package update

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func validManifestForTest() Manifest {
	return Manifest{
		Schema:       ChannelSchemaV2,
		Channel:      DevelopmentChannel,
		Version:      "dev-0123456789ab",
		SourceCommit: "0123456789abcdef0123456789abcdef01234567",
		Payload: Payload{
			URL:    "https://github.com/ordaxsystems/ordax-os/releases/download/creator-dev/ordax-creator-dev-windows-amd64.zip",
			SHA256: strings.Repeat("a", 64),
			Size:   1234,
		},
		Entrypoints: Entrypoints{
			Inspect: "1-Inspect-OrdaXUSB.cmd",
			Status:  "ordax-creator-physical-test.exe",
		},
	}
}

func legacyManifestBytesForTest(trust string) []byte {
	m := validManifestForTest()
	legacy := legacyManifestV1{
		Schema: m.Schema, Channel: m.Channel, Version: m.Version,
		SourceCommit: m.SourceCommit, Payload: m.Payload,
		Entrypoints: legacyEntrypointsV1{Inspect: m.Entrypoints.Inspect, Trust: trust, Status: m.Entrypoints.Status},
	}
	legacy.Schema = ChannelSchemaV1
	data, _ := json.Marshal(legacy)
	return data
}

func TestValidateManifestAcceptsBoundDevelopmentChannel(t *testing.T) {
	if err := ValidateManifest(validManifestForTest()); err != nil {
		t.Fatal(err)
	}
}

func TestDecodeManifestAcceptsLegacyProtocolWithoutExposingTrustCapability(t *testing.T) {
	manifest, err := decodeManifest(legacyManifestBytesForTest("2-Legacy-Protocol-Compatibility.cmd"))
	if err != nil {
		t.Fatal(err)
	}
	if manifest.Schema != ChannelSchemaV1 {
		t.Fatalf("schema = %q", manifest.Schema)
	}
	if manifest.Entrypoints.Inspect == "" || manifest.Entrypoints.Status == "" {
		t.Fatal("consumer entrypoints were lost")
	}
}

func TestDecodeManifestRejectsUnsafeLegacyCompatibilityEntrypoint(t *testing.T) {
	if _, err := decodeManifest(legacyManifestBytesForTest("../evil.cmd")); err == nil {
		t.Fatal("unsafe legacy compatibility entrypoint unexpectedly accepted")
	}
}

func TestValidateManifestRejectsExternalPayload(t *testing.T) {
	m := validManifestForTest()
	m.Payload.URL = "https://example.invalid/creator.zip"
	if err := ValidateManifest(m); err == nil {
		t.Fatal("external payload URL unexpectedly accepted")
	}
}

func TestValidateManifestRejectsTraversalEntrypoint(t *testing.T) {
	m := validManifestForTest()
	m.Entrypoints.Inspect = "../evil.cmd"
	if err := ValidateManifest(m); err == nil {
		t.Fatal("path traversal entrypoint unexpectedly accepted")
	}
}

func TestDecodeV2RejectsPublisherTrustField(t *testing.T) {
	m := validManifestForTest()
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	text := strings.TrimSuffix(string(data), "}") + `,"trust":"forbidden"}`
	if _, err := decodeManifest([]byte(text)); err == nil {
		t.Fatal("publisher trust field unexpectedly accepted in v2 consumer manifest")
	}
}

func writeZip(t *testing.T, entries map[string]string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "payload.zip")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	archive := zip.NewWriter(file)
	for name, body := range entries {
		entry, err := archive.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := entry.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := archive.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestExtractPayloadRequiresAllEntrypoints(t *testing.T) {
	m := validManifestForTest()
	zipPath := writeZip(t, map[string]string{
		m.Entrypoints.Inspect: "@echo inspect\r\n",
	})
	dest := filepath.Join(t.TempDir(), "out")
	if err := os.Mkdir(dest, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := extractPayload(zipPath, dest, m); err == nil || !strings.Contains(err.Error(), "status entrypoint") {
		t.Fatalf("missing status entrypoint error = %v", err)
	}
}

func TestExtractPayloadRejectsTraversal(t *testing.T) {
	path := filepath.Join(t.TempDir(), "payload.zip")
	var data bytes.Buffer
	archive := zip.NewWriter(&data)
	entry, err := archive.Create("../evil.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = entry.Write([]byte("evil"))
	if err := archive.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(t.TempDir(), "out")
	if err := os.Mkdir(dest, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := extractPayload(path, dest, validManifestForTest()); err == nil || !strings.Contains(err.Error(), "unsafe archive path") {
		t.Fatalf("traversal error = %v", err)
	}
}

func TestCurrentRejectsDirectoryOutsideVersionSlot(t *testing.T) {
	root := t.TempDir()
	current := Installed{
		Schema:       "prototype-ordax.creator-installed/1",
		Channel:      DevelopmentChannel,
		Version:      "dev-0123456789ab",
		SourceCommit: "0123456789abcdef0123456789abcdef01234567",
		Directory:    filepath.Join(root, "elsewhere"),
	}
	data, _ := json.Marshal(current)
	if err := os.WriteFile(filepath.Join(root, "current.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Current(root); err == nil {
		t.Fatal("external current directory unexpectedly accepted")
	}
}


func TestSafeRelativeRejectsAmbiguousWindowsPathsOnAllHosts(t *testing.T) {
	for _, name := range []string{
		"C:payload.exe", `C:\payload.exe`, `dir\payload.exe`,
		"payload.exe:stream", "nested/NUL", "nested/CON.txt",
		"nested/com1.exe", "nested/LPT9", "nested/COM¹",
		"trailing.", "nested/trailing ", "a//b", "./payload.exe",
		"../payload.exe", "/absolute/payload.exe",
	} {
		if safeRelative(name) {
			t.Errorf("unsafe path %q accepted", name)
		}
	}
	for _, name := range []string{"payload.zip", "bin/payload.exe", "1-Inspect-OrdaXUSB.cmd"} {
		if !safeRelative(name) {
			t.Errorf("canonical relative path %q rejected", name)
		}
	}
}

func TestValidateManifestRejectsNonCanonicalCreatorDevAssetURLs(t *testing.T) {
	base := "https://github.com/ordaxsystems/ordax-os/releases/download/creator-dev/"
	for _, raw := range []string{
		base,
		base + "payload.exe",
		base + "subfolder/payload.zip",
		base + "payload.zip?redirect=true",
		base + "payload.zip?",
		base + "payload%2Ezip",
		base + "NUL.zip",
		"https://github.com.evil.example/ordaxsystems/ordax-os/releases/download/creator-dev/payload.zip",
	} {
		m := validManifestForTest()
		m.Payload.URL = raw
		if err := ValidateManifest(m); err == nil {
			t.Errorf("noncanonical release asset URL %q accepted", raw)
		}
	}
}

func TestExtractPayloadRejectsCaseFoldedDuplicatePaths(t *testing.T) {
	zipPath := writeZip(t, map[string]string{
		"payload/Component.exe": "first",
		"payload/component.exe": "second",
	})
	destination := filepath.Join(t.TempDir(), "out")
	if err := os.Mkdir(destination, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := extractPayload(zipPath, destination, validManifestForTest()); err == nil ||
		!strings.Contains(err.Error(), "duplicate archive path") {
		t.Fatalf("case-folded duplicate error = %v", err)
	}
}
