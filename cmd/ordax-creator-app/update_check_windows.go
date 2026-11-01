//go:build windows

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	ownerUpdateManifestURL = "https://github.com/ordaxsystems/ordax-os/releases/download/creator-owner-prototype/creator-owner-update.json"
	// The signed-source owner bundle includes system images and may exceed
	// 512 MiB. This is only a bounded read-only manifest check, never a
	// download/installation entitlement or a physical-write authorization.
	maxOwnerUpdateBundleBytes int64 = 2 << 30
)

var (
	updateMu           sync.Mutex
	currentUpdateState updateUIState
)

type ownerUpdateManifest struct {
	Schema               string `json:"$schema"`
	Channel              string `json:"channel"`
	Version              string `json:"version"`
	SourceCommit         string `json:"source_commit"`
	Artifact             string `json:"artifact"`
	DownloadURL          string `json:"download_url"`
	SHA256               string `json:"sha256"`
	Size                 int64  `json:"size"`
	AutomaticInAppUpdate bool   `json:"automatic_in_app_update"`
}

type updateUIState struct {
	Checking     bool
	Checked      bool
	Manual       bool
	Available    bool
	Version      string
	SourceCommit string
	Error        string
}

func validLowerHexString(value string, size int) bool {
	if len(value) != size || value != strings.ToLower(value) {
		return false
	}
	for _, ch := range value {
		if (ch < '0' || ch > '9') && (ch < 'a' || ch > 'f') {
			return false
		}
	}
	return true
}

func validateOwnerUpdateManifest(manifest ownerUpdateManifest) error {
	if manifest.Schema != "prototype-ordax.creator-owner-bundle-update/1" {
		return errors.New("Creator update manifest schema is incompatible")
	}
	if manifest.Channel != "owner-prototype" {
		return errors.New("Creator update channel is unexpected")
	}
	if !validLowerHexString(manifest.SourceCommit, 40) {
		return errors.New("Creator update source_commit is invalid")
	}
	if !validCreatorVersion(manifest.Version) {
		return errors.New("Creator update version is invalid")
	}
	if manifest.Artifact != "OrdaX-Creator-Owner-Prototype.zip" {
		return errors.New("Creator update artifact is unexpected")
	}
	if !validLowerHexString(manifest.SHA256, 64) || manifest.Size <= 0 || manifest.Size > maxOwnerUpdateBundleBytes {
		return errors.New("Creator update bundle binding is invalid")
	}
	parsed, err := url.Parse(manifest.DownloadURL)
	if err != nil || parsed.Scheme != "https" || parsed.Host != "github.com" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return errors.New("Creator update URL is invalid")
	}
	if parsed.Path != "/ordaxsystems/ordax-os/releases/download/creator-owner-prototype/OrdaX-Creator-Owner-Prototype.zip" {
		return errors.New("Creator update bundle is outside the OrdaX release channel")
	}
	if manifest.AutomaticInAppUpdate {
		return errors.New("automatic in-app bundle update is not supported by this Creator")
	}
	return nil
}

