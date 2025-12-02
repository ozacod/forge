package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/ozacod/cpx/internal/config"
	"github.com/spf13/cobra"
)

// NewDocCmd creates the doc command
func NewDocCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doc",
		Short: "Generate documentation",
		Long:  "Generate documentation using Doxygen. Use --open to open in browser after generation.",
		RunE:  runDoc,
	}

	cmd.Flags().Bool("open", false, "Open documentation in browser")

	return cmd
}

func runDoc(cmd *cobra.Command, args []string) error {
	open, _ := cmd.Flags().GetBool("open")
	return generateDocs(open)
}

// Doc is kept for backward compatibility (if needed)
func Doc(args []string) {
	// This function is deprecated - use NewDocCmd instead
	// Kept for compatibility during migration
}

func generateDocs(openBrowser bool) error {
	// Check if Doxygen is available
	if _, err := exec.LookPath("doxygen"); err != nil {
		return fmt.Errorf("doxygen not found. Please install it first:\n  macOS: brew install doxygen\n  Ubuntu: sudo apt install doxygen")
	}

	cfg, err := config.LoadProject(DefaultCfgFile)
	if err != nil {
		return err
	}

	fmt.Printf("%s Generating documentation...%s\n", Cyan, Reset)

	// Create Doxyfile if it doesn't exist
	if _, err := os.Stat("Doxyfile"); os.IsNotExist(err) {
		doxyContent := fmt.Sprintf(`PROJECT_NAME           = "%s"
PROJECT_NUMBER         = "%s"
OUTPUT_DIRECTORY       = docs
INPUT                  = src include
RECURSIVE              = YES
EXTRACT_ALL            = YES
GENERATE_HTML          = YES
GENERATE_LATEX         = NO
HTML_OUTPUT            = html
USE_MDFILE_AS_MAINPAGE = README.md
`, cfg.Package.Name, cfg.Package.Version)

		if err := os.WriteFile("Doxyfile", []byte(doxyContent), 0644); err != nil {
			return fmt.Errorf("failed to create Doxyfile: %w", err)
		}
		fmt.Printf("    Created Doxyfile\n")
	}

	// Run Doxygen
	cmd := exec.Command("doxygen")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("doxygen failed: %w", err)
	}

	indexPath := "docs/html/index.html"
	fmt.Printf("%s Documentation generated at %s%s\n", Green, indexPath, Reset)

	if openBrowser {
		var openCmd string
		switch runtime.GOOS {
		case "darwin":
			openCmd = "open"
		case "linux":
			openCmd = "xdg-open"
		case "windows":
			openCmd = "start"
		}

		if openCmd != "" {
			exec.Command(openCmd, indexPath).Start()
		}
	}

	return nil
}
