package main

import (
	"os"
	"path/filepath"

	"github.com/cainseing/drop-cli/internal/app"
	"github.com/cainseing/drop-cli/internal/config"
	"github.com/cainseing/drop-cli/internal/update"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

func main() {
	cobra.OnInitialize(func() {
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				// If we can't find home, we skip config loading rather than failing
				return
			}
			viper.AddConfigPath(filepath.Join(home, ".config", "drop"))
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
		}
		viper.SetDefault("api_url", "https://api.getdrop.dev/")
		viper.SetEnvPrefix("drop")
		viper.AutomaticEnv()
		_ = viper.ReadInConfig()
	})

	rootCmd := app.CreateRootCmd()

	update.CheckForUpdates(config.Version, false)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
