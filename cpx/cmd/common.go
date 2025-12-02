package cmd

import (
	"fmt"
	"os"
)

// Colors for terminal output
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Bold    = "\033[1m"
)

// Version is the cpx version
const Version = "1.0.0"

// DefaultServer is the default server URL
const DefaultServer = "https://cpxcpp.vercel.app"

// DefaultCfgFile is the default config file name
const DefaultCfgFile = "cpx.yaml"

// LockFile is the lock file name
const LockFile = "cpx.lock"

// ExitWithError prints an error message and exits with status 1
func ExitWithError(err error) {
	fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
	os.Exit(1)
}
