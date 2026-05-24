package app

import (
	"fmt"
	"io"
	"os"

	"github.com/cainseing/drop-cli/internal/api"
	"github.com/cainseing/drop-cli/internal/config"
	"github.com/cainseing/drop-cli/internal/output"
	"github.com/cainseing/drop-cli/internal/update"
	"github.com/spf13/cobra"
)

var (
	ttl        int
	reads      int
	signed     bool
	shouldCopy bool
)

func CreateRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "drop [secret]",
		Short: "Zero-knowledge secret sharing",
		Args:  cobra.MaximumNArgs(1),
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Run: func(cmd *cobra.Command, args []string) {
			var input []byte
			stat, err := os.Stdin.Stat()

			if err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
				input, _ = io.ReadAll(os.Stdin)
			} else if len(args) > 0 {
				input = []byte(args[0])
			}

			if len(input) == 0 {
				_ = cmd.Help()
				return
			}

			api.HandleCreateCommand(input, ttl, reads, signed, shouldCopy)
		},
	}

	rootCmd.Flags().IntVarP(&ttl, "ttl", "t", 5, "Expiry in minutes")
	rootCmd.Flags().IntVarP(&reads, "reads", "r", 1, "Maximum number of times drop can be read")
	rootCmd.Flags().BoolVarP(&signed, "signed", "s", false, "Sign the drop")
	rootCmd.Flags().BoolVarP(&shouldCopy, "copy", "c", false, "Copy token to clipboard")

	rootCmd.AddCommand(
		createGetCmd(),
		createPurgeCmd(),
		createVersionCommand(),
		createIdentityCmd(),
	)

	applyHelpFunc(rootCmd)

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
				output.PrintError("Unsupported provider. Use 'github' or 'gitlab'.", nil)
				return
			}

			cfg := config.LoadUserConfig()
			cfg.Username = username
			cfg.Provider = provider

			if err := config.SaveUserConfig(cfg); err != nil {
				output.PrintError("Failed to save identity", err)
				return
			}

			output.PrintSuccess(fmt.Sprintf("Signed as @%s via %s", username, provider), "Identity saved")
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
			if update.CheckForUpdates(config.Version, true) {
				return
			}

			fmt.Fprintln(os.Stderr)
			output.PrintPropertyToStderr("VERSION", config.Version)
			fmt.Fprintln(os.Stderr)
		},
	}
}
