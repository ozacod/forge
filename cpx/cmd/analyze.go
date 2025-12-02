package cmd

import (
	"github.com/ozacod/cpx/internal/quality"
	"github.com/spf13/cobra"
)

// vcpkgAdapter implements quality.VcpkgSetup interface
type analyzeVcpkgAdapter struct {
	setupEnv func() error
	getPath  func() (string, error)
}

func (v *analyzeVcpkgAdapter) SetupVcpkgEnv() error {
	return v.setupEnv()
}

func (v *analyzeVcpkgAdapter) GetVcpkgPath() (string, error) {
	return v.getPath()
}

var analyzeSetupVcpkgEnvFunc func() error
var analyzeGetVcpkgPathFunc func() (string, error)

// NewAnalyzeCmd creates the analyze command
func NewAnalyzeCmd(setupVcpkgEnv func() error, getVcpkgPath func() (string, error)) *cobra.Command {
	analyzeSetupVcpkgEnvFunc = setupVcpkgEnv
	analyzeGetVcpkgPathFunc = getVcpkgPath

	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Run comprehensive code analysis and generate HTML report",
		Long:  "Run comprehensive code analysis using cppcheck, clang-tidy, and flawfinder. Generates a combined HTML report (analyze.html).",
		RunE:  runAnalyze,
		Args:  cobra.ArbitraryArgs,
	}

	cmd.Flags().String("output", "analyze.html", "Output HTML file path")
	cmd.Flags().Bool("skip-cppcheck", false, "Skip Cppcheck analysis")
	cmd.Flags().Bool("skip-lint", false, "Skip clang-tidy analysis")
	cmd.Flags().Bool("skip-flawfinder", false, "Skip Flawfinder analysis")

	return cmd
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	output, _ := cmd.Flags().GetString("output")
	skipCppcheck, _ := cmd.Flags().GetBool("skip-cppcheck")
	skipLint, _ := cmd.Flags().GetBool("skip-lint")
	skipFlawfinder, _ := cmd.Flags().GetBool("skip-flawfinder")

	// Get remaining args as target directories (default to current directory)
	targets := args
	if len(targets) == 0 {
		targets = []string{"."}
	}

	vcpkg := &analyzeVcpkgAdapter{
		setupEnv: analyzeSetupVcpkgEnvFunc,
		getPath:  analyzeGetVcpkgPathFunc,
	}

	return quality.RunComprehensiveAnalysis(output, skipCppcheck, skipLint, skipFlawfinder, targets, vcpkg)
}

