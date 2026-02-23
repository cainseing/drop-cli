package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	green = lipgloss.Color("#9bddca")
	red   = lipgloss.Color("#FF6961")

	successStyle = lipgloss.NewStyle().
			Background(green).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 1).
			Bold(true)
	errorStyle = lipgloss.NewStyle().
			Background(red).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 1).
			Bold(true)
)

func PrintSuccess(header string, message string) {
	fmt.Println()

	if header == "" && message == "" {
		return
	}

	fmt.Printf("%s%s\n", successStyle.Render("✔"), successStyle.Render(header))
	fmt.Printf("%s", message)
}

func PrintError(header string, message string, err error) {
	fmt.Println()

	if err == nil && header == "" && message == "" {
		return
	}

	if header != "" {
		fmt.Printf("%s%s\n\n", errorStyle.Render("✘"), errorStyle.Render(header))
	}

	if message != "" {
		fmt.Printf("%s - ", message)
	}

	if err == nil {
		return
	}

	fmt.Printf("%s\n", err.Error())
}
