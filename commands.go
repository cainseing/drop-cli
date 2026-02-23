package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func createRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "drop [secret]",
		Short: "Secure, zero-knowledge, secret sharing CLI",
		Args:  cobra.MaximumNArgs(1),
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Run: func(cmd *cobra.Command, args []string) {
			var input []byte
			stat, _ := os.Stdin.Stat()

			if (stat.Mode() & os.ModeCharDevice) == 0 {
				input, _ = io.ReadAll(os.Stdin)
			} else if len(args) > 0 {
				input = []byte(args[0])
			}

			if len(input) == 0 {
				cmd.Help()
				return
			}

			handleCreateCommand(input, ttl, reads)
		},
	}

	rootCmd.Flags().IntVarP(&ttl, "ttl", "t", 5, "Expiry in minutes")
	rootCmd.Flags().IntVarP(&reads, "reads", "r", 1, "Maximum number of times drop can be read")
	rootCmd.AddCommand(createGetCmd(), createPurgeCmd(), createVersionCommand())
	return rootCmd
}

func createGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get [token]",
		Short: "Fetch drop",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			handleGetCommand(args[0])
		},
	}
}

func createPurgeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "purge [token]",
		Short: "Purge a drop",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			handlePurgeCommand(args[0])
		},
	}
}

// func createConfigCmd() *cobra.Command {
// 	var configCmd = &cobra.Command{Use: "config", Short: "Manage CLI settings"}
// 	var setUrlCmd = &cobra.Command{
// 		Use:  "url [url]",
// 		Args: cobra.ExactArgs(1),
// 		Run: func(cmd *cobra.Command, args []string) {
// 			viper.Set("api_url", args[0])
// 			home, _ := os.UserHomeDir()
// 			configPath := filepath.Join(home, ".config", "drop")
// 			os.MkdirAll(configPath, 0755)
// 			if err := viper.WriteConfig(); err != nil {
// 				viper.SafeWriteConfig()
// 			}
// 			fmt.Printf("%s%s[config]%s API Uplink set to: %s%s%s\n", ColorDim, ColorGreen, ColorReset, ColorBold, args[0], ColorReset)
// 		},
// 	}
// 	configCmd.AddCommand(setUrlCmd)
// 	return configCmd
// }

func createVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "View running version",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			PrintInfo("Version", fmt.Sprintf("\n%s\n", dim.Render(version)))
		},
	}
}
