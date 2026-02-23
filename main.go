package main

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const version = "v0.2.1-beta"
const protocolVersion = "1"
const MaxBlobSize = 1 * 1024 * 1024 // 1MB
const MaxTTLMinutes = 10080         // 7 Days

var (
	ttl     int
	cfgFile string
	reads   int
)

var (
	dim        = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	secret     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	errorLabel = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	errorText  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8888"))
)

func main() {
	cobra.OnInitialize(func() {
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			home, _ := os.UserHomeDir()
			viper.AddConfigPath(filepath.Join(home, ".config", "drop"))
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
		}
		viper.SetDefault("api_url", "https://api.getdrop.dev/")
		viper.SetEnvPrefix("drop")
		viper.AutomaticEnv()
		_ = viper.ReadInConfig()
	})

	rootCmd := createRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
