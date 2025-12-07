package cli

import (
	"fmt"

	"github.com/ozacod/cpx/internal/app/cli/tui"
	"github.com/spf13/cobra"
)

var searchRunVcpkgCommandFunc func([]string) error
var searchGetVcpkgPath func() (string, error)

// NewSearchCmd creates the search command
func NewSearchCmd(runVcpkgCommand func([]string) error, getVcpkgPath func() (string, error)) *cobra.Command {
	searchRunVcpkgCommandFunc = runVcpkgCommand
	searchGetVcpkgPath = getVcpkgPath

	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search for libraries interactively",
		Long:  "Search for libraries using an interactive TUI. Select packages to add them to your project.",
		RunE:  runSearch,
		Args:  cobra.MaximumNArgs(1),
	}

	return cmd
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := ""
	if len(args) > 0 {
		query = args[0]
	}

	if err := requireVcpkgProject("cpx search"); err != nil {
		return err
	}

	vcpkgPath, err := searchGetVcpkgPath()
	if err != nil {
		return fmt.Errorf("failed to get vcpkg path: %w", err)
	}

	return tui.RunSearch(query, vcpkgPath, searchRunVcpkgCommandFunc)
}

// Search handles the search command - passes through to vcpkg
func Search(args []string, runVcpkgCommand func([]string) error) {
	// This function is deprecated - use NewSearchCmd instead
	// Kept for compatibility during migration
}
