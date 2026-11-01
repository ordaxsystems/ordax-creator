package cachepath

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureVersionsCreatesAndReusesCanonicalDirectory(t *testing.T) {
	root := t.TempDir()
	first, err := EnsureVersions(root)
	if err != nil {
		t.Fatal(err)
	}
	if first != filepath.Join(root, "versions") {
		t.Fatalf("unexpected canonical versions directory: %q", first)
	}
	again, err := EnsureVersions(root)
	if err != nil || again != first {
		t.Fatalf("reusing real directory failed: %q, %v", again, err)
	}
	found, err := Versions(root)
	if err != nil || found != first {
		t.Fatalf("offline verification failed: %q, %v", found, err)
	}
}

func TestVersionsMissingDoesNotCreateCache(t *testing.T) {
	root := t.TempDir()
	if _, err := Versions(root); !os.IsNotExist(err) {
		t.Fatalf("missing cache must return os.ErrNotExist, got %v", err)
	}
	if _, err := os.Lstat(filepath.Join(root, "versions")); !os.IsNotExist(err) {
		t.Fatalf("offline read unexpectedly created cache: %v", err)
	}
}

func TestVersionsRejectsSymlinkAndRegularFile(t *testing.T) {
	for _, tc := range []struct {
		name string
		setup func(*testing.T, string)
	}{
		{"file", func(t *testing.T, target string) {
			if err := os.WriteFile(target, []byte("not a directory"), 0o600); err != nil {
				t.Fatal(err)
			}
		}},
		{"symlink", func(t *testing.T, target string) {
			if err := os.Symlink(t.TempDir(), target); err != nil {
				t.Skipf("symlinks unavailable on this host: %v", err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			tc.setup(t, filepath.Join(root, "versions"))
			if _, err := Versions(root); err == nil {
				t.Fatal("unsafe versions cache accepted by offline reader")
			}
			if _, err := EnsureVersions(root); err == nil {
				t.Fatal("unsafe versions cache accepted by installer")
			}
		})
	}
}
