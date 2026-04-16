package main

import (
	"os"
	"strings"

	"github.com/cainseing/drop-cli/internal/app"
	"github.com/cainseing/drop-cli/internal/config"
	"github.com/cainseing/drop-cli/internal/update"
)

func main() {
	update.CheckForUpdates(config.Version, false)

	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "drop_") {
		token := strings.TrimPrefix(os.Args[1], "drop_")
		os.Args = []string{os.Args[0], "get", token}
	}

	rootCmd := app.CreateRootCmd()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
