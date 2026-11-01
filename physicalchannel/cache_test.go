package physicalchannel

import (
	"crypto/ed25519"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCurrentRequiresSignedCachedEnvelopeAndFreshBindings(t *testing.T) {
	root := t.TempDir()
	trustBytes, private, trustSHA := testTrust(t)
	commit := "0123456789abcdef0123456789abcdef01234567"
	candidate := sequencedCandidate(t, root, commit, 1)
	cacheSignedCandidate(t, root, candidate, private)

	installed, err := Current(root, trustBytes, trustSHA)
	if err != nil {
		t.Fatal(err)
	}
	if installed.SourceCommit != commit || installed.Directory != candidate.Directory {
		t.Fatalf("unexpected cached candidate: %+v", installed)
	}
	sequence, err := ReleaseSequence(installed)
	if err != nil {
		t.Fatal(err)
	}
	if sequence != 1 {
		t.Fatalf("release sequence=%d want=1", sequence)
	}

	if err := os.WriteFile(filepath.Join(candidate.Directory, "provenance.json"), []byte("tampered"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Current(root, trustBytes, trustSHA); err == nil || !strings.Contains(err.Error(), "changed") {
		t.Fatalf("tampered cached candidate unexpectedly accepted: %v", err)
	}
}

func TestCurrentRejectsTamperedCachedEnvelope(t *testing.T) {
	trustBytes, private, trustSHA := testTrust(t)
	manifest := validManifest()
	payload, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	signature := ed25519.Sign(private, payload)
	signature[0] ^= 0xff
	envelopeBytes, err := json.Marshal(Envelope{
		Schema: EnvelopeSchema,
		Payload: payload,
		Signature: signature,
		KeyID: "ordax-prototype-release-v1",
	})
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, currentEnvelopeName), envelopeBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Current(root, trustBytes, trustSHA); err == nil || !strings.Contains(err.Error(), "signature") {
		t.Fatalf("tampered cached envelope unexpectedly accepted: %v", err)
	}
}

func TestPhysicalChannelRejectsRedirectedVersionsRootForCachedAndOnlineCandidate(t *testing.T) {
	root := t.TempDir()
	trust, private, trustSHA := testTrust(t)
	candidate := sequencedCandidate(t, root, "0123456789abcdef0123456789abcdef01234567", 1)
	cacheSignedCandidate(t, root, candidate, private)
	if _, err := Current(root, trust, trustSHA); err != nil {
		t.Fatalf("legitimate cached candidate rejected before symlink: %v", err)
	}
	versions := filepath.Join(root, "versions")
	backing := filepath.Join(t.TempDir(), "versions")
	if err := os.Rename(versions, backing); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(backing, versions); err != nil {
		t.Skipf("symlinks unavailable on this host: %v", err)
	}
	if _, err := Current(root, trust, trustSHA); err == nil || !strings.Contains(err.Error(), "versions cache") {
		t.Fatalf("physical offline cache followed redirected versions root: %v", err)
	}
	envelope := signedEnvelope(t, validManifest(), private)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(envelope)
	}))
	defer server.Close()
	if _, _, err := Acquire(server.Client(), root, server.URL, trust, trustSHA); err == nil || !strings.Contains(err.Error(), "versions cache") {
		t.Fatalf("physical online acquisition followed redirected versions root: %v", err)
	}
	if _, _, err := AcquireCached(server.Client(), root, server.URL, trust, trustSHA); err == nil {
		t.Fatal("physical signed cache promotion accepted redirected versions root")
	}
}
