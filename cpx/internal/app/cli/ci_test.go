package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ozacod/cpx/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeriveTargetConfig(t *testing.T) {
	tests := []struct {
		name           string
		targetName     string
		expectedTarget config.CITarget
	}{
		{
			name:       "Linux AMD64",
			targetName: "linux-amd64",
			expectedTarget: config.CITarget{
				Name:       "linux-amd64",
				Dockerfile: "Dockerfile.linux-amd64",
				Image:      "cpx-linux-amd64",
				Triplet:    "x64-linux",
				Platform:   "linux/amd64",
			},
		},
		{
			name:       "Linux ARM64",
			targetName: "linux-arm64",
			expectedTarget: config.CITarget{
				Name:       "linux-arm64",
				Dockerfile: "Dockerfile.linux-arm64",
				Image:      "cpx-linux-arm64",
				Triplet:    "arm64-linux",
				Platform:   "linux/arm64",
			},
		},
		{
			name:       "Linux AMD64 MUSL",
			targetName: "linux-amd64-musl",
			expectedTarget: config.CITarget{
				Name:       "linux-amd64-musl",
				Dockerfile: "Dockerfile.linux-amd64-musl",
				Image:      "cpx-linux-amd64-musl",
				Triplet:    "x64-linux",
				Platform:   "linux/amd64",
			},
		},
		{
			name:       "Windows AMD64",
			targetName: "windows-amd64",
			expectedTarget: config.CITarget{
				Name:       "windows-amd64",
				Dockerfile: "Dockerfile.windows-amd64",
				Image:      "cpx-windows-amd64",
				Triplet:    "x64-mingw-static",
				Platform:   "linux/amd64",
			},
		},
		{
			name:       "macOS ARM64",
			targetName: "macos-arm64",
			expectedTarget: config.CITarget{
				Name:       "macos-arm64",
				Dockerfile: "Dockerfile.macos-arm64",
				Image:      "cpx-macos-arm64",
				Triplet:    "arm64-osx",
				Platform:   "linux/arm64",
			},
		},
		{
			name:       "macOS AMD64",
			targetName: "macos-amd64",
			expectedTarget: config.CITarget{
				Name:       "macos-amd64",
				Dockerfile: "Dockerfile.macos-amd64",
				Image:      "cpx-macos-amd64",
				Triplet:    "x64-osx",
				Platform:   "linux/amd64",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deriveTargetConfig(tt.targetName)
			assert.Equal(t, tt.expectedTarget.Name, result.Name)
			assert.Equal(t, tt.expectedTarget.Dockerfile, result.Dockerfile)
			assert.Equal(t, tt.expectedTarget.Image, result.Image)
			assert.Equal(t, tt.expectedTarget.Triplet, result.Triplet)
			assert.Equal(t, tt.expectedTarget.Platform, result.Platform)
		})
	}
}

func TestSaveCIConfig(t *testing.T) {
	tmpDir := t.TempDir()
	ciPath := filepath.Join(tmpDir, "cpx.ci")

	// Create test config
	ciConfig := &config.CIConfig{
		Targets: []config.CITarget{
			{
				Name:       "linux-amd64",
				Dockerfile: "Dockerfile.linux-amd64",
				Image:      "cpx-linux-amd64",
				Triplet:    "x64-linux",
				Platform:   "linux/amd64",
			},
		},
		Build: config.CIBuild{
			Type:         "Release",
			Optimization: "2",
			Jobs:         0,
		},
		Output: ".bin/ci",
	}

	// Save config
	err := config.SaveCI(ciConfig, ciPath)
	require.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(ciPath)
	require.NoError(t, err)

	// Load it back
	loadedConfig, err := config.LoadCI(ciPath)
	require.NoError(t, err)

	// Verify content
	assert.Len(t, loadedConfig.Targets, 1)
	assert.Equal(t, "linux-amd64", loadedConfig.Targets[0].Name)
	assert.Equal(t, "Release", loadedConfig.Build.Type)
	assert.Equal(t, ".bin/ci", loadedConfig.Output)
}