func fetchOwnerUpdateManifest() (ownerUpdateManifest, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest(http.MethodGet, ownerUpdateManifestURL+"?ordax_nocache="+fmt.Sprint(time.Now().UnixNano()), nil)
	if err != nil {
		return ownerUpdateManifest{}, err
	}
	req.Header.Set("User-Agent", "OrdaX-Creator-Owner-Updater/4")
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := client.Do(req)
	if err != nil {
		return ownerUpdateManifest{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ownerUpdateManifest{}, fmt.Errorf("Creator update source returned HTTP %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10+1))
	if err != nil {
		return ownerUpdateManifest{}, err
	}
	if len(data) == 0 || len(data) > 64<<10 {
		return ownerUpdateManifest{}, errors.New("Creator update manifest has invalid size")
	}
	var manifest ownerUpdateManifest
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&manifest); err != nil {
		return ownerUpdateManifest{}, err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return ownerUpdateManifest{}, errors.New("Creator update manifest contains trailing data")
	}
	if err := validateOwnerUpdateManifest(manifest); err != nil {
		return ownerUpdateManifest{}, err
	}
	return manifest, nil
}

func ownerUpdateIsForward(currentSource, candidateSource string) (bool, error) {
	if currentSource == candidateSource {
		return false, nil
	}
	if !validLowerHexString(currentSource, 40) || !validLowerHexString(candidateSource, 40) {
		return false, errors.New("Creator update lineage cannot be validated")
	}
	compareURL := "https://api.github.com/repos/ordaxsystems/ordax-os/compare/" + currentSource + "..." + candidateSource + "?per_page=1"
	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest(http.MethodGet, compareURL, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", "OrdaX-Creator-Owner-Updater/4")
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Creator anti-downgrade validation returned HTTP %d", resp.StatusCode)
	}
	var comparison struct {
		Status   string `json:"status"`
		AheadBy  int    `json:"ahead_by"`
		BehindBy int    `json:"behind_by"`
	}
	decoder := json.NewDecoder(io.LimitReader(resp.Body, 16<<20))
	if err := decoder.Decode(&comparison); err != nil {
		return false, fmt.Errorf("decode Creator lineage comparison: %w", err)
	}
	switch comparison.Status {
	case "ahead":
		if comparison.AheadBy > 0 && comparison.BehindBy == 0 {
			return true, nil
		}
	case "behind":
		if comparison.BehindBy > 0 && comparison.AheadBy == 0 {
			return false, nil
		}
	case "identical":
		return false, nil
	}
	return false, fmt.Errorf("published Creator history diverges from the current version; update blocked (status=%s)", comparison.Status)
}

func beginUpdateCheck(manual bool) {
	if writeInProgress() {
		return
	}
	_, currentSource, owner := ownerPrototypeBuildInfo()
	if !owner {
		if manual {
			messageBox(creatorT(msgUpdateDevChannelBody), windowTitle, mbOK|mbIconInformation)
		}
		return
	}

	updateMu.Lock()
	if currentUpdateState.Checking {
		updateMu.Unlock()
		return
	}
	currentUpdateState = updateUIState{Checking: true, Manual: manual}
	updateMu.Unlock()
	renderUpdateUI()

	go func() {
		manifest, err := fetchOwnerUpdateManifest()
		state := updateUIState{Manual: manual}
		if err != nil {
			state.Error = err.Error()
		} else {
			state.Version = manifest.Version
			state.SourceCommit = manifest.SourceCommit
			forward, forwardErr := ownerUpdateIsForward(currentSource, manifest.SourceCommit)
			if forwardErr != nil && manifest.SourceCommit != currentSource {
				state.Error = forwardErr.Error()
			} else {
				state.Available = forward
				state.Checked = true
			}
		}
		updateMu.Lock()
		currentUpdateState = state
		updateMu.Unlock()
		procPostMessageW.Call(mainWindow, wmAppUpdateDone, 0, 0)
	}()
}

// The installed build identity belongs to the executable, not the owner
// release's read-only update notice. Keep both versions and source commits
// distinct in the manual instructions; no in-app update is performed.
func ownerUpdateManualMessage(installedVersion, installedSource string, state updateUIState) string {
	body := creatorT(msgUpdateBundleManualBody)
	if !state.Available || state.Error != "" ||
		!validCreatorVersion(installedVersion) || !validLowerHexString(installedSource, 40) ||
		!validCreatorVersion(state.Version) || !validLowerHexString(state.SourceCommit, 40) {
		return body
	}
	details := creatorTValues(msgUpdateVersionComparison, map[string]string{
		"installed":       installedVersion,
		"installed_commit": shortSourceCommit(installedSource),
		"available":       state.Version,
		"available_commit": shortSourceCommit(state.SourceCommit),
	})
	return details + "\n\n" + body
}

func showOwnerUpdateManualMessage(state updateUIState) {
	installedVersion, installedSource, _ := ownerPrototypeBuildInfo()
	messageBox(ownerUpdateManualMessage(installedVersion, installedSource, state),
		creatorT(msgUpdateAvailableTitle), mbOK|mbIconInformation)
}

// The dedicated update control represents the outcome of the *check*, not
// an installation. Completed checks may be repeated, and an unsuccessful
// network/trust check must offer retry without declaring the build current.
func ownerUpdateAction(state updateUIState) (creatorMessageID, bool) {
	switch {
	case state.Checking:
		return msgUpdateChecking, false
	case state.Error != "":
		return msgUpdateRetryAction, true
	case state.Available:
		return msgUpdateAvailableAction, true
	case state.Checked:
		return msgUpdateCurrentAction, true
	default:
		return msgUpdateControl, true
	}
}

// An asynchronous manual result must never open a modal dialog over a
// physical write or USB refresh. The result remains in the button state so
// the user can review or recheck when the system is idle.
func ownerUpdateMayShowDialog(state updateUIState, writing, refreshing bool) bool {
	return state.Manual && !state.Checking && !writing && !refreshing
}

func renderUpdateUI() {
	if updateButton == 0 {
		return
	}
	updateMu.Lock()
	state := currentUpdateState
	updateMu.Unlock()
	if writeInProgress() {
		enable(updateButton, false)
		return
	}
	action, enabled := ownerUpdateAction(state)
	setText(updateButton, creatorT(action))
	enable(updateButton, enabled)
}

func renderUpdateDone() {
	updateMu.Lock()
	state := currentUpdateState
	updateMu.Unlock()
	renderUpdateUI()
	if !ownerUpdateMayShowDialog(state, writeInProgress(), refreshInProgress) {
		return
	}

	if state.Available && state.Error == "" {
		// Never replace the installed build label with the candidate version.
		// The button owns availability, and the manual dialog owns the
		// comparison; installation is unchanged until externally updated.
		showOwnerUpdateManualMessage(state)
		return
	}
	if state.Error != "" {
		messageBox(creatorT(msgUpdateFailedBody), windowTitle, mbOK|mbIconError)
		return
	}
	if state.Checked {
		messageBox(creatorT(msgUpdateCurrentBody), windowTitle, mbOK|mbIconInformation)
	}
}

func handleUpdateButton() {
	if writeInProgress() {
		return
	}
	updateMu.Lock()
	state := currentUpdateState
	updateMu.Unlock()
	if state.Available && state.Error == "" {
		showOwnerUpdateManualMessage(state)
		return
	}
	beginUpdateCheck(true)
}
