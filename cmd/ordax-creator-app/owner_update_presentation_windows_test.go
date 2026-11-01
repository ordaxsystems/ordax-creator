//go:build windows

package main

import (
	"strings"
	"testing"
)

func validOwnerUpdateManifestForTest() ownerUpdateManifest {
	return ownerUpdateManifest{
		Schema:       "prototype-ordax.creator-owner-bundle-update/1",
		Channel:      "owner-prototype",
		Version:      "0.6.2-dev",
		SourceCommit: strings.Repeat("b", 40),
		Artifact:     "OrdaX-Creator-Owner-Prototype.zip",
		DownloadURL:  "https://github.com/ordaxsystems/ordax-os/releases/download/creator-owner-prototype/OrdaX-Creator-Owner-Prototype.zip",
		SHA256:       strings.Repeat("c", 64),
		Size:         889240014, // Regression: official 2026-11-01 owner bundle exceeds the old 512 MiB limit.
	}
}

func TestOfficialOwnerBundleSizeAcceptedWithoutRelaxingManifestPolicy(t *testing.T) {
	for _, tc := range []struct {
		name string
		size int64
		wantValid bool
	}{
		{"published-size-over-512-mib", 889240014, true},
		{"at-maximum", maxOwnerUpdateBundleBytes, true},
		{"zero", 0, false},
		{"negative", -1, false},
		{"over-maximum", maxOwnerUpdateBundleBytes + 1, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			manifest := validOwnerUpdateManifestForTest()
			manifest.Size = tc.size
			if err := validateOwnerUpdateManifest(manifest); (err == nil) != tc.wantValid {
				t.Fatalf("bundle size=%d valid=%v error=%v, want valid=%v", tc.size, err == nil, err, tc.wantValid)
			}
		})
	}
	for _, mutate := range []func(*ownerUpdateManifest){
		func(m *ownerUpdateManifest) { m.Channel = "development" },
		func(m *ownerUpdateManifest) { m.SourceCommit = "untrusted" },
		func(m *ownerUpdateManifest) { m.SHA256 = strings.Repeat("z", 64) },
		func(m *ownerUpdateManifest) { m.DownloadURL = "http://github.com/ordaxsystems/ordax-os/releases/download/creator-owner-prototype/OrdaX-Creator-Owner-Prototype.zip" },
		func(m *ownerUpdateManifest) { m.AutomaticInAppUpdate = true },
	} {
		manifest := validOwnerUpdateManifestForTest()
		mutate(&manifest)
		if err := validateOwnerUpdateManifest(manifest); err == nil {
			t.Fatalf("unsafe owner manifest accepted: %+v", manifest)
		}
	}
}

func TestOwnerManualUpdateMessageSeparatesInstalledAndAvailableForBothLocales(t *testing.T) {
	previous := currentCreatorLocale()
	t.Cleanup(func() { setCreatorLocale(string(previous)) })
	candidate := updateUIState{
		Available: true,
		Version: "0.6.2-dev",
		SourceCommit: strings.Repeat("b", 40),
	}
	for _, locale := range []creatorLocale{creatorLocalePTBR, creatorLocaleENUS} {
		setCreatorLocale(string(locale))
		got := ownerUpdateManualMessage("0.6.1-dev", strings.Repeat("a", 40), candidate)
		for _, phrase := range []string{"0.6.1-dev (aaaaaaaa)", "0.6.2-dev (bbbbbbbb)", "OrdaX-Creator-Owner-Prototype.zip"} {
			if !strings.Contains(got, phrase) {
				t.Fatalf("locale %s lost version/bundle identity %q: %q", locale, phrase, got)
			}
		}
		if strings.Contains(got, "{installed}") || strings.Contains(got, "{available}") {
			t.Fatalf("locale %s has unexpanded version placeholders: %q", locale, got)
		}
	}

	candidate.Error = "lineage check failed"
	got := ownerUpdateManualMessage("0.6.1-dev", strings.Repeat("a", 40), candidate)
	if strings.Contains(got, "bbbbbbbb") {
		t.Fatalf("failed lineage was presented as validated available update: %q", got)
	}
}

func TestAvailableUpdateNeverReplacesInstalledVersionLabel(t *testing.T) {
	oldLabel := creatorVersionText
	updateMu.Lock()
	previous := currentUpdateState
	currentUpdateState = updateUIState{
		Available: true,
		Version: "0.6.2-dev",
		SourceCommit: strings.Repeat("b", 40),
		Manual: false,
	}
	updateMu.Unlock()
	t.Cleanup(func() {
		creatorVersionText = oldLabel
		updateMu.Lock()
		currentUpdateState = previous
		updateMu.Unlock()
	})

	installed := "OrdaX Creator • 0.6.1-dev • aaaaaaaa"
	creatorVersionText = installed
	renderUpdateDone() // No window, network, download or physical device access.
	if creatorVersionText != installed {
		t.Fatalf("advertised update masqueraded as installed: got %q, want %q", creatorVersionText, installed)
	}
}
