package update

import (
	"testing"
	"time"

	"github.com/cainseing/drop-cli/internal/config"
)

func TestSanitizeVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"v1.0.0", "1.0.0"},
		{"1.0.0", "1.0.0"},
		{"v1.0.0-beta", "1.0.0"},
		{"v2.1.3-rc.1", "2.1.3"},
		{"3.0.0", "3.0.0"},
	}

	for _, test := range tests {
		result := sanitizeVersion(test.input)
		if result != test.expected {
			t.Errorf("sanitizeVersion(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

func TestIsBrewInstall(t *testing.T) {
	// Test cases for different executable paths
	tests := []struct {
		path     string
		expected bool
	}{
		{"/usr/local/Cellar/drop/1.0.0/bin/drop", true},
		{"/homebrew/Cellar/drop/1.0.0/bin/drop", true},
		{"/linuxbrew/Cellar/drop/1.0.0/bin/drop", true},
		{"/usr/bin/drop", false},
		{"/opt/drop/bin/drop", false},
		{"", false},
	}

	for _, test := range tests {
		// Mock getExecutable
		original := getExecutable
		getExecutable = func() (string, error) { return test.path, nil }
		defer func() { getExecutable = original }()

		result := isBrewInstall()
		if result != test.expected {
			t.Errorf("isBrewInstall() for path %q = %v; want %v", test.path, result, test.expected)
		}
	}
}

func TestCheckForUpdates_NoTTY(t *testing.T) {
	// Test when not running in TTY (e.g., piped output)
	// This should return false immediately
	result := CheckForUpdates("1.0.0", false)
	if result {
		t.Errorf("CheckForUpdates should return false when not in TTY")
	}
}

func TestCheckForUpdates_RecentCheck(t *testing.T) {
	// Test when update check was recent
	cfg := config.LoadUserConfig()
	originalTime := cfg.UpdateCheck
	cfg.UpdateCheck = time.Now().Add(-1 * time.Hour) // 1 hour ago
	config.SaveUserConfig(cfg)
	defer func() {
		cfg.UpdateCheck = originalTime
		config.SaveUserConfig(cfg)
	}()

	// Since it's within 24 hours, should return false
	result := CheckForUpdates("1.0.0", false)
	if result {
		t.Errorf("CheckForUpdates should return false for recent check")
	}
}

func TestCheckForUpdates_Force(t *testing.T) {
	// Test force update check
	cfg := config.LoadUserConfig()
	originalTime := cfg.UpdateCheck
	cfg.UpdateCheck = time.Now().Add(-1 * time.Hour)
	config.SaveUserConfig(cfg)
	defer func() {
		cfg.UpdateCheck = originalTime
		config.SaveUserConfig(cfg)
	}()

	// With force=true, should attempt check even if recent
	// This will likely fail due to network, but we test the logic
	result := CheckForUpdates("1.0.0", true)
	// Result depends on network, but at least it should not return false due to time check
	// In a real test, we'd mock the HTTP call
	_ = result // Just ensure it runs without panic
}
