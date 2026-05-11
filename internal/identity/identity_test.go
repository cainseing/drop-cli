package identity

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

type rewriteTransport struct {
	base   http.RoundTripper
	target string
}

func (t rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	proxyReq := req.Clone(req.Context())
	if proxyReq.URL.Host == "github.com" || proxyReq.URL.Host == "gitlab.com" {
		proxyReq.URL.Scheme = "http"
		proxyReq.URL.Host = t.target
	}
	return t.base.RoundTrip(proxyReq)
}

func TestIsValidUsername(t *testing.T) {
	valid := []string{"alice", "bob-123", "user_name"}
	invalid := []string{"alice@example.com", "bob$", "user name", ""}

	for _, username := range valid {
		if !isValidUsername(username) {
			t.Fatalf("expected valid username for %q", username)
		}
	}

	for _, username := range invalid {
		if isValidUsername(username) {
			t.Fatalf("expected invalid username for %q", username)
		}
	}
}

func TestBuildURL(t *testing.T) {
	url, err := buildURL(GitHub, "alice")
	if err != nil {
		t.Fatalf("buildURL failed: %v", err)
	}
	if url != "https://github.com/alice.keys" {
		t.Fatalf("unexpected GitHub URL: %s", url)
	}

	url, err = buildURL(GitLab, "alice")
	if err != nil {
		t.Fatalf("buildURL failed: %v", err)
	}
	if url != "https://gitlab.com/alice.keys" {
		t.Fatalf("unexpected GitLab URL: %s", url)
	}

	_, err = buildURL(Provider("unsupported"), "alice")
	if err == nil {
		t.Fatalf("expected unsupported provider error")
	}
}

func TestFetchKeysValidation(t *testing.T) {
	_, err := FetchKeys(GitHub, "")
	if err == nil {
		t.Fatal("expected error for empty username")
	}
}

func TestFetchKeysUnsupportedProvider(t *testing.T) {
	_, err := FetchKeys(Provider("unsupported"), "alice")
	if err == nil {
		t.Fatal("expected error for unsupported provider")
	}
}

func TestParseKeys(t *testing.T) {
	body := `# comment line

ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKjgQYAaN0w4tACiObT3Z9xBr1eK1Eh435D/3b8cRSB8 test@example.com
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyno0QZuFziVR50S13NjDuYMZMwTmB51Y6yxOhww/4RBl4alds+Y/DRgB1+JXqZmCekkGq1W/OO9bGiEFjzp9zdrXqk7wihR2a+oVQdVE6CCw1hC1E7A200AiknGSZxWFhakPsKNetp8Yo4LDPxSyr/PAVepjj2wpwTYeeoFXkVu5nMuEqOI8WeBkx2v8XqGGgIZhWYjkkCdudkOmsbDQ1RZjHmT86bX5W9MHS14D9yeMV2IjGs6QW8S5fsIF7VWFlvRaMwsE4Oypg/Xt0uAZdg4rPOdeIbVja5nwYSFc4wO29td71mr3oqDwlUDjKeVLEZBn/tqvgclm8neRF3VpP cainseing@MacBook-Pro-M2.localdomain
invalid-key
`

	keys, err := parseKeys(strings.NewReader(body))
	if err != nil {
		t.Fatalf("parseKeys failed: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("expected 2 valid keys, got %d", len(keys))
	}
	if keys[0].Type() != ssh.KeyAlgoED25519 {
		t.Fatalf("expected first key type %q, got %q", ssh.KeyAlgoED25519, keys[0].Type())
	}
	if keys[1].Type() != ssh.KeyAlgoRSA {
		t.Fatalf("expected second key type %q, got %q", ssh.KeyAlgoRSA, keys[1].Type())
	}
}

func TestParseKeysNoValidEntries(t *testing.T) {
	body := "# comment\nnot-a-key\n"
	_, err := parseKeys(strings.NewReader(body))
	if err == nil {
		t.Fatal("expected error when no valid SSH keys are present")
	}
}

func TestFetchKeysFromServer(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/alice.keys" {
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
		fmt.Fprintln(w, "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKjgQYAaN0w4tACiObT3Z9xBr1eK1Eh435D/3b8cRSB8 test@example.com")
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	origClient := client
	client = &http.Client{
		Timeout:   5 * time.Second,
		Transport: rewriteTransport{base: http.DefaultTransport, target: ts.Listener.Addr().String()},
	}
	defer func() { client = origClient }()

	keys, err := FetchKeys(GitHub, "alice")
	if err != nil {
		t.Fatalf("FetchKeys failed: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(keys))
	}
	if keys[0].Type() != ssh.KeyAlgoED25519 {
		t.Fatalf("expected ed25519 key, got %s", keys[0].Type())
	}
}
