package cli

import (
	"github.com/ozacod/cpx/internal/pkg/build"
	"github.com/spf13/cobra"
)

var benchSetupVcpkgEnvFunc func() error

// BenchCmd creates the bench command
func BenchCmd(setupVcpkgEnv func() error) *cobra.Command {
	benchSetupVcpkgEnvFunc = setupVcpkgEnv

	cmd := &cobra.Command{
		Use:   "bench",
		Short: "Build and run benchmarks",
		Long:  "Build the project benchmarks and run them.",
		Example: `  cpx bench            # Build + run all benchmarks
  cpx bench --verbose  # Show verbose output`,
		RunE: runBench,
	}

	cmd.Flags().BoolP("verbose", "v", false, "Show verbose build output")

	return cmd
}

func runBench(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")

	return build.RunBenchmarks(verbose, benchSetupVcpkgEnvFunc)
}
