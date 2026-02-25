// Package topics provides Go performance optimization demonstrations.
package topics

import (
	"fmt"
	"runtime"
)

// =============================================================================
// STACK VS HEAP ALLOCATION
// =============================================================================
//
// This file demonstrates the difference between stack and heap allocation
// and when each is used.
//
// WHAT'S THE STACK?
// - Fast, limited size (~MBs)
// - Automatic cleanup (just move stack pointer)
// - Great for short-lived data
//
// WHAT'S THE HEAP?
// - Slower, larger (~GBs)
// - Requires garbage collection
// - For data that outlives its function

// CreateLargeStructOnStack creates a large struct and returns it by value.
//
// WHAT'S HAPPENING:
// - LargeStruct lives on stack (fast, no GC)
// - But copy cost is still paid (1KB copied to caller)
// - Compiler CAN optimize this with RVO, but not guaranteed
//
// KEY TAKEAWAY: Even stack-allocated large structs have copy overhead.
func CreateLargeStructOnStack() LargeStruct {
	s := LargeStruct{
		Field1: 1,
		Field2: 2,
		Field3: 3,
		Field4: 4,
		Field5: 5,
		Field6: 6,
		Field7: 7,
		Field8: 8,
	}
	for i := range len(s.Data) {
		s.Data[i] = int64(i)
	}
	return s
}

// CreateLargeStructOnHeap creates a large struct on the HEAP.
//
// WHAT'S HAPPENING:
// - Go must allocate 1KB on the heap (slower than stack)
// - This allocation triggers the garbage collector
// - Data survives after function returns (caller gets the pointer)
//
// WHY DO WE NEED HEAP?
// - The returned pointer &s must be valid after createLargeStructOnHeap() ends
// - Stack data is automatically cleaned up when function returns
// - So we NEED heap to persist the data beyond the function call!
//
// KEY TAKEAWAY: Heap = when data must outlive its creating function.
func CreateLargeStructOnHeap() *LargeStruct {
	s := LargeStruct{
		Field1: 1,
		Field2: 2,
		Field3: 3,
		Field4: 4,
		Field5: 5,
		Field6: 6,
		Field7: 7,
		Field8: 8,
	}
	for i := range len(s.Data) {
		s.Data[i] = int64(i)
	}
	return &s
}

// =============================================================================
// DEMONSTRATION
// =============================================================================

// RunStackVsHeapDemo demonstrates stack vs heap allocation.
func RunStackVsHeapDemo() {
	fmt.Println("================================================================================")
	fmt.Println("                    STACK VS HEAP ALLOCATION DEMONSTRATION                     ")
	fmt.Println("================================================================================")
	fmt.Println()

	fmt.Println("THE STACK:")
	fmt.Println("  - Fast allocation (just move a pointer)")
	fmt.Println("  - Automatic cleanup (no GC needed)")
	fmt.Println("  - Limited size (~MB)")
	fmt.Println("  - Great for short-lived data")
	fmt.Println()
	fmt.Println("THE HEAP:")
	fmt.Println("  - Slower allocation (requires finding free space)")
	fmt.Println("  - Requires garbage collection")
	fmt.Println("  - Much larger (~GB)")
	fmt.Println("  - For data that outlives its function")

	// Get current runtime info
	fmt.Println()
	fmt.Println("=== Runtime Information ===")
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("NumCPU: %d\n", runtime.NumCPU())
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("GC cycles: %d\n", m.NumGC)

	// Key insights
	fmt.Println()
	fmt.Println("=== Key Insights ===")
	fmt.Println("✓ Stack allocation is ~10-100x faster than heap")
	fmt.Println("✓ Small allocations may not trigger GC at all")
	fmt.Println("✓ Escape analysis happens at compile time")
	fmt.Println("✓ Large objects (> 64KB) go directly to heap")
	fmt.Println("✓ Use pprof to identify heap allocations: go tool pprof")
	fmt.Println()

	fmt.Println("================================================================================")
}
