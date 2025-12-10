package cli

import (
	"github.com/spf13/cobra"
)

var listRunVcpkgCommandFunc func([]string) error

// ListCmd creates the list command
func ListCmd(runVcpkgCommand func([]string) error) *cobra.Command {
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

func runList(_ *cobra.Command, args []string) error {
	// Directly pass all arguments to vcpkg list command
	vcpkgArgs := []string{"list"}
	vcpkgArgs = append(vcpkgArgs, args...)

	return listRunVcpkgCommandFunc(vcpkgArgs)
}
