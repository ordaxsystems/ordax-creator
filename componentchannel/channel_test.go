package componentchannel

import (
	"bytes"
	"crypto/ed25519"
	"errors"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testTrust(t *testing.T) ([]byte, ed25519.PrivateKey, string) {
	t.Helper()
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	trust := TrustAnchor{Schema: TrustSchema, KeyID: "ordax-prototype-release-v1", PublicKeyB64: base64.StdEncoding.EncodeToString(public)}
	data, err := json.Marshal(trust)
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(data)
	return data, private, hex.EncodeToString(digest[:])
}

func validManifest() Manifest {
	return Manifest{
		Schema:            ManifestSchema,
		Purpose:           Purpose,
		SourceRepository:  SourceRepository,
		SourceCommit:      "1111111111111111111111111111111111111111",
		Version:           "1.0.0",
		ReleaseSequence:   1,
		CreatedFromRecipe: Recipe,
		Bundle: Bundle{
			URL:    "https://github.com/ordaxsystems/ordax-os/releases/download/creator-components/ordax-creator-components-windows-amd64.zip",
			SHA256: strings.Repeat("a", 64),
			Size:   1234,
		},
		File: FileBinding{Name: expectedFile, SHA256: strings.Repeat("b", 64), Size: 123},
	}
}

func signedEnvelope(t *testing.T, manifest Manifest, private ed25519.PrivateKey) []byte {
	t.Helper()
	payload, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	envelope := Envelope{Schema: EnvelopeSchema, Payload: payload, Signature: ed25519.Sign(private, payload), KeyID: "ordax-prototype-release-v1"}
	data, err := json.Marshal(envelope)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func installManifest(t *testing.T, root string, manifest Manifest, body []byte) Manifest {
	t.Helper()
	directory := filepath.Join(root, "versions", manifest.SourceCommit)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, expectedFile), body, 0o700); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(body)
	manifest.File.SHA256 = hex.EncodeToString(digest[:])
	manifest.File.Size = int64(len(body))
	return manifest
}

func TestVerifyEnvelopeAcceptsCanonicalPurpose(t *testing.T) {
	trust, private, trustSHA := testTrust(t)
	manifest, err := VerifyEnvelope(signedEnvelope(t, validManifest(), private), trust, trustSHA)
	if err != nil {
		t.Fatal(err)
	}
	if manifest.Purpose != Purpose || manifest.ReleaseSequence != 1 {
		t.Fatalf("unexpected manifest: %+v", manifest)
	}
}

func TestVerifyEnvelopeRejectsWrongPurposeSignatureAndTrust(t *testing.T) {
	trust, private, trustSHA := testTrust(t)
	wrongPurpose := validManifest()
	wrongPurpose.Purpose = "creator-physical-windows-amd64"
	if _, err := VerifyEnvelope(signedEnvelope(t, wrongPurpose, private), trust, trustSHA); err == nil || !strings.Contains(err.Error(), "identity/purpose") {
		t.Fatalf("wrong purpose error=%v", err)
	}

	envelopeBytes := signedEnvelope(t, validManifest(), private)
	var envelope Envelope
	if err := json.Unmarshal(envelopeBytes, &envelope); err != nil {
		t.Fatal(err)
	}
	envelope.Signature[0] ^= 0xff
	tampered, _ := json.Marshal(envelope)
	if _, err := VerifyEnvelope(tampered, trust, trustSHA); err == nil || !strings.Contains(err.Error(), "signature") {
		t.Fatalf("tampered signature error=%v", err)
	}

	if _, err := VerifyEnvelope(envelopeBytes, trust, strings.Repeat("0", 64)); err == nil || !strings.Contains(err.Error(), "pinned") {
		t.Fatalf("wrong trust hash error=%v", err)
	}
}

func TestVerifyEnvelopeRejectsInvalidSequence(t *testing.T) {
	trust, private, trustSHA := testTrust(t)
	manifest := validManifest()
	manifest.ReleaseSequence = 0
	if _, err := VerifyEnvelope(signedEnvelope(t, manifest, private), trust, trustSHA); err == nil || !strings.Contains(err.Error(), "release_sequence") {
		t.Fatalf("invalid sequence error=%v", err)
	}
}

