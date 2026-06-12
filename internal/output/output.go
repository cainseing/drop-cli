package output

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

const (
	hexEmerald = "#10b981"
	hexLight   = "#34d399"
	hexLighter = "#6ee7b7"
	hexDim     = "#64748b"
	hexWhite   = "#f1f5f9"
	hexRed     = "#f87171"
)

// ansi wraps s in a 24-bit ANSI foreground sequence, or returns s plain when
// colour output is disabled (NO_COLOR, non-TTY, etc.).
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

// Styled is a simple styled-text type whose .Render() method is called by
// service.go, providing backward compatibility without lipgloss.
type Styled struct {
	hex  string
	bold bool
}

func (s Styled) Render(text string) string {
	return ansi(s.hex, s.bold, text)
}

// Exported Styled vars used directly in service.go.
var (
	DimText = Styled{hex: hexDim}
	Secret  = Styled{hex: hexWhite}
)

const labelWidth = 10

func label(s string) string {
	padded := fmt.Sprintf("%-*s", labelWidth, s)
	return ansi(hexDim, false, padded)
}

// PrintProperty writes a label/value pair to stdout.
func PrintProperty(lbl, value string) {
	printPropertyTo(os.Stdout, lbl, value)
}

// PrintPropertyToStderr writes a label/value pair to stderr.
func PrintPropertyToStderr(lbl, value string) {
	printPropertyTo(os.Stderr, lbl, value)
}

func printPropertyTo(out *os.File, lbl, value string) {
	if lbl == "" {
		return
	}
	fmt.Fprintf(out, "  %s  %s\n", label(lbl), value)
}

// PrintSuccess prints a success status line.
func PrintSuccess(title, description string) {
	fmt.Println()
	fmt.Printf("  %s  %s\n", label("STATUS"), RenderVerified(title))
	if description != "" {
		fmt.Printf("  %s\n", ansi(hexDim, false, description))
	}
}

// PrintInfo prints an informational status line.
func PrintInfo(title, description string) {
	fmt.Println()
	fmt.Printf("  %s  %s\n", label("INFO"), ansi(hexEmerald, true, title))
	if description != "" {
		fmt.Printf("  %s\n", ansi(hexDim, false, description))
	}
}

// PrintError writes an error block to stderr.
func PrintError(message string, err error) {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "  %s\n", ansi(hexRed, true, "✗ ERROR"))

	var context string
	switch {
	case message != "" && err != nil:
		context = fmt.Sprintf("%s (%s)", message, err.Error())
	case message != "":
		context = message
	case err != nil:
		context = err.Error()
	}

	if context != "" {
		fmt.Fprintf(os.Stderr, "  %s\n", ansi(hexDim, false, ucFirst(context)))
	}
}

// RenderToken renders the drop token with the branded prefix.
func RenderToken(token string) string {
	prefix := ansi(hexEmerald, true, "drop_")
	body := ansi(hexLighter, false, strings.TrimPrefix(token, "drop_"))
	return "\n" + prefix + body
}

// RenderVerified renders a verified checkmark string.
func RenderVerified(text string) string {
	return ansi(hexLight, true, "✓ "+text)
}

func ucFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
