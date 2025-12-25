# Sanitizer Examples

This directory contains example C++ files that demonstrate violations that sanitizers can detect.

## Files

### Sanitizer Examples
- `sanitizer_examples.cpp` - Comprehensive examples for all sanitizers
- `asan_example.cpp` - AddressSanitizer examples (buffer overflow, use-after-free)
- `tsan_example.cpp` - ThreadSanitizer examples (data races)
- `msan_example.cpp` - MemorySanitizer examples (uninitialized memory)
- `ubsan_example.cpp` - UndefinedBehaviorSanitizer examples (undefined behavior)

## Usage

### AddressSanitizer (ASan)
Detects memory errors: buffer overflows, use-after-free, double-free, memory leaks.

```bash
cpx check --asan
./build/asan_example
```

**Example violations:**
- Stack buffer overflow: `arr[10]` when array size is 5
- Use after free: accessing deleted pointer
- Double free: calling `delete` twice on same pointer
- Memory leak: allocating memory without freeing

### ThreadSanitizer (TSan)
Detects data races in multithreaded code.

```bash
cpx check --tsan
./build/tsan_example
```

**Example violations:**
- Data race: multiple threads accessing shared variable without synchronization
- Race condition: concurrent modifications to shared data structures

### MemorySanitizer (MSan)
Detects uninitialized memory reads.

```bash
cpx check --msan
./build/msan_example
```

**Example violations:**
- Reading uninitialized variables
- Reading uninitialized array elements
- Reading uninitialized struct members

### UndefinedBehaviorSanitizer (UBSan)
Detects undefined behavior in C++ code.

```bash
cpx check --ubsan
./build/ubsan_example
```

**Example violations:**
- Signed integer overflow
- Division by zero
- Shift out of bounds
- Array index out of bounds
- Null pointer dereference
- Misaligned pointer access

## Notes

- Sanitizers significantly slow down execution (2-10x slower)
- Use sanitizers during development and testing, not in production
- Only one sanitizer can be used at a time
- Some sanitizers require specific compiler flags and runtime libraries

