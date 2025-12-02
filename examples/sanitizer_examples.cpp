// Sanitizer Examples - Demonstrates violations that each sanitizer can detect
// This file is for educational purposes only

#include <iostream>
#include <thread>
#include <vector>
#include <cstring>
#include <cstdlib>

// ============================================================================
// AddressSanitizer (ASan) Examples
// ============================================================================

// Example 1: Buffer overflow (stack)
void asan_buffer_overflow() {
    int arr[10];
    arr[15] = 42;  // Out of bounds write - ASan will catch this
    std::cout << "Buffer overflow: " << arr[15] << std::endl;
}

// Example 2: Use after free
void asan_use_after_free() {
    int* ptr = new int(42);
    delete ptr;
    *ptr = 100;  // Use after free - ASan will catch this
    std::cout << "Use after free: " << *ptr << std::endl;
}

// Example 3: Double free
void asan_double_free() {
    int* ptr = new int(42);
    delete ptr;
    delete ptr;  // Double free - ASan will catch this
}

// Example 4: Memory leak (if compiled with leak detection)
void asan_memory_leak() {
    int* ptr = new int[1000];
    // Forgot to delete - ASan with leak detection will report this
    // delete[] ptr;
}

// ============================================================================
// ThreadSanitizer (TSan) Examples
// ============================================================================

int shared_counter = 0;  // Shared variable without synchronization

// Example 5: Data race
void tsan_data_race_increment() {
    for (int i = 0; i < 100000; ++i) {
        shared_counter++;  // Data race - TSan will catch this
    }
}

void tsan_data_race_example() {
    std::thread t1(tsan_data_race_increment);
    std::thread t2(tsan_data_race_increment);
    t1.join();
    t2.join();
    std::cout << "Counter value: " << shared_counter << std::endl;
}

// Example 6: Race condition on vector
std::vector<int> shared_vec;

void tsan_vector_race() {
    for (int i = 0; i < 1000; ++i) {
        shared_vec.push_back(i);  // Data race - TSan will catch this
    }
}

void tsan_vector_race_example() {
    std::thread t1(tsan_vector_race);
    std::thread t2(tsan_vector_race);
    t1.join();
    t2.join();
}

// ============================================================================
// MemorySanitizer (MSan) Examples
// ============================================================================

// Example 7: Uninitialized memory read
void msan_uninitialized_read() {
    int x;  // Uninitialized
    if (x > 0) {  // MSan will catch this
        std::cout << "Uninitialized read: " << x << std::endl;
    }
}

// Example 8: Uninitialized array
void msan_uninitialized_array() {
    int arr[10];
    // arr[0] is never initialized
    std::cout << "Uninitialized array: " << arr[0] << std::endl;  // MSan will catch this
}

// Example 9: Uninitialized struct member
struct MyStruct {
    int value;
    char name[10];
};

void msan_uninitialized_struct() {
    MyStruct s;
    // s.value is never initialized
    std::cout << "Uninitialized struct: " << s.value << std::endl;  // MSan will catch this
}

// ============================================================================
// UndefinedBehaviorSanitizer (UBSan) Examples
// ============================================================================

// Example 10: Signed integer overflow
void ubsan_signed_overflow() {
    int x = INT_MAX;
    x++;  // Signed overflow - UBSan will catch this
    std::cout << "Signed overflow: " << x << std::endl;
}

// Example 11: Null pointer dereference
void ubsan_null_pointer() {
    int* ptr = nullptr;
    *ptr = 42;  // Null pointer dereference - UBSan will catch this
}

// Example 12: Division by zero
void ubsan_division_by_zero() {
    int x = 10;
    int y = 0;
    int result = x / y;  // Division by zero - UBSan will catch this
    std::cout << "Division: " << result << std::endl;
}

// Example 13: Shift out of bounds
void ubsan_shift_out_of_bounds() {
    int x = 1;
    int shift = 100;  // Too large for int
    int result = x << shift;  // Shift out of bounds - UBSan will catch this
    std::cout << "Shift: " << result << std::endl;
}

// Example 14: Array index out of bounds (undefined behavior)
void ubsan_array_bounds() {
    int arr[5] = {1, 2, 3, 4, 5};
    int index = 10;
    int value = arr[index];  // Out of bounds - UBSan will catch this
    std::cout << "Array access: " << value << std::endl;
}

// Example 15: Misaligned pointer access
void ubsan_misaligned_pointer() {
    char data[10] = "123456789";
    int* ptr = reinterpret_cast<int*>(data + 1);  // Misaligned
    *ptr = 42;  // Misaligned access - UBSan will catch this (on some platforms)
}

// Example 16: Invalid cast
void ubsan_invalid_cast() {
    float f = 3.14f;
    int* ptr = reinterpret_cast<int*>(&f);
    // Accessing float as int - potential undefined behavior
    std::cout << "Invalid cast: " << *ptr << std::endl;
}

// ============================================================================
// Main function - Uncomment the examples you want to test
// ============================================================================

int main() {
    std::cout << "Sanitizer Examples\n";
    std::cout << "==================\n\n";
    
    // Uncomment the example you want to test:
    
    // AddressSanitizer examples:
    // asan_buffer_overflow();
    // asan_use_after_free();
    // asan_double_free();
    // asan_memory_leak();
    
    // ThreadSanitizer examples:
    // tsan_data_race_example();
    // tsan_vector_race_example();
    
    // MemorySanitizer examples:
    // msan_uninitialized_read();
    // msan_uninitialized_array();
    // msan_uninitialized_struct();
    
    // UndefinedBehaviorSanitizer examples:
    // ubsan_signed_overflow();
    // ubsan_null_pointer();
    // ubsan_division_by_zero();
    // ubsan_shift_out_of_bounds();
    // ubsan_array_bounds();
    // ubsan_misaligned_pointer();
    // ubsan_invalid_cast();
    
    std::cout << "No examples executed. Uncomment examples in main() to test.\n";
    return 0;
}

