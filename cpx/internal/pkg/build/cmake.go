package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
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
func ConfigureCMake(buildDir, buildType, cxxFlags string, verbose bool, setupVcpkgEnv func() error) error {
	// Set VCPKG_ROOT from cpx config if not already set
	if err := setupVcpkgEnv(); err != nil {
		return err
	}

	// Determine absolute path for shared vcpkg_installed directory
	// Determine absolute path for shared vcpkg_installed directory in .cache
	cwd, _ := os.Getwd()
	vcpkgInstalledDir := filepath.Join(cwd, ".cache", "vcpkg_installed")
	vcpkgInstallArg := "-DVCPKG_INSTALLED_DIR=" + vcpkgInstalledDir

	// Check if CMakePresets.json exists, use preset if available
	if _, err := os.Stat("CMakePresets.json"); err == nil {
		// Use "default" preset (VCPKG_ROOT is now set from config)
		// We explicitly pass -B to override the preset's binaryDir if it differs
		// We explicitly pass VCPKG_INSTALLED_DIR to share dependencies between build configs
		cmd := exec.Command("cmake", "--preset=default", "-B", buildDir, vcpkgInstallArg)
		// Ensure all vcpkg environment variables are in command environment
		cmd.Env = os.Environ()
		// Debug: Show environment variables being passed to CMake
		if os.Getenv("CPX_DEBUG") != "" {
			fmt.Printf("%s[DEBUG] CMake environment (preset):%s\n", "\033[36m", "\033[0m")
			for _, env := range cmd.Env {
				if strings.HasPrefix(env, "VCPKG_") {
					fmt.Printf("  %s\n", env)
				}
			}
		}
		if err := runCMakeConfigure(cmd, verbose); err != nil {
			return fmt.Errorf("cmake configure failed (preset 'default'): %w", err)
		}
	} else {
		// Fallback to traditional cmake configure
		cmakeArgs := []string{"-B", buildDir, "-DCMAKE_BUILD_TYPE=" + buildType, vcpkgInstallArg}

		if cxxFlags != "" {
			cmakeArgs = append(cmakeArgs, "-DCMAKE_CXX_FLAGS="+cxxFlags)
		}

		cmd := exec.Command("cmake", cmakeArgs...)
		// Ensure all vcpkg environment variables are in command environment
		cmd.Env = os.Environ()
		// Debug: Show environment variables being passed to CMake
		if os.Getenv("CPX_DEBUG") != "" {
			fmt.Printf("%s[DEBUG] CMake environment (traditional):%s\n", "\033[36m", "\033[0m")
			for _, env := range cmd.Env {
				if strings.HasPrefix(env, "VCPKG_") {
					fmt.Printf("  %s\n", env)
				}
			}
		}
		if err := runCMakeConfigure(cmd, verbose); err != nil {
			return fmt.Errorf("cmake configure failed: %w", err)
		}
	}

	return nil
}

