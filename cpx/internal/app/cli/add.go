package cli

import (
	"github.com/spf13/cobra"
)

var addRunVcpkgCommandFunc func([]string) error

// NewAddCmd creates the add command
func NewAddCmd(runVcpkgCommand func([]string) error) *cobra.Command {
	addRunVcpkgCommandFunc = runVcpkgCommand

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a dependency",
		Long:  "Add a dependency. Passes through to vcpkg add command.",
		RunE:  runAdd,
		Args:  cobra.MinimumNArgs(1),
	}

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) error {
	if err := requireVcpkgProject("cpx add"); err != nil {
		return err
	}

	// Directly pass all arguments to vcpkg add command
	// cpx add <args> -> vcpkg add <args>
	vcpkgArgs := []string{"add"}
	vcpkgArgs = append(vcpkgArgs, args...)

	return addRunVcpkgCommandFunc(vcpkgArgs)
}

// Add handles the add command - passes through to vcpkg
func Add(args []string, runVcpkgCommand func([]string) error) {
	// This function is deprecated - use NewAddCmd instead
	// Kept for compatibility during migration
}
