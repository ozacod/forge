package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ozacod/cpx/internal/pkg/build"
	"github.com/spf13/cobra"
)

var setupVcpkgEnvFunc func() error

// BuildCmd creates the build command
func BuildCmd(setupVcpkgEnv func() error) *cobra.Command {
	setupVcpkgEnvFunc = setupVcpkgEnv

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Compile the project",
		Long: `Compile the project. Automatically detects project type:
  - vcpkg/CMake projects: Uses CMake with vcpkg toolchain
  - Bazel projects: Uses bazel build`,
		Example: `  cpx build              # Debug build (default)
  cpx build --release    # Release build (-O2)
  cpx build -O3          # Maximum optimization
  cpx build -j 8         # Use 8 parallel jobs
  cpx build --clean      # Clean rebuild
  cpx build --watch      # Watch for changes and rebuild`,
		RunE: runBuild,
	}

	cmd.Flags().BoolP("release", "r", false, "Release build (-O2). Default is debug")
	cmd.Flags().Bool("debug", false, "Debug build (-O0). Default; kept for compatibility")
	cmd.Flags().IntP("jobs", "j", 0, "Parallel jobs for build (0 = auto)")
	cmd.Flags().String("target", "", "Specific target to build")
	cmd.Flags().BoolP("clean", "c", false, "Clean build directory before building")
	cmd.Flags().StringP("opt", "O", "", "Override optimization level: 0,1,2,3,s,fast")
	cmd.Flags().BoolP("watch", "w", false, "Watch for file changes and rebuild automatically")
	cmd.Flags().Bool("verbose", false, "Show full build output")

	return cmd
}

func runBuild(cmd *cobra.Command, args []string) error {
	release, _ := cmd.Flags().GetBool("release")
	jobs, _ := cmd.Flags().GetInt("jobs")
	target, _ := cmd.Flags().GetString("target")
	clean, _ := cmd.Flags().GetBool("clean")
	optLevel, _ := cmd.Flags().GetString("opt")
	watch, _ := cmd.Flags().GetBool("watch")
	verbose, _ := cmd.Flags().GetBool("verbose")

	projectType := DetectProjectType()

	switch projectType {
	case ProjectTypeBazel:
		if watch {
			fmt.Printf("%sWatch mode not yet supported for Bazel projects%s\n", Yellow, Reset)
			return nil
		}
		return runBazelBuild(release, target, clean, verbose, optLevel)
	case ProjectTypeMeson:
		if watch {
			fmt.Printf("%sWatch mode not yet supported for Meson projects%s\n", Yellow, Reset)
			return nil
		}
		return runMesonBuild(release, target, clean, verbose, optLevel)
	case ProjectTypeVcpkg:
		if watch {
			return build.WatchAndBuild(release, jobs, target, optLevel, verbose, setupVcpkgEnvFunc)
		}
		return build.BuildProject(release, jobs, target, clean, optLevel, verbose, setupVcpkgEnvFunc)
	default:
		// Fall back to CMake build even without vcpkg.json
		if watch {
			return build.WatchAndBuild(release, jobs, target, optLevel, verbose, setupVcpkgEnvFunc)
		}
		return build.BuildProject(release, jobs, target, clean, optLevel, verbose, setupVcpkgEnvFunc)
	}
}