// BuildProject builds the project using CMake
func BuildProject(release bool, jobs int, target string, clean bool, optLevel string, verbose bool, setupVcpkgEnv func() error) error {
	// Set VCPKG_ROOT from cpx config if not already set
	if err := setupVcpkgEnv(); err != nil {
		return err
	}

	// Get project name from CMakeLists.txt (optional, for display only)
	projectName := GetProjectNameFromCMakeLists()
	if projectName == "" {
		projectName = "project"
	}

	// Determine build output directory based on optimization/release
	outDirName := "debug"
	if optLevel != "" {
		outDirName = "O" + optLevel
	} else if release {
		outDirName = "release"
	}

	// Use hidden cache directory for build artifacts
	// .cache/native/<variant>
	cacheBuildDir := filepath.Join(".cache", "native", outDirName)
	// Final executables go to .bin/native/<variant>
	finalBuildDir := filepath.Join(".bin", "native", outDirName)

	if clean {
		if verbose {
			fmt.Printf("%s  Cleaning build directory...%s\n", colorCyan, colorReset)
		}
		os.RemoveAll(cacheBuildDir)
		os.RemoveAll(finalBuildDir)
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(cacheBuildDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache build dir: %w", err)
	}

	// Determine build type and optimization
	buildType, cxxFlags := DetermineBuildType(release, optLevel)
	optLabel := cxxFlags
	if optLabel == "" {
		optLabel = "default (CMake)"
	}

	fmt.Printf("\n%s▸ Build%s %s %s(%s)%s %s[opt: %s]%s\n",
		colorCyan, colorReset, projectName, colorGray, buildType, colorReset,
		colorGray, optLabel, colorReset)

	// Check if configure is needed (if cache doesn't exist or CMakeCache.txt logic)
	// Actually we check CMakeCache.txt inside cacheBuildDir
	needsConfigure := false
	if _, err := os.Stat(filepath.Join(cacheBuildDir, "CMakeCache.txt")); os.IsNotExist(err) {
		needsConfigure = true
	}

	// Determine total steps
	totalSteps := 1
	currentStep := 0
	if needsConfigure {
		totalSteps = 3 // configure + build + copy
	} else {
		totalSteps = 2 // build + copy
	}

	if needsConfigure {
		currentStep++
		if verbose {
			fmt.Printf("%s  • Configuring CMake%s\n", colorCyan, colorReset)
		} else {
			fmt.Printf("\r\033[2K%s[%d/%d]%s Configuring...", colorCyan, currentStep, totalSteps, colorReset)
		}

		// Configure CMake
		// We use -B to specify the build directory in .cache

		// Determine absolute path for shared vcpkg_installed directory
		cwd, _ := os.Getwd()
		vcpkgInstalledDir := filepath.Join(cwd, ".cache", "native", "vcpkg_installed")
		vcpkgInstallArg := "-DVCPKG_INSTALLED_DIR=" + vcpkgInstalledDir

		// Check if CMakePresets.json exists, use preset if available
		if _, err := os.Stat("CMakePresets.json"); err == nil {
			// Use "default" preset (VCPKG_ROOT is now set from config)
			// Pass -B explicitly to override preset binaryDir if needed, or ensure it goes to our cache
			// Also pass VCPKG_INSTALLED_DIR to force shared vcpkg location
			cmd := exec.Command("cmake", "--preset=default", "-B", cacheBuildDir, vcpkgInstallArg)
			cmd.Env = os.Environ()
			if err := runCMakeConfigure(cmd, verbose); err != nil {
				fmt.Println()
				return fmt.Errorf("cmake configure failed (preset 'default'): %w", err)
			}
		} else {
			// Fallback to traditional cmake configure
			cmd := exec.Command("cmake", "-B", cacheBuildDir, "-DCMAKE_BUILD_TYPE="+buildType, vcpkgInstallArg)
			if cxxFlags != "" {
				cmd.Args = append(cmd.Args, "-DCMAKE_CXX_FLAGS="+cxxFlags)
			}
			cmd.Env = os.Environ()
			if err := runCMakeConfigure(cmd, verbose); err != nil {
				fmt.Println()
				return fmt.Errorf("cmake configure failed: %w", err)
			}
		}

		if !verbose {
			fmt.Printf("\r\033[2K%s[%d/%d]%s Configured ✓\n", colorCyan, currentStep, totalSteps, colorReset)
		}
	}

	// Build
	buildStart := time.Now()
	buildArgs := []string{"--build", cacheBuildDir, "--config", buildType}

	if jobs > 0 {
		buildArgs = append(buildArgs, "--parallel", fmt.Sprintf("%d", jobs))
	} else {
		buildArgs = append(buildArgs, "--parallel", fmt.Sprintf("%d", runtime.NumCPU()))
	}

	if target != "" {
		buildArgs = append(buildArgs, "--target", target)
	}

	currentStep++
	if err := runCMakeBuild(buildArgs, verbose, currentStep, totalSteps); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Copy artifacts to final build directory
	if err := os.MkdirAll(finalBuildDir, 0755); err != nil {
		return fmt.Errorf("failed to create final build dir: %w", err)
	}

	executables, err := FindExecutables(cacheBuildDir)
	if err == nil {
		for _, exe := range executables {
			dest := filepath.Join(finalBuildDir, filepath.Base(exe))
			// Copy file
			input, err := os.ReadFile(exe)
			if err != nil {
				continue
			}
			os.WriteFile(dest, input, 0755)
		}
	}

	fmt.Printf("%s  ✔ Build complete%s %s[%s]%s\n", colorGreen, colorReset, colorGray, time.Since(buildStart).Round(10*time.Millisecond), colorReset)
	fmt.Printf("  Artifacts in: %s/\n\n", finalBuildDir)
	return nil
}
