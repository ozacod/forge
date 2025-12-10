package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// UpgradeCmd creates the upgrade command
func UpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade cpx to the latest version",
		Long:  "Upgrade cpx to the latest version from GitHub releases.",
		RunE:  runUpgrade,
	}

	return cmd
}

func runUpgrade(_ *cobra.Command, args []string) error {
	Upgrade(args)
	return nil
}

func Upgrade(_ []string) {
	fmt.Printf("%s Checking for updates...%s\n", Cyan, Reset)

	// Get latest version from GitHub releases API
	resp, err := http.Get("https://api.github.com/repos/ozacod/cpx/releases/latest")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s Failed to check for updates: %v\n", Red, Reset, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("%s  No releases found. This may be the first version.%s\n", Yellow, Reset)
		fmt.Printf("   Repository: https://github.com/ozacod/cpx\n")
		return
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "%sError:%s Failed to check for updates (status %d): %s\n", Red, Reset, resp.StatusCode, string(body))
		os.Exit(1)
	}

	var release struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s Failed to parse release info: %v\n", Red, Reset, err)
		os.Exit(1)
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := Version

	if latestVersion == currentVersion {
		fmt.Printf("%s You're already running the latest version (%s)%s\n", Green, currentVersion, Reset)
		return
	}

	fmt.Printf("%s New version available: %s  %s%s\n", Yellow, currentVersion, latestVersion, Reset)
	fmt.Printf("   Release: %s\n", release.HTMLURL)

	// Determine platform and architecture
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	var binaryName string
	switch goos {
	case "darwin":
		binaryName = fmt.Sprintf("cpx-darwin-%s", goarch)
	case "linux":
		binaryName = fmt.Sprintf("cpx-linux-%s", goarch)
	case "windows":
		binaryName = fmt.Sprintf("cpx-windows-%s.exe", goarch)
	default:
		fmt.Fprintf(os.Stderr, "%sError:%s Unsupported platform: %s\n", Red, Reset, goos)
		os.Exit(1)
	}

	downloadURL := fmt.Sprintf("https://github.com/ozacod/cpx/releases/download/%s/%s", release.TagName, binaryName)
	fmt.Printf("%s Downloading %s...%s\n", Cyan, binaryName, Reset)

	// Download the new binary
	resp, err = http.Get(downloadURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s Failed to download: %v\n", Red, Reset, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "%sError:%s Download failed with status %d\n", Red, Reset, resp.StatusCode)
		os.Exit(1)
	}

	binaryData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s Failed to read download: %v\n", Red, Reset, err)
		os.Exit(1)
	}

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s Failed to get executable path: %v\n", Red, Reset, err)
		os.Exit(1)
	}
	execPath, _ = filepath.EvalSymlinks(execPath)

	// Write to temp file first
	tempPath := execPath + ".new"
	if err := os.WriteFile(tempPath, binaryData, 0755); err != nil {
		// Try writing to temp directory instead
		tempPath = filepath.Join(os.TempDir(), "cpx-new")
		if err := os.WriteFile(tempPath, binaryData, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "%sError:%s Failed to write binary: %v\n", Red, Reset, err)
			os.Exit(1)
		}
		fmt.Printf("%s Downloaded to %s%s\n", Green, tempPath, Reset)
		fmt.Printf("\nTo complete the upgrade, run:\n")
		fmt.Printf("  sudo mv %s %s\n", tempPath, execPath)
		return
	}

	// Remove old binary and rename new one
	os.Remove(execPath)
	if err := os.Rename(tempPath, execPath); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s Failed to replace binary: %v\n", Red, Reset, err)
		fmt.Printf("\nTo complete manually, run:\n")
		fmt.Printf("  sudo mv %s %s\n", tempPath, execPath)
		os.Exit(1)
	}

	fmt.Printf("%s Successfully upgraded to %s!%s\n", Green, latestVersion, Reset)
	fmt.Printf("  Run %scpx version%s to verify.\n", Cyan, Reset)
}
