package main

import (
	"math/rand"
	"testing"
	"unsafe"
)

// BenchmarkStructSizeComparison shows the memory footprint of aligned vs unaligned structs.
// Run with: go test -v -run=NONE -bench=BenchmarkStructSize
// Expected output shows UnalignedStruct (48 bytes) vs AlignedStruct (32 bytes)
func BenchmarkStructSize(b *testing.B) {
	// These benchmarks just measure sizeof - the compiler may optimize away unused values
	_ = unsafe.Sizeof(UnalignedStruct{})
	_ = unsafe.Sizeof(AlignedStruct{})
	_ = unsafe.Sizeof(PoorlyPaddedStruct{})
	_ = unsafe.Sizeof(MixedTypesAligned{})
	_ = unsafe.Sizeof(MixedTypesUnaligned{})
	_ = unsafe.Sizeof(UnalignedInts{})
	_ = unsafe.Sizeof(AlignedInts{})
}

// BenchmarkProcessUnaligned tests processing with poorly aligned struct.
// WHAT IT DEMONSTRATES:
// - UnalignedStruct has 48 bytes (24 data + 24 padding)
// - Passing by value copies all 48 bytes
// - Cache line efficiency is poor due to padding
//
// EXPECTED RESULTS:
// - Slower than aligned version due to larger copy size
// - More memory bandwidth used
func BenchmarkProcessUnaligned(b *testing.B) {
	s := UnalignedStruct{
		Field1: 1,
		Field2: 2,
		Field3: 3,
		Field4: 4,
		Field5: 5,
		Field6: 6,
	}
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_ = ProcessUnaligned(s)
	}
}

// BenchmarkProcessAligned tests processing with well-aligned struct.
// WHAT IT DEMONSTRATES:
// - AlignedStruct has 32 bytes (24 data + 8 padding)
// - Passing by value copies 16 bytes LESS than unaligned
// - Better cache line utilization
//
// EXPECTED RESULTS:
// - Faster than unaligned version
// - Less memory bandwidth
func BenchmarkProcessAligned(b *testing.B) {
	s := AlignedStruct{
		Field1: 1,
		Field2: 2,
		Field3: 3,
		Field4: 4,
		Field5: 5,
		Field6: 6,
	}
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_ = ProcessAligned(s)
	}
}

// BenchmarkProcessUnalignedPtr tests processing unaligned struct by pointer.
// WHAT IT DEMONSTRATES:
// - Pointer passes only 8 bytes regardless of struct size
// - Indirection cost but smaller data transfer
func BenchmarkProcessUnalignedPtr(b *testing.B) {
	s := UnalignedStruct{
		Field1: 1,
		Field2: 2,
		Field3: 3,
		Field4: 4,
		Field5: 5,
		Field6: 6,
	}
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_ = ProcessUnalignedPtr(&s)
	}
}

// BenchmarkProcessAlignedPtr tests processing aligned struct by pointer.
// WHAT IT DEMONSTRATES:
// - Pointer passing avoids copy overhead entirely
// - Best of both worlds: good layout + no copy
func BenchmarkProcessAlignedPtr(b *testing.B) {
	s := AlignedStruct{
		Field1: 1,
		Field2: 2,
		Field3: 3,
		Field4: 4,
		Field5: 5,
		Field6: 6,
	}
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_ = ProcessAlignedPtr(&s)
	}
}

// BenchmarkUnalignedIntsVsAlignedInts demonstrates the classic int8/int64 interleaving problem.
// This is the most common real-world alignment issue.
//
// WHY THIS MATTERS:
// Many protobuf or database structs have this problem!
// Check your structs with unsafe.Sizeof()
func BenchmarkUnalignedIntsVsAlignedInts(b *testing.B) {
	// Show size difference - compiler may optimize away
	_ = unsafe.Sizeof(UnalignedInts{})
	_ = unsafe.Sizeof(AlignedInts{})
}

// =============================================================================
// COMPREHENSIVE ALIGNMENT BENCHMARKS - Slice-Based with Cache Effects
// =============================================================================

// Benchmark constants - large enough to exceed L1/L2 cache
const (
	testBenchSliceSize   = 1000000 // 1M elements to exceed L1/L2 cache
	testBenchCacheLine   = 64      // Typical cache line size
	testBenchL1CacheSize = 32 * 1024
)

// createUnalignedSlice creates a large slice of UnalignedStruct
func createUnalignedSlice(size int) []UnalignedStruct {
	data := make([]UnalignedStruct, size)
	for i := range data {
		data[i] = UnalignedStruct{
			Field1: int8(i % 256),
			Field2: int64(i),
			Field3: int8((i + 1) % 256),
			Field4: int64(i + 1),
			Field5: int8((i + 2) % 256),
			Field6: int64(i + 2),
		}
	}
	return data
}

