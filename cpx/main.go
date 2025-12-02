package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ozacod/cpx/cmd"
	"github.com/ozacod/cpx/cmd/root"
	"github.com/ozacod/cpx/internal/config"
	"github.com/ozacod/cpx/internal/git"
	"github.com/ozacod/cpx/internal/template"
	"github.com/ozacod/cpx/internal/vcpkg"

	"gopkg.in/yaml.v3"
)

const (
	Version        = cmd.Version
	DefaultServer  = cmd.DefaultServer
	DefaultCfgFile = cmd.DefaultCfgFile
	LockFile       = cmd.LockFile
)

var vcpkgClient *vcpkg.Client

func getVcpkgClient() (*vcpkg.Client, error) {
	if vcpkgClient == nil {
		var err error
		vcpkgClient, err = vcpkg.NewClient()
		if err != nil {
			return nil, err
		}
	}
	return vcpkgClient, nil
}

func setupVcpkgEnv() error {
	client, err := getVcpkgClient()
	if err != nil {
		return err
	}

	err = client.SetupEnv()
	if err != nil {
		return err
	}

	if os.Getenv("FORGE_DEBUG") != "" {
		fmt.Printf("%s[DEBUG] VCPKG Environment:%s\n", cmd.Cyan, cmd.Reset)
		fmt.Printf("  VCPKG_ROOT=%s\n", os.Getenv("VCPKG_ROOT"))
		fmt.Printf("  VCPKG_FEATURE_FLAGS=%s\n", os.Getenv("VCPKG_FEATURE_FLAGS"))
		fmt.Printf("  VCPKG_DISABLE_REGISTRY_UPDATE=%s\n", os.Getenv("VCPKG_DISABLE_REGISTRY_UPDATE"))
	}

	return nil
}

const (
	Reset   = cmd.Reset
	Red     = cmd.Red
	Green   = cmd.Green
	Yellow  = cmd.Yellow
	Blue    = cmd.Blue
	Magenta = cmd.Magenta
	Cyan    = cmd.Cyan
	Bold    = cmd.Bold
)

type CpxConfig = config.ProjectConfig

func getVcpkgPath() (string, error) {
	client, err := getVcpkgClient()
	if err != nil {
		return "", err
	}
	return client.GetPath()
}

func runVcpkgCommand(args []string) error {
	client, err := getVcpkgClient()
	if err != nil {
		return err
	}
	return client.RunCommand(args)
}

