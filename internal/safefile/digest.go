package safefile

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

// DigestBoundRegular hashes a cached artifact only after verifying its expected
// size and stable file identity on the open descriptor. The read is bounded by
// the signed size, not by a potentially corrupted on-disk size. Consumers
// compare the returned digest with their independently verified signed binding.
func DigestBoundRegular(path string, expectedSize int64) (string, int64, error) {
	if expectedSize <= 0 || expectedSize == int64(^uint64(0)>>1) {
		return "", 0, errors.New("invalid signed Creator artifact size")
	}
	entry, err := os.Lstat(path)
	if err != nil {
		return "", 0, err
	}
	if !entry.Mode().IsRegular() || entry.Size() != expectedSize {
		return "", 0, fmt.Errorf("signed Creator artifact is unsafe or changed size: %s", path)
	}
	file, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil {
		return "", 0, err
	}
	if !opened.Mode().IsRegular() || !os.SameFile(entry, opened) || opened.Size() != expectedSize {
		return "", 0, fmt.Errorf("signed Creator artifact changed identity or size while opening: %s", path)
	}
	digest := sha256.New()
	n, err := io.Copy(digest, io.LimitReader(file, expectedSize+1))
	if err != nil {
		return "", 0, err
	}
	if n != expectedSize {
		return "", 0, fmt.Errorf("signed Creator artifact changed size while reading: %s", path)
	}
	return hex.EncodeToString(digest.Sum(nil)), n, nil
}