// createAlignedSlice creates a large slice of AlignedStruct
func createAlignedSlice(size int) []AlignedStruct {
	data := make([]AlignedStruct, size)
	for i := range data {
		data[i] = AlignedStruct{
			Field2: int64(i),
			Field4: int64(i + 1),
			Field6: int64(i + 2),
			Field1: int8(i % 256),
			Field3: int8((i + 1) % 256),
			Field5: int8((i + 2) % 256),
		}
	}
	return data
}

// =============================================================================
// SEQUENTIAL ACCESS BENCHMARKS
// =============================================================================

// BenchmarkSequentialUnaligned benchmarks sequential access through UnalignedStruct slice
func BenchmarkSequentialUnaligned(b *testing.B) {
	data := createUnalignedSlice(testBenchSliceSize)

	b.SetBytes(int64(testBenchSliceSize * int(unsafe.Sizeof(UnalignedStruct{}))))

	var sum int64
	for b.Loop() {
		for j := range data {
			sum += data[j].Field2 + data[j].Field4 + data[j].Field6
		}
	}
	_ = sum // Prevent optimization
}

// BenchmarkSequentialAligned benchmarks sequential access through AlignedStruct slice
func BenchmarkSequentialAligned(b *testing.B) {
	data := createAlignedSlice(testBenchSliceSize)
	b.ResetTimer()
	b.SetBytes(int64(testBenchSliceSize * int(unsafe.Sizeof(AlignedStruct{}))))

	var sum int64
	for b.Loop() {
		for j := range data {
			sum += data[j].Field2 + data[j].Field4 + data[j].Field6
		}
	}
	_ = sum // Prevent optimization
}

// BenchmarkSequentialUnalignedInts benchmarks sequential access through UnalignedInts slice
func BenchmarkSequentialUnalignedInts(b *testing.B) {
	data := make([]UnalignedInts, testBenchSliceSize)
	for i := range data {
		data[i] = UnalignedInts{
			A: int8(i % 256),
			B: int64(i),
			C: int8((i + 1) % 256),
			D: int64(i + 1),
			E: int8((i + 2) % 256),
			F: int64(i + 2),
		}
	}
	b.ResetTimer()
	b.SetBytes(int64(testBenchSliceSize * int(unsafe.Sizeof(UnalignedInts{}))))

	var sum int64
	for b.Loop() {
		for j := range data {
			sum += data[j].B + data[j].D + data[j].F
		}
	}
	_ = sum
}

// BenchmarkSequentialAlignedInts benchmarks sequential access through AlignedInts slice
func BenchmarkSequentialAlignedInts(b *testing.B) {
	data := make([]AlignedInts, testBenchSliceSize)
	for i := range data {
		data[i] = AlignedInts{
			B: int64(i),
			D: int64(i + 1),
			F: int64(i + 2),
			A: int8(i % 256),
			C: int8((i + 1) % 256),
			E: int8((i + 2) % 256),
		}
	}
	b.ResetTimer()
	b.SetBytes(int64(testBenchSliceSize * int(unsafe.Sizeof(AlignedInts{}))))

	var sum int64
	for b.Loop() {
		for j := range data {
			sum += data[j].B + data[j].D + data[j].F
		}
	}
	_ = sum
}

// =============================================================================
// RANDOM ACCESS BENCHMARKS (Cache-Unfriendly Patterns)
// =============================================================================

// BenchmarkRandomUnaligned benchmarks random access through UnalignedStruct slice
func BenchmarkRandomUnaligned(b *testing.B) {
	data := createUnalignedSlice(testBenchSliceSize)
	// Create random indices for random access pattern
	indices := make([]int, testBenchSliceSize)
	rng := rand.New(rand.NewSource(42))
	for i := range indices {
		indices[i] = rng.Intn(len(data))
	}

	var sum int64
	for b.Loop() {
		for j := range indices {
			sum += data[indices[j]].Field2 + data[indices[j]].Field4 + data[indices[j]].Field6
		}
	}
	_ = sum
}

// BenchmarkRandomAligned benchmarks random access through AlignedStruct slice
func BenchmarkRandomAligned(b *testing.B) {
	data := createAlignedSlice(testBenchSliceSize)
	// Create random indices for random access pattern
	indices := make([]int, testBenchSliceSize)
	rng := rand.New(rand.NewSource(42))
	for i := range indices {
		indices[i] = rng.Intn(len(data))
	}
	b.ResetTimer()

	var sum int64
	for b.Loop() {
		for j := range indices {
			sum += data[indices[j]].Field2 + data[indices[j]].Field4 + data[indices[j]].Field6
		}
	}
	_ = sum
}

