package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetConfigPath(t *testing.T) {
	tmpHome := t.TempDir()
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)
	if err := os.Setenv("HOME", tmpHome); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}

	got := GetConfigPath()
	want := filepath.Join(tmpHome, ".config", "drop", "config.json")
	if got != want {
		t.Fatalf("GetConfigPath() = %q, want %q", got, want)
	}
}

func TestLoadUserConfigMissingFile(t *testing.T) {
	tmpHome := t.TempDir()
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)
	if err := os.Setenv("HOME", tmpHome); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}

	cfg := LoadUserConfig()
	if cfg.Username != "" || cfg.Provider != "" || !cfg.UpdateCheck.IsZero() {
		t.Fatalf("expected empty config for missing file, got %+v", cfg)
	}
}

func TestLoadUserConfigInvalidJSON(t *testing.T) {
	tmpHome := t.TempDir()
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)
	if err := os.Setenv("HOME", tmpHome); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}

	path := GetConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}
	if err := os.WriteFile(path, []byte("{invalid json}"), 0o644); err != nil {
		t.Fatalf("failed to write invalid config file: %v", err)
	}

	cfg := LoadUserConfig()
	if cfg.Username != "" || cfg.Provider != "" || !cfg.UpdateCheck.IsZero() {
		t.Fatalf("expected empty config for invalid JSON, got %+v", cfg)
	}
}

func TestSaveAndLoadUserConfig(t *testing.T) {
	tmpHome := t.TempDir()
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)
	if err := os.Setenv("HOME", tmpHome); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}

	cfg := &UserConfig{
		Username:    "alice",
		Provider:    "github",
		UpdateCheck: time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC),
	}

	if err := SaveUserConfig(cfg); err != nil {
		t.Fatalf("SaveUserConfig() error: %v", err)
	}

	loaded := LoadUserConfig()
	if loaded.Username != cfg.Username {
		t.Fatalf("Username mismatch: want %q, got %q", cfg.Username, loaded.Username)
	}
	if loaded.Provider != cfg.Provider {
		t.Fatalf("Provider mismatch: want %q, got %q", cfg.Provider, loaded.Provider)
	}
	if !loaded.UpdateCheck.Equal(cfg.UpdateCheck) {
		t.Fatalf("UpdateCheck mismatch: want %v, got %v", cfg.UpdateCheck, loaded.UpdateCheck)
	}
}

func TestPublicConfigConstants(t *testing.T) {
	if Version == "" {
		t.Fatal("expected Version to be set")
	}
	if ProtocolVersion != "2" {
		t.Fatalf("expected ProtocolVersion to be %q, got %q", "2", ProtocolVersion)
	}
	if MaxBlobSize != 1*1024*1024 {
		t.Fatalf("expected MaxBlobSize to be %d, got %d", 1*1024*1024, MaxBlobSize)
	}
	if MaxTTLMinutes != 10080 {
		t.Fatalf("expected MaxTTLMinutes to be %d, got %d", 10080, MaxTTLMinutes)
	}
	if ApiURL != "https://api.getdrop.dev/" {
		t.Fatalf("expected ApiURL to be %q, got %q", "https://api.getdrop.dev/", ApiURL)
	}
}
