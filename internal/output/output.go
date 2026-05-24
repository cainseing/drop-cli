package output

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
)

var (
	accent  = lipgloss.AdaptiveColor{Light: "#10b981", Dark: "#10b981"}
	success = lipgloss.Color("#34d399")
	danger  = lipgloss.Color("#f87171")
	dim     = lipgloss.Color("#666666")

	LabelStyle = lipgloss.NewStyle().Foreground(dim).Width(12)
	ValueStyle = lipgloss.NewStyle().Foreground(accent).Bold(true)
	DimText    = lipgloss.NewStyle().Foreground(dim)

	CopyStyle = lipgloss.NewStyle().Foreground(accent).Italic(true)

	Secret = lipgloss.NewStyle().
		Padding(0, 0).
		MarginTop(1).
		Foreground(lipgloss.Color("#FFFFFF"))

	Token = lipgloss.NewStyle().
		Padding(0, 0).
		MarginTop(1).
		Foreground(lipgloss.Color("#FFFFFF"))

	labelPadding = 13

	StatusVerified = lipgloss.NewStyle().Foreground(success).Bold(true)
	StatusError    = lipgloss.NewStyle().Foreground(danger).Bold(true)
)

func RenderVerified(text string) string {
	return StatusVerified.Render("✓ " + text)
}

func RenderToken(token string) string {
	prefix := lipgloss.NewStyle().Foreground(lipgloss.Color("#10b981")).Bold(true).Render("drop_")
	body := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Render(strings.TrimPrefix(token, "drop_"))
	return "\n" + prefix + body
}

func PrintProperty(label string, value string) {
	printPropertyTo(os.Stdout, label, value)
}

func PrintPropertyToStderr(label string, value string) {
	printPropertyTo(os.Stderr, label, value)
}

func printPropertyTo(out *os.File, label string, value string) {
	if label == "" {
		return
	}

	fmt.Fprintf(out, "%s %s\n", LabelStyle.Render(label), value)
}

func PrintSuccess(title string, description string) {
	fmt.Println()
	fmt.Printf("%s %s\n", LabelStyle.Render("STATUS"), StatusVerified.Render(title))
	if description != "" {
		fmt.Println(DimText.PaddingLeft(labelPadding).Render(description))
	}
}

func PrintInfo(title string, description string) {
	fmt.Println()
	fmt.Printf("%s %s\n", LabelStyle.Render("INFO"), ValueStyle.Render(title))
	if description != "" {
		fmt.Println(DimText.PaddingLeft(labelPadding).Render(description))
	}
}

func PrintError(message string, err error) {
	fmt.Fprintln(os.Stderr)

	fmt.Fprintf(os.Stderr, "%s\n",
		StatusError.Render("✗ ERROR"),
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

func ucFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
