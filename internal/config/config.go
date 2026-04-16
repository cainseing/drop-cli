package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const Version = "v0.4.1-beta"
const ProtocolVersion = "1"
const MaxBlobSize = 1 * 1024 * 1024 // 1MB
const MaxTTLMinutes = 10080         // 7 Days
const ApiURL = "https://api.getdrop.dev/"

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
	data, err := os.ReadFile(path)
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
	os.MkdirAll(filepath.Dir(path), 0755)

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
