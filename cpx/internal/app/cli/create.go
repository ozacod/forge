package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ozacod/cpx/internal/pkg/git"
	"github.com/ozacod/cpx/internal/pkg/template"
	"github.com/ozacod/cpx/pkg/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var createLoadConfigFunc func(string) (*CpxConfig, error)
var createGetVcpkgPathFunc func() (string, error)
var createSetupVcpkgProjectFunc func(string, string, bool, []string) error
var createGenerateVcpkgProjectFilesFromConfigFunc func(string, *CpxConfig, string, bool) error

// NewCreateCmd creates the create command
func NewCreateCmd(loadConfig func(string) (*CpxConfig, error), getVcpkgPath func() (string, error), setupVcpkgProject func(string, string, bool, []string) error, generateVcpkgProjectFilesFromConfig func(string, *CpxConfig, string, bool) error) *cobra.Command {
	createLoadConfigFunc = loadConfig
	createGetVcpkgPathFunc = getVcpkgPath
	createSetupVcpkgProjectFunc = setupVcpkgProject
	createGenerateVcpkgProjectFilesFromConfigFunc = generateVcpkgProjectFilesFromConfig

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new project from cpx.yaml template",
		Long:  "Create a new project from cpx.yaml template. The project name is required as the first argument.",
		RunE:  runCreate,
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.Flags().StringP("template", "t", "", "Path to cpx.yaml template file")
	cmd.Flags().Bool("lib", false, "Create a library project")

	return cmd
}

func runCreate(cmd *cobra.Command, args []string) error {
	templatePath, _ := cmd.Flags().GetString("template")
	isLib, _ := cmd.Flags().GetBool("lib")

	// Get project name from args
	var projectName string
	for _, arg := range args {
		// Skip flags
		if strings.HasPrefix(arg, "-") {
			continue
		}
		// Check for lib/exe keywords
		switch arg {
		case "lib", "library":
			isLib = true
			continue
		case "exe", "bin":
			isLib = false
			continue
		}
		// First non-flag, non-keyword is the project name
		if projectName == "" {
			projectName = arg
		}
	}

	if projectName == "" {
		return fmt.Errorf("project name required")
	}

	return createProject(projectName, templatePath, isLib, createLoadConfigFunc, createGetVcpkgPathFunc, createSetupVcpkgProjectFunc, createGenerateVcpkgProjectFilesFromConfigFunc)
}

// Create is kept for backward compatibility (if needed)
func Create(args []string, loadConfig func(string) (*CpxConfig, error), getVcpkgPath func() (string, error), setupVcpkgProject func(string, string, bool, []string) error, generateVcpkgProjectFilesFromConfig func(string, *CpxConfig, string, bool) error) {
	// This function is deprecated - use NewCreateCmd instead
	// Kept for compatibility during migration
}

// CpxConfig represents the cpx.yaml structure
// This is a local type for create command - main.go has its own CpxConfig
type CpxConfig struct {
	Package struct {
		Name        string   `yaml:"name"`
		Version     string   `yaml:"version"`
		CppStandard int      `yaml:"cpp_standard"`
		Authors     []string `yaml:"authors,omitempty"`
		Description string   `yaml:"description,omitempty"`
	} `yaml:"package"`
	Build struct {
		SharedLibs  bool   `yaml:"shared_libs"`
		ClangFormat string `yaml:"clang_format"`
		BuildType   string `yaml:"build_type,omitempty"`
		CxxFlags    string `yaml:"cxx_flags,omitempty"`
	} `yaml:"build"`
	VCS struct {
		Type string `yaml:"type,omitempty"` // "git" or "none"
	} `yaml:"vcs,omitempty"`
	PackageManager struct {
		Type string `yaml:"type,omitempty"` // "vcpkg" or "none"
	} `yaml:"package_manager,omitempty"`
	Testing struct {
		Framework string `yaml:"framework"`
	} `yaml:"testing"`
	Hooks struct {
		PreCommit []string `yaml:"precommit,omitempty"`
		PrePush   []string `yaml:"prepush,omitempty"`
	} `yaml:"hooks,omitempty"`
	Features     map[string]config.FeatureConfig `yaml:"features,omitempty"`
	Dependencies []string                        `yaml:"dependencies,omitempty"`
}

