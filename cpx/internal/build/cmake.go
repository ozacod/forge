package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// GetProjectNameFromCMakeLists extracts project name from CMakeLists.txt in current directory
func GetProjectNameFromCMakeLists() string {
	return GetProjectNameFromCMakeListsInDir(".")
}

// GetProjectNameFromCMakeListsInDir extracts project name from CMakeLists.txt in the given directory
func GetProjectNameFromCMakeListsInDir(dir string) string {
	cmakeListsPath := filepath.Join(dir, "CMakeLists.txt")
	data, err := os.ReadFile(cmakeListsPath)
	if err != nil {
		return ""
	}

	// Look for: project(PROJECT_NAME ...)
	re := regexp.MustCompile(`project\s*\(\s*([^\s\)]+)`)
	matches := re.FindStringSubmatch(string(data))
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

// DetermineBuildType determines the CMake build type and CXX flags based on release flag and optimization level.
// Returns (buildType, cxxFlags)
func DetermineBuildType(release bool, optLevel string) (string, string) {
	buildType := "Debug"
	cxxFlags := ""

	if release {
		buildType = "Release"
	}

	// Handle optimization level
	switch optLevel {
	case "0":
		cxxFlags = "-O0"
		buildType = "Debug"
	case "1":
		cxxFlags = "-O1"
		buildType = "RelWithDebInfo"
	case "2":
		cxxFlags = "-O2"
		buildType = "Release"
	case "3":
		cxxFlags = "-O3"
		buildType = "Release"
	case "s":
		cxxFlags = "-Os"
		buildType = "MinSizeRel"
	case "fast":
		cxxFlags = "-Ofast"
		buildType = "Release"
	}

	return buildType, cxxFlags
}

// ConfigureCMake configures CMake for the project
func ConfigureCMake(buildDir, buildType, cxxFlags string, setupVcpkgEnv func() error) error {
	// Set VCPKG_ROOT from cpx config if not already set
	if err := setupVcpkgEnv(); err != nil {
		return err
	}

	// Check if CMakePresets.json exists, use preset if available
	if _, err := os.Stat("CMakePresets.json"); err == nil {
		// Use "default" preset (VCPKG_ROOT is now set from config)
		cmd := exec.Command("cmake", "--preset=default")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		// Ensure all vcpkg environment variables are in command environment
		cmd.Env = os.Environ()
		// Debug: Show environment variables being passed to CMake
		if os.Getenv("FORGE_DEBUG") != "" {
			fmt.Printf("%s[DEBUG] CMake environment (preset):%s\n", "\033[36m", "\033[0m")
			for _, env := range cmd.Env {
				if strings.HasPrefix(env, "VCPKG_") {
					fmt.Printf("  %s\n", env)
				}
			}
		}
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("cmake configure failed (preset 'default'): %w", err)
		}
	} else {
		// Fallback to traditional cmake configure
		cmakeArgs := []string{"-B", buildDir, "-DCMAKE_BUILD_TYPE=" + buildType}

		if cxxFlags != "" {
			cmakeArgs = append(cmakeArgs, "-DCMAKE_CXX_FLAGS="+cxxFlags)
		}

		cmd := exec.Command("cmake", cmakeArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		// Ensure all vcpkg environment variables are in command environment
		cmd.Env = os.Environ()
		// Debug: Show environment variables being passed to CMake
		if os.Getenv("FORGE_DEBUG") != "" {
			fmt.Printf("%s[DEBUG] CMake environment (traditional):%s\n", "\033[36m", "\033[0m")
			for _, env := range cmd.Env {
				if strings.HasPrefix(env, "VCPKG_") {
					fmt.Printf("  %s\n", env)
				}
			}
		}
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("cmake configure failed: %w", err)
		}
	}

	return nil
}

// BuildProject builds the project using CMake
func BuildProject(release bool, jobs int, target string, clean bool, optLevel string, setupVcpkgEnv func() error) error {
	// Set VCPKG_ROOT from cpx config if not already set
	if err := setupVcpkgEnv(); err != nil {
		return err
	}

	// Get project name from CMakeLists.txt (optional, for display only)
	projectName := GetProjectNameFromCMakeLists()
	if projectName == "" {
		projectName = "project"
	}

	buildDir := "build"

	// Clean if requested
	if clean {
		fmt.Printf("%s Cleaning build directory...%s\n", "\033[36m", "\033[0m")
		os.RemoveAll(buildDir)
	}

	// Determine build type and optimization
	buildType, cxxFlags := DetermineBuildType(release, optLevel)
	optInfo := ""
	if cxxFlags != "" {
		optInfo = fmt.Sprintf(" [%s]", cxxFlags)
	}

	fmt.Printf("%s¨ Building '%s' (%s%s)...%s\n", "\033[36m", projectName, buildType, optInfo, "\033[0m")

	// Configure CMake if needed or if clean was done
	needsConfigure := clean
	if _, err := os.Stat(filepath.Join(buildDir, "CMakeCache.txt")); os.IsNotExist(err) {
		needsConfigure = true
	}

	if needsConfigure {
		fmt.Printf("%s™  Configuring CMake...%s\n", "\033[36m", "\033[0m")
		if err := ConfigureCMake(buildDir, buildType, cxxFlags, setupVcpkgEnv); err != nil {
			return err
		}
	}

	// Build
	fmt.Printf("%s Compiling...%s\n", "\033[36m", "\033[0m")
	buildArgs := []string{"--build", buildDir, "--config", buildType}

	if jobs > 0 {
		buildArgs = append(buildArgs, "--parallel", fmt.Sprintf("%d", jobs))
	} else {
		buildArgs = append(buildArgs, "--parallel", fmt.Sprintf("%d", runtime.NumCPU()))
	}

	if target != "" {
		buildArgs = append(buildArgs, "--target", target)
	}

	buildCmd := exec.Command("cmake", buildArgs...)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Printf("%s Build complete!%s\n", "\033[32m", "\033[0m")
	return nil
}
