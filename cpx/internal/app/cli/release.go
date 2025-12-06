package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/ozacod/cpx/pkg/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// NewReleaseCmd creates the release command
func NewReleaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Bump version number",
		Long:  "Bump version number (major, minor, or patch). Defaults to patch if not specified.",
		RunE:  runRelease,
		Args:  cobra.MaximumNArgs(1),
	}

	return cmd
}

func runRelease(cmd *cobra.Command, args []string) error {
	bumpType := "patch"
	if len(args) > 0 {
		bumpType = args[0]
	}
	return bumpVersion(bumpType)
}

// Release is kept for backward compatibility (if needed)
func Release(args []string) {
	// This function is deprecated - use NewReleaseCmd instead
	// Kept for compatibility during migration
}

func bumpVersion(bumpType string) error {
	cfg, err := config.LoadProject(DefaultCfgFile)
	if err != nil {
		return err
	}

	version := cfg.Package.Version
	if version == "" {
		version = "0.1.0"
	}

	// Parse version
	parts := strings.Split(strings.TrimPrefix(version, "v"), ".")
	if len(parts) < 3 {
		parts = append(parts, make([]string, 3-len(parts))...)
	}

	major, minor, patch := 0, 0, 0
	fmt.Sscanf(parts[0], "%d", &major)
	fmt.Sscanf(parts[1], "%d", &minor)
	fmt.Sscanf(parts[2], "%d", &patch)

	switch bumpType {
	case "major":
		major++
		minor = 0
		patch = 0
	case "minor":
		minor++
		patch = 0
	case "patch":
		patch++
	default:
		return fmt.Errorf("invalid bump type: %s (use major, minor, or patch)", bumpType)
	}

	newVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	cfg.Package.Version = newVersion

	fmt.Printf("%s Bumping version: %s  %s%s\n", Cyan, version, newVersion, Reset)

	if err := saveProjectConfig(cfg); err != nil {
		return err
	}

	fmt.Printf("%s Version updated to %s%s\n", Green, newVersion, Reset)
	return nil
}

func saveProjectConfig(cfg *config.ProjectConfig) error {
	// Create a copy without dependencies (dependencies are in vcpkg.json)
	configCopy := *cfg
	configCopy.Dependencies = nil // Don't save dependencies to cpx.yaml

	data, err := yaml.Marshal(&configCopy)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Add header comment
	header := "# cpx.yaml - C++ Project Configuration\n# Dependencies are managed in vcpkg.json (use 'vcpkg add port <package>')\n\n"
	data = append([]byte(header), data...)

	if err := os.WriteFile(DefaultCfgFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

