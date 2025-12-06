package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewUpdateCmd creates the update command
func NewUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update dependencies to latest versions",
		Long:  "Update dependencies to latest versions. Use 'vcpkg upgrade' to update vcpkg packages.",
		RunE:  runUpdate,
		Args:  cobra.MaximumNArgs(1),
	}

	cmd.Flags().StringP("server", "s", DefaultServer, "Server URL")

	return cmd
}

func runUpdate(cmd *cobra.Command, args []string) error {
	var libName string
	if len(args) > 0 {
		libName = args[0]
	}
	return updateDependencies(libName)
}

// Update is kept for backward compatibility (if needed)
func Update(args []string) {
	// This function is deprecated - use NewUpdateCmd instead
	// Kept for compatibility during migration
}

func updateDependencies(specificLib string) error {
	// Dependencies are managed in vcpkg.json, not cpx.yaml
	// Read dependencies from vcpkg.json
	deps, err := getDependenciesFromVcpkgJson(".")
	if err != nil {
		return fmt.Errorf("failed to read vcpkg.json: %w\n   Run 'cpx create <project>' first or 'vcpkg new --application' to initialize", err)
	}

	if len(deps) == 0 {
		fmt.Printf("%s No dependencies to update%s\n", Green, Reset)
		return nil
	}

	fmt.Printf("%s Checking for updates...%s\n", Cyan, Reset)
	fmt.Printf("%s  Use 'vcpkg upgrade' to update vcpkg packages%s\n", Yellow, Reset)
	fmt.Printf("   Dependencies in vcpkg.json:\n")
	for _, dep := range deps {
		if specificLib != "" && dep != specificLib {
			continue
		}
		fmt.Printf("    %s\n", dep)
	}

	return nil
}

func getDependenciesFromVcpkgJson(projectDir string) ([]string, error) {
	vcpkgJsonPath := "vcpkg.json"
	if projectDir != "." {
		vcpkgJsonPath = fmt.Sprintf("%s/vcpkg.json", projectDir)
	}

	data, err := os.ReadFile(vcpkgJsonPath)
	if err != nil {
		return nil, err
	}

	var vcpkgConfig struct {
		Dependencies []interface{} `json:"dependencies"`
	}

	if err := json.Unmarshal(data, &vcpkgConfig); err != nil {
		return nil, fmt.Errorf("failed to parse vcpkg.json: %w", err)
	}

	var dependencies []string
	for _, dep := range vcpkgConfig.Dependencies {
		var depStr string
		switch v := dep.(type) {
		case string:
			depStr = v
		case map[string]interface{}:
			if name, ok := v["name"].(string); ok {
				depStr = name
			}
		}
		if depStr != "" {
			dependencies = append(dependencies, depStr)
		}
	}

	return dependencies, nil
}
