package identity

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type Provider string

const (
	GitHub Provider = "github"
	GitLab Provider = "gitlab"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

func FetchKeys(provider Provider, username string) ([]ssh.PublicKey, error) {
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	if !isValidUsername(username) {
		return nil, fmt.Errorf("invalid username: contains disallowed characters")
	}

	url, err := buildURL(provider, username)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to reach %s: %w", provider, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user '%s' not found on %s", username, provider)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("provider returned error: %s", resp.Status)
	}

	return parseKeys(resp.Body)
}

func isValidUsername(username string) bool {
	// Allow alphanumeric, dash, and underscore only
	pattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return pattern.MatchString(username)
}

func buildURL(provider Provider, username string) (string, error) {
	if provider == GitHub {
		return fmt.Sprintf("https://github.com/%s.keys", username), nil
	}

	if provider == GitLab {
		return fmt.Sprintf("https://gitlab.com/%s.keys", username), nil
	}

	return "", fmt.Errorf("unsupported identity provider: %s", provider)
}

func parseKeys(body io.Reader) ([]ssh.PublicKey, error) {
	var keys []ssh.PublicKey
	scanner := bufio.NewScanner(body)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(line))
		if err != nil {
			continue
		}

		keys = append(keys, pubKey)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed reading stream: %w", err)
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("no valid SSH keys found in provider profile")
	}

	return keys, nil
}
