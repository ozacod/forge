// Semgrep Example - Security and Bug Detection
// Run: forge semgrep
// This file demonstrates code patterns that Semgrep can detect

#include <iostream>
#include <cstring>
#include <cstdio>
#include <cstdlib>

// Example 1: Command Injection (Semgrep rule: python.lang.security.audit.subprocess-shell-true)
void semgrep_command_injection() {
    char user_input[100];
    std::cin >> user_input;
    
    // Dangerous: user input in system command
    char cmd[200];
    sprintf(cmd, "ls %s", user_input);  // Semgrep will flag this!
    system(cmd);  // Command injection risk
}

// Example 2: Buffer Overflow (Semgrep rule: c.lang.security.buffer-overflow)
void semgrep_buffer_overflow() {
    char buffer[10];
    char input[100];
    std::cin >> input;
    
    // Dangerous: no bounds checking
    strcpy(buffer, input);  // Semgrep will flag this!
    std::cout << buffer << std::endl;
}

// Example 3: Use of Dangerous Function (Semgrep rule: c.lang.security.dangerous-function-use)
void semgrep_dangerous_function() {
    char buffer[100];
    char* input = "some string";
    
    // Dangerous: strcpy doesn't check bounds
    strcpy(buffer, input);  // Semgrep will suggest strncpy or safer alternatives
    
    // Dangerous: gets() is unsafe
    // gets(buffer);  // Semgrep will flag this!
}

// Example 4: Hardcoded Secrets (Semgrep rule: secrets)
void semgrep_hardcoded_secret() {
    // Semgrep will detect hardcoded secrets
    const char* api_key = "sk-1234567890abcdef";  // Semgrep will flag this!
    const char* password = "admin123";  // Semgrep will flag this!
    
    std::cout << "API Key: " << api_key << std::endl;
}

// Example 5: SQL Injection Pattern (if using database)
void semgrep_sql_injection() {
    char user_input[100];
    std::cin >> user_input;
    
    // Dangerous: user input in SQL query
    char query[200];
    sprintf(query, "SELECT * FROM users WHERE name = '%s'", user_input);  // Semgrep will flag this!
    // execute_query(query);
}

// Example 6: Weak Cryptography (Semgrep rule: c.lang.security.weak-crypto)
void semgrep_weak_crypto() {
    // Using weak hash functions
    // MD5 is cryptographically broken
    // Semgrep will suggest using SHA-256 or better
}

// Example 7: Race Condition (Semgrep rule: c.lang.security.race-condition)
int global_counter = 0;

void semgrep_race_condition() {
    // Accessing global variable without synchronization
    global_counter++;  // Semgrep may flag this in multi-threaded context
}

// Example 8: Null Pointer Dereference (Semgrep rule: c.lang.bugs.null-pointer-dereference)
void semgrep_null_pointer() {
    int* ptr = nullptr;
    
    // Dangerous: potential null dereference
    if (ptr != nullptr) {
        *ptr = 42;  // This is safe
    }
    
    // Dangerous: forgot to check
    // *ptr = 42;  // Semgrep will flag this!
}

// Example 9: Memory Leak (Semgrep rule: c.lang.bugs.memory-leak)
void semgrep_memory_leak() {
    int* ptr = new int[100];
    // Forgot to delete - Semgrep will flag this!
    // delete[] ptr;
}

// Example 10: Use After Free (Semgrep rule: c.lang.bugs.use-after-free)
void semgrep_use_after_free() {
    int* ptr = new int(42);
    delete ptr;
    // *ptr = 100;  // Semgrep will flag this!
}

int main() {
    std::cout << "Semgrep Example: Security and Bug Detection\n";
    std::cout << "===========================================\n\n";
    std::cout << "This file contains code patterns that Semgrep can detect.\n";
    std::cout << "Run 'forge semgrep' to scan for issues.\n\n";
    
    std::cout << "Common Semgrep detections:\n";
    std::cout << "- Command injection vulnerabilities\n";
    std::cout << "- Buffer overflows\n";
    std::cout << "- Dangerous function usage (strcpy, gets, etc.)\n";
    std::cout << "- Hardcoded secrets and credentials\n";
    std::cout << "- SQL injection patterns\n";
    std::cout << "- Weak cryptography\n";
    std::cout << "- Race conditions\n";
    std::cout << "- Null pointer dereferences\n";
    std::cout << "- Memory leaks\n";
    std::cout << "- Use after free\n";
    
    return 0;
}

