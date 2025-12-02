package cmd

import (
	"github.com/ozacod/cpx/internal/quality"
	"github.com/spf13/cobra"
)

// vcpkgAdapter implements quality.VcpkgSetup interface
type vcpkgAdapter struct {
	setupEnv func() error
	getPath  func() (string, error)
}

func (v *vcpkgAdapter) SetupVcpkgEnv() error {
	return v.setupEnv()
}

func (v *vcpkgAdapter) GetVcpkgPath() (string, error) {
	return v.getPath()
}

var lintSetupVcpkgEnvFunc func() error
var lintGetVcpkgPathFunc func() (string, error)

// NewLintCmd creates the lint command
func NewLintCmd(setupVcpkgEnv func() error, getVcpkgPath func() (string, error)) *cobra.Command {
	lintSetupVcpkgEnvFunc = setupVcpkgEnv
	lintGetVcpkgPathFunc = getVcpkgPath

	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Run clang-tidy static analysis",
		Long:  "Run clang-tidy static analysis. Use --fix to automatically fix issues.",
		RunE:  runLint,
	}

	cmd.Flags().Bool("fix", false, "Automatically fix issues")

	return cmd
}

func runLint(cmd *cobra.Command, args []string) error {
	fix, _ := cmd.Flags().GetBool("fix")

	vcpkg := &vcpkgAdapter{
		setupEnv: lintSetupVcpkgEnvFunc,
		getPath:  lintGetVcpkgPathFunc,
	}
	return quality.LintCode(fix, vcpkg)
}

// Lint is kept for backward compatibility (if needed)
func Lint(args []string, setupVcpkgEnv func() error, getVcpkgPath func() (string, error)) {
	// This function is deprecated - use NewLintCmd instead
	// Kept for compatibility during migration
}
