// Semgrep Example - Go Security and Bug Detection
// Run: forge semgrep
// This file demonstrates Go code patterns that Semgrep can detect

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Example 1: Command Injection (Semgrep rule: go.lang.security.command-injection)
func semgrepCommandInjection(userInput string) {
	// Dangerous: user input in command
	cmd := exec.Command("sh", "-c", "ls "+userInput) // Semgrep will flag this!
	cmd.Run()
}

// Example 2: SQL Injection (Semgrep rule: go.lang.security.sql-injection)
func semgrepSQLInjection(userInput string) {
	query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", userInput) // Semgrep will flag this!
	// db.Query(query)
	fmt.Println(query)
}

// Example 3: Path Traversal (Semgrep rule: go.lang.security.path-traversal)
func semgrepPathTraversal(userInput string) {
	file, err := os.Open("/data/" + userInput) // Semgrep will flag this!
	if err != nil {
		return
	}
	defer file.Close()
}

// Example 4: Hardcoded Secrets (Semgrep rule: secrets)
func semgrepHardcodedSecret() {
	apiKey := "sk-1234567890abcdef" // Semgrep will flag this!
	password := "admin123"          // Semgrep will flag this!
	fmt.Println("API Key:", apiKey)
	fmt.Println("Password:", password)
}

// Example 5: Weak Cryptography (Semgrep rule: go.lang.security.weak-crypto)
func semgrepWeakCrypto() {
	// Using MD5 which is cryptographically broken
	// Semgrep will suggest using crypto/sha256 or better
}

// Example 6: Insecure Random (Semgrep rule: go.lang.security.insecure-random)
func semgrepInsecureRandom() {
	// Using math/rand instead of crypto/rand
	// Semgrep will flag this for security-sensitive operations
}

// Example 7: Race Condition (Semgrep rule: go.lang.bugs.race-condition)
var globalCounter int

func semgrepRaceCondition() {
	// Accessing global variable without synchronization
	globalCounter++ // Semgrep will flag this!
}

// Example 8: Error Handling (Semgrep rule: go.lang.bugs.missing-error-check)
func semgrepMissingErrorCheck() {
	file, _ := os.Open("file.txt") // Semgrep will flag ignored error!
	defer file.Close()
}

// Example 9: Use of Dangerous syscall (Semgrep rule: go.lang.security.dangerous-syscall)
func semgrepDangerousSyscall() {
	syscall.Exec("/bin/sh", []string{"sh"}, nil) // Semgrep will flag this!
}

// Example 10: Insecure TLS (Semgrep rule: go.lang.security.insecure-tls)
func semgrepInsecureTLS() {
	// Using InsecureSkipVerify
	// Semgrep will flag this!
}

func main() {
	fmt.Println("Semgrep Example: Go Security and Bug Detection")
	fmt.Println("=============================================")
	fmt.Println()
	fmt.Println("This file contains Go code patterns that Semgrep can detect.")
	fmt.Println("Run 'forge semgrep' to scan for issues.")
	fmt.Println()
	fmt.Println("Common Semgrep detections in Go:")
	fmt.Println("- Command injection vulnerabilities")
	fmt.Println("- SQL injection patterns")
	fmt.Println("- Path traversal vulnerabilities")
	fmt.Println("- Hardcoded secrets and credentials")
	fmt.Println("- Weak cryptography")
	fmt.Println("- Insecure random number generation")
	fmt.Println("- Race conditions")
	fmt.Println("- Missing error checks")
	fmt.Println("- Dangerous syscalls")
	fmt.Println("- Insecure TLS configuration")
}
