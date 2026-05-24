package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cainseing/drop-cli/internal/config"
	"github.com/charmbracelet/lipgloss"
)

var getExecutable = os.Executable

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

	emerald := lipgloss.Color("#10b981")
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	versionStyle := lipgloss.NewStyle().Foreground(emerald).Bold(true)
	cmdStyle := lipgloss.NewStyle().Foreground(emerald).Italic(true)
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(emerald).
		Padding(0, 1)

	var cmd string
	if isBrewInstall() {
		cmd = "brew upgrade drop"
	} else {
		cmd = "curl -sL getdrop.dev/install.sh | bash"
	}

	message := fmt.Sprintf("Update Available: %s → %s\nRun %s to update.",
		dimStyle.Render(currentVersion),
		versionStyle.Render(release.TagName),
		cmdStyle.Render(cmd))

	fmt.Println()
	fmt.Println(box.Render(message))
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
