// Package topics provides Go performance optimization demonstrations.
package topics

import (
	"fmt"
	"time"
)

// =============================================================================
// SLICE ESCAPE ANALYSIS
// =============================================================================
//
// This file demonstrates when slices escape to heap and when they stay
// on the stack.
//
// WHAT'S A SLICE?
// - A slice has 3 parts: pointer (data location), length, capacity
// - The pointer points to an underlying array (which CAN be on heap)

// globalSlice is a global variable used to demonstrate escape analysis.
// Global variables force anything assigned to them to escape to heap.
var globalSlice []int

// ProcessSliceWithEscape demonstrates when slices ESCAPE to heap.
//
// WHY ESCAPE HERE?
// - We assign s to globalSlice (global variable)
// - Global variables live for the entire program run
// - Therefore s's data MUST survive function scope → HEAP allocation!
//
// ANALOGY: Writing something in a notebook vs. publishing a book.
//          Global = published book (can't be taken back!)
//
// KEY TAKEAWAY: Assigning to globals/returning/storing = escape to heap.
func ProcessSliceWithEscape(n int) int {
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}
	// Assignment to global escapes to heap
	// The slice header AND its data both escape
	globalSlice = s
	return len(s)
}

// ProcessSliceNoEscape demonstrates when slices STAY ON STACK.
//
// WHY NO ESCAPE?
// - Slice s is only used locally within this function
// - It's never returned, stored, or shared with anything that escapes
// - Go's escape analysis: "This slice dies with the function" → stack allocation!
//
// BENEFIT: No heap allocation, no GC pressure
// COST: None! This is ideal.
//
// ANALOGY: Working on scratch paper - throw it away when done, no filing needed.
//
// KEY TAKEAWAY: Keep data local = stack = fast! Don't return/store unless needed.
func ProcessSliceNoEscape(n int) int {
	s := make([]int, n)
	sum := 0
	for i := range s {
		s[i] = i
		sum += s[i]
	}
	return sum
}

// =============================================================================
// DEMONSTRATION
// =============================================================================

// RunSliceEscapeDemo demonstrates escape analysis with slices.
func RunSliceEscapeDemo() {
	fmt.Println("================================================================================")
	fmt.Println("                    SLICE ESCAPE ANALYSIS DEMONSTRATION                        ")
	fmt.Println("================================================================================")
	fmt.Println()

	fmt.Println("WHAT IS ESCAPE ANALYSIS?")
	fmt.Println("Go determines at compile time whether data can stay on the stack")
	fmt.Println("or must be moved to the heap (where it 'escapes' the function).")
	fmt.Println()
	fmt.Println("STAY ON STACK (no escape):")
	fmt.Println("  ✓ Fast allocation (just move stack pointer)")
	fmt.Println("  ✓ No GC pressure")
	fmt.Println("  ✓ Automatic cleanup")
	fmt.Println()
	fmt.Println("ESCAPE TO HEAP:")
	fmt.Println("  ✗ Requires allocation")
	fmt.Println("  ✗ GC must track and collect")
	fmt.Println("  ✗ Slower than stack")

	// Demonstrate escape with timing
	fmt.Println()
	fmt.Println("=== Escape Demonstration (Timing Test) ===")

	// Run timing test to show difference
	for _, size := range []int{100, 1000, 10000} {
		start := time.Now()
		var escapeSum, noEscapeSum int

		// Warm up
		for i := range 1000 {
			if i%2 == 0 {
				escapeSum += ProcessSliceWithEscape(size)
			} else {
				noEscapeSum += ProcessSliceNoEscape(size)
			}
		}
		elapsed := time.Since(start)

		_ = escapeSum
		_ = noEscapeSum

		fmt.Printf("Slice size %5d: %v for 1000 iterations\n", size, elapsed)
	}

	// Guidelines
	fmt.Println()
	fmt.Println("=== Guidelines ===")
	fmt.Println("✓ Keep data local to avoid escape")
	fmt.Println("✓ Avoid assigning to global variables")
	fmt.Println("✓ Don't return slices unnecessarily")
	fmt.Println("✓ Use //go:noinline to prevent optimization if testing")
	fmt.Println()

	fmt.Println("================================================================================")
}
