# Sanitizer and Semgrep Examples

This directory contains example C++ and Go files that demonstrate violations that sanitizers and Semgrep can detect.

## Files

### Sanitizer Examples
- `sanitizer_examples.cpp` - Comprehensive examples for all sanitizers
- `asan_example.cpp` - AddressSanitizer examples (buffer overflow, use-after-free)
- `tsan_example.cpp` - ThreadSanitizer examples (data races)
- `msan_example.cpp` - MemorySanitizer examples (uninitialized memory)
- `ubsan_example.cpp` - UndefinedBehaviorSanitizer examples (undefined behavior)

### Semgrep Examples
- `semgrep_example.cpp` - C++ security and bug patterns that Semgrep detects
- `semgrep_example.go` - Go security and bug patterns that Semgrep detects

## Usage

### AddressSanitizer (ASan)
Detects memory errors: buffer overflows, use-after-free, double-free, memory leaks.

```bash
forge check --asan
./build/asan_example
```

**Example violations:**
- Stack buffer overflow: `arr[10]` when array size is 5
- Use after free: accessing deleted pointer
- Double free: calling `delete` twice on same pointer
- Memory leak: allocating memory without freeing

### ThreadSanitizer (TSan)
Detects data races in multi-threaded code.

```bash
forge check --tsan
./build/tsan_example
```

**Example violations:**
- Data race: multiple threads accessing shared variable without synchronization
- Race condition: concurrent modifications to shared data structures

### MemorySanitizer (MSan)
Detects uninitialized memory reads.

```bash
forge check --msan
./build/msan_example
```

**Example violations:**
- Reading uninitialized variables
- Reading uninitialized array elements
- Reading uninitialized struct members

### UndefinedBehaviorSanitizer (UBSan)
Detects undefined behavior in C++ code.

```bash
forge check --ubsan
./build/ubsan_example
```

**Example violations:**
- Signed integer overflow
- Division by zero
- Shift out of bounds
- Array index out of bounds
- Null pointer dereference
- Misaligned pointer access

### Semgrep
Detects security vulnerabilities and bugs in C++ and Go code.

```bash
forge semgrep
```

**Example detections:**
- Command injection vulnerabilities
- Buffer overflows
- Dangerous function usage (strcpy, gets, etc.)
- Hardcoded secrets and credentials
- SQL injection patterns
- Weak cryptography
- Race conditions
- Null pointer dereferences
- Memory leaks
- Use after free
- Path traversal
- Missing error checks

## Notes

- Sanitizers significantly slow down execution (2-10x slower)
- Use sanitizers during development and testing, not in production
- Only one sanitizer can be used at a time
- Some sanitizers require specific compiler flags and runtime libraries
- Semgrep is free and open-source (Community Edition)
- Semgrep can scan multiple languages simultaneously
- Semgrep works best with proper configuration files (`.semgrep.yml`)

