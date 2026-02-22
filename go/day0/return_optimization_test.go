package main

import (
	"testing"
)

// BenchmarkReturnAddByValue tests return-by-value with RVO.
// WHAT IT DEMONSTRATES:
// - Return Value Optimization (RVO): compiler allocates return space in caller
// - No copy of the large struct - it's built directly where it's needed
// - Stack-only operation, no heap allocation
//
// EXPECTED RESULTS:
// - 0 heap allocations per operation
// - As fast as local variables (just stack operations)
// - Best of both worlds: clean code + great performance!
//
// WHY THIS MATTERS:
// RVO is one of Go's best optimizations. Return by value when you can!
func BenchmarkReturnAddByValue(b *testing.B) {
	a := LargeStruct{Field1: 1, Field2: 2}
	bVal := LargeStruct{Field3: 3, Field4: 4}
	for i := range len(a.Data) {
		a.Data[i] = int64(i)
		bVal.Data[i] = int64(i * 2)
	}
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_ = returnAddByValue(a, bVal)
	}
}

// BenchmarkReturnAddByPointer tests returning a pointer to local data.
// WHAT IT DEMONSTRATES:
// - The LOCAL VARIABLE ESCAPES to heap because we return its address
// - Go proves the data outlives the function → must allocate on heap
// - This is the KEY benchmark showing heap vs stack difference!
//
// EXPECTED RESULTS:
// - 1 heap allocation per operation (the 1KB struct)
// - Slower than stack due to allocation + GC overhead
// - More variable performance due to GC pauses
//
// WHY THIS MATTERS:
// This is the "expensive" pattern. Every call allocates 1KB on heap!
// In real code, doing this in a tight loop = GC pressure = stuttering.
//
// TAKEAWAY: Prefer returning by value when possible. Let RVO help you!
func BenchmarkReturnAddByPointer(b *testing.B) {
	a := LargeStruct{Field1: 1, Field2: 2}
	bVal := LargeStruct{Field3: 3, Field4: 4}
	for i := range len(a.Data) {
		a.Data[i] = int64(i)
		bVal.Data[i] = int64(i * 2)
	}
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_ = returnAddByPointer(a, bVal)
	}
}