func runBazelBuild(release bool, target string, clean bool, verbose bool, optLevel string) error {
	// Clean if requested
	if clean {
		fmt.Printf("%sCleaning Bazel build...%s\n", Cyan, Reset)
		cleanCmd := execCommand("bazel", "clean")
		cleanCmd.Stdout = os.Stdout
		cleanCmd.Stderr = os.Stderr
		if err := cleanCmd.Run(); err != nil {
			return fmt.Errorf("bazel clean failed: %w", err)
		}
		// Also remove build directory
		os.RemoveAll("build")
	}

	// Build args
	bazelArgs := []string{"build"}

	// Handle optimization level - optLevel takes precedence over release flag
	optLabel := "debug"
	switch optLevel {
	case "0":
		bazelArgs = append(bazelArgs, "--copt=-O0", "-c", "dbg")
		optLabel = "-O0 (debug)"
	case "1":
		bazelArgs = append(bazelArgs, "--copt=-O1", "-c", "opt")
		optLabel = "-O1"
	case "2":
		bazelArgs = append(bazelArgs, "--copt=-O2", "-c", "opt")
		optLabel = "-O2"
	case "3":
		bazelArgs = append(bazelArgs, "--copt=-O3", "-c", "opt")
		optLabel = "-O3"
	case "s":
		bazelArgs = append(bazelArgs, "--copt=-Os", "-c", "opt")
		optLabel = "-Os (size)"
	case "fast":
		bazelArgs = append(bazelArgs, "--copt=-Ofast", "-c", "opt")
		optLabel = "-Ofast"
	default:
		// No explicit opt level, use release/debug config
		if release {
			bazelArgs = append(bazelArgs, "--config=release")
			optLabel = "release"
		} else {
			bazelArgs = append(bazelArgs, "--config=debug")
			optLabel = "debug"
		}
	}

	// Add target or default to //...
	if target != "" {
		bazelArgs = append(bazelArgs, target)
	} else {
		bazelArgs = append(bazelArgs, "//...")
	}

	fmt.Printf("%sBuilding with Bazel [%s]...%s\n", Cyan, optLabel, Reset)
	if verbose {
		fmt.Printf("  Running: bazel %v\n", bazelArgs)
	} else {
		// Suppress progress bars for cleaner output (like vcpkg)
		// Use hidden symlinks (.bazel-bin, .bazel-out, etc.)
		bazelArgs = append(bazelArgs, "--noshow_progress", "--symlink_prefix=.bazel-")
	}

	buildCmd := execCommand("bazel", bazelArgs...)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("bazel build failed: %w", err)
	}

	// Determine output directory based on config
	outDirName := "debug"
	if optLevel != "" {
		outDirName = "O" + optLevel
	} else if release {
		outDirName = "release"
	}
	outputDir := filepath.Join(".bin", "native", outDirName)

	// Copy artifacts to build/<config>/ directory
	// Remove existing build artifacts for this config first
	os.RemoveAll(outputDir)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create build directory: %w", err)
	}

	// Copy executables and libraries from bazel-bin to build/<config>/
	fmt.Printf("%sCopying artifacts to %s/...%s\n", Cyan, outputDir, Reset)

	// Create a script to copy with the correct output directory variable
	script := fmt.Sprintf(`
		# Find the bazel-bin symlink (prefer hidden .bazel-bin)
		BAZEL_BIN=""
		if [ -L ".bazel-bin" ] || [ -d ".bazel-bin" ]; then
			BAZEL_BIN=".bazel-bin"
		elif [ -L "bazel-bin" ] || [ -d "bazel-bin" ]; then
			BAZEL_BIN="bazel-bin"
		elif [ -L ".bin" ] || [ -d ".bin" ]; then
			BAZEL_BIN=".bin"
		fi

		if [ -z "$BAZEL_BIN" ]; then
			echo "No bazel-bin found"
			exit 0
		fi

		# Copy executables from src/ directory (where cc_binary targets are placed)
		find -L "$BAZEL_BIN/src" -maxdepth 1 -type f -perm +111 ! -name "*.params" ! -name "*.sh" ! -name "*.cppmap" ! -name "*.repo_mapping" ! -name "*runfiles*" ! -name "*.d" -exec cp -f {} %[1]s/ \; 2>/dev/null || true

		# Also copy from root of bazel-bin (for root aliases)
		find -L "$BAZEL_BIN" -maxdepth 1 -type f -perm +111 ! -name "*.params" ! -name "*.sh" ! -name "*.cppmap" ! -name "*.repo_mapping" ! -name "*runfiles*" ! -name "*.d" -exec cp -f {} %[1]s/ \; 2>/dev/null || true

		# Copy libraries from src/
		find -L "$BAZEL_BIN/src" -maxdepth 1 -type f \( -name "*.a" -o -name "*.so" -o -name "*.dylib" \) -exec cp -f {} %[1]s/ \; 2>/dev/null || true

		# Make copied files writable (Bazel creates read-only files)
		chmod -R u+w %[1]s/ 2>/dev/null || true

		# List what was copied
		ls %[1]s/ 2>/dev/null || true
	`, outputDir)

	copyCmd := execCommand("bash", "-c", script)
	copyCmd.Stdout = os.Stdout
	copyCmd.Stderr = os.Stderr
	copyCmd.Run() // Ignore errors - may have no artifacts

	fmt.Printf("%s✓ Build successful%s\n", Green, Reset)
	fmt.Printf("  Artifacts in: %s/\n", outputDir)
	return nil
}

