package safefile

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDigestBoundRegularRequiresExactSignedSizeAndContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "signed-component.exe")
	body := []byte("exact signed component fixture")
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(body)
	got, n, err := DigestBoundRegular(path, int64(len(body)))
	if err != nil || n != int64(len(body)) || got != hex.EncodeToString(sum[:]) {
		t.Fatalf("correct signed file rejected: hash=%q size=%d err=%v", got, n, err)
	}
	for _, expected := range []int64{-1, 0, int64(len(body) - 1), int64(len(body) + 1), int64(^uint64(0) >> 1)} {
		if _, _, err := DigestBoundRegular(path, expected); err == nil {
			t.Fatalf("incorrect signed size %d was accepted", expected)
		}
	}
	if err := os.WriteFile(path, []byte(strings.Repeat("x", len(body))), 0o600); err != nil {
		t.Fatal(err)
	}
	got, n, err = DigestBoundRegular(path, int64(len(body)))
	if err != nil || n != int64(len(body)) || got == hex.EncodeToString(sum[:]) {
		t.Fatalf("signed SHA-256 check cannot detect same-size modification: hash=%q n=%d err=%v", got, n, err)
	}
}

func TestDigestBoundRegularRejectsSymlinkDirectoryAndOversizedSparseFile(t *testing.T) {
	root := t.TempDir()
	if _, _, err := DigestBoundRegular(root, 1); err == nil {
		t.Fatal("directory accepted as executable")
	}
	path := filepath.Join(root, "signed-release")
	if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "symlink")
	if err := os.Symlink(path, link); err != nil {
		t.Logf("symlink unavailable: %v", err)
	} else if _, _, err := DigestBoundRegular(link, 1); err == nil {
		t.Fatal("symlink accepted as a signed executable")
	}
	if err := os.Truncate(path, 8<<20); err != nil {
		t.Skipf("sparse file unavailable: %v", err)
	}
	if _, _, err := DigestBoundRegular(path, 1<<20); err == nil || !strings.Contains(err.Error(), "changed size") {
		t.Fatalf("oversized sparse file bypassed size bound: %v", err)
	}
}
