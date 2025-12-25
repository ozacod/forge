// Package meson provides Meson build system Docker integration.
package meson

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	build "github.com/ozacod/cpx/internal/pkg/build/interfaces"
	"github.com/ozacod/cpx/internal/pkg/utils/colors"
)

// RunDockerBuild implements the DockerBuilder interface for Meson builds.
func (b *Builder) RunDockerBuild(ctx context.Context, opts build.DockerBuildOptions) error {
	absProjectRoot, err := filepath.Abs(opts.ProjectRoot)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for project root: %w", err)
	}

	absOutputDir, err := filepath.Abs(filepath.Join(opts.ProjectRoot, opts.OutputDir))
	if err != nil {
		return fmt.Errorf("failed to get absolute path for output directory: %w", err)
	}

	// Create persistent build directory
	hostBuildDir := filepath.Join(opts.ProjectRoot, ".cache", "ci", opts.TargetName)
	if err := os.MkdirAll(hostBuildDir, 0755); err != nil {
		return fmt.Errorf("failed to create build directory: %w", err)
	}
	absBuildDir, err := filepath.Abs(hostBuildDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for build directory: %w", err)
	}

	// Determine build type
	buildType := "release"
	if opts.BuildType == "Debug" || opts.BuildType == "debug" {
		buildType = "debug"
	}

	// Create subprojects directory
	hostSubprojectsDir := filepath.Join(opts.ProjectRoot, "subprojects")
	if err := os.MkdirAll(hostSubprojectsDir, 0755); err != nil {
		return fmt.Errorf("failed to create subprojects directory: %w", err)
	}
	absSubprojectsDir, err := filepath.Abs(hostSubprojectsDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for subprojects directory: %w", err)
	}

	// Environment exports
	var envExports string
	if len(opts.Env) > 0 {
		envExports = "# User-defined environment variables\n"
		for k, v := range opts.Env {
			envExports += fmt.Sprintf("export %s=\"%s\"\n", k, v)
		}
	}

	// Build Meson arguments
	setupArgs := []string{"--buildtype=" + buildType}
	setupArgs = append(setupArgs, opts.MesonArgs...)

	testSection := ""
	if opts.RunTests {
		testSection = `
echo "  Running tests..."
meson test -C /tmp/builddir -v
`
	}

	benchSection := ""
	if opts.RunBenchmarks {
		benchSection = `
echo "  Running benchmarks..."
meson test -C /tmp/builddir --benchmark -v
`
	}

	buildScript := fmt.Sprintf(`#!/bin/bash
set -e
%smkdir -p /tmp/builddir
echo "  Configuring Meson..."
if [ ! -f /tmp/builddir/build.ninja ]; then
    meson setup /tmp/builddir %s
else
    echo "  Build directory already configured, skipping setup."
fi
echo "  Building..."
meson compile -C /tmp/builddir
echo "  Copying artifacts..."
mkdir -p /workspace/out/%s
if [ -d "/tmp/builddir/src" ]; then
    find /tmp/builddir/src -maxdepth 1 -type f -perm +111 ! -name "*.so" ! -name "*.dylib" ! -name "*.a" ! -name "*.p" ! -name "*_test" -exec cp {} /workspace/out/%s/ \; 2>/dev/null || true
fi
find /tmp/builddir -maxdepth 1 -type f -perm +111 ! -name "*.so" ! -name "*.dylib" ! -name "*.a" ! -name "*.p" ! -name "build.ninja" ! -name "*.json" -exec cp {} /workspace/out/%s/ \; 2>/dev/null || true
find /tmp/builddir -maxdepth 2 -type f \( -name "*.a" -o -name "*.so" -o -name "*.dylib" \) -exec cp {} /workspace/out/%s/ \; 2>/dev/null || true
ls -la /workspace/out/%s/ 2>/dev/null || echo "  (no artifacts found)"
echo "  Build complete!"
%s%s
`, envExports, strings.Join(setupArgs, " "), opts.TargetName, opts.TargetName, opts.TargetName, opts.TargetName, opts.TargetName, testSection, benchSection)

	fmt.Printf("  %s Running Meson build in Docker container...%s\n", colors.Cyan, colors.Reset)

	dockerArgs := []string{"run", "--rm"}
	if opts.Platform != "" {
		dockerArgs = append(dockerArgs, "--platform", opts.Platform)
	}

	dockerArgs = append(dockerArgs,
		"-v", absProjectRoot+":/workspace:ro",
		"-v", absBuildDir+":/tmp/builddir",
		"-v", absSubprojectsDir+":/workspace/subprojects",
		"-v", absOutputDir+":/workspace/out",
		"-w", "/workspace",
		opts.ImageName,
		"bash", "-c", buildScript)

	cmd := exec.Command("docker", dockerArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker meson build failed: %w", err)
	}

	return nil
}

// Compile-time check that Builder implements DockerBuilder
var _ build.DockerBuilder = (*Builder)(nil)
