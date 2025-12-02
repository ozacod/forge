package quality

import (
	"fmt"
	"os"
	"os/exec"
)

// RunFlawfinder runs Flawfinder security analysis for C/C++
func RunFlawfinder(minLevel int, csv, html bool, output string, dataflow, quiet, singleline bool, context int, targets []string) error {
	// Check if flawfinder is available
	if _, err := exec.LookPath("flawfinder"); err != nil {
		return fmt.Errorf("flawfinder not found. Please install it first:\n  pip install flawfinder\n  or\n  brew install flawfinder\n  or\n  apt-get install flawfinder (Debian/Ubuntu)")
	}

	// Validate output file for HTML/CSV
	if (html || csv) && output == "" {
		return fmt.Errorf("--output file is required when using --html or --csv flags")
	}

	fmt.Printf("%s Running Flawfinder analysis...%s\n", Cyan, Reset)

	// Filter targets to only include git-tracked files (respect .gitignore)
	filteredTargets, err := FilterGitTrackedFiles(targets)
	if err != nil {
		// If git is not available or not in a git repo, use original targets
		fmt.Printf("%s Warning: Not in a git repository or git not available. Scanning all files.%s\n", Yellow, Reset)
		filteredTargets = targets
	} else if len(filteredTargets) == 0 {
		return fmt.Errorf("no git-tracked C/C++ files found to scan")
	}

	// Build flawfinder command
	flawfinderArgs := []string{}

	// Add min level
	if minLevel >= 0 && minLevel <= 5 {
		flawfinderArgs = append(flawfinderArgs, "-m", fmt.Sprintf("%d", minLevel))
	}

	// Output format
	if csv {
		flawfinderArgs = append(flawfinderArgs, "-C")
	} else if html {
		flawfinderArgs = append(flawfinderArgs, "-H")
	}

	// Dataflow analysis
	if dataflow {
		flawfinderArgs = append(flawfinderArgs, "-D")
	}

	// Quiet mode
	if quiet {
		flawfinderArgs = append(flawfinderArgs, "--quiet")
	}

	// Single line output
	if singleline {
		flawfinderArgs = append(flawfinderArgs, "--singleline")
	}

	// Context lines
	if context > 0 {
		flawfinderArgs = append(flawfinderArgs, "-c", fmt.Sprintf("%d", context))
	}

	// Add filtered target files
	flawfinderArgs = append(flawfinderArgs, filteredTargets...)

	cmd := exec.Command("flawfinder", flawfinderArgs...)

	// Handle output file for HTML/CSV
	if output != "" {
		file, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer file.Close()
		cmd.Stdout = file
		fmt.Printf("%s Writing output to: %s%s\n", Cyan, output, Reset)
	} else {
		cmd.Stdout = os.Stdout
	}

	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Flawfinder returns non-zero on findings, which is normal
		if output != "" {
			fmt.Printf("%s  Flawfinder found potential issues (saved to %s)%s\n", Yellow, output, Reset)
		} else {
			fmt.Printf("%s  Flawfinder found potential issues%s\n", Yellow, Reset)
		}
		return nil
	}

	if output != "" {
		fmt.Printf("%s Analysis complete! Report saved to: %s%s\n", Green, output, Reset)
	} else {
		fmt.Printf("%s No issues found!%s\n", Green, Reset)
	}
	return nil
}