// BenchmarkStridedUnaligned benchmarks strided access (every Nth element)
func BenchmarkStridedUnaligned(b *testing.B) {
	data := createUnalignedSlice(testBenchSliceSize)
	stride := 64 // One cache line apart
	b.ResetTimer()

	var sum int64
	for b.Loop() {
		for j := 0; j < len(data); j += stride {
			sum += data[j].Field2 + data[j].Field4 + data[j].Field6
		}
	}
	_ = sum
}

// BenchmarkStridedAligned benchmarks strided access (every Nth element)
func BenchmarkStridedAligned(b *testing.B) {
	data := createAlignedSlice(testBenchSliceSize)
	stride := 64 // One cache line apart
	b.ResetTimer()

	var sum int64
	for b.Loop() {
		for j := 0; j < len(data); j += stride {
			sum += data[j].Field2 + data[j].Field4 + data[j].Field6
		}
	}
	_ = sum
}

// =============================================================================
// VALUE VS POINTER PASSING BENCHMARKS
// =============================================================================

// BenchmarkValueUnaligned benchmarks passing UnalignedStruct by value
func BenchmarkValueUnaligned(b *testing.B) {
	data := createUnalignedSlice(testBenchSliceSize)
	b.ResetTimer()

	var sum int64
	for b.Loop() {
		for j := range data {
			sum += ProcessUnaligned(data[j])
		}
	}
	_ = sum
}

// BenchmarkValueAligned benchmarks passing AlignedStruct by value
func BenchmarkValueAligned(b *testing.B) {
	data := createAlignedSlice(testBenchSliceSize)
	b.ResetTimer()

	var sum int64
	for b.Loop() {
		for j := range data {
			sum += ProcessAligned(data[j])
		}
	}
	_ = sum
}

// BenchmarkPointerUnaligned benchmarks passing UnalignedStruct by pointer
func BenchmarkPointerUnaligned(b *testing.B) {
	data := createUnalignedSlice(testBenchSliceSize)
	b.ResetTimer()

	var sum int64
	for b.Loop() {
		for j := range data {
			sum += ProcessUnalignedPtr(&data[j])
		}
	}
	_ = sum
}

// BenchmarkPointerAligned benchmarks passing AlignedStruct by pointer
func BenchmarkPointerAligned(b *testing.B) {
	data := createAlignedSlice(testBenchSliceSize)
	b.ResetTimer()

	var sum int64
	for b.Loop() {
		for j := range data {
			sum += ProcessAlignedPtr(&data[j])
		}
	}
	_ = sum
}

// =============================================================================
// MIXED TYPES BENCHMARKS
// =============================================================================

// BenchmarkMixedTypesAligned benchmarks access through MixedTypesAligned slice
func BenchmarkMixedTypesAligned(b *testing.B) {
	data := make([]MixedTypesAligned, testBenchSliceSize)
	for i := range data {
		val := int64(i)
		data[i] = MixedTypesAligned{
			Pointer: &val,
			Float:   float64(i),
			Counter: val,
			Count:   int32(i),
			Flag:    float32(i),
			Short:   int16(i),
			Char:    int16(i + 1),
			Byte:    byte(i % 256),
			Bool:    i%2 == 0,
		}
	}
	b.ResetTimer()
	b.SetBytes(int64(testBenchSliceSize * int(unsafe.Sizeof(MixedTypesAligned{}))))

	var sum int64
	for b.Loop() {
		for j := range data {
			sum += data[j].Counter
		}
	}
	_ = sum
}

// BenchmarkMixedTypesUnaligned benchmarks access through MixedTypesUnaligned slice
func BenchmarkMixedTypesUnaligned(b *testing.B) {
	data := make([]MixedTypesUnaligned, testBenchSliceSize)
	for i := range data {
		val := int64(i)
		data[i] = MixedTypesUnaligned{
			Bool:    i%2 == 0,
			Byte:    byte(i % 256),
			Char:    int16(i + 1),
			Short:   int16(i),
			Flag:    float32(i),
			Count:   int32(i),
			Counter: val,
			Float:   float64(i),
			Pointer: &val,
		}
	}
	b.ResetTimer()
	b.SetBytes(int64(testBenchSliceSize * int(unsafe.Sizeof(MixedTypesUnaligned{}))))

	var sum int64
	for b.Loop() {
		for j := range data {
			sum += data[j].Counter
		}
	}
	_ = sum
}

// =============================================================================
// HELPER FUNCTION FOR DEMO
// =============================================================================

// Used by RunAlignmentDemo for quick timing test
func init() {
	rand.Seed(42)
}
