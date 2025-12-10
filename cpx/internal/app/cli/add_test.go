package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunVcpkgAdd(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectsError bool
		mockError    error
	}{
		{
			name:         "Successful vcpkg add",
			args:         []string{"zlib"},
			expectsError: false,
			mockError:    nil,
		},
		{
			name:         "Failed vcpkg add",
			args:         []string{"nonexistent"},
			expectsError: true,
			mockError:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup the global variable
			addRunVcpkgCommandFunc = func(args []string) error {
				return tt.mockError
			}

			// Capture output
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			runErr := runVcpkgAdd(tt.args)

			// Restore stdout
			err := w.Close()
			if err != nil {
				t.Fatalf("Failed to close pipe: %v", err)
			}
			os.Stdout = old

			// Read captured output
			var buf2 bytes.Buffer
			if _, err := buf2.ReadFrom(r); err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}

			if tt.expectsError {
				assert.Error(t, runErr)
			} else {
				assert.NoError(t, runErr)
			}
		})
	}
}

func TestRunBazelAdd(t *testing.T) {
	tests := []struct {
		name                 string
		args                 []string
		bcrPath              string
		mockVersions         map[string]string
		mockErrors           map[string]error
		expectsError         bool
		expectedDependencies []string
	}{
		{
			name:    "Successful bazel add",
			args:    []string{"com_google_googletest"},
			bcrPath: "/test/bcr",
			mockVersions: map[string]string{
				"com_google_googletest": "1.14.0",
			},
			mockErrors:           nil,
			expectsError:         false,
			expectedDependencies: []string{"com_google_googletest"},
		},
		{
			name:    "Failed bazel add - module not found",
			args:    []string{"nonexistent_module"},
			bcrPath: "/test/bcr",
			mockVersions: map[string]string{
				"nonexistent_module": "",
			},
			mockErrors: map[string]error{
				"nonexistent_module": assert.AnError,
			},
			expectsError: false, // Should not return error, just skip the module
		},
		{
			name:         "Failed bazel add - no bcr path",
			args:         []string{"com_google_googletest"},
			bcrPath:      "",
			expectsError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore original functions
			oldGetBcrPathFunc := addGetBcrPathFunc
			oldGetLatestVersionFunc := bazelGetLatestVersionFunc
			oldAddDependencyFunc := bazelAddDependencyFunc
			defer func() {
				addGetBcrPathFunc = oldGetBcrPathFunc
				bazelGetLatestVersionFunc = oldGetLatestVersionFunc
				bazelAddDependencyFunc = oldAddDependencyFunc
			}()

			// Setup the mock for getBcrPath
			addGetBcrPathFunc = func() string {
				return tt.bcrPath
			}

			// Setup the mock for GetLatestVersion
			bazelGetLatestVersionFunc = func(bcrPath, moduleName string) (string, error) {
				if tt.mockErrors != nil && tt.mockErrors[moduleName] != nil {
					return "", tt.mockErrors[moduleName]
				}
				if version, ok := tt.mockVersions[moduleName]; ok {
					return version, nil
				}
				return "", assert.AnError
			}

			// Setup the mock for AddDependency (uses the real implementation)
			bazelAddDependencyFunc = func(modulePath, depName, version string) error {
				// For testing, we'll use the real AddDependency since it just writes to a file
				return nil // Just succeed for test purposes
			}

			// Create a temporary directory for MODULE.bazel
			tmpDir := t.TempDir()
			oldWd, err := os.Getwd()
			require.NoError(t, err)
			defer func(dir string) {
				err := os.Chdir(dir)
				if err != nil {
					t.Fatalf("Failed to restore working directory: %v", err)
				}
			}(oldWd)
			require.NoError(t, os.Chdir(tmpDir))

			// Create MODULE.bazel file
			moduleContent := `module(name = "test")

bazel_dep(name = "rules_cc", version = "0.0.1")
`
			require.NoError(t, os.WriteFile("MODULE.bazel", []byte(moduleContent), 0644))

			// Capture output
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			runErr := runBazelAdd(tt.args)

			// Restore stdout
			if err := w.Close(); err != nil {
				t.Fatalf("Failed to close pipe: %v", err)
			}
			os.Stdout = old

			// Read captured output
			var buf2 bytes.Buffer
			if _, err := buf2.ReadFrom(r); err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}
			output := buf2.String()

			if tt.expectsError {
				assert.Error(t, runErr)
			} else {
				if tt.bcrPath != "" {
					assert.NoError(t, runErr)
				}

				// Verify output contains expected messages for successful cases
				if tt.bcrPath != "" && len(tt.expectedDependencies) > 0 {
					for _, dep := range tt.expectedDependencies {
						if tt.mockVersions[dep] != "" {
							assert.Contains(t, output, dep)
							assert.Contains(t, output, tt.mockVersions[dep])
						}
					}
				}
			}
		})
	}
}

func TestPrintVcpkgUsageInfo(t *testing.T) {
	tests := []struct {
		name     string
		pkgName  string
		mockHTTP bool
	}{
		{
			name:     "Print usage info for existing package",
			pkgName:  "zlib",
			mockHTTP: true,
		},
		{
			name:     "Print usage info for nonexistent package",
			pkgName:  "nonexistent",
			mockHTTP: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// This is a simple test - the actual HTTP call would need more complex mocking
			// For now, we just test that the function doesn't crash
			printVcpkgUsageInfo(tt.pkgName)

			// Restore stdout
			err := w.Close()
			if err != nil {
				t.Fatalf("Failed to close pipe: %v", err)
			}
			os.Stdout = old

			// Read captured output to avoid "declared and not used" error
			var buf4 bytes.Buffer
			_, err = buf4.ReadFrom(r)
			if err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}

			// We can't easily test the HTTP response without more complex mocking
			// This test just ensures the function runs without crashing
		})
	}
}

func TestPrintBazelUsageInfo(t *testing.T) {
	tests := []struct {
		name    string
		pkgName string
	}{
		{
			name:    "Print bazel usage info",
			pkgName: "com_google_googletest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			printBazelUsageInfo(tt.pkgName)

			// Restore stdout
			err := w.Close()
			if err != nil {
				t.Fatalf("Close failed: %v", err)
			}
			os.Stdout = old

			// Read captured output
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}
			output := buf.String()

			assert.Contains(t, output, "USAGE INFO FOR")
			assert.Contains(t, output, tt.pkgName)
			assert.Contains(t, output, "BUILD.bazel")
			assert.Contains(t, output, "@"+tt.pkgName)
		})
	}
}

func TestPrintMesonUsageInfo(t *testing.T) {
	tests := []struct {
		name    string
		pkgName string
	}{
		{
			name:    "Print meson usage info",
			pkgName: "zlib",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			printMesonUsageInfo(tt.pkgName)

			// Restore stdout
			err := w.Close()
			if err != nil {
				t.Fatalf("Close failed: %v", err)
			}
			os.Stdout = old

			// Read captured output
			var buf bytes.Buffer
			_, err = buf.ReadFrom(r)
			if err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}
			output := buf.String()

			assert.Contains(t, output, "USAGE INFO FOR")
			assert.Contains(t, output, tt.pkgName)
			assert.Contains(t, output, "meson.build")
			assert.Contains(t, output, "dependency")
		})
	}
}
