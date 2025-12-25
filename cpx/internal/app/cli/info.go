package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ozacod/cpx/internal/pkg/build/bazel"
	build "github.com/ozacod/cpx/internal/pkg/build/interfaces"
	"github.com/ozacod/cpx/internal/pkg/build/meson"
	"github.com/ozacod/cpx/internal/pkg/build/vcpkg"
	"github.com/ozacod/cpx/internal/pkg/utils/colors"
	"github.com/spf13/cobra"
)

func InfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <package>",
		Short: "Show detailed library information",
		Long:  "Show detailed library information for a vcpkg package.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInfo(cmd, args)
		},
		Args: cobra.MinimumNArgs(1),
	}

	cmd.Flags().Bool("json", false, "Output in JSON format")

	return cmd
}

// PackageInfo represents the structure of vcpkg x-package-info output
type PackageInfo struct {
	Results map[string]struct {
		Name         string `json:"name"`
		Version      string `json:"version-semver"`
		VersionDate  string `json:"version-date"`
		VersionStr   string `json:"version-string"`
		Description  any    `json:"description"`
		Homepage     string `json:"homepage"`
		License      string `json:"license"`
		Dependencies []any  `json:"dependencies"`
		Features     map[string]struct {
			Description string `json:"description"`
		} `json:"features"`
	} `json:"results"`
}

func runInfo(cmd *cobra.Command, args []string) error {
	jsonOutput, _ := cmd.Flags().GetBool("json")
	packageName := args[0]

	projectType := DetectProjectType()

	var builder build.BuildSystem
	var err error

	switch projectType {
	case ProjectTypeBazel:
		builder = bazel.New()
	case ProjectTypeMeson:
		builder = meson.New()
	default: // vcpkg/cmake
		builder = vcpkg.New()
	}

	info, err := builder.DependencyInfo(context.Background(), packageName)
	if err != nil {
		return err
	}

	if jsonOutput {
		// Output as JSON
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(info)
	}

	// Print formatted output
	fmt.Printf("%s📦 %s%s %s%s%s\n", colors.Bold, colors.Cyan, info.Name, colors.Yellow, info.Version, colors.Reset)

	if info.Description != "" {
		// Handle multi-line description
		lines := strings.Split(info.Description, "\n")
		for _, line := range lines {
			fmt.Printf("   %s\n", line)
		}
	}

	if info.Homepage != "" {
		fmt.Printf("\n%s🔗 Homepage:%s %s\n", colors.Bold, colors.Reset, info.Homepage)
	}

	if info.License != "" {
		fmt.Printf("%s📄 License:%s  %s\n", colors.Bold, colors.Reset, info.License)
	}

	// Dependencies
	if len(info.Dependencies) > 0 {
		fmt.Printf("\n%s📚 Dependencies:%s\n", colors.Bold, colors.Reset)
		for _, dep := range info.Dependencies {
			fmt.Printf("   • %s\n", dep)
		}
	}

	return nil
}
