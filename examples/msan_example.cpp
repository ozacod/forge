// MemorySanitizer Example - Uninitialized Memory
// Compile: forge check --msan
// Run: ./build/msan_example

#include <iostream>

int main() {
    std::cout << "MemorySanitizer Example: Uninitialized Memory\n";
    std::cout << "==============================================\n\n";
    
    // Uninitialized variable
    int x;  // Never initialized!
    
    std::cout << "Reading uninitialized variable x...\n";
    if (x > 0) {  // MSan will catch this!
        std::cout << "x is positive: " << x << std::endl;
    } else {
        std::cout << "x is not positive: " << x << std::endl;
    }
    
    // Uninitialized array
    int arr[5];
    std::cout << "\nReading uninitialized array element...\n";
    std::cout << "arr[0] = " << arr[0] << std::endl;  // MSan will catch this!
    
    // Partially initialized
    int arr2[5] = {1, 2};  // Only first 2 elements initialized
    std::cout << "\nReading partially initialized array...\n";
    std::cout << "arr2[0] = " << arr2[0] << " (OK)\n";
    std::cout << "arr2[2] = " << arr2[2] << " (MSan will catch this!)\n";
    
    return 0;
}

