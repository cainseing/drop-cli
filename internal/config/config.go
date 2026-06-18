package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

var Version = "v0.5.1-beta"
const ProtocolVersion = "2"
const MaxBlobSize = 1 * 1024 * 1024 // 1MB
const MaxTTLMinutes = 10080         // 7 Days
var ApiURL = "https://api.getdrop.dev/"

type UserConfig struct {
	Username    string    `json:"username"`
	Provider    string    `json:"provider"`
	UpdateCheck time.Time `json:"update_check"`
}

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "drop", "config.json")
}

func LoadUserConfig() *UserConfig {
	path := GetConfigPath()
	data, err := os.ReadFile(path) // #nosec G304 -- path is a fixed location under user home dir
	if err != nil {
		return &UserConfig{}
	}

	var cfg UserConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return &UserConfig{}
	}

	return &cfg
}

func SaveUserConfig(cfg *UserConfig) error {
	path := GetConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
