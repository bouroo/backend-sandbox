package benchmarks

import (
	"testing"

	"day0/topics"
)

// =============================================================================
// STRUCT ALIGNMENT BENCHMARKS
// =============================================================================

func BenchmarkProcessUnaligned(b *testing.B) {
	s := topics.UnalignedStruct{
		Field1: 1,
		Field2: 2,
		Field3: 3,
		Field4: 4,
		Field5: 5,
		Field6: 6,
	}

	for b.Loop() {
		_ = topics.ProcessUnaligned(s)
	}
}

func BenchmarkProcessAligned(b *testing.B) {
	s := topics.AlignedStruct{
		Field1: 1,
		Field2: 2,
		Field3: 3,
		Field4: 4,
		Field5: 5,
		Field6: 6,
	}

	for b.Loop() {
		_ = topics.ProcessAligned(s)
	}
}

func BenchmarkProcessUnalignedPtr(b *testing.B) {
	s := &topics.UnalignedStruct{
		Field1: 1,
		Field2: 2,
		Field3: 3,
		Field4: 4,
		Field5: 5,
		Field6: 6,
	}

	for b.Loop() {
		_ = topics.ProcessUnalignedPtr(s)
	}
}

func BenchmarkProcessAlignedPtr(b *testing.B) {
	s := &topics.AlignedStruct{
		Field1: 1,
		Field2: 2,
		Field3: 3,
		Field4: 4,
		Field5: 5,
		Field6: 6,
	}

	for b.Loop() {
		_ = topics.ProcessAlignedPtr(s)
	}
}

func BenchmarkMixedTypesAligned(b *testing.B) {
	var i int64 = 42
	s := topics.MixedTypesAligned{
		Pointer: &i,
		Float:   3.14,
		Counter: 100,
		Count:   50,
		Flag:    2.5,
		Short:   10,
		Char:    5,
		Byte:    1,
		Bool:    true,
	}

	for b.Loop() {
		_ = s.Counter + int64(s.Count)
	}
}

func BenchmarkMixedTypesUnaligned(b *testing.B) {
	var i int64 = 42
	s := topics.MixedTypesUnaligned{
		Bool:    true,
		Byte:    1,
		Char:    5,
		Short:   10,
		Flag:    2.5,
		Count:   50,
		Counter: 100,
		Float:   3.14,
		Pointer: &i,
	}

	for b.Loop() {
		_ = s.Counter + int64(s.Count)
	}
}

// =============================================================================
// PASS BY VALUE VS POINTER BENCHMARKS
// =============================================================================

func BenchmarkAddByValue(b *testing.B) {
	a := topics.LargeStruct{Field1: 1, Field2: 2}
	c := topics.LargeStruct{Field1: 3, Field2: 4}

	for b.Loop() {
		_ = topics.AddByValue(a, c)
	}
}

func BenchmarkAddByPointer(b *testing.B) {
	a := &topics.LargeStruct{Field1: 1, Field2: 2}
	c := &topics.LargeStruct{Field1: 3, Field2: 4}

	for b.Loop() {
		_ = topics.AddByPointer(a, c)
	}
}

// =============================================================================
// RECEIVER TYPES BENCHMARKS
// =============================================================================

func BenchmarkIncrementByValue(b *testing.B) {
	c := topics.Counter{}

	for b.Loop() {
		_ = c.IncrementByValue()
	}
}

func BenchmarkIncrementByPointer(b *testing.B) {
	c := &topics.Counter{}

	for b.Loop() {
		_ = c.IncrementByPointer()
	}
}

func BenchmarkProcessByValue(b *testing.B) {
	dp := topics.DataProcessor{}
	dp.Field1 = 1
	dp.Field2 = 2
	dp.Field3 = 3
	dp.Field4 = 4
	dp.Field5 = 5
	dp.Field6 = 6
	dp.Field7 = 7
	dp.Field8 = 8
	for i := range dp.Data {
		dp.Data[i] = int64(i)
	}

	for b.Loop() {
		_ = dp.ProcessByValue()
	}
}

func BenchmarkProcessByPointer(b *testing.B) {
	dp := &topics.DataProcessor{}
	dp.Field1 = 1
	dp.Field2 = 2
	dp.Field3 = 3
	dp.Field4 = 4
	dp.Field5 = 5
	dp.Field6 = 6
	dp.Field7 = 7
	dp.Field8 = 8
	for i := range dp.Data {
		dp.Data[i] = int64(i)
	}

	for b.Loop() {
		_ = dp.ProcessByPointer()
	}
}

// =============================================================================
// RETURN VALUE OPTIMIZATION BENCHMARKS
// =============================================================================

func BenchmarkReturnAddByValue(b *testing.B) {
	a := topics.LargeStruct{Field1: 1, Field2: 2}
	c := topics.LargeStruct{Field1: 3, Field2: 4}

	for b.Loop() {
		_ = topics.ReturnAddByValue(a, c)
	}
}

func BenchmarkReturnAddByPointer(b *testing.B) {
	a := topics.LargeStruct{Field1: 1, Field2: 2}
	c := topics.LargeStruct{Field1: 3, Field2: 4}

	for b.Loop() {
		_ = topics.ReturnAddByPointer(a, c)
	}
}

// =============================================================================
// SLICE ESCAPE ANALYSIS BENCHMARKS
// =============================================================================

func BenchmarkProcessSliceWithEscape(b *testing.B) {

	for b.Loop() {
		_ = topics.ProcessSliceWithEscape(1000)
	}
}

func BenchmarkProcessSliceNoEscape(b *testing.B) {

	for b.Loop() {
		_ = topics.ProcessSliceNoEscape(1000)
	}
}

// =============================================================================
// STACK VS HEAP ALLOCATION BENCHMARKS
// =============================================================================

func BenchmarkCreateLargeStructOnStack(b *testing.B) {

	for b.Loop() {
		_ = topics.CreateLargeStructOnStack()
	}
}

func BenchmarkCreateLargeStructOnHeap(b *testing.B) {

	for b.Loop() {
		_ = topics.CreateLargeStructOnHeap()
	}
}
