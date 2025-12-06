package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ProjectConfig represents the cpx.yaml structure
type ProjectConfig struct {
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
		PreCommit []string `yaml:"precommit,omitempty"` // e.g., ["fmt", "lint"]
		PrePush   []string `yaml:"prepush,omitempty"`   // e.g., ["test"]
	} `yaml:"hooks,omitempty"`
	Features     map[string]FeatureConfig `yaml:"features,omitempty"`
	Dependencies []string                 `yaml:"dependencies,omitempty"` // Deprecated: Dependencies are now in vcpkg.json
}

// FeatureConfig represents feature-specific configuration
type FeatureConfig struct {
	Dependencies map[string]map[string]interface{} `yaml:"dependencies,omitempty"`
}

// CIConfig represents the cpx.ci structure for cross-compilation
type CIConfig struct {
	Targets []CITarget `yaml:"targets"`
	Build   CIBuild    `yaml:"build"`
	Output  string     `yaml:"output"`
}

// CITarget represents a cross-compilation target
type CITarget struct {
	Name       string `yaml:"name"`
	Dockerfile string `yaml:"dockerfile"`
	Image      string `yaml:"image"`
	Triplet    string `yaml:"triplet"`
	Platform   string `yaml:"platform"`
}

// CIBuild represents CI build configuration
type CIBuild struct {
	Type         string   `yaml:"type"`
	Optimization string   `yaml:"optimization"`
	Jobs         int      `yaml:"jobs"`
	CMakeArgs    []string `yaml:"cmake_args"`
	BuildArgs    []string `yaml:"build_args"`
}

// LoadProject loads the project configuration from cpx.yaml
func LoadProject(path string) (*ProjectConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	var config ProjectConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	return &config, nil
}

// SaveProject saves the project configuration to cpx.yaml
func SaveProject(config *ProjectConfig, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}

	return nil
}

// LoadCI loads the CI configuration from cpx.ci
func LoadCI(path string) (*CIConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config CIConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse cpx.ci: %w", err)
	}

	// Set defaults
	if config.Output == "" {
		config.Output = "out"
	}
	if config.Build.Type == "" {
		config.Build.Type = "Release"
	}
	if config.Build.Optimization == "" {
		config.Build.Optimization = "2"
	}

	return &config, nil
}
