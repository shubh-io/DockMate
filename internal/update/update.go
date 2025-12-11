package update

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/shubh-io/dockmate/pkg/version"
)

// getLatestReleaseTag fetches the latest release tag name from GitHub for the given repo (owner/repo)
// This uses a small shell pipeline to keep the implementation compact
func getLatestReleaseTag(repo string) (string, error) {
	// build the shell command; use the compact pipeline suggested by the user
	cmd := fmt.Sprintf("curl -s https://api.github.com/repos/%s/releases/latest | grep '\"tag_name\":' | sed -E 's/.*\"([^\"]+)\".*/\\1/'", repo)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	tag := strings.TrimSpace(string(out))
	if tag == "" {
		return "", fmt.Errorf("empty tag from GitHub for %s", repo)
	}
	return tag, nil
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

	fmt.Printf("New release available: %s → %s\n", current, latestTag)
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

	fmt.Println("Updated successfully!")
}
