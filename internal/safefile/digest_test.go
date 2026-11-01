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

func TestDigestInventoryRegularAllowsEmptyFilesWithoutWeakeningSignedPolicy(t *testing.T) {
	root := t.TempDir()
	empty := filepath.Join(root, "empty-file")
	if err := os.WriteFile(empty, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	expected := sha256.Sum256(nil)
	digest, size, err := DigestInventoryRegular(empty, 0)
	if err != nil || size != 0 || digest != hex.EncodeToString(expected[:]) {
		t.Fatalf("empty Dev cache file unexpectedly rejected: sha=%s size=%d err=%v", digest, size, err)
	}
	if _, _, err := DigestBoundRegular(empty, 0); err == nil {
		t.Fatal("signed artifact unexpectedly accepted empty file")
	}
	for _, size := range []int64{-1, 1, int64(^uint64(0) >> 1)} {
		if _, _, err := DigestInventoryRegular(empty, size); err == nil {
			t.Fatalf("invalid Dev cache size %d accepted", size)
		}
	}
	if _, _, err := DigestInventoryRegular(root, 0); err == nil {
		t.Fatal("directory accepted as empty regular cache file")
	}
	link := filepath.Join(root, "linked-empty")
	if err := os.Symlink(empty, link); err == nil {
		if _, _, err := DigestInventoryRegular(link, 0); err == nil {
			t.Fatal("symlink to empty cache file accepted")
		}
	} else {
		t.Logf("symlinks unavailable on host: %v", err)
	}
}

func TestDigestInventoryRegularRejectsReplacedOrUnexpectedlySizedFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "receipt-source")
	body := []byte("legitimate development receipt")
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
	expected := sha256.Sum256(body)
	got, size, err := DigestInventoryRegular(path, int64(len(body)))
	if err != nil || got != hex.EncodeToString(expected[:]) || size != int64(len(body)) {
		t.Fatalf("development cache digest mismatch: sha=%s size=%d err=%v", got, size, err)
	}
	if err := os.Truncate(path, 8<<20); err != nil {
		t.Fatal(err)
	}
	if _, _, err := DigestInventoryRegular(path, int64(len(body))); err == nil || !strings.Contains(err.Error(), "changed size") {
		t.Fatalf("oversized sparse cache file passed bounded identity check: %v", err)
	}
}
