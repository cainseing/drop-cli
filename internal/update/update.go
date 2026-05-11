package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cainseing/drop-cli/internal/config"
	"github.com/pterm/pterm"
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

	if latest <= current {
		return false
	}

	message := ""
	if isBrewInstall() {
		message = fmt.Sprintf("Update Available: %s -> %s\nRun %s to update.",
			pterm.Gray(currentVersion),
			pterm.LightGreen(release.TagName),
			pterm.LightCyan("brew upgrade drop"))
	} else {
		message = fmt.Sprintf("Update Available: %s -> %s\nRun %s to update.",
			pterm.Gray(currentVersion),
			pterm.LightGreen(release.TagName),
			pterm.LightCyan("curl -sL getdrop.dev/install.sh | bash"))
	}

	fmt.Println()
	pterm.DefaultBox.
		WithTitle(pterm.LightYellow(" New Version ")).
		WithBoxStyle(pterm.NewStyle(pterm.FgLightYellow)).
		Println(message)
	fmt.Println()

	return true
}

func sanitizeVersion(v string) string {
	v = strings.TrimPrefix(v, "v")
	return strings.Split(v, "-")[0]
}

func isBrewInstall() bool {
	exe, err := getExecutable()
	if err != nil {
		return false
	}

	return strings.Contains(exe, "/homebrew/") || strings.Contains(exe, "/linuxbrew/") || strings.Contains(exe, "/usr/local/Cellar/")
}
