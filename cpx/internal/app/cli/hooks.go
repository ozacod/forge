package cli

import (
	"github.com/ozacod/cpx/pkg/config"
	"github.com/ozacod/cpx/internal/pkg/git"
	"github.com/spf13/cobra"
)

var hooksLoadConfigFunc func(string) (*config.ProjectConfig, error)

// NewHooksCmd creates the hooks command
func NewHooksCmd(loadConfig func(string) (*config.ProjectConfig, error)) *cobra.Command {
	hooksLoadConfigFunc = loadConfig

	cmd := &cobra.Command{
		Use:   "hooks",
		Short: "Install git hooks",
		Long: `Install git hooks for code quality and automation:
   pre-commit   - Format code and run linters before commit
   pre-push     - Run tests and security checks before push
   commit-msg   - Validate commit message format
   post-merge   - Update dependencies if vcpkg.json changed`,
	}

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install git hooks",
		Long:  "Install git hooks configured in cpx.yaml.",
		RunE:  runHooksInstall,
	}
	cmd.AddCommand(installCmd)

	return cmd
}

func runHooksInstall(cmd *cobra.Command, args []string) error {
	return git.InstallHooks(hooksLoadConfigFunc, DefaultCfgFile)
}

// Hooks is kept for backward compatibility (if needed)
func Hooks(args []string, loadConfig func(string) (*config.ProjectConfig, error)) {
	// This function is deprecated - use NewHooksCmd instead
	// Kept for compatibility during migration
}
