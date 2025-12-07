package cli

import (
	"github.com/ozacod/cpx/internal/pkg/build"
	"github.com/spf13/cobra"
)

var runSetupVcpkgEnvFunc func() error

// NewRunCmd creates the run command
func NewRunCmd(setupVcpkgEnv func() error) *cobra.Command {
	runSetupVcpkgEnvFunc = setupVcpkgEnv

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Build and run the project",
		Long:  "Build the project (configuring CMake if needed) and run the chosen executable. Arguments after -- are passed to the binary.",
		Example: `  cpx run                 # Debug build by default
  cpx run --release        # Release build, then run
  cpx run --target app -- --flag value`,
		RunE: runRun,
	}

	cmd.Flags().Bool("release", false, "Build in release mode (-O2). Default is debug")
	cmd.Flags().String("target", "", "Executable target to run (useful if multiple)")
	cmd.Flags().Bool("verbose", false, "Show full CMake/Ninja output during build")

	return cmd
}

func runRun(cmd *cobra.Command, args []string) error {
	release, _ := cmd.Flags().GetBool("release")
	target, _ := cmd.Flags().GetString("target")
	verbose, _ := cmd.Flags().GetBool("verbose")

	return build.RunProject(release, target, args, verbose, runSetupVcpkgEnvFunc)
}

// Run is kept for backward compatibility (if needed)
func Run(args []string, setupVcpkgEnv func() error) {
	// This function is deprecated - use NewRunCmd instead
	// Kept for compatibility during migration
}
