package cmd

import (
	"github.com/spf13/cobra"
)

var searchRunVcpkgCommandFunc func([]string) error

// NewSearchCmd creates the search command
func NewSearchCmd(runVcpkgCommand func([]string) error) *cobra.Command {
	searchRunVcpkgCommandFunc = runVcpkgCommand

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for libraries",
		Long:  "Search for libraries. Passes through to vcpkg search command.",
		RunE:  runSearch,
		Args:  cobra.MinimumNArgs(1),
	}

	return cmd
}

func runSearch(cmd *cobra.Command, args []string) error {
	// Directly pass all arguments to vcpkg search command
	vcpkgArgs := []string{"search"}
	vcpkgArgs = append(vcpkgArgs, args...)

	return searchRunVcpkgCommandFunc(vcpkgArgs)
}

// Search handles the search command - passes through to vcpkg
func Search(args []string, runVcpkgCommand func([]string) error) {
	// This function is deprecated - use NewSearchCmd instead
	// Kept for compatibility during migration
}
