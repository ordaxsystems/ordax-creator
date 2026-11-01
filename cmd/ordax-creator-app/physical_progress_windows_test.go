//go:build windows

package main

import (
	"encoding/json"
	"math"
	"math/big"
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


func TestPhysicalProgressPercentExactForLargeByteCounters(t *testing.T) {
	for _, tc := range []struct {
		name       string
		completed  int64
		total      int64
		start      int
		span       int
	}{
		{"writing-half-small", 50, 100, 15, 45},
		{"verifying-half-small", 50, 100, 62, 28},
		{"writing-half-maxint64", math.MaxInt64 / 2, math.MaxInt64, 15, 45},
		{"verifying-half-maxint64", math.MaxInt64 / 2, math.MaxInt64, 62, 28},
		{"writing-full-maxint64", math.MaxInt64, math.MaxInt64, 15, 45},
		{"verifying-full-maxint64", math.MaxInt64, math.MaxInt64, 62, 28},
		{"writing-near-full-maxint64", math.MaxInt64 - 1, math.MaxInt64, 15, 45},
		{"verifying-near-full-maxint64", math.MaxInt64 - 1, math.MaxInt64, 62, 28},
		{"writing-quarter-maxint64", math.MaxInt64 / 4, math.MaxInt64, 15, 45},
		{"writing-near-zero", 1, math.MaxInt64, 15, 45},
		{"verifying-one-byte", 1, 1, 62, 28},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Independent arbitrary-precision reference, not susceptible to
			// the signed int64 multiplication overflow in the former GUI.
			product := new(big.Int).Mul(big.NewInt(tc.completed), big.NewInt(int64(tc.span)))
			want := tc.start + int(product.Div(product, big.NewInt(tc.total)).Int64())
			got := physicalProgressPercent(tc.completed, tc.total, tc.start, tc.span)
			if got != want || got < tc.start || got > tc.start+tc.span {
				t.Fatalf("scaled progress=%d, exact=%d, range=[%d,%d]", got, want, tc.start, tc.start+tc.span)
			}
		})
	}
	for _, tc := range []struct {
		completed int64
		total     int64
		span      int
	}{
		{0, 0, 45},
		{0, math.MaxInt64, 45},
		{-1, math.MaxInt64, 28},
		{math.MaxInt64, 0, 45},
		{math.MaxInt64, 1, 0},
	} {
		if got := physicalProgressPercent(tc.completed, tc.total, 15, tc.span); got != 15 {
			t.Fatalf("zero/unknown/invalid total must remain at phase start: %d", got)
		}
	}
}

func TestPhysicalProgressPresentationScalesHugeValidSnapshots(t *testing.T) {
	for _, tc := range []struct {
		phase string
		completed int64
		want int
	}{
		{"writing", math.MaxInt64 / 2, 37},
		{"writing", math.MaxInt64 - 1, 59},
		{"writing", math.MaxInt64, 60},
		{"verifying", math.MaxInt64 / 2, 75},
		{"verifying", math.MaxInt64 - 1, 89},
		{"verifying", math.MaxInt64, 90},
	} {
		doc := physicalProgressDocument{Schema: physicalProgressSchema, Phase: tc.phase, CompletedBytes: tc.completed, TotalBytes: math.MaxInt64}
		_, _, got, ok := physicalProgressPresentation(doc)
		if !ok || got != tc.want {
			t.Fatalf("phase %q bytes %d/%d gave %d (ok=%v), want %d", tc.phase, tc.completed, doc.TotalBytes, got, ok, tc.want)
		}
	}
}

func TestReadPhysicalProgressRejectsIncoherentByteCounters(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	for _, tc := range []struct {
		name string
		completed int64
		total int64
		valid bool
	}{
		{"starting-zero", 0, 0, true},
		{"unknown-total-but-positive-bytes", 1, 0, false},
		{"more-completed-than-total", 100, 99, false},
		{"negative-completed", -1, 10, false},
		{"negative-total", 0, -1, false},
		{"large-valid", math.MaxInt64 - 1, math.MaxInt64, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body, err := json.Marshal(physicalProgressDocument{
				Schema: physicalProgressSchema, Phase: "writing",
				CompletedBytes: tc.completed, TotalBytes: tc.total,
			})
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, body, 0o600); err != nil {
				t.Fatal(err)
			}
			_, err = readPhysicalProgress(path)
			if (err == nil) != tc.valid {
				t.Fatalf("counters %d/%d accepted=%v, want=%v: %v", tc.completed, tc.total, err == nil, tc.valid, err)
			}
		})
	}
}
