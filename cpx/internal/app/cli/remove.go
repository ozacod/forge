package cli

import (
	"github.com/spf13/cobra"
)

var removeRunVcpkgCommandFunc func([]string) error

// NewRemoveCmd creates the remove command
func NewRemoveCmd(runVcpkgCommand func([]string) error) *cobra.Command {
	removeRunVcpkgCommandFunc = runVcpkgCommand

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a dependency",
		Long:  "Remove a dependency. Passes through to vcpkg remove command.",
		Aliases: []string{"rm"},
		RunE:   runRemove,
		Args:   cobra.MinimumNArgs(1),
	}

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) error {
	// Directly pass all arguments to vcpkg remove command
	// cpx remove <args> -> vcpkg remove <args>
	vcpkgArgs := []string{"remove"}
	vcpkgArgs = append(vcpkgArgs, args...)

	return removeRunVcpkgCommandFunc(vcpkgArgs)
}

// Remove handles the remove command - passes through to vcpkg
func Remove(args []string, runVcpkgCommand func([]string) error) {
	// This function is deprecated - use NewRemoveCmd instead
	// Kept for compatibility during migration
}
