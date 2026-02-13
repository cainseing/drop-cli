package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const version = "0.1.0-beta"
const protocolVersion = "1"

var (
	ttl     int
	cfgFile string
	reads   int
)

const MaxBlobSize = 1 * 1024 * 1024 // 1MB
const MaxTTLMinutes = 10080         // 7 Days

var (
	accent      = lipgloss.NewStyle().Foreground(lipgloss.Color("#d47fd4"))
	dim         = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	errorPrefix = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3333")).Bold(true)
	highlight   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFD1")).Bold(true)
	secret      = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	success     = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFD1")).Bold(true)
	errorLabel  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	errorText   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8888"))
)

func main() {
	cobra.OnInitialize(initConfig)

	rootCmd := createRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initConfig() {
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
}

func printError(label string, err error) {
	fmt.Fprintf(os.Stderr, "\n%s %s", errorPrefix.Render("×"), errorPrefix.Render("ERROR"))

	if label != "" && err != nil {
		fmt.Fprintf(os.Stderr, " %s %s\n", dim.Render("─"), errorLabel.Render(label))
		fmt.Fprintf(os.Stderr, "  %s\n\n", errorText.Render(err.Error()))
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "\n  %s\n\n", errorText.Render(err.Error()))
	} else {
		fmt.Fprintf(os.Stderr, "\n  %s\n\n", errorText.Render(label))
	}
}
