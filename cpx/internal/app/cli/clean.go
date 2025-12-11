package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// CleanCmd creates the clean command
func CleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Remove build artifacts",
		Long: `Remove build artifacts. Automatically detects project type:
  - Bazel: runs 'bazel clean' and removes symlinks (.bin, .out, bazel-*)
  - Meson: removes builddir/
  - CMake/vcpkg: removes build/

Use --all to also remove additional generated files.`,
		Example: `  cpx clean         # Clean build artifacts
  cpx clean --all   # Also remove all generated files`,
		RunE: runClean,
	}

	cmd.Flags().Bool("all", false, "Also remove generated files")

	return cmd
}

func runClean(cmd *cobra.Command, _ []string) error {
	all, _ := cmd.Flags().GetBool("all")

	projectType := DetectProjectType()

	switch projectType {
	case ProjectTypeBazel:
		return cleanBazel(all)
	case ProjectTypeMeson:
		return cleanMeson(all)
	default:
		// CMake/vcpkg or unknown - clean generic build directory
		return cleanCMake(all)
	}
}

func cleanBazel(all bool) error {
	fmt.Printf("%sCleaning Bazel project...%s\n", Cyan, Reset)

	// Run bazel clean
	cleanCmd := execCommand("bazel", "clean")
	cleanCmd.Stdout = os.Stdout
	cleanCmd.Stderr = os.Stderr
	if err := cleanCmd.Run(); err != nil {
		fmt.Printf("%s⚠ bazel clean failed (may not be initialized)%s\n", Yellow, Reset)
	} else {
		fmt.Printf("%s✓ Ran bazel clean%s\n", Green, Reset)
	}

	// Remove common build output directory
	removeDir("build")

	// Remove Bazel symlinks
	bazelSymlinks := []string{".bin", ".out", ".testlogs"}
	for _, symlink := range bazelSymlinks {
		if _, err := os.Lstat(symlink); err == nil {
			fmt.Printf("%s  Removing %s...%s\n", Cyan, symlink, Reset)
			os.RemoveAll(symlink)
		}
	}

	// Remove bazel-* symlinks (bazel-bin, bazel-out, bazel-testlogs, bazel-<project>)
	entries, err := os.ReadDir(".")
	if err == nil {
		for _, entry := range entries {
			matched, _ := filepath.Match("bazel-*", entry.Name())
			if matched {
				fmt.Printf("%s  Removing %s...%s\n", Cyan, entry.Name(), Reset)
				os.RemoveAll(entry.Name())
			}
		}
	}

	if all {
		// Remove additional Bazel artifacts
		removeDir(".bazel")
		removeDir("external")
	}

	fmt.Printf("%s✓ Bazel project cleaned%s\n", Green, Reset)
	return nil
}

func cleanMeson(all bool) error {
	fmt.Printf("%sCleaning Meson project...%s\n", Cyan, Reset)

	// Remove builddir
	removeDir("builddir")

	// Remove common build output directory
	removeDir("build")

	if all {
		// Remove additional Meson artifacts
		removeDir("subprojects/packagecache")

		// Remove build-* directories
		entries, err := os.ReadDir(".")
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					matched, _ := filepath.Match("build-*", entry.Name())
					if matched {
						fmt.Printf("%s  Removing %s...%s\n", Cyan, entry.Name(), Reset)
						os.RemoveAll(entry.Name())
					}
				}
			}
		}
	}

	fmt.Printf("%s✓ Meson project cleaned%s\n", Green, Reset)
	return nil
}

func cleanCMake(all bool) error {
	fmt.Printf("%sCleaning CMake/vcpkg project...%s\n", Cyan, Reset)

	// Remove bin directory (artifacts)
	removeDir(filepath.Join(".bin", "native"))

	// Remove intermediate build directories (keep vcpkg_installed unless --all)
	// We iterate common variants instead of blowing away .cache/native
	variants := []string{"debug", "release", "O0", "O1", "O2", "O3", "Os", "Ofast"}
	for _, v := range variants {
		removeDir(filepath.Join(".cache", "native", v))
	}

	if all {
		// Clean everything including vcpkg dependencies and CI artifacts
		dirsToRemove := []string{
			filepath.Join(".cache", "native"),
			filepath.Join(".cache", "ci"),
			filepath.Join(".bin", "ci"),
			"out",
			"cmake-build-debug",
			"cmake-build-release",
		}
		for _, dir := range dirsToRemove {
			removeDir(dir)
		}

		// Remove build-* directories
		entries, err := os.ReadDir(".")
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					matched, _ := filepath.Match("build-*", entry.Name())
					if matched {
						fmt.Printf("%s  Removing %s...%s\n", Cyan, entry.Name(), Reset)
						os.RemoveAll(entry.Name())
					}
				}
			}
		}
	}

	fmt.Printf("%s✓ CMake project cleaned%s\n", Green, Reset)
	return nil
}

func removeDir(path string) {
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("%s  Removing %s...%s\n", Cyan, path, Reset)
		if err := os.RemoveAll(path); err != nil {
			fmt.Printf("%s⚠ Failed to remove %s: %v%s\n", Yellow, path, err, Reset)
		}
	}
}