func runMesonBuild(release bool, target string, clean bool, verbose bool, optLevel string) error {
	buildDir := "builddir"

	// Determine build type and optimization from flags
	buildType := "debug"
	optimization := "0" // Meson optimization: 0, 1, 2, 3, s
	optLabel := "debug"

	switch optLevel {
	case "0":
		buildType = "debug"
		optimization = "0"
		optLabel = "-O0 (debug)"
	case "1":
		buildType = "debugoptimized"
		optimization = "1"
		optLabel = "-O1"
	case "2":
		buildType = "release"
		optimization = "2"
		optLabel = "-O2"
	case "3":
		buildType = "release"
		optimization = "3"
		optLabel = "-O3"
	case "s":
		buildType = "minsize"
		optimization = "s"
		optLabel = "-Os (size)"
	case "fast":
		// Meson doesn't have -Ofast directly, use -O3 with custom flags
		buildType = "release"
		optimization = "3"
		optLabel = "-Ofast"
	default:
		// No explicit opt level, use release/debug
		if release {
			buildType = "release"
			optimization = "2"
			optLabel = "release"
		} else {
			buildType = "debug"
			optimization = "0"
			optLabel = "debug"
		}
	}

	// Clean if requested or if optimization changed
	if clean {
		fmt.Printf("%sCleaning Meson build...%s\n", Cyan, Reset)
		os.RemoveAll(buildDir)
	}

	// Check if build directory exists (needs setup)
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		fmt.Printf("%sSetting up Meson build directory [%s]...%s\n", Cyan, optLabel, Reset)
		setupArgs := []string{"setup", buildDir}
		setupArgs = append(setupArgs, "--buildtype="+buildType)
		setupArgs = append(setupArgs, "--optimization="+optimization)
		if optLevel == "fast" {
			// Add -ffast-math for -Ofast equivalent
			setupArgs = append(setupArgs, "-Dc_args=-ffast-math", "-Dcpp_args=-ffast-math")
		}
		setupCmd := execCommand("meson", setupArgs...)
		setupCmd.Stdout = os.Stdout
		setupCmd.Stderr = os.Stderr
		if err := setupCmd.Run(); err != nil {
			return fmt.Errorf("meson setup failed: %w", err)
		}
	} else {
		// Build directory exists, reconfigure if optimization changed
		fmt.Printf("%sReconfiguring Meson [%s]...%s\n", Cyan, optLabel, Reset)
		reconfigArgs := []string{"configure", buildDir}
		reconfigArgs = append(reconfigArgs, "--buildtype="+buildType)
		reconfigArgs = append(reconfigArgs, "--optimization="+optimization)
		if optLevel == "fast" {
			reconfigArgs = append(reconfigArgs, "-Dc_args=-ffast-math", "-Dcpp_args=-ffast-math")
		}
		reconfigCmd := execCommand("meson", reconfigArgs...)
		reconfigCmd.Stdout = os.Stdout
		reconfigCmd.Stderr = os.Stderr
		// Ignore reconfigure errors - may fail if no changes needed
		reconfigCmd.Run()
	}

	// Build
	fmt.Printf("%sBuilding with Meson...%s\n", Cyan, Reset)
	compileArgs := []string{"compile", "-C", buildDir}
	if target != "" {
		compileArgs = append(compileArgs, target)
	}
	if verbose {
		compileArgs = append(compileArgs, "-v")
	}
	buildCmd := execCommand("meson", compileArgs...)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("meson compile failed: %w", err)
	}

	// Determine output directory based on config
	outDirName := "debug"
	if optLevel != "" {
		outDirName = "O" + optLevel
	} else if release {
		outDirName = "release"
	}
	outputDir := filepath.Join(".bin", "native", outDirName)

	// Copy artifacts to output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	fmt.Printf("%sCopying artifacts to %s/...%s\n", Cyan, outputDir, Reset)
	copyCmd := execCommand("bash", "-c", fmt.Sprintf(`
		# Meson places executables in subdirectories (src/, bench/, etc.)
		# Search in builddir/src/ first (main executables)
		if [ -d "builddir/src" ]; then
			find builddir/src -maxdepth 1 -type f -perm +111 ! -name "*.p" ! -name "*_test" -exec cp {} %[1]s/ \; 2>/dev/null || true
		fi

		# Also check builddir root for executables
		find builddir -maxdepth 1 -type f -perm +111 ! -name "*.p" ! -name "*_test" -exec cp {} %[1]s/ \; 2>/dev/null || true

		# Copy libraries from builddir and subdirectories
		find builddir -maxdepth 2 -type f \( -name "*.a" -o -name "*.so" -o -name "*.dylib" \) -exec cp {} %[1]s/ \; 2>/dev/null || true

		# List what was copied
		ls %[1]s/ 2>/dev/null || true
	`, outputDir))
	copyCmd.Stdout = os.Stdout
	copyCmd.Stderr = os.Stderr
	copyCmd.Run()

	fmt.Printf("%s✓ Build successful%s\n", Green, Reset)
	fmt.Printf("  Artifacts in: %s/\n", outputDir)
	return nil
}
