// ThreadSanitizer Example - Data Race
// Compile: forge check --tsan
// Run: ./build/tsan_example

#include <iostream>
#include <thread>
#include <vector>

int counter = 0;  // Shared variable without synchronization

void increment() {
    for (int i = 0; i < 100000; ++i) {
        counter++;  // Data race! TSan will catch this
    }
}

int main() {
    std::cout << "ThreadSanitizer Example: Data Race\n";
    std::cout << "==================================\n\n";
    
    std::cout << "Starting two threads that increment a shared counter...\n";
    std::cout << "Initial counter: " << counter << std::endl;
    
    std::thread t1(increment);
    std::thread t2(increment);
    
    t1.join();
    t2.join();
    
    std::cout << "Final counter: " << counter << std::endl;
    std::cout << "Expected: 200000, but may be less due to race condition\n";
    std::cout << "TSan will report the data race!\n";
    
    return 0;
}