func createProject(projectName, templatePath string, isLib bool, loadConfig func(string) (*CpxConfig, error), getVcpkgPath func() (string, error), setupVcpkgProject func(string, string, bool, []string) error, generateVcpkgProjectFilesFromConfig func(string, *CpxConfig, string, bool) error) error {
	// Validate project name
	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`).MatchString(projectName) {
		return fmt.Errorf("invalid project name '%s': must start with letter and contain only letters, numbers, underscores, or hyphens", projectName)
	}

	// Check if directory already exists
	if _, err := os.Stat(projectName); err == nil {
		return fmt.Errorf("directory '%s' already exists", projectName)
	}

	// Create the new directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", projectName, err)
	}

	fmt.Printf("%s Creating project '%s'...%s\n", Cyan, projectName, Reset)

	// Load cpx.yaml template
	var cfg *CpxConfig
	if templatePath != "" {
		// Check if it's a template name (no path separators, no .yaml extension)
		actualTemplatePath := templatePath
		if !strings.Contains(templatePath, string(filepath.Separator)) && !strings.Contains(templatePath, "/") && !strings.HasSuffix(templatePath, ".yaml") && !strings.HasSuffix(templatePath, ".yml") {
			// It's a template name, download from GitHub templates/ folder
			tempDir := filepath.Join(os.TempDir(), "cpx-templates")
			if err := os.MkdirAll(tempDir, 0755); err != nil {
				return fmt.Errorf("failed to create temp directory: %w", err)
			}
			actualTemplatePath = filepath.Join(tempDir, templatePath+".yaml")

			// Download from GitHub
			fmt.Printf("%s Downloading template '%s' from GitHub...%s\n", Cyan, templatePath, Reset)
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
		// No template specified - try to download default.yaml from GitHub
		tempDir := filepath.Join(os.TempDir(), "cpx-templates")
		if err := os.MkdirAll(tempDir, 0755); err != nil {
			return fmt.Errorf("failed to create temp directory: %w", err)
		}
		defaultTemplatePath := filepath.Join(tempDir, "default.yaml")

		// Try to download from GitHub
		fmt.Printf("%s Downloading default template from GitHub...%s\n", Cyan, Reset)
		if err := template.DownloadFromGitHub("default.yaml", defaultTemplatePath); err != nil {
			// If download fails, use hardcoded defaults
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
			// Successfully downloaded, now load it
			cfg, err = loadConfig(defaultTemplatePath)
			if err != nil {
				return fmt.Errorf("failed to load default template: %w", err)
			}
			cfg.Package.Name = projectName
		}
	}

	// Determine if it's a library from config or flag
	if cfg.Build.SharedLibs {
		isLib = true
	}

	// Initialize git repository
	fmt.Printf("%s Initializing git repository...%s\n", Cyan, Reset)
	cmd := exec.Command("git", "init")
	cmd.Dir = projectName
	if err := cmd.Run(); err != nil {
		// Git init failure is not critical, just warn
		fmt.Printf("%s  Warning: Failed to initialize git repository: %v%s\n", Yellow, err, Reset)
	} else {
		fmt.Printf("%s Initialized git repository%s\n", Green, Reset)
	}

	fmt.Printf("%s Created project '%s'%s\n", Green, projectName, Reset)
	fmt.Printf("   Directory: %s\n", projectName)

	// Set up vcpkg integration
	dependencies := cfg.Dependencies
	if dependencies == nil {
		dependencies = []string{}
	}
	if err := setupVcpkgProject(projectName, projectName, isLib, dependencies); err != nil {
		return fmt.Errorf("failed to set up vcpkg project: %w", err)
	}

	// Generate project files with vcpkg integration
	fmt.Printf("\n%s Generating project files...%s\n", Cyan, Reset)
	if err := generateVcpkgProjectFilesFromConfig(projectName, cfg, projectName, isLib); err != nil {
		return fmt.Errorf("failed to generate project files: %w", err)
	}

	// Write cpx.yaml to the project so hooks can be installed
	configCopy := *cfg
	configCopy.Dependencies = nil // Don't save dependencies to cpx.yaml
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
		fmt.Printf("\n%s Installing git hooks...%s\n", Cyan, Reset)
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir) // Restore original directory

		if err := os.Chdir(projectName); err == nil {
			if err := git.InstallHooks(func(path string) (*config.ProjectConfig, error) {
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
				// Convert Features map
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
