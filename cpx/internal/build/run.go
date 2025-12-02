package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// RunProject builds and runs the project
func RunProject(release bool, target string, execArgs []string, setupVcpkgEnv func() error) error {
	// Set VCPKG_ROOT from cpx config if not already set
	if err := setupVcpkgEnv(); err != nil {
		return err
	}

	// Get project name from CMakeLists.txt (optional, for display only)
	projectName := GetProjectNameFromCMakeLists()
	if projectName == "" {
		projectName = "project"
	}

	buildType, _ := DetermineBuildType(release, "")

	fmt.Printf("%s¨ Building '%s' (%s)...%s\n", "\033[36m", projectName, buildType, "\033[0m")

	// Configure CMake if needed
	buildDir := "build"
	if _, err := os.Stat(filepath.Join(buildDir, "CMakeCache.txt")); os.IsNotExist(err) {
		fmt.Printf("%s™  Configuring CMake...%s\n", "\033[36m", "\033[0m")

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
			cmd := exec.Command("cmake", "-B", buildDir, "-DCMAKE_BUILD_TYPE="+buildType)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("cmake configure failed: %w", err)
			}
		}
	}

	// Build
	fmt.Printf("%s Compiling...%s\n", "\033[36m", "\033[0m")
	buildCmd := exec.Command("cmake", "--build", buildDir, "--config", buildType)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Find and run executable
	execName := projectName
	if runtime.GOOS == "windows" {
		execName += ".exe"
	}

	execPath := filepath.Join(buildDir, execName)
	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		// Try to find executable in build directory
		entries, err := os.ReadDir(buildDir)
		if err != nil {
			return fmt.Errorf("failed to read build directory: %w", err)
		}

		found := false
		for _, entry := range entries {
			if !entry.IsDir() {
				info, err := entry.Info()
				if err != nil {
					continue
				}
				// Check if it's executable (Unix) or .exe (Windows)
				if runtime.GOOS == "windows" {
					if strings.HasSuffix(entry.Name(), ".exe") {
						execPath = filepath.Join(buildDir, entry.Name())
						found = true
						break
					}
				} else {
					if info.Mode()&0111 != 0 && !strings.Contains(entry.Name(), "_test") && !strings.Contains(entry.Name(), "_tests") {
						execPath = filepath.Join(buildDir, entry.Name())
						found = true
						break
					}
				}
			}
		}

		if !found {
			return fmt.Errorf("executable not found in %s. Make sure the project builds an executable", buildDir)
		}
	}

	fmt.Printf("%s€ Running '%s'...%s\n", "\033[36m", execPath, "\033[0m")
	runCmd := exec.Command(execPath, execArgs...)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	runCmd.Stdin = os.Stdin
	return runCmd.Run()
}
