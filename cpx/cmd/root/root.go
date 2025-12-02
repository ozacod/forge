package root

import (
	"fmt"
	"os"

	"github.com/ozacod/cpx/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cpx",
	Short: "C++ Project Generator (like Cargo for Rust)",
	Long: `cpx - C++ Project Generator (like Cargo for Rust)

A modern build tool and project generator for C++ projects using CMake and vcpkg.`,
	Version: cmd.Version,
	// Don't show usage on errors by default
	SilenceUsage:  true,
	SilenceErrors: false,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", cmd.Red, cmd.Reset, err)
		os.Exit(1)
	}
}

// GetRootCmd returns the root command (for testing or extending)
func GetRootCmd() *cobra.Command {
	return rootCmd
}

