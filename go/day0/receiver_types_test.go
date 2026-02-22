package main

import (
	"testing"
)

// BenchmarkIncrementByValue tests value receiver method.
// WHAT IT DEMONSTRATES:
// - Value receiver = copy of struct for each call
// - For small Counter (8 bytes), copy cost is tiny
// - Original is never modified
//
// EXPECTED RESULTS:
// - 0 heap allocations
// - Slightly slower than pointer (copy overhead)
// - But simpler reasoning: no aliasing possible
//
// WHY THIS MATTERS:
// For small structs, value receivers are fine and prevent subtle bugs.
// The copy is cheap for small data!
func BenchmarkIncrementByValue(b *testing.B) {
	c := Counter{value: 0}
	for i := 0; b.Loop(); i++ {
		_ = c.IncrementByValue()
	}
}

// BenchmarkIncrementByPointer tests pointer receiver method.
// WHAT IT DEMONSTRATES:
// - Pointer receiver = no copy, just 8-byte pointer passed
// - Modifies original directly
// - Slight risk of race conditions if shared across goroutines
//
// EXPECTED RESULTS:
// - 0 heap allocations
// - Potentially faster (no copy)
// - But needs careful synchronization in concurrent code
//
// WHY THIS MATTERS:
// Pointer receivers are efficient but bring concurrency considerations.
// Choose based on whether you need mutation and thread safety.
func BenchmarkIncrementByPointer(b *testing.B) {
	c := Counter{value: 0}
	for i := 0; b.Loop(); i++ {
		_ = c.IncrementByPointer()
	}
}

// BenchmarkProcessByValue tests value receiver with large struct.
// WHAT IT DEMONSTRATES:
// - Value receiver on 1KB struct = 1KB copy EVERY call!
// - This is expensive in tight loops
// - The KEY comparison against pointer receiver
//
// EXPECTED RESULTS:
// - 0 heap allocations
// - 1KB copy per call = slow
// - This is what NOT to do with large structs
//
// WHY THIS MATTERS:
// This shows why receiver type matters for large types.
// Always use pointer receiver for LargeStruct!
func BenchmarkProcessByValue(b *testing.B) {
	dp := DataProcessor{}
	for i := range len(dp.Data) {
		dp.Data[i] = int64(i)
	}
	for i := 0; b.Loop(); i++ {
		_ = dp.ProcessByValue()
	}
}

// BenchmarkProcessByPointer tests pointer receiver with large struct.
// WHAT IT DEMONSTRATES:
// - Pointer receiver = just 8-byte pointer passed
// - No copy of the 1KB data
// - Much more efficient for large structs
//
// EXPECTED RESULTS:
// - 0 heap allocations
// - Much faster than value receiver
// - Should see significant speedup vs ProcessByValue
//
// WHY THIS MATTERS:
// This is the correct pattern for large structs!
// The small pointer overhead is worth avoiding the 1KB copy.
func BenchmarkProcessByPointer(b *testing.B) {
	dp := DataProcessor{}
	for i := range len(dp.Data) {
		dp.Data[i] = int64(i)
	}
	for i := 0; b.Loop(); i++ {
		_ = dp.ProcessByPointer()
	}
}
