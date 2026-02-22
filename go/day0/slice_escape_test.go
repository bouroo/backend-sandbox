package main

import (
	"testing"
)

// BenchmarkProcessSliceWithEscape tests slice that escapes to heap.
// WHAT IT DEMONSTRATES:
// - Slice assigned to global = escape to heap
// - Both slice header AND underlying array go to heap
// - GC must track and eventually collect this memory
//
// EXPECTED RESULTS:
// - Heap allocations for both the slice header and data
// - More GC pressure than non-escaping version
// - Performance degrades as heap grows
//
// WHY THIS MATTERS:
// Global variables are the "silent killer" of performance.
// They force everything they touch onto the heap!
func BenchmarkProcessSliceWithEscape(b *testing.B) {
	for i := 0; b.Loop(); i++ {
		_ = processSliceWithEscape(1000)
	}
}

// BenchmarkProcessSliceNoEscape tests slice that stays on stack.
// WHAT IT DEMONSTRATES:
// - Slice only used locally = no escape
// - Go keeps everything on stack = super fast
// - No heap allocation, no GC pressure
//
// EXPECTED RESULTS:
// - 0 heap allocations
// - Much faster than escaping version
// - Ideal performance pattern
//
// WHY THIS MATTERS:
// This is the "golden path" - keep data local, avoid escaping!
func BenchmarkProcessSliceNoEscape(b *testing.B) {
	for i := 0; b.Loop(); i++ {
		_ = processSliceNoEscape(1000)
	}
}
