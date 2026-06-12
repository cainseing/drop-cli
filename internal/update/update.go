package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cainseing/drop-cli/internal/config"
	"github.com/fatih/color"
)

var getExecutable = os.Executable

// ansiRe matches ANSI escape sequences for stripping in width calculations.
var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

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

// visibleLen returns the printable rune count of s, ignoring ANSI escape codes.
func visibleLen(s string) int {
	return len([]rune(ansiRe.ReplaceAllString(s, "")))
}

func CheckForUpdates(currentVersion string, force bool) bool {
	fi, _ := os.Stdout.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return false
	}

	cfg := config.LoadUserConfig()

	if !force && time.Since(cfg.UpdateCheck) < 24*time.Hour {
		return false
	}

	cfg.UpdateCheck = time.Now()
	_ = config.SaveUserConfig(cfg)

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/cainseing/drop-cli/releases/latest")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return false
	}

	latest := sanitizeVersion(release.TagName)
	current := sanitizeVersion(currentVersion)

	if compareVersions(latest, current) <= 0 {
		return false
	}

	var cmd string
	if isBrewInstall() {
		cmd = "brew upgrade drop"
	} else {
		cmd = "curl -sL getdrop.dev/install.sh | bash"
	}

	line1 := fmt.Sprintf("  Update Available: %s → %s",
		ansi("#64748b", false, currentVersion),
		ansi("#10b981", true, release.TagName),
	)
	line2 := fmt.Sprintf("  Run %s to update.", ansi("#10b981", false, cmd))

	width := max(visibleLen(line1), visibleLen(line2)) + 2
	bar := strings.Repeat("─", width)

	fmt.Println()
	fmt.Println(ansi("#10b981", false, "╭"+bar+"╮"))
	fmt.Println(ansi("#10b981", false, "│") + line1 + strings.Repeat(" ", width-visibleLen(line1)) + ansi("#10b981", false, "│"))
	fmt.Println(ansi("#10b981", false, "│") + line2 + strings.Repeat(" ", width-visibleLen(line2)) + ansi("#10b981", false, "│"))
	fmt.Println(ansi("#10b981", false, "╰"+bar+"╯"))
	fmt.Println()

	return true
}

func sanitizeVersion(v string) string {
	v = strings.TrimPrefix(v, "v")
	return strings.Split(v, "-")[0]
}

func compareVersions(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	for i := range max(len(aParts), len(bParts)) {
		var aVal, bVal int
		if i < len(aParts) {
			aVal, _ = strconv.Atoi(aParts[i])
		}
		if i < len(bParts) {
			bVal, _ = strconv.Atoi(bParts[i])
		}
		if aVal != bVal {
			if aVal > bVal {
				return 1
			}
			return -1
		}
	}
	return 0
}

func isBrewInstall() bool {
	exe, err := getExecutable()
	if err != nil {
		return false
	}

	return strings.Contains(exe, "/homebrew/") || strings.Contains(exe, "/linuxbrew/") || strings.Contains(exe, "/usr/local/Cellar/")
}
