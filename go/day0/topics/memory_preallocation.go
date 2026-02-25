// Package topics provides Go performance optimization demonstrations.
package topics

import (
	"fmt"
	"time"
)

// =============================================================================
// MEMORY PREALLOCATION
// =============================================================================
//
// This file demonstrates the performance benefits of preallocating memory
// for slices and maps, avoiding repeated reallocations during growth.
//
// ANALOGY:
// - Dynamic growth: Renting a bigger apartment every time your family grows
// - Preallocation: Buying a big house from the start
//
// BENEFITS:
// - Avoids multiple memory allocations
// - Reduces GC pressure
// - Improves performance in hot paths
// - Better memory locality

// =============================================================================
// SLICE PREALLOCATION
// =============================================================================

// DynamicSlice demonstrates GROWTH WITHOUT preallocation.
//
// WHAT'S HAPPENING:
// - Start with empty slice
// - Each append may trigger reallocation when capacity is exceeded
// - Go doubles capacity each time (1, 2, 4, 8, 16, ...)
// - Each reallocation: allocate new memory, copy old data, free old memory
//
// USE WHEN: You don't know the final size and it's small
func DynamicSlice(n int) []int {
	s := []int{}
	for i := range n {
		s = append(s, i)
	}
	return s
}

// PreallocatedSlice demonstrates GROWTH WITH preallocation.
//
// WHAT'S HAPPENING:
// - Allocate exact capacity upfront
// - No reallocations needed during growth
// - Single allocation, single copy (if any)
//
// USE WHEN: You know or can estimate the final size
func PreallocatedSlice(n int) []int {
	s := make([]int, 0, n) // Preallocate capacity n
	for i := range n {
		s = append(s, i)
	}
	return s
}

// PreallocatedSliceExact demonstrates using make with exact size.
//
// WHAT'S HAPPENING:
// - Allocate exact size AND capacity
// - Can use index assignment instead of append
// - Most efficient when you know the exact size
//
// USE WHEN: You know the exact final size
func PreallocatedSliceExact(n int) []int {
	s := make([]int, n) // Allocate size n
	for i := range n {
		s[i] = i
	}
	return s
}

// =============================================================================
// MAP PREALLOCATION
// =============================================================================

// DynamicMap demonstrates growth without preallocation.
//
// WHAT'S HAPPENING:
// - Start with empty map
// - Each insertion may trigger rehashing when load factor exceeds threshold
// - Rehashing: allocate new bucket array, redistribute all entries
//
// USE WHEN: You don't know the number of entries
func DynamicMap(n int) map[string]int {
	m := map[string]int{}
	for i := range n {
		m[fmt.Sprintf("key%d", i)] = i
	}
	return m
}

// PreallocatedMap demonstrates growth with preallocation.
//
// WHAT'S HAPPENING:
// - Preallocate with expected size using make(map[type]type, size)
// - Gives the runtime a hint for initial bucket count
// - Reduces rehashing during growth
//
// USE WHEN: You know or can estimate the number of entries
func PreallocatedMap(n int) map[string]int {
	m := make(map[string]int, n) // Preallocate for n entries
	for i := range n {
		m[fmt.Sprintf("key%d", i)] = i
	}
	return m
}

// =============================================================================
// PERFORMANCE COMPARISON
// =============================================================================

// benchmarkSliceComparison compares dynamic vs preallocated slice growth.
func benchmarkSliceComparison(size int) (dynamicTime, preallocatedTime time.Duration) {
	// Warm up
	DynamicSlice(100)
	PreallocatedSlice(100)

	// Test dynamic growth
	start := time.Now()
	for range 1000 {
		_ = DynamicSlice(size)
	}
	dynamicTime = time.Since(start)

	// Test preallocated growth
	start = time.Now()
	for range 1000 {
		_ = PreallocatedSlice(size)
	}
	preallocatedTime = time.Since(start)

	return
}

