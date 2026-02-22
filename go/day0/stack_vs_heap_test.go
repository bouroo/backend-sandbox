package main

import (
	"testing"
)

// BenchmarkCreateLargeStructOnStack tests creating large struct by value.
// WHAT IT DEMONSTRATES:
// - LargeStruct created and returned by value
// - Compiler MAY apply RVO, may not - depends on complexity
// - Copy cost is visible either way
//
// EXPECTED RESULTS:
// - 0 heap allocations (with RVO)
// - Copy cost still paid (1KB moved around)
// - Performance depends on RVO application
func BenchmarkCreateLargeStructOnStack(b *testing.B) {
	for i := 0; b.Loop(); i++ {
		_ = createLargeStructOnStack()
	}
}

// BenchmarkCreateLargeStructOnHeap tests creating large struct on heap.
// WHAT IT DEMONSTRATES:
// - Every call allocates 1KB on the heap
// - No RVO because we're returning a pointer
// - Clear demonstration of allocation cost
//
// EXPECTED RESULTS:
// - 1 heap allocation per call (1KB)
// - Slower than stack version
// - Shows the baseline cost of heap allocation
//
// WHY THIS MATTERS:
// This is what heap allocation looks like. Benchmark it against stack!
func BenchmarkCreateLargeStructOnHeap(b *testing.B) {
	for i := 0; b.Loop(); i++ {
		_ = createLargeStructOnHeap()
	}
}
