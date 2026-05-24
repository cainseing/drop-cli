package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	helpLogoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#10b981"))
	helpTitleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#10b981")).Bold(true)
	helpSectionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#10b981"))
	helpNameStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	helpDimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
)

var logo = "      _\n" +
	"   __| |_ __ ___  _ __\n" +
	"  / _` | '__/ _ \\| '_ \\\n" +
	" | (_| | | | (_) | |_) |\n" +
	"  \\__,_|_|  \\___/| .__/\n" +
	"                 |_|"

func applyHelpFunc(cmd *cobra.Command) {
	cmd.SetHelpFunc(renderHelp)
	for _, sub := range cmd.Commands() {
		applyHelpFunc(sub)
	}
}

func renderHelp(cmd *cobra.Command, _ []string) {
	fmt.Println()

	if !cmd.HasParent() {
		for _, line := range strings.Split(logo, "\n") {
			fmt.Printf("  %s\n", helpLogoStyle.Render(line))
		}
		fmt.Println()
		fmt.Printf("  %s\n", helpDimStyle.Render(cmd.Short))
	} else {
		fmt.Printf("  %s  %s\n", helpTitleStyle.Render(cmd.Name()), helpDimStyle.Render(cmd.Short))
	}

	fmt.Println()

	fmt.Printf("  %s\n", helpSectionStyle.Render("Usage"))
	fmt.Printf("    %s\n", helpDimStyle.Render(cmd.UseLine()))
	if cmd.HasAvailableSubCommands() {
		fmt.Printf("    %s\n", helpDimStyle.Render(cmd.CommandPath()+" [command]"))
	}
	fmt.Println()

	if cmd.HasAvailableSubCommands() {
		fmt.Printf("  %s\n", helpSectionStyle.Render("Commands"))

		maxLen := 0
		var available []*cobra.Command
		for _, c := range cmd.Commands() {
			if c.IsAvailableCommand() {
				available = append(available, c)
				if len(c.Name()) > maxLen {
					maxLen = len(c.Name())
				}
			}
		}
		for _, c := range available {
			name := helpNameStyle.Render(fmt.Sprintf("%-*s", maxLen, c.Name()))
			fmt.Printf("    %s  %s\n", name, helpDimStyle.Render(c.Short))
		}
		fmt.Println()
	}

	if cmd.HasAvailableFlags() {
		fmt.Printf("  %s\n", helpSectionStyle.Render("Flags"))
		for _, line := range strings.Split(strings.TrimRight(cmd.Flags().FlagUsages(), "\n"), "\n") {
			if strings.TrimSpace(line) != "" {
				fmt.Printf("  %s\n", helpDimStyle.Render(line))
			}
		}
		fmt.Println()
	}

	if cmd.HasAvailableSubCommands() {
		fmt.Printf("  %s\n\n", helpDimStyle.Render("Use \"drop [command] --help\" for more information about a command."))
	}
}
