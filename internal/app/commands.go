package app

import (
	"fmt"
	"io"
	"os"

	"github.com/cainseing/drop-cli/internal/api"
	"github.com/cainseing/drop-cli/internal/config"
	"github.com/cainseing/drop-cli/internal/update"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	ttl    int
	reads  int
	signed bool
)

func CreateRootCmd() *cobra.Command {
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

			api.HandleCreateCommand(input, ttl, reads, signed)
		},
	}

	rootCmd.Flags().IntVarP(&ttl, "ttl", "t", 5, "Expiry in minutes")
	rootCmd.Flags().IntVarP(&reads, "reads", "r", 1, "Maximum number of times drop can be read")
	rootCmd.Flags().BoolVarP(&signed, "signed", "s", false, "Sign the drop")

	rootCmd.AddCommand(
		createGetCmd(),
		createPurgeCmd(),
		createVersionCommand(),
		createIdentityCmd(),
	)

	return rootCmd
}

func createGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get [token]",
		Short: "Fetch drop",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			api.HandleGetCommand(args[0])
		},
	}
}

func createPurgeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "purge [token]",
		Short: "Purge a drop",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			api.HandlePurgeCommand(args[0])
		},
	}
}

func createIdentityCmd() *cobra.Command {
	var identityCmd = &cobra.Command{
		Use:   "identity",
		Short: "Manage your signing identity",
	}

	var setCmd = &cobra.Command{
		Use:   "set [github|gitlab] [username]",
		Short: "Assign a public profile to your signed drops",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			provider, username := args[0], args[1]

			if provider != "github" && provider != "gitlab" {
				pterm.Error.Println("Unsupported provider. Use 'github' or 'gitlab'.")
				return
			}

			// Load existing config to preserve other fields (like LastUpdateCheck)
			cfg := config.LoadUserConfig()
			cfg.Username = username
			cfg.Provider = provider

			if err := config.SaveUserConfig(cfg); err != nil {
				pterm.Error.Printf("Failed to save identity: %v\n", err)
				return
			}

			pterm.Success.Printf("Identity saved! Drops will now be signed as @%s via %s\n", username, provider)
		},
	}

	identityCmd.AddCommand(setCmd)
	return identityCmd
}

func createVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "View running version",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			result := update.CheckForUpdates(config.Version, true)

			if result {
				return
			}

			message := fmt.Sprintf("%s", pterm.Gray(config.Version))

			fmt.Println()
			pterm.DefaultBox.
				WithTitle(pterm.LightYellow(" Version ")).
				WithBoxStyle(pterm.NewStyle(pterm.FgLightYellow)).
				Println(message)
			fmt.Println()
		},
	}
}
