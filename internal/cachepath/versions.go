// Package cachepath defines the Creator's shared verified cache-directory
// boundary. It owns no release, signature, artifact, or authorization policy.
package cachepath

import (
	"errors"
	"os"
	"path/filepath"
)

// Versions returns the existing versions directory only if its final path
// component is a real directory, never a symlink or a regular file.
func Versions(root string) (string, error) {
	dir := filepath.Join(root, "versions")
	info, err := os.Lstat(dir)
	if err != nil {
		return "", err
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return "", errors.New("Creator versions cache must be a real directory")
	}
	return dir, nil
}

// EnsureVersions creates the versions cache when absent, then enforces the
// same boundary used by the offline readers. MkdirAll alone is insufficient:
// it succeeds if an existing versions path is a symlink to a directory.
func EnsureVersions(root string) (string, error) {
	dir := filepath.Join(root, "versions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return Versions(root)
}
