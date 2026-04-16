package display

import (
	"fmt"
	"os"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

var (
	subtle  = lipgloss.AdaptiveColor{Light: "#D9D9D9", Dark: "#383838"}
	accent  = lipgloss.AdaptiveColor{Light: "#00ADD8", Dark: "#00ADD8"}
	success = lipgloss.Color("#A3BE8C")
	danger  = lipgloss.Color("#BF616A")
	dim     = lipgloss.Color("#666666")

	LabelStyle = lipgloss.NewStyle().Foreground(dim).Width(12)
	ValueStyle = lipgloss.NewStyle().Foreground(accent).Bold(true)
	DimText    = lipgloss.NewStyle().Foreground(dim)

	CopyStyle = lipgloss.NewStyle().Foreground(accent).Italic(true)

	Secret = lipgloss.NewStyle().
		Foreground(subtle).
		Padding(0, 0).
		MarginTop(1).
		Foreground(lipgloss.Color("#FFFFFF"))

	Token = lipgloss.NewStyle().
		Foreground(subtle).
		Padding(0, 0).
		MarginTop(1).
		Foreground(lipgloss.Color("#FFFFFF"))

	StatusVerified = lipgloss.NewStyle().Foreground(success).Bold(true)
	StatusError    = lipgloss.NewStyle().Foreground(danger).Bold(true)
)

func PrintPropertyToStderr(label string, value string) {
	if label == "" {
		return
	}

	fmt.Fprintf(os.Stderr, "%s %s\n", LabelStyle.Render(label), value)
}

func PrintProperty(label string, value string) {
	if label == "" {
		return
	}

	fmt.Printf("%s %s\n", LabelStyle.Render(label), value)
}

func PrintSuccess(title string, description string) {
	fmt.Println()
	fmt.Printf("%s %s\n", LabelStyle.Render("STATUS"), StatusVerified.Render(title))
	if description != "" {
		fmt.Println(DimText.PaddingLeft(13).Render(description))
	}
}

func PrintInfo(title string, description string) {
	fmt.Println()
	fmt.Printf("%s %s\n", LabelStyle.Render("INFO"), ValueStyle.Render(title))
	if description != "" {
		fmt.Println(DimText.PaddingLeft(13).Render(description))
	}
}

func PrintError(message string, err error) {
	fmt.Fprintln(os.Stderr)

	fmt.Fprintf(os.Stderr, "%s\n",
		StatusError.Bold(true).Render("ERROR"),
	)

	var context string
	if message != "" && err != nil {
		context = fmt.Sprintf("%s (%s)", message, err.Error())
	} else if message != "" {
		context = message
	} else if err != nil {
		context = err.Error()
	}

	if context != "" {
		fmt.Fprintf(os.Stderr, "%s\n", DimText.Render(ucFirst(context)))
	}
}

func getTerminalWidth() int {
	width, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil || width <= 0 {
		return 80
	}
	return width
}

func ucFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
