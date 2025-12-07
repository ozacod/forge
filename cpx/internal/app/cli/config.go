package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ozacod/cpx/pkg/config"
	"github.com/spf13/cobra"
)

// NewConfigCmd creates the config command
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage cpx configuration",
		Long:  "Manage cpx configuration. Use subcommands to get or set configuration values.",
		RunE:  runConfigShow,
	}

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get config value",
		Long:  "Get a configuration value by key.",
		RunE:  runConfigGet,
		Args:  cobra.ExactArgs(1),
	}
	cmd.AddCommand(getCmd)

	setVcpkgRootCmd := &cobra.Command{
		Use:   "set-vcpkg-root",
		Short: "Set vcpkg root directory",
		Long:  "Set the vcpkg root directory path.",
		RunE:  runConfigSetVcpkgRoot,
		Args:  cobra.ExactArgs(1),
	}
	cmd.AddCommand(setVcpkgRootCmd)

	return cmd
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	return showConfig()
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	return getConfig(args[0])
}

func runConfigSetVcpkgRoot(cmd *cobra.Command, args []string) error {
	return setVcpkgRoot(args[0])
}

func showConfig() error {
	configPath, err := config.GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	cfg, err := config.LoadGlobal()
	if err != nil {
		fmt.Printf("%sCpx Configuration%s\n", Bold, Reset)
		fmt.Printf("  Config file: %s\n", configPath)
		fmt.Printf("  %sError: %s%s\n", Red, err, Reset)
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Printf("%sCpx Configuration%s\n", Bold, Reset)
	fmt.Printf("  Config file: %s\n", configPath)
	fmt.Printf("  vcpkg_root: %s\n", cfg.VcpkgRoot)
	return nil
}

func getConfig(key string) error {
	cfg, err := config.LoadGlobal()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	switch key {
	case "vcpkg_root", "vcpkg-root":
		fmt.Println(cfg.VcpkgRoot)
		return nil
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
}

func setVcpkgRoot(path string) error {
	// Validate path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}

	// Check if it looks like a vcpkg directory
	vcpkgExe := filepath.Join(path, "vcpkg")
	if runtime.GOOS == "windows" {
		vcpkgExe += ".exe"
	}
	if _, err := os.Stat(vcpkgExe); os.IsNotExist(err) {
		fmt.Printf("%s Warning: %s does not appear to be a vcpkg directory%s\n", Yellow, path, Reset)
		fmt.Printf("  (vcpkg executable not found at %s)\n", vcpkgExe)
	}

	cfg, err := config.LoadGlobal()
	if err != nil {
		cfg = &config.GlobalConfig{}
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	cfg.VcpkgRoot = absPath

	if err := config.SaveGlobal(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("%s Set vcpkg_root to %s%s\n", Green, absPath, Reset)
	return nil
}

// Config is kept for backward compatibility (if needed)
func Config(args []string) {
	// This function is deprecated - use NewConfigCmd instead
	// Kept for compatibility during migration
}
