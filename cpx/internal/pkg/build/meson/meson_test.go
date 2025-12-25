package meson

import (
	"context"
	"os"
	"os/exec"
	"testing"

	build "github.com/ozacod/cpx/internal/pkg/build/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHelperProcess isn't a real test. It's used as a helper process
// for mocking exec.Command.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}

func TestBuild(t *testing.T) {
	// Mock execCommand
	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()

	var capturedArgs [][]string

	execCommand = func(name string, arg ...string) *exec.Cmd {
		args := append([]string{name}, arg...)
		capturedArgs = append(capturedArgs, args)

		cs := []string{"-test.run=TestHelperProcess", "--", name}
		cs = append(cs, arg...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
		return cmd
	}

	// Use temp dir
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(oldWd) }()
	require.NoError(t, os.Chdir(tmpDir))

	builder := New()

	// Test Debug Build
	capturedArgs = nil
	err = builder.Build(context.Background(), build.BuildOptions{
		Release: false,
	})
	assert.NoError(t, err)

	require.Len(t, capturedArgs, 3) // setup, compile, copy
	// meson setup
	assert.Equal(t, "meson", capturedArgs[0][0])
	assert.Equal(t, "setup", capturedArgs[0][1])
	assert.Contains(t, capturedArgs[0], "--buildtype=debug")
	// meson compile
	assert.Equal(t, "meson", capturedArgs[1][0])
	assert.Equal(t, "compile", capturedArgs[1][1])

	// Test Release Build
	// Note: builddir already exists, so setup will be SKIPPED unless we clean or reuse code that detects changes.
	// But in our mock, builddir dir is created by the actual os.Stat checks calling real FS, or does it?
	// In the real code: if _, err := os.Stat(buildDir); os.IsNotExist(err)
	// Since the previous run (debug build) would have executed "setup", but specifically checking if "builddir" exists...
	// Wait, our mock execution DOES NOT create the directory `builddir`.
	// So `os.Stat(buildDir)` will return NotExist every time in this test environment unless we create it.

	// Let's verify what happens.
	// If `builddir` is not created, `Build` sees it missing and calls `setup` again.
	// So for Release build, we expect `setup` to be called again with release flags.

	capturedArgs = nil
	err = builder.Build(context.Background(), build.BuildOptions{
		Release: true,
		Clean:   true, // Clean will remove builddir anyway
	})
	assert.NoError(t, err)

	// With clean=true:
	// 1. builder.Clean is called (execs meson clean usually? No, Clean logic calls removeDir)
	// builder.Clean calls internal removeDir, which uses os.RemoveAll.
	// So builddir is removed.

	// So calls expected: setup, compile, copy
	// Note: capturedArgs might contain args from Clean if Clean calls any specific command.
	// builder.Clean calls `removeDir`... does it call `meson clean`?
	// Checking code... Clean calls `meson clean` if it can? No, `os.RemoveAll("builddir")`.

	// Wait, checking Step 1396 content for `Clean` method.
	// It calls `removeDir("builddir")`. No `exec` command for clean in Meson builder, unlike Bazel.
	// Ah, step 1410 shows the new code for Build.
	// Clean implementation in Meson builder (Step 1396 lines 203+):
	// It just removes directories.

	// So capturedArgs should just have Set up, Compile, Copy

	require.Len(t, capturedArgs, 3)
	assert.Equal(t, "setup", capturedArgs[0][1])
	assert.Contains(t, capturedArgs[0], "--buildtype=release")
}

func TestRun(t *testing.T) {
	// Mock execCommand
	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()

	var capturedArgs [][]string

	execCommand = func(name string, arg ...string) *exec.Cmd {
		args := append([]string{name}, arg...)
		capturedArgs = append(capturedArgs, args)

		cs := []string{"-test.run=TestHelperProcess", "--", name}
		cs = append(cs, arg...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
		return cmd
	}

	// Use temp dir
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(oldWd) }()
	require.NoError(t, os.Chdir(tmpDir))

	// Create meson.build
	require.NoError(t, os.WriteFile("meson.build", []byte("project('test', 'cpp')"), 0644))

	// Create executable in builddir/src
	// (Note: Build logic usually creates builddir, but we are mocking.
	// Run logic checks for executable existence on disk.)
	srcDir := "builddir/src"
	err = os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(srcDir+"/myapp", []byte("#!/bin/sh\necho hello"), 0755)
	require.NoError(t, err)

	builder := New()

	err = builder.Run(context.Background(), build.RunOptions{
		Release: false,
		Target:  "myapp",
		Verbose: true,
	})
	assert.NoError(t, err)

	// Calls expected:
	// 1. Build calls Setup (since builddir exists but maybe not fully configured? No, we just created directories)
	// Build logic checks `os.Stat(buildDir)`. We created `builddir/src`, so `builddir` exists.
	// So Build calls `meson configure`.
	// 2. Build calls `meson compile`
	// 3. Build calls `bash` copy script
	// 4. Run calls `./builddir/src/myapp`

	require.GreaterOrEqual(t, len(capturedArgs), 4)

	// Check for run execution
	lastCmd := capturedArgs[len(capturedArgs)-1]
	assert.Equal(t, "builddir/src/myapp", lastCmd[0])
}

func TestTest(t *testing.T) {
	// Mock execCommand
	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()

	var capturedArgs [][]string

	execCommand = func(name string, arg ...string) *exec.Cmd {
		args := append([]string{name}, arg...)
		capturedArgs = append(capturedArgs, args)

		cs := []string{"-test.run=TestHelperProcess", "--", name}
		cs = append(cs, arg...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
		return cmd
	}

	// Create temp dir
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(oldWd) }()
	require.NoError(t, os.Chdir(tmpDir))

	// Create meson.build and builddir so Test logic thinks it's built
	require.NoError(t, os.WriteFile("meson.build", []byte("project('test', 'cpp')"), 0644))
	require.NoError(t, os.MkdirAll("builddir", 0755))

	builder := New()

	err = builder.Test(context.Background(), build.TestOptions{
		Verbose: true,
		Filter:  "mytest",
	})
	assert.NoError(t, err)

	require.Len(t, capturedArgs, 1) // meson test
	assert.Equal(t, "meson", capturedArgs[0][0])
	assert.Equal(t, "test", capturedArgs[0][1])
	assert.Contains(t, capturedArgs[0], "mytest")
}

func TestBench(t *testing.T) {
	// Mock execCommand
	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()

	var capturedArgs [][]string

	execCommand = func(name string, arg ...string) *exec.Cmd {
		args := append([]string{name}, arg...)
		capturedArgs = append(capturedArgs, args)

		cs := []string{"-test.run=TestHelperProcess", "--", name}
		cs = append(cs, arg...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
		return cmd
	}

	// Create temp dir
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(oldWd) }()
	require.NoError(t, os.Chdir(tmpDir))

	// Create meson.build and builddir so Bench logic thinks it's built
	require.NoError(t, os.WriteFile("meson.build", []byte("project('test', 'cpp')"), 0644))
	require.NoError(t, os.MkdirAll("builddir", 0755))

	// Mock existence of benchmark executable
	benchDir := "builddir/bench/myapp_bench"
	require.NoError(t, os.MkdirAll(benchDir, 0755))

	builder := New()

	// Test Bench with specific target
	err = builder.Bench(context.Background(), build.BenchOptions{
		Verbose: true,
		Target:  "myapp_bench",
	})
	assert.NoError(t, err)

	// Since builddir exists, Build() calls (re)configure, then compile, then copy...
	// Wait, builder.Bench checks if builddir exists. If it does, it calls b.Build only if IT DOES NOT exist?
	// Checking `meson.go`:
	// 	if _, err := os.Stat("builddir"); os.IsNotExist(err) { if err := b.Build(...) ... }
	// So if builddir exists, Build is NOT called.

	// Then it looks for bench executable.
	// benchPath = "builddir/bench/myapp_bench" (we created it as dir, but os.Stat counts dir too? ExecCommand might fail if it's a dir, but we mock it)

	require.Len(t, capturedArgs, 1)
	assert.Equal(t, "builddir/bench/myapp_bench", capturedArgs[0][0])
}

func TestClean(t *testing.T) {
	// Use temp dir
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(oldWd) }()
	require.NoError(t, os.Chdir(tmpDir))

	// Create Meson project files
	require.NoError(t, os.WriteFile("meson.build", []byte("project('test', 'cpp')"), 0644))

	tests := []struct {
		name          string
		all           bool
		expectRemoved []string
	}{
		{
			name:          "Clean without all flag",
			all:           false,
			expectRemoved: []string{"builddir", "build"},
		},
		{
			name:          "Clean with all flag",
			all:           true,
			expectRemoved: []string{"builddir", "build", "subprojects/packagecache", "build-release"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Recreate dirs for each test
			_ = os.MkdirAll("builddir", 0755)
			_ = os.MkdirAll("build", 0755)
			_ = os.MkdirAll("subprojects/packagecache", 0755)
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
