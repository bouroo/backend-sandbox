package main

import (
	"testing"
)

// BenchmarkAddByValue tests pass-by-value with LargeStruct (1KB).
// WHAT IT DEMONSTRATES:
// - Passing large structs BY VALUE means copying the entire 1KB
// - All data stays on the stack - no heap allocation, no GC pressure
// - Copy cost exists but is predictable (stack-to-stack copy is fast)
//
// EXPECTED RESULTS:
// - 0 heap allocations per operation
// - May be slower than pointer due to copy overhead
// - BUT: No GC pressure = more consistent performance over time
//
// WHY THIS MATTERS:
// For hot paths with large structs, this copy cost adds up.
// Consider: Is the data small enough that copy is cheaper than pointer indirection?
func BenchmarkAddByValue(b *testing.B) {
	a := LargeStruct{Field1: 1, Field2: 2}
	bVal := LargeStruct{Field3: 3, Field4: 4}
	for i := range len(a.Data) {
		a.Data[i] = int64(i)
		bVal.Data[i] = int64(i * 2)
	}
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_ = addByValue(a, bVal)
	}
}

// BenchmarkAddByPointer tests pass-by-pointer with LargeStruct.
// WHAT IT DEMONSTRATES:
// - Passing by pointer only copies the 8-byte pointer
// - The pointers are on the stack, pointing to heap data
// - No copy of the actual 1KB data
//
// EXPECTED RESULTS:
// - 0 heap allocations (data was already on heap from setup)
// - Slightly faster than value due to avoiding 1KB copy
// - BUT: Potential cache misses if data not local
//
// WHY THIS MATTERS:
// Pointers avoid the copy but add indirection overhead.
// The real cost is hidden in cache behavior and whether data is already in CPU cache.
func BenchmarkAddByPointer(b *testing.B) {
	a := LargeStruct{Field1: 1, Field2: 2}
	bVal := LargeStruct{Field3: 3, Field4: 4}
	for i := range len(a.Data) {
		a.Data[i] = int64(i)
		bVal.Data[i] = int64(i * 2)
	}
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_ = addByPointer(&a, &bVal)
	}
}
