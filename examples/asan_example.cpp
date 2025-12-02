// AddressSanitizer Example - Buffer Overflow
// Compile: forge check --asan
// Run: ./build/asan_example

#include <iostream>

int main() {
    std::cout << "AddressSanitizer Example: Buffer Overflow\n";
    std::cout << "==========================================\n\n";
    
    // Stack buffer overflow
    int arr[5] = {1, 2, 3, 4, 5};
    
    std::cout << "Writing to arr[10] (out of bounds)...\n";
    arr[10] = 99;  // ASan will catch this!
    
    std::cout << "Value at arr[10]: " << arr[10] << std::endl;
    
    // Use after free
    std::cout << "\nUse after free example:\n";
    int* ptr = new int(42);
    std::cout << "Allocated: " << *ptr << std::endl;
    delete ptr;
    std::cout << "After delete, trying to access: " << *ptr << std::endl;  // ASan will catch this!
    
    return 0;
}