// benchmarkMapComparison compares dynamic vs preallocated map growth.
func benchmarkMapComparison(size int) (dynamicTime, preallocatedTime time.Duration) {
	// Warm up
	DynamicMap(100)
	PreallocatedMap(100)

	// Test dynamic growth
	start := time.Now()
	for range 1000 {
		_ = DynamicMap(size)
	}
	dynamicTime = time.Since(start)

	// Test preallocated growth
	start = time.Now()
	for range 1000 {
		_ = PreallocatedMap(size)
	}
	preallocatedTime = time.Since(start)

	return
}

// =============================================================================
// DEMONSTRATION
// =============================================================================

// RunMemoryPreallocationDemo demonstrates the performance impact of preallocation.
func RunMemoryPreallocationDemo() {
	fmt.Println("================================================================================")
	fmt.Println("                MEMORY PREALLOCATION DEMONSTRATION                              ")
	fmt.Println("================================================================================")
	fmt.Println()

	fmt.Println("WHAT IS PREALLOCATION?")
	fmt.Println("Preallocating memory for slices and maps before adding elements.")
	fmt.Println()
	fmt.Println("DYNAMIC GROWTH (without preallocation):")
	fmt.Println("  - Start with small capacity")
	fmt.Println("  - Trigger reallocation when capacity exceeded")
	fmt.Println("  - Each reallocation: copy all data to new location")
	fmt.Println("  - More allocations = more GC pressure")
	fmt.Println()
	fmt.Println("PREALLOCATION (with capacity hint):")
	fmt.Println("  - Allocate capacity upfront")
	fmt.Println("  - No reallocations during growth")
	fmt.Println("  - Single allocation, better performance")
	fmt.Println()

	// Demonstrate slice growth
	fmt.Println("=== SLICE GROWTH COMPARISON ===")
	for _, size := range []int{100, 1000, 10000} {
		dynamic, preallocated := benchmarkSliceComparison(size)
		speedup := float64(dynamic.Nanoseconds()) / float64(preallocated.Nanoseconds())
		fmt.Printf("Size %6d: Dynamic: %12v, Preallocated: %12v, Speedup: %.2fx\n",
			size, dynamic, preallocated, speedup)
	}
	fmt.Println()

	// Demonstrate map growth
	fmt.Println("=== MAP GROWTH COMPARISON ===")
	for _, size := range []int{100, 1000, 10000} {
		dynamic, preallocated := benchmarkMapComparison(size)
		speedup := float64(dynamic.Nanoseconds()) / float64(preallocated.Nanoseconds())
		fmt.Printf("Size %6d: Dynamic: %12v, Preallocated: %12v, Speedup: %.2fx\n",
			size, dynamic, preallocated, speedup)
	}
	fmt.Println()

	// Guidelines
	fmt.Println("=== GUIDELINES ===")
	fmt.Println("PREALLOCATE SLICES when:")
	fmt.Println("  - You know or can estimate the final size")
	fmt.Println("  - Working in tight loops (hot paths)")
	fmt.Println("  - Building up a slice incrementally")
	fmt.Println()
	fmt.Println("PREALLOCATE MAPS when:")
	fmt.Println("  - You know the approximate number of entries")
	fmt.Println("  - Inserting many items at once")
	fmt.Println("  - Performance is critical")
	fmt.Println()
	fmt.Println("DON'T PREALLOCATE when:")
	fmt.Println("  - Size is unknown and could be very large")
	fmt.Println("  - Memory is constrained")
	fmt.Println("  - Code readability matters more than micro-optimization")
	fmt.Println()

	// Syntax examples
	fmt.Println("=== SYNTAX EXAMPLES ===")
	fmt.Println("// Preallocate slice with capacity")
	fmt.Println("s := make([]int, 0, 100)  // length=0, capacity=100")
	fmt.Println()
	fmt.Println("// Preallocate slice with exact size")
	fmt.Println("s := make([]int, 100)     // length=100, capacity=100")
	fmt.Println()
	fmt.Println("// Preallocate map with size hint")
	fmt.Println("m := make(map[string]int, 100)  // preallocate for 100 entries")
	fmt.Println()

	fmt.Println("================================================================================")
}
