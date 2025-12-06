package template

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadFromGitHub downloads a template file from the GitHub repository
func DownloadFromGitHub(templateName, localPath string) error {
	// Create directory if it doesn't exist
	templatesDir := filepath.Dir(localPath)
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Download from GitHub raw content
	url := fmt.Sprintf("https://raw.githubusercontent.com/ozacod/cpx/master/templates/%s", templateName)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download template: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download template: HTTP %d", resp.StatusCode)
	}

	// Write to local file
	out, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write template file: %w", err)
	}

	return nil
}
