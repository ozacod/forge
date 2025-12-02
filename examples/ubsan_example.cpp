// UndefinedBehaviorSanitizer Example - Undefined Behavior
// Compile: forge check --ubsan
// Run: ./build/ubsan_example

#include <iostream>
#include <climits>

int main() {
    std::cout << "UndefinedBehaviorSanitizer Example: Undefined Behavior\n";
    std::cout << "======================================================\n\n";
    
    // Signed integer overflow
    std::cout << "1. Signed integer overflow:\n";
    int x = INT_MAX;
    std::cout << "   x = " << x << std::endl;
    x++;  // UBSan will catch this!
    std::cout << "   x++ = " << x << " (undefined behavior)\n\n";
    
    // Division by zero
    std::cout << "2. Division by zero:\n";
    int a = 10;
    int b = 0;
    // Uncomment to trigger:
    // int result = a / b;  // UBSan will catch this!
    // std::cout << "   10 / 0 = " << result << std::endl;
    std::cout << "   (Commented out to avoid crash)\n\n";
    
    // Shift out of bounds
    std::cout << "3. Shift out of bounds:\n";
    int value = 1;
    int shift = 100;  // Too large
    int shifted = value << shift;  // UBSan will catch this!
    std::cout << "   1 << 100 = " << shifted << " (undefined behavior)\n\n";
    
    // Array index out of bounds (undefined behavior)
    std::cout << "4. Array index out of bounds:\n";
    int arr[5] = {1, 2, 3, 4, 5};
    int index = 10;
    int val = arr[index];  // UBSan will catch this!
    std::cout << "   arr[10] = " << val << " (undefined behavior)\n\n";
    
    // Null pointer dereference
    std::cout << "5. Null pointer dereference:\n";
    int* ptr = nullptr;
    // Uncomment to trigger:
    // *ptr = 42;  // UBSan will catch this!
    std::cout << "   (Commented out to avoid crash)\n";
    
    return 0;
}

