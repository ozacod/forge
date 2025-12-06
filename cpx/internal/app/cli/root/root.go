package root

import (
	"fmt"
	"os"

	"github.com/ozacod/cpx/internal/app/cli"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cpx",
	Short: "Cargo-like DX for modern C++ projects",
	Long: `cpx - Cargo-like DX for modern C++

Generate, build, lint, test, and ship CMake/vcpkg-based C++ projects with sensible defaults and cross-compilation ready Docker targets.`,
	Version: cli.Version,
	// Don't show usage on errors by default
	SilenceUsage:  true,
	SilenceErrors: false,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", cli.Red, cli.Reset, err)
		os.Exit(1)
	}
}

// GetRootCmd returns the root command (for testing or extending)
func GetRootCmd() *cobra.Command {
	return rootCmd
}
