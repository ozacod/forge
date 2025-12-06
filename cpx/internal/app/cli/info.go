package cli

import (
	"github.com/spf13/cobra"
)

var infoRunVcpkgCommandFunc func([]string) error

// NewInfoCmd creates the info command
func NewInfoCmd(runVcpkgCommand func([]string) error) *cobra.Command {
	infoRunVcpkgCommandFunc = runVcpkgCommand

	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show detailed library information",
		Long:  "Show detailed library information. Passes through to vcpkg show command.",
		RunE:  runInfo,
		Args:  cobra.MinimumNArgs(1),
	}

	return cmd
}

func runInfo(cmd *cobra.Command, args []string) error {
	// Directly pass all arguments to vcpkg show command
	// cpx info <package> -> vcpkg show <package>
	vcpkgArgs := []string{"show"}
	vcpkgArgs = append(vcpkgArgs, args...)

	return infoRunVcpkgCommandFunc(vcpkgArgs)
}

// Info handles the info command - passes through to vcpkg
func Info(args []string, runVcpkgCommand func([]string) error) {
	// This function is deprecated - use NewInfoCmd instead
	// Kept for compatibility during migration
}
