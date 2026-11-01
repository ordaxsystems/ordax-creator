//go:build windows

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPhysicalProgressPresentationUsesByteProgressForWriteAndVerify(t *testing.T) {
	write := physicalProgressDocument{Schema: physicalProgressSchema, Phase: "writing", CompletedBytes: 50, TotalBytes: 100}
	_, _, writePercent, ok := physicalProgressPresentation(write)
	if !ok {
		t.Fatal("writing phase must be recognized")
	}
	if writePercent != 37 {
		t.Fatalf("writing percent=%d want=37", writePercent)
	}

	verify := physicalProgressDocument{Schema: physicalProgressSchema, Phase: "verifying", CompletedBytes: 50, TotalBytes: 100}
	_, _, verifyPercent, ok := physicalProgressPresentation(verify)
	if !ok {
		t.Fatal("verifying phase must be recognized")
	}
	if verifyPercent != 76 {
		t.Fatalf("verifying percent=%d want=76", verifyPercent)
	}
}

func TestPhysicalProgressPresentationCompletesAtOneHundred(t *testing.T) {
	_, _, percent, ok := physicalProgressPresentation(physicalProgressDocument{Schema: physicalProgressSchema, Phase: "complete", CompletedBytes: 1, TotalBytes: 1})
	if !ok || percent != 100 {
		t.Fatalf("complete progress=(%d,%v) want=(100,true)", percent, ok)
	}
}

func TestGuidedPhysicalProgressCopyKeepsLocalizedSafetyInstructionVisible(t *testing.T) {
	statusInput := creatorT(msgProgressWritingStatus)
	hintInput := creatorTValues(msgProgressWritingBytes, map[string]string{
		"completed": "1.0",
		"total":     "2.0",
	})
	status, hint := guidedPhysicalProgressCopy(statusInput, hintInput)
	if !strings.Contains(status, creatorT(msgCreatingEyebrow)) {
		t.Fatalf("status=%q must carry localized guided creation stage", status)
	}
	if !strings.Contains(hint, creatorT(msgCreatingTitle)) {
		t.Fatalf("hint=%q must keep localized removal warning visible", hint)
	}
	if !strings.Contains(hint, hintInput) {
		t.Fatalf("hint=%q must preserve localized backend byte progress %q", hint, hintInput)
	}
}

func TestPhysicalProgressSafetyCompositionWorksInBundledLocales(t *testing.T) {
	previous := currentCreatorLocale()
	t.Cleanup(func() { setCreatorLocale(string(previous)) })

	for _, locale := range []creatorLocale{creatorLocalePTBR, creatorLocaleENUS} {
		setCreatorLocale(string(locale))
		statusInput := creatorT(msgProgressWritingStatus)
		hintInput := creatorTValues(msgProgressWritingBytes, map[string]string{"completed": "1.0", "total": "2.0"})
		status, hint := guidedPhysicalProgressCopy(statusInput, hintInput)
		if !strings.Contains(status, creatorMessageFor(locale, msgCreatingEyebrow, nil)) || !strings.Contains(hint, creatorMessageFor(locale, msgCreatingTitle, nil)) {
			t.Fatalf("locale %s lost guided safety copy: status=%q hint=%q", locale, status, hint)
		}
	}
}

func TestReadPhysicalProgressUsesBoundedRegularFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "progress.json")
	body := []byte(`{"$schema":"prototype-ordax.creator-physical-progress/1","phase":"writing","completed_bytes":1,"total_bytes":2}`)
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
	document, err := readPhysicalProgress(path)
	if err != nil || document.Phase != "writing" {
		t.Fatalf("safe progress snapshot rejected: %+v, %v", document, err)
	}
	if err := os.Truncate(path, 4097); err != nil {
		t.Fatal(err)
	}
	if _, err := readPhysicalProgress(path); err == nil {
		t.Fatal("oversized progress evidence was accepted")
	}
}

func TestReadPhysicalProgressRejectsRedirectedFile(t *testing.T) {
	root := t.TempDir()
	realFile := filepath.Join(root, "real-progress.json")
	if err := os.WriteFile(realFile, []byte(`{"$schema":"prototype-ordax.creator-physical-progress/1","phase":"complete","completed_bytes":2,"total_bytes":2}`), 0o600); err != nil {
		t.Fatal(err)
	}
	linked := filepath.Join(root, "progress.json")
	if err := os.Symlink(realFile, linked); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if _, err := readPhysicalProgress(linked); err == nil {
		t.Fatal("symlinked progress evidence was accepted")
	}
}
