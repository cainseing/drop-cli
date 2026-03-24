package display

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

var (
	green = lipgloss.Color("#9bddca")
	blue  = lipgloss.Color("#5292e6")
	red   = lipgloss.Color("#FF6961")

	successStyle = lipgloss.NewStyle().
			Background(green).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 1).
			Bold(true)
	infoStyle = lipgloss.NewStyle().
			Background(blue).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 1).
			Bold(true)
	errorStyle = lipgloss.NewStyle().
			Background(red).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 1).
			Bold(true)

	Secret      = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	ErrorLabel  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	ErrorText   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8888")).MarginLeft(2)
	MessageText = lipgloss.NewStyle().MarginLeft(2)
)

func PrintSuccess(header string, message string) {
	fmt.Println()

	if header == "" && message == "" {
		return
	}

	fmt.Println(successStyle.Render(fmt.Sprintf("✔ %s", header)))
	if message != "" {
		fmt.Println(MessageText.Render(message))
	}
}

func PrintInfo(header string, message string) {
	fmt.Println()

	if header == "" && message == "" {
		return
	}

	fmt.Println(infoStyle.Render(fmt.Sprintf("ℹ %s", header)))
	if message != "" {
		fmt.Println(MessageText.Render(message))
	}
}

func PrintError(header string, message string, err error) {
	fmt.Fprintln(os.Stderr)

	if err == nil && header == "" && message == "" {
		return
	}

	if header != "" {
		fmt.Fprintln(os.Stderr, errorStyle.Render(fmt.Sprintf("✘ %s", header)))
		fmt.Fprintln(os.Stderr)
	}

	if message != "" && err != nil {
		fmt.Fprintln(os.Stderr, ErrorText.Render(fmt.Sprintf("%s: %s", message, err.Error())))
	} else if message != "" {
		fmt.Fprintln(os.Stderr, ErrorText.Render(message))
	} else if err != nil {
		fmt.Fprintln(os.Stderr, ErrorText.Render(err.Error()))
	}
}
