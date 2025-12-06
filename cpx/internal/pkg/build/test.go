package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// RunTests runs the project tests
func RunTests(verbose bool, filter string, setupVcpkgEnv func() error) error {
	// Set VCPKG_ROOT from cpx config if not already set
	if err := setupVcpkgEnv(); err != nil {
		return err
	}

	projectName := GetProjectNameFromCMakeLists()
	if projectName == "" {
		return fmt.Errorf("failed to get project name from CMakeLists.txt")
	}
	fmt.Printf("%s Running tests for '%s'...%s\n", "\033[36m", projectName, "\033[0m")

	buildDir := "build"

	// Configure CMake if needed
	if _, err := os.Stat(filepath.Join(buildDir, "CMakeCache.txt")); os.IsNotExist(err) {
		fmt.Printf("%s  Configuring CMake...%s\n", "\033[36m", "\033[0m")
		// Check if CMakePresets.json exists, use preset if available
		if _, err := os.Stat("CMakePresets.json"); err == nil {
			// Use "default" preset (VCPKG_ROOT is now set from config)
			cmd := exec.Command("cmake", "--preset=default")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			// Ensure VCPKG_ROOT is in command environment
			cmd.Env = os.Environ()
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("cmake configure failed (preset 'default'): %w", err)
			}
		} else {
			// Fallback to traditional cmake configure
			cmd := exec.Command("cmake", "-B", buildDir)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("cmake configure failed: %w", err)
			}
		}
	}

	// Build tests
	fmt.Printf("%s Building tests...%s\n", "\033[36m", "\033[0m")
	buildCmd := exec.Command("cmake", "--build", buildDir, "--target", projectName+"_tests")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build tests: %w", err)
	}

	// Run tests with CTest
	fmt.Printf("%s Running tests...%s\n", "\033[36m", "\033[0m")
	ctestArgs := []string{"--test-dir", buildDir}

	if verbose {
		ctestArgs = append(ctestArgs, "--verbose")
	}

	if filter != "" {
		ctestArgs = append(ctestArgs, "--output-on-failure", "-R", filter)
	} else {
		ctestArgs = append(ctestArgs, "--output-on-failure")
	}

	ctestCmd := exec.Command("ctest", ctestArgs...)
	ctestCmd.Stdout = os.Stdout
	ctestCmd.Stderr = os.Stderr

	if err := ctestCmd.Run(); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}

	fmt.Printf("%s All tests passed!%s\n", "\033[32m", "\033[0m")
	return nil
}