func TestVerifyInstalledRehashesAndRejectsExtraFiles(t *testing.T) {
	root := t.TempDir()
	manifest := installManifest(t, root, validManifest(), []byte("inspection-backend"))
	directory := filepath.Join(root, "versions", manifest.SourceCommit)
	if err := VerifyInstalled(directory, manifest); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "extra.exe"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := VerifyInstalled(directory, manifest); err == nil || !strings.Contains(err.Error(), "exactly") {
		t.Fatalf("extra file error=%v", err)
	}
	if err := os.Remove(filepath.Join(directory, "extra.exe")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, expectedFile), []byte("tampered"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := VerifyInstalled(directory, manifest); err == nil || !strings.Contains(err.Error(), "changed") {
		t.Fatalf("tamper error=%v", err)
	}
}

func TestRollbackPolicyRejectsOlderAndSequenceReuse(t *testing.T) {
	root := t.TempDir()
	trust, private, trustSHA := testTrust(t)
	current := validManifest()
	current.SourceCommit = "2222222222222222222222222222222222222222"
	current.ReleaseSequence = 2
	current = installManifest(t, root, current, []byte("current"))
	if err := os.WriteFile(filepath.Join(root, currentEnvelopeName), signedEnvelope(t, current, private), 0o600); err != nil {
		t.Fatal(err)
	}

	older := validManifest()
	older.SourceCommit = "1111111111111111111111111111111111111111"
	older.ReleaseSequence = 1
	if err := checkRollback(root, older, trust, trustSHA); err == nil || !strings.Contains(err.Error(), "rollback rejected") {
		t.Fatalf("older error=%v", err)
	}

	reused := validManifest()
	reused.SourceCommit = "3333333333333333333333333333333333333333"
	reused.ReleaseSequence = 2
	if err := checkRollback(root, reused, trust, trustSHA); err == nil || !strings.Contains(err.Error(), "reused") {
		t.Fatalf("reuse error=%v", err)
	}

	newer := validManifest()
	newer.SourceCommit = "4444444444444444444444444444444444444444"
	newer.ReleaseSequence = 3
	if err := checkRollback(root, newer, trust, trustSHA); err != nil {
		t.Fatalf("newer candidate rejected: %v", err)
	}
}

func TestSignedComponentRejectsSymlinkedVersionsRootOfflineAndOnline(t *testing.T) {
	root := t.TempDir()
	external := t.TempDir()
	if err := os.Symlink(external, filepath.Join(root, "versions")); err != nil {
		t.Skipf("symlinks unavailable on this host: %v", err)
	}
	trust, private, trustSHA := testTrust(t)
	manifest := validManifest()
	body := []byte("signed-inspection-component")
	digest := sha256.Sum256(body)
	manifest.File.SHA256 = hex.EncodeToString(digest[:])
	manifest.File.Size = int64(len(body))
	dir := filepath.Join(external, manifest.SourceCommit)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, expectedFile), body, 0o700); err != nil {
		t.Fatal(err)
	}
	envelope := signedEnvelope(t, manifest, private)
	if err := os.WriteFile(filepath.Join(root, currentEnvelopeName), envelope, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Current(root, trust, trustSHA); err == nil || !strings.Contains(err.Error(), "versions cache") {
		t.Fatalf("offline signed component followed redirected versions root: %v", err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(envelope)
	}))
	defer server.Close()
	if _, _, err := AcquireCached(server.Client(), root, server.URL, trust, trustSHA); err == nil {
		t.Fatal("online signed component cache followed redirected versions root")
	}
}


type signedComponentTestTransport func(*http.Request) (*http.Response, error)

func (fn signedComponentTestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

// The UI uses a separate, strictly verified Current fallback on network
// failures. An unconfirmed signed component may never be reported as active.
func TestSignedComponentCachedActivationRequiresSecondStableEnvelope(t *testing.T) {
	for _, mode := range []string{
		"stable", "network-failure", "http-failure",
		"changed-signed-artifact", "invalid-second-signature", "pointer-write-failure",
	} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			trust, private, trustSHA := testTrust(t)
			manifest := installManifest(t, root, validManifest(), []byte("valid-signed-inspection-component"))
			first := signedEnvelope(t, manifest, private)
			second := first
			switch mode {
			case "changed-signed-artifact":
				rebound := manifest
				rebound.File.SHA256 = strings.Repeat("a", 64)
				second = signedEnvelope(t, rebound, private)
			case "invalid-second-signature":
				var altered Envelope
				if err := json.Unmarshal(first, &altered); err != nil {
					t.Fatal(err)
				}
				altered.Signature[0] ^= 0xff
				var err error
				second, err = json.Marshal(altered)
				if err != nil {
					t.Fatal(err)
				}
			}
			requests := 0
			client := &http.Client{Transport: signedComponentTestTransport(func(req *http.Request) (*http.Response, error) {
				requests++
				if req.URL.Host != "github.com" || req.URL.Path != "/ordaxsystems/ordax-os/releases/download/creator-components/creator-component-envelope.json" || requests > 2 {
					return nil, errors.New("unexpected signed component test URL or third request")
				}
				if requests == 2 {
					switch mode {
					case "network-failure":
						return nil, errors.New("synthetic signed component network failure")
					case "http-failure":
						return &http.Response{StatusCode: http.StatusServiceUnavailable, Body: io.NopCloser(bytes.NewReader(nil)), Header: make(http.Header), Request: req}, nil
					case "pointer-write-failure":
						if err := os.Mkdir(filepath.Join(root, currentEnvelopeName), 0o700); err != nil {
							return nil, err
						}
					}
				}
				data := first
				if requests == 2 {
					data = second
				}
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(data)), Header: make(http.Header), Request: req}, nil
			})}
			got, changed, err := AcquireCached(client, root, DefaultEnvelopeURL, trust, trustSHA)
			if requests != 2 {
				t.Fatalf("expected two independent envelope fetches, got %d, err=%v", requests, err)
			}
			if mode == "stable" {
				if err != nil || changed || got.SourceCommit != manifest.SourceCommit {
					t.Fatalf("stable signed component changed=%t installed=%+v err=%v", changed, got, err)
				}
				if _, err := Current(root, trust, trustSHA); err != nil {
					t.Fatalf("valid component cache pointer did not activate: %v", err)
				}
			} else {
				if err == nil || changed || got.SourceCommit != "" {
					t.Fatalf("unconfirmed signed component was returned as success: changed=%t installed=%+v err=%v", changed, got, err)
				}
				if _, err := Current(root, trust, trustSHA); err == nil {
					t.Fatal("unconfirmed signed component became offline current")
				}
			}
		})
	}
}