func main() {
	rootCmd := root.GetRootCmd()

	// Register all commands
	rootCmd.AddCommand(cmd.NewBuildCmd(setupVcpkgEnv))
	rootCmd.AddCommand(cmd.NewRunCmd(setupVcpkgEnv))
	rootCmd.AddCommand(cmd.NewTestCmd(setupVcpkgEnv))
	rootCmd.AddCommand(cmd.NewCleanCmd())
	rootCmd.AddCommand(cmd.NewCreateCmd(
		func(path string) (*cmd.CpxConfig, error) {
			cfg, err := loadConfig(path)
			if err != nil {
				return nil, err
			}
			// Convert main.CpxConfig to cmd.CpxConfig
			result := &cmd.CpxConfig{}
			result.Package.Name = cfg.Package.Name
			result.Package.Version = cfg.Package.Version
			result.Package.CppStandard = cfg.Package.CppStandard
			result.Package.Authors = cfg.Package.Authors
			result.Package.Description = cfg.Package.Description
			result.Build.SharedLibs = cfg.Build.SharedLibs
			result.Build.ClangFormat = cfg.Build.ClangFormat
			result.Build.BuildType = cfg.Build.BuildType
			result.Build.CxxFlags = cfg.Build.CxxFlags
			result.Testing.Framework = cfg.Testing.Framework
			result.Hooks.PreCommit = cfg.Hooks.PreCommit
			result.Hooks.PrePush = cfg.Hooks.PrePush
			if cfg.Features != nil {
				result.Features = make(map[string]config.FeatureConfig)
				for k, v := range cfg.Features {
					result.Features[k] = config.FeatureConfig{
						Dependencies: v.Dependencies,
					}
				}
			}
			return result, nil
		},
		getVcpkgPath,
		setupVcpkgProject,
		func(targetDir string, cfg *cmd.CpxConfig, projectName string, isLib bool) error {
			// Convert cmd.CpxConfig to main.CpxConfig
			mainCfg := &CpxConfig{}
			mainCfg.Package.Name = cfg.Package.Name
			mainCfg.Package.Version = cfg.Package.Version
			mainCfg.Package.CppStandard = cfg.Package.CppStandard
			mainCfg.Package.Authors = cfg.Package.Authors
			mainCfg.Package.Description = cfg.Package.Description
			mainCfg.Build.SharedLibs = cfg.Build.SharedLibs
			mainCfg.Build.ClangFormat = cfg.Build.ClangFormat
			mainCfg.Build.BuildType = cfg.Build.BuildType
			mainCfg.Build.CxxFlags = cfg.Build.CxxFlags
			mainCfg.Testing.Framework = cfg.Testing.Framework
			mainCfg.Hooks.PreCommit = cfg.Hooks.PreCommit
			mainCfg.Hooks.PrePush = cfg.Hooks.PrePush
			if cfg.Features != nil {
				mainCfg.Features = make(map[string]config.FeatureConfig)
				for k, v := range cfg.Features {
					mainCfg.Features[k] = config.FeatureConfig{
						Dependencies: v.Dependencies,
					}
				}
			}
			return generateVcpkgProjectFilesFromConfig(targetDir, mainCfg, projectName, isLib)
		}))
	rootCmd.AddCommand(cmd.NewAddCmd(runVcpkgCommand))
	rootCmd.AddCommand(cmd.NewRemoveCmd(runVcpkgCommand))
	rootCmd.AddCommand(cmd.NewListCmd(runVcpkgCommand))
	rootCmd.AddCommand(cmd.NewSearchCmd(runVcpkgCommand))
	rootCmd.AddCommand(cmd.NewInfoCmd(runVcpkgCommand))
	rootCmd.AddCommand(cmd.NewFmtCmd())
	rootCmd.AddCommand(cmd.NewLintCmd(setupVcpkgEnv, getVcpkgPath))
	rootCmd.AddCommand(cmd.NewFlawfinderCmd())
	rootCmd.AddCommand(cmd.NewCppcheckCmd())
	rootCmd.AddCommand(cmd.NewAnalyzeCmd(setupVcpkgEnv, getVcpkgPath))
	rootCmd.AddCommand(cmd.NewCheckCmd(setupVcpkgEnv))
	rootCmd.AddCommand(cmd.NewDocCmd())
	rootCmd.AddCommand(cmd.NewReleaseCmd())
	rootCmd.AddCommand(cmd.NewUpgradeCmd())
	rootCmd.AddCommand(cmd.NewConfigCmd())
	rootCmd.AddCommand(cmd.NewCICmd())
	rootCmd.AddCommand(cmd.NewHooksCmd(func(path string) (*config.ProjectConfig, error) {
		cfg, err := loadConfig(path)
		if err != nil {
			return nil, err
		}
		// Convert CpxConfig to config.ProjectConfig
		projectCfg := &config.ProjectConfig{}
		projectCfg.Package.Name = cfg.Package.Name
		projectCfg.Package.Version = cfg.Package.Version
		projectCfg.Package.CppStandard = cfg.Package.CppStandard
		projectCfg.Package.Authors = cfg.Package.Authors
		projectCfg.Package.Description = cfg.Package.Description
		projectCfg.Build.SharedLibs = cfg.Build.SharedLibs
		projectCfg.Build.ClangFormat = cfg.Build.ClangFormat
		projectCfg.Build.BuildType = cfg.Build.BuildType
		projectCfg.Build.CxxFlags = cfg.Build.CxxFlags
		projectCfg.Testing.Framework = cfg.Testing.Framework
		projectCfg.Hooks.PreCommit = cfg.Hooks.PreCommit
		projectCfg.Hooks.PrePush = cfg.Hooks.PrePush
		if cfg.Features != nil {
			projectCfg.Features = make(map[string]config.FeatureConfig)
			for k, v := range cfg.Features {
				projectCfg.Features[k] = config.FeatureConfig{
					Dependencies: v.Dependencies,
				}
			}
		}
		return projectCfg, nil
	}))
	rootCmd.AddCommand(cmd.NewUpdateCmd())

	// Handle vcpkg passthrough for unknown commands
	// Check if command exists before executing
	if len(os.Args) > 1 {
		command := os.Args[1]
		// Skip version/help flags - cobra handles these
		if command != "-v" && command != "--version" && command != "version" &&
			command != "-h" && command != "--help" && command != "help" {
			// Check if it's a known command
			found := false
			for _, c := range rootCmd.Commands() {
				if c.Name() == command || contains(c.Aliases, command) {
					found = true
					break
				}
			}
			// If not found, try vcpkg passthrough
			if !found {
				if err := runVcpkgCommand(os.Args[1:]); err != nil {
					fmt.Fprintf(os.Stderr, "%sError:%s Failed to run vcpkg command: %v\n", Red, Reset, err)
					fmt.Fprintf(os.Stderr, "Make sure vcpkg is installed and configured: cpx config set-vcpkg-root <path>\n")
					os.Exit(1)
				}
				return
			}
		}
	}

	// Execute root command
	root.Execute()
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// printUsage is no longer needed - cobra handles help automatically

// cmdCreate is no longer needed - use cmd.NewCreateCmd instead

func createProject(projectName, templatePath string, isLib bool) error {
	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`).MatchString(projectName) {
		return fmt.Errorf("invalid project name '%s': must start with letter and contain only letters, numbers, underscores, or hyphens", projectName)
	}

	if _, err := os.Stat(projectName); err == nil {
		return fmt.Errorf("directory '%s' already exists", projectName)
	}

	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", projectName, err)
	}

	fmt.Printf("%s� Creating project '%s'...%s\n", Cyan, projectName, Reset)

	var cfg *CpxConfig
	if templatePath != "" {
		actualTemplatePath := templatePath
		if !strings.Contains(templatePath, string(filepath.Separator)) && !strings.Contains(templatePath, "/") && !strings.HasSuffix(templatePath, ".yaml") && !strings.HasSuffix(templatePath, ".yml") {
			tempDir := filepath.Join(os.TempDir(), "cpx-templates")
			if err := os.MkdirAll(tempDir, 0755); err != nil {
				return fmt.Errorf("failed to create temp directory: %w", err)
			}
			actualTemplatePath = filepath.Join(tempDir, templatePath+".yaml")

			fmt.Printf("%s� Downloading template '%s' from GitHub...%s\n", Cyan, templatePath, Reset)
			if err := template.DownloadFromGitHub(templatePath+".yaml", actualTemplatePath); err != nil {
				return fmt.Errorf("failed to download template '%s' from GitHub: %w", templatePath, err)
			}
			fmt.Printf("%s Using template: %s%s\n", Cyan, templatePath, Reset)
		} else {
			fmt.Printf("%s Using template: %s%s\n", Cyan, templatePath, Reset)
		}

		var err error
		cfg, err = loadConfig(actualTemplatePath)
		if err != nil {
			return fmt.Errorf("failed to load template file '%s': %w", actualTemplatePath, err)
		}
		cfg.Package.Name = projectName
	} else {
		tempDir := filepath.Join(os.TempDir(), "cpx-templates")
		if err := os.MkdirAll(tempDir, 0755); err != nil {
			return fmt.Errorf("failed to create temp directory: %w", err)
		}
		defaultTemplatePath := filepath.Join(tempDir, "default.yaml")

		fmt.Printf("%s� Downloading default template from GitHub...%s\n", Cyan, Reset)
		if err := template.DownloadFromGitHub("default.yaml", defaultTemplatePath); err != nil {
			fmt.Printf("%s  Could not load default template, using built-in defaults...%s\n", Yellow, Reset)
			cfg = &CpxConfig{}
			cfg.Package.Name = projectName
			cfg.Package.Version = "0.1.0"
			cfg.Package.CppStandard = 17
			cfg.Build.SharedLibs = false
			cfg.Build.ClangFormat = "Google"
			cfg.Testing.Framework = "googletest"
			cfg.Hooks.PreCommit = []string{"fmt", "lint"}
			cfg.Hooks.PrePush = []string{"test"}
		} else {
			cfg, err = loadConfig(defaultTemplatePath)
			if err != nil {
				return fmt.Errorf("failed to load default template: %w", err)
			}
			cfg.Package.Name = projectName
		}
	}

	if cfg.Build.SharedLibs {
		isLib = true
	}

	fmt.Printf("%s Initializing git repository...%s\n", Cyan, Reset)
	cmd := exec.Command("git", "init")
	cmd.Dir = projectName
	if err := cmd.Run(); err != nil {
		fmt.Printf("%s  Warning: Failed to initialize git repository: %v%s\n", Yellow, err, Reset)
	} else {
		fmt.Printf("%s Initialized git repository%s\n", Green, Reset)
	}

	fmt.Printf("%s Created project '%s'%s\n", Green, projectName, Reset)
	fmt.Printf("   Directory: %s\n", projectName)

	dependencies := cfg.Dependencies
	if dependencies == nil {
		dependencies = []string{}
	}
	if err := setupVcpkgProject(projectName, projectName, isLib, dependencies); err != nil {
		return fmt.Errorf("failed to set up vcpkg project: %w", err)
	}

	fmt.Printf("\n%s Generating project files...%s\n", Cyan, Reset)
	if err := generateVcpkgProjectFilesFromConfig(projectName, cfg, projectName, isLib); err != nil {
		return fmt.Errorf("failed to generate project files: %w", err)
	}

	configCopy := *cfg
	configCopy.Dependencies = nil
	data, err := yaml.Marshal(&configCopy)
	if err == nil {
		header := "# cpx.yaml - C++ Project Configuration\n# Dependencies are managed in vcpkg.json (use 'vcpkg add port <package>')\n\n"
		data = append([]byte(header), data...)
		cpxYamlPath := filepath.Join(projectName, DefaultCfgFile)
		if err := os.WriteFile(cpxYamlPath, data, 0644); err == nil {
			fmt.Printf("%s   cpx.yaml%s\n", Green, Reset)
		}
	}

	// Install hooks if configured
	if len(cfg.Hooks.PreCommit) > 0 || len(cfg.Hooks.PrePush) > 0 {
		fmt.Printf("\n%s� Installing git hooks...%s\n", Cyan, Reset)
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir) // Restore original directory

		if err := os.Chdir(projectName); err == nil {
			if err := git.InstallHooks(func(path string) (*config.ProjectConfig, error) {
				cfg, err := loadConfig(path)
				if err != nil {
					return nil, err
				}
				projectCfg := &config.ProjectConfig{}
				projectCfg.Package.Name = cfg.Package.Name
				projectCfg.Package.Version = cfg.Package.Version
				projectCfg.Package.CppStandard = cfg.Package.CppStandard
				projectCfg.Package.Authors = cfg.Package.Authors
				projectCfg.Package.Description = cfg.Package.Description
				projectCfg.Build.SharedLibs = cfg.Build.SharedLibs
				projectCfg.Build.ClangFormat = cfg.Build.ClangFormat
				projectCfg.Build.BuildType = cfg.Build.BuildType
				projectCfg.Build.CxxFlags = cfg.Build.CxxFlags
				projectCfg.Testing.Framework = cfg.Testing.Framework
				projectCfg.Hooks.PreCommit = cfg.Hooks.PreCommit
				projectCfg.Hooks.PrePush = cfg.Hooks.PrePush
				if cfg.Features != nil {
					projectCfg.Features = make(map[string]config.FeatureConfig)
					for k, v := range cfg.Features {
						projectCfg.Features[k] = config.FeatureConfig{
							Dependencies: v.Dependencies,
						}
					}
				}
				return projectCfg, nil
			}, DefaultCfgFile); err != nil {
				// Non-fatal error, just warn
				fmt.Printf("%s  Warning: Failed to install hooks: %v%s\n", Yellow, err, Reset)
			}
		}
	}

	fmt.Printf("\n%s Project '%s' ready!%s\n\n", Green, projectName, Reset)
	fmt.Printf("Next steps:\n")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Printf("  %scpx build%s       # Compile the project\n", Cyan, Reset)
	fmt.Printf("  %scpx run%s         # Build and run\n", Cyan, Reset)

	return nil
}

func setupVcpkgProject(targetDir, _ string, _ bool, dependencies []string) error {
	vcpkgPath, err := getVcpkgPath()
	if err != nil {
		return fmt.Errorf("vcpkg not configured: %w\n   Run: cpx config set-vcpkg-root <path>", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(targetDir); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	fmt.Printf("%s Initializing vcpkg.json...%s\n", Cyan, Reset)

	cmd := exec.Command(vcpkgPath, "new", "--application")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	for i, env := range cmd.Env {
		if strings.HasPrefix(env, "VCPKG_ROOT=") {
			cmd.Env = append(cmd.Env[:i], cmd.Env[i+1:]...)
			break
		}
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize vcpkg.json: %w", err)
	}

	if len(dependencies) > 0 {
		fmt.Printf("%s Adding dependencies from template...%s\n", Cyan, Reset)
		for _, dep := range dependencies {
			if dep == "" {
				continue
			}
			fmt.Printf("   Adding %s...\n", dep)
			// vcpkg add requires "port" or "artifact" as the second argument
			// We're adding ports (packages), so use "port"
			addCmd := exec.Command(vcpkgPath, "add", "port", dep)
			addCmd.Stdout = os.Stdout
			addCmd.Stderr = os.Stderr
			addCmd.Env = cmd.Env // Use same environment
			if err := addCmd.Run(); err != nil {
				fmt.Printf("%s  Warning: Failed to add dependency '%s': %v%s\n", Yellow, dep, err, Reset)
				// Continue with other dependencies even if one fails
			}
		}
	}

	return nil
}

// generateVcpkgProjectFilesFromConfig generates CMake files with vcpkg integration from config struct
func generateVcpkgProjectFilesFromConfig(targetDir string, cfg *CpxConfig, projectName string, isLib bool) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	cppStandard := cfg.Package.CppStandard
	if cppStandard == 0 {
		cppStandard = 17
	}

	projectVersion := cfg.Package.Version
	if projectVersion == "" {
		projectVersion = "0.1.0"
	}

	// Get dependencies from vcpkg.json, not cpx.yaml
	dependencies, err := getDependenciesFromVcpkgJsonLocal(targetDir)
	if err != nil {
		// If vcpkg.json doesn't exist or can't be read, use empty list
		dependencies = []string{}
	}

	// Create directories
	dirs := []string{
		"include/" + projectName,
		"src",
		"tests",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(targetDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Generate CMakeLists.txt with vcpkg integration
	cmakeLists := generateVcpkgCMakeLists(projectName, cppStandard, dependencies, !isLib, cfg.Testing.Framework != "" && cfg.Testing.Framework != "none", cfg.Testing.Framework, projectVersion)
	if err := os.WriteFile(filepath.Join(targetDir, "CMakeLists.txt"), []byte(cmakeLists), 0644); err != nil {
		return fmt.Errorf("failed to write CMakeLists.txt: %w", err)
	}

	// Generate CMakePresets.json
	// Assumes VCPKG_ROOT environment variable is set when cpx build is called
	cmakePresets := generateCMakePresets()
	if err := os.WriteFile(filepath.Join(targetDir, "CMakePresets.json"), []byte(cmakePresets), 0644); err != nil {
		return fmt.Errorf("failed to write CMakePresets.json: %w", err)
	}

	// Generate version.hpp
	versionHpp := generateVersionHpp(projectName, projectVersion)
	if err := os.WriteFile(filepath.Join(targetDir, "include/"+projectName+"/version.hpp"), []byte(versionHpp), 0644); err != nil {
		return fmt.Errorf("failed to write version.hpp: %w", err)
	}

	// Generate header file
	libHeader := generateLibHeader(projectName)
	if err := os.WriteFile(filepath.Join(targetDir, "include/"+projectName+"/"+projectName+".hpp"), []byte(libHeader), 0644); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Generate source files
	if !isLib {
		mainCpp := generateMainCpp(projectName, dependencies)
		if err := os.WriteFile(filepath.Join(targetDir, "src/main.cpp"), []byte(mainCpp), 0644); err != nil {
			return fmt.Errorf("failed to write main.cpp: %w", err)
		}
	}

	libSource := generateLibSource(projectName, dependencies)
	if err := os.WriteFile(filepath.Join(targetDir, "src/"+projectName+".cpp"), []byte(libSource), 0644); err != nil {
		return fmt.Errorf("failed to write source: %w", err)
	}

	// Generate README
	readme := generateVcpkgReadme(projectName, dependencies, cppStandard, isLib)
	if err := os.WriteFile(filepath.Join(targetDir, "README.md"), []byte(readme), 0644); err != nil {
		return fmt.Errorf("failed to write README: %w", err)
	}

	// Generate .gitignore
	gitignore := generateGitignore()
	if err := os.WriteFile(filepath.Join(targetDir, ".gitignore"), []byte(gitignore), 0644); err != nil {
		return fmt.Errorf("failed to write .gitignore: %w", err)
	}

	// Generate .clang-format
	clangFormatStyle := cfg.Build.ClangFormat
	if clangFormatStyle == "" {
		clangFormatStyle = "Google"
	}
	clangFormat := generateClangFormat(clangFormatStyle)
	if err := os.WriteFile(filepath.Join(targetDir, ".clang-format"), []byte(clangFormat), 0644); err != nil {
		return fmt.Errorf("failed to write .clang-format: %w", err)
	}

	// Generate test files if testing framework is enabled
	if cfg.Testing.Framework != "" && cfg.Testing.Framework != "none" {
		// Generate tests/CMakeLists.txt
		testCMake := generateTestCMake(projectName, dependencies, cfg.Testing.Framework)
		if err := os.WriteFile(filepath.Join(targetDir, "tests/CMakeLists.txt"), []byte(testCMake), 0644); err != nil {
			return fmt.Errorf("failed to write tests/CMakeLists.txt: %w", err)
		}

		// Generate tests/test_main.cpp
		testMain := generateTestMain(projectName, dependencies, cfg.Testing.Framework)
		if err := os.WriteFile(filepath.Join(targetDir, "tests/test_main.cpp"), []byte(testMain), 0644); err != nil {
			return fmt.Errorf("failed to write tests/test_main.cpp: %w", err)
		}
	}

	// Generate cpx.ci with empty targets
	cpxCI := generateCpxCI()
	if err := os.WriteFile(filepath.Join(targetDir, "cpx.ci"), []byte(cpxCI), 0644); err != nil {
		return fmt.Errorf("failed to write cpx.ci: %w", err)
	}

	return nil
}

// removeDependenciesFromYaml removes the dependencies section from YAML content
// getDependenciesFromVcpkgJson reads dependencies from vcpkg.json
// getDependenciesFromVcpkgJsonLocal is a local helper for createProject
func getDependenciesFromVcpkgJsonLocal(projectDir string) ([]string, error) {
	vcpkgJsonPath := filepath.Join(projectDir, "vcpkg.json")

	// Check if vcpkg.json exists
	if _, err := os.Stat(vcpkgJsonPath); os.IsNotExist(err) {
		return []string{}, nil
	}

	// Read vcpkg.json
	data, err := os.ReadFile(vcpkgJsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read vcpkg.json: %w", err)
	}

	var vcpkgJson map[string]interface{}
	if err := json.Unmarshal(data, &vcpkgJson); err != nil {
		return nil, fmt.Errorf("failed to parse vcpkg.json: %w", err)
	}

	// Extract dependencies
	deps, ok := vcpkgJson["dependencies"].([]interface{})
	if !ok {
		return []string{}, nil
	}

	dependencies := make([]string, 0, len(deps))
	for _, dep := range deps {
		if depStr, ok := dep.(string); ok {
			dependencies = append(dependencies, depStr)
		}
	}

	return dependencies, nil
}

func loadConfig(path string) (*CpxConfig, error) {
	return config.LoadProject(path)
}
