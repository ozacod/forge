package vcpkg

import (
	"context"
	"os"
	"testing"

	build "github.com/ozacod/cpx/internal/pkg/build/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClean(t *testing.T) {
	// Use temp dir
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(oldWd) }()
	require.NoError(t, os.Chdir(tmpDir))

	// Create CMake project files
	require.NoError(t, os.WriteFile("CMakeLists.txt", []byte("cmake_minimum_required(VERSION 3.16)"), 0644))

	tests := []struct {
		name          string
		all           bool
		expectRemoved []string
	}{
		{
			name:          "Clean without all flag",
			all:           false,
			expectRemoved: []string{".bin/native"},
		},
		{
			name:          "Clean with all flag",
			all:           true,
			expectRemoved: []string{".bin/native", ".bin/ci", "out", "cmake-build-debug", "build-release"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Recreate dirs for each test
			_ = os.MkdirAll(".bin/native", 0755)
			_ = os.MkdirAll(".bin/ci", 0755)
			_ = os.MkdirAll("out", 0755)
			_ = os.MkdirAll("cmake-build-debug", 0755)
			_ = os.MkdirAll("build-release", 0755)

			builder := New()
			err := builder.Clean(context.Background(), build.CleanOptions{All: tt.all})
			assert.NoError(t, err)

			// Verify expected directories were removed
			for _, dir := range tt.expectRemoved {
				_, err = os.Stat(dir)
				assert.True(t, os.IsNotExist(err), "%s should be removed", dir)
			}
		})
	}
}
