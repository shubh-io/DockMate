package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/shubh-io/dockmate/pkg/version"
)

// Check if dockmate is installed via Homebrew
func isHomebrewInstall() bool {
	// First, check if brew exists and whether dockmate is registered under brew
	if _, err := exec.LookPath("brew"); err == nil {
		// Prefer explicit tap name; fallback to plain formula name
		if err := exec.Command("brew", "list", "--versions", "shubh-io/tap/dockmate").Run(); err == nil {
			return true
		}
		if err := exec.Command("brew", "list", "--versions", "dockmate").Run(); err == nil {
			return true
		}
		// Fallback: compare executable path to brew prefix
		exe, err := os.Executable()
		if err == nil {
			prefixOut, pErr := exec.Command("brew", "--prefix").Output()
			if pErr == nil {
				prefix := strings.TrimSpace(string(prefixOut))
				exeLower := strings.ToLower(exe)
				// Common brew locations
				if strings.HasPrefix(exeLower, strings.ToLower(prefix)) ||
					strings.Contains(exeLower, "/cellar/dockmate") ||
					strings.Contains(exeLower, "/opt/homebrew") ||
					strings.Contains(exeLower, "/usr/local/cellar") ||
					strings.Contains(exeLower, ".linuxbrew") ||
					strings.Contains(exeLower, "/home/linuxbrew") {
					return true
				}
			}
		}
	}

	// As a last resort, path heuristics without brew available
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	exePath := strings.ToLower(exe)
	homebrewHints := []string{
		"/linuxbrew/",
		"/home/linuxbrew/",
		"/homebrew/",
		"/opt/homebrew/",
		"/usr/local/cellar/",
		"cellar/dockmate",
		".linuxbrew",
	}
	for _, h := range homebrewHints {
		if strings.Contains(exePath, strings.ToLower(h)) {
			return true
		}
	}
	return false
}

// getLatestReleaseTag fetches the latest release tag name from GitHub for the given repo (owner/repo)
// This uses a small shell pipeline to keep the implementation compact
func getLatestReleaseTag(repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.Unmarshal(body, &release); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if strings.TrimSpace(release.TagName) == "" {
		return "", fmt.Errorf("no tag name found in release")
	}

	return release.TagName, nil
}

// trims whitespace and leading 'v' or 'V'
func normalizeTag(tag string) string {
	tag = strings.TrimSpace(tag)
	if strings.HasPrefix(tag, "v") || strings.HasPrefix(tag, "V") {
		return tag[1:]
	}
	return tag
}

// compareSemver compares two simple dot-separated semantic versions (major.minor.patch).
// returns -1 if a<b, 0 if equal, +1 if a>b.
// non-numeric parts are compared lexically!
func compareSemver(a, b string) int {
	a = normalizeTag(a)
	b = normalizeTag(b)
	if a == b {
		return 0
	}
	// split a into parts - eg: "1.2.3" -> ["1","2","3"]
	a_splited := strings.Split(a, ".")
	// split b into parts - eg: "1.2.0" -> ["1","2","0"]
	b_splited := strings.Split(b, ".")

	// compare each part
	n := len(a_splited)
	// if length of b is greater than a, then we compare by b's length
	if len(b_splited) > n {
		n = len(b_splited)
	}

	for i := 0; i < n; i++ {
		var a_value, b_value string
		if i < len(a_splited) {
			a_value = a_splited[i]
		}
		if i < len(b_splited) {
			b_value = b_splited[i]
		}
		if a_value == b_value {
			continue
		}
		// attempting numeric compare for best accuracy
		ai, aErr := strconv.Atoi(a_value)
		bi, bErr := strconv.Atoi(b_value)
		if aErr == nil && bErr == nil {
			if ai < bi {
				return -1
			}
			if ai > bi {
				return 1
			}
			// equal numerically, continue
			continue
		}
		// fallback to lexical comparison if either segment isn't a plain integer
		if cmp := strings.Compare(a_value, b_value); cmp != 0 {
			return cmp
		}
	}
	return 0
}

func UpdateCommand() {
	fmt.Println("Checking for updates...")

	// Check if installed via Homebrew FIRST
	if isHomebrewInstall() {
		fmt.Println("⚠️ Detected: dockmate is installed via Homebrew")
		fmt.Println("")
		fmt.Println("To update, please run:")
		fmt.Println("  brew upgrade shubh-io/tap/dockmate")
		fmt.Println("")
		fmt.Println("Current version:", version.Dockmate_Version)
		return
	}

	// Ensure we have the current version constant available
	current := version.Dockmate_Version

	latestTag, err := getLatestReleaseTag(version.Repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not check latest release: %v\n", err)
		return
	}

	// compare normalized tags (striped 'v')
	cmp := compareSemver(current, latestTag)
	if cmp >= 0 {
		fmt.Printf("Already up-to-date (current: %s, latest: %s)\n", current, latestTag)
		return
	}

	fmt.Printf("New release available! : %s → %s\n", current, latestTag)
	fmt.Println("Re-running installer to update...")

	// Try piped install first
	cmd := exec.Command("bash", "-c", "curl -fsSL https://raw.githubusercontent.com/shubh-io/dockmate/main/install.sh | bash")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Piped install failed, trying fallback method...")

		// Fallback- download script, run it, then delete it lol
		downloadCmd := exec.Command("bash", "-c", "curl -fsSL https://raw.githubusercontent.com/shubh-io/dockmate/main/install.sh -o install.sh")
		if err := downloadCmd.Run(); err != nil {
			fmt.Printf("Failed to download install script: %v\n", err)
			return
		}

		runCmd := exec.Command("bash", "install.sh")
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		if err := runCmd.Run(); err != nil {
			fmt.Printf("Update failed: %v\n", err)
			// Still try to clean up
			os.Remove("install.sh")
			return
		}

		// Clean up the script file
		if err := os.Remove("install.sh"); err != nil {
			fmt.Printf("Warning: could not remove install.sh: %v\n", err)
		}
	}

	fmt.Println("")
	fmt.Println("Updated successfully!")
	fmt.Println("Please Run 'dockmate version' to verify")
	fmt.Println("⚠️ Note: You may need to run 'hash -r' if version doesn't update")
}
