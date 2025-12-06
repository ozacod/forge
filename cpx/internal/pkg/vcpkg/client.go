package vcpkg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ozacod/cpx/pkg/config"
)

// Client handles vcpkg operations
type Client struct {
	globalConfig *config.GlobalConfig
}

// NewClient creates a new vcpkg client
func NewClient() (*Client, error) {
	globalConfig, err := config.LoadGlobal()
	if err != nil {
		return nil, fmt.Errorf("failed to load global config: %w", err)
	}
	return &Client{globalConfig: globalConfig}, nil
}

// SetupEnv sets VCPKG_ROOT and VCPKG_FEATURE_FLAGS environment variables from cpx config
// This ensures CMake presets can find vcpkg and uses manifest mode consistently
func (c *Client) SetupEnv() error {
	// Set VCPKG_ROOT if not already set and we have it in config
	if os.Getenv("VCPKG_ROOT") == "" {
		if c.globalConfig.VcpkgRoot == "" {
			return fmt.Errorf("vcpkg_root not set in config. Run: cpx config set-vcpkg-root <path>")
		}
		if err := os.Setenv("VCPKG_ROOT", c.globalConfig.VcpkgRoot); err != nil {
			return fmt.Errorf("failed to set VCPKG_ROOT: %w", err)
		}
	}

	// Set VCPKG_FEATURE_FLAGS=manifests if not already set
	// This ensures consistent behavior between host and Docker builds
	// Manifest mode is required when using vcpkg.json
	if os.Getenv("VCPKG_FEATURE_FLAGS") == "" {
		if err := os.Setenv("VCPKG_FEATURE_FLAGS", "manifests"); err != nil {
			return fmt.Errorf("failed to set VCPKG_FEATURE_FLAGS: %w", err)
		}
	}

	// Disable registry update check to speed up builds
	// Packages will still be cached and reused
	if os.Getenv("VCPKG_DISABLE_REGISTRY_UPDATE") == "" {
		if err := os.Setenv("VCPKG_DISABLE_REGISTRY_UPDATE", "1"); err != nil {
			return fmt.Errorf("failed to set VCPKG_DISABLE_REGISTRY_UPDATE: %w", err)
		}
	}

	return nil
}

// GetPath returns the path to the vcpkg executable
func (c *Client) GetPath() (string, error) {
	vcpkgRoot := c.globalConfig.VcpkgRoot

	// If not set in config, check environment variable as fallback
	if vcpkgRoot == "" {
		if envRoot := os.Getenv("VCPKG_ROOT"); envRoot != "" {
			vcpkgRoot = envRoot
		}
	}

	if vcpkgRoot == "" {
		return "", fmt.Errorf("vcpkg_root not set in config. Run: cpx config set-vcpkg-root <path>")
	}

	// Convert to absolute path
	absVcpkgRoot, err := filepath.Abs(vcpkgRoot)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute vcpkg root path: %w", err)
	}

	vcpkgPath := filepath.Join(absVcpkgRoot, "vcpkg")
	if runtime.GOOS == "windows" {
		vcpkgPath += ".exe"
	}

	if _, err := os.Stat(vcpkgPath); os.IsNotExist(err) {
		return "", fmt.Errorf("vcpkg not found at %s. Make sure vcpkg is installed and bootstrapped", vcpkgPath)
	}

	return vcpkgPath, nil
}

// RunCommand runs a vcpkg command
func (c *Client) RunCommand(args []string) error {
	vcpkgPath, err := c.GetPath()
	if err != nil {
		return err
	}

	cmd := exec.Command(vcpkgPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	// Remove VCPKG_ROOT from environment to use the one from config
	cmd.Env = os.Environ()
	for i, env := range cmd.Env {
		if strings.HasPrefix(env, "VCPKG_ROOT=") {
			cmd.Env = append(cmd.Env[:i], cmd.Env[i+1:]...)
			break
		}
	}
	return cmd.Run()
}

