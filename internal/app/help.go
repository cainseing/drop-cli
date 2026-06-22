package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var Logo = "      _\n" +
	"   __| |_ __ ___  _ __\n" +
	"  / _` | '__/ _ \\| '_ \\\n" +
	" | (_| | | | (_) | |_) |\n" +
	"  \\__,_|_|  \\___/| .__/\n" +
	"                 |_|"

// ansi returns s wrapped in a 24-bit ANSI foreground sequence, plain when
// colour is disabled.
func ansi(hex string, bold bool, s string) string {
	if color.NoColor {
		return s
	}
	r, g, b := hexRGB(hex)
	prefix := fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	if bold {
		prefix = "\033[1m" + prefix
	}
	return prefix + s + "\033[0m"
}

func hexRGB(h string) (uint8, uint8, uint8) {
	h = strings.TrimPrefix(h, "#")
	if len(h) != 6 {
		return 255, 255, 255
	}
	ri, _ := strconv.ParseUint(h[0:2], 16, 8)
	gi, _ := strconv.ParseUint(h[2:4], 16, 8)
	bi, _ := strconv.ParseUint(h[4:6], 16, 8)
	return uint8(ri), uint8(gi), uint8(bi)
}

const (
	Emerald = "#10b981"
	Dim     = "#64748b"
	White   = "#f1f5f9"
)

const (
	emerald = Emerald
	dim     = Dim
	white   = White
)

func RenderText(hex string, bold bool, s string) string {
	return ansi(hex, bold, s)
}

func applyHelpFunc(cmd *cobra.Command) {
	cmd.SetHelpFunc(renderHelp)
	for _, sub := range cmd.Commands() {
		applyHelpFunc(sub)
	}
}

func renderHelp(cmd *cobra.Command, _ []string) {
	fmt.Println()

	if !cmd.HasParent() {
		for _, line := range strings.Split(Logo, "\n") {
			fmt.Printf("  %s\n", ansi(emerald, false, line))
		}
		fmt.Println()
		fmt.Printf("  %s\n", ansi(dim, false, cmd.Short))
	} else {
		fmt.Printf("  %s  %s\n", ansi(emerald, true, cmd.Name()), ansi(dim, false, cmd.Short))
	}

	fmt.Println()

	fmt.Printf("  %s\n", ansi(emerald, true, "Usage"))
	fmt.Printf("    %s\n", ansi(dim, false, cmd.UseLine()))
	if cmd.HasAvailableSubCommands() {
		fmt.Printf("    %s\n", ansi(dim, false, cmd.CommandPath()+" [command]"))
	}
	fmt.Println()

	if cmd.HasAvailableSubCommands() {
		fmt.Printf("  %s\n", ansi(emerald, true, "Commands"))

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
			name := ansi(dim, false, fmt.Sprintf("%-*s", maxLen, c.Name()))
			fmt.Printf("    %s  %s\n", name, ansi(dim, false, c.Short))
		}
		fmt.Println()
	}

	if cmd.HasAvailableFlags() {
		fmt.Printf("  %s\n", ansi(emerald, true, "Flags"))
		for _, line := range strings.Split(strings.TrimRight(cmd.Flags().FlagUsages(), "\n"), "\n") {
			if strings.TrimSpace(line) != "" {
				fmt.Printf("  %s\n", ansi(dim, false, line))
			}
		}
		fmt.Println()
	}

	if cmd.HasAvailableSubCommands() {
		fmt.Printf("  %s\n\n", ansi(dim, false, "Use \"drop [command] --help\" for more information about a command."))
	}
}
