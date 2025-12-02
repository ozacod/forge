package cmd

import (
	"github.com/spf13/cobra"
)

var listRunVcpkgCommandFunc func([]string) error

// NewListCmd creates the list command
func NewListCmd(runVcpkgCommand func([]string) error) *cobra.Command {
	listRunVcpkgCommandFunc = runVcpkgCommand

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available libraries",
		Long:  "List available libraries. Passes through to vcpkg list command.",
		RunE:  runList,
		Args:  cobra.ArbitraryArgs,
	}

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	// Directly pass all arguments to vcpkg list command
	vcpkgArgs := []string{"list"}
	vcpkgArgs = append(vcpkgArgs, args...)

	return listRunVcpkgCommandFunc(vcpkgArgs)
}

// List handles the list command - passes through to vcpkg
func List(args []string, runVcpkgCommand func([]string) error) {
	// This function is deprecated - use NewListCmd instead
	// Kept for compatibility during migration
}
