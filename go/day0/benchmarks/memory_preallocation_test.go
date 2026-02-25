package benchmarks

import (
	"testing"

	"day0/topics"
)

// =============================================================================
// SLICE BENCHMARKS
// =============================================================================

// BenchmarkDynamicSliceSmall benchmarks dynamic slice growth for small sizes.
func BenchmarkDynamicSliceSmall(b *testing.B) {
	for b.Loop() {
		_ = topics.DynamicSlice(100)
	}
}

// BenchmarkDynamicSliceMedium benchmarks dynamic slice growth for medium sizes.
func BenchmarkDynamicSliceMedium(b *testing.B) {
	for b.Loop() {
		_ = topics.DynamicSlice(1000)
	}
}

// BenchmarkDynamicSliceLarge benchmarks dynamic slice growth for large sizes.
func BenchmarkDynamicSliceLarge(b *testing.B) {
	for b.Loop() {
		_ = topics.DynamicSlice(10000)
	}
}

// BenchmarkPreallocatedSliceSmall benchmarks preallocated slice growth for small sizes.
func BenchmarkPreallocatedSliceSmall(b *testing.B) {
	for b.Loop() {
		_ = topics.PreallocatedSlice(100)
	}
}

// BenchmarkPreallocatedSliceMedium benchmarks preallocated slice growth for medium sizes.
func BenchmarkPreallocatedSliceMedium(b *testing.B) {
	for b.Loop() {
		_ = topics.PreallocatedSlice(1000)
	}
}

// BenchmarkPreallocatedSliceLarge benchmarks preallocated slice growth for large sizes.
func BenchmarkPreallocatedSliceLarge(b *testing.B) {
	for b.Loop() {
		_ = topics.PreallocatedSlice(10000)
	}
}

// BenchmarkPreallocatedSliceExactSmall benchmarks exact size preallocation for small sizes.
func BenchmarkPreallocatedSliceExactSmall(b *testing.B) {
	for b.Loop() {
		_ = topics.PreallocatedSliceExact(100)
	}
}

// BenchmarkPreallocatedSliceExactMedium benchmarks exact size preallocation for medium sizes.
func BenchmarkPreallocatedSliceExactMedium(b *testing.B) {
	for b.Loop() {
		_ = topics.PreallocatedSliceExact(1000)
	}
}

// BenchmarkPreallocatedSliceExactLarge benchmarks exact size preallocation for large sizes.
func BenchmarkPreallocatedSliceExactLarge(b *testing.B) {
	for b.Loop() {
		_ = topics.PreallocatedSliceExact(10000)
	}
}

// =============================================================================
// MAP BENCHMARKS
// =============================================================================

// BenchmarkDynamicMapSmall benchmarks dynamic map growth for small sizes.
func BenchmarkDynamicMapSmall(b *testing.B) {
	for b.Loop() {
		_ = topics.DynamicMap(100)
	}
}

// BenchmarkDynamicMapMedium benchmarks dynamic map growth for medium sizes.
func BenchmarkDynamicMapMedium(b *testing.B) {
	for b.Loop() {
		_ = topics.DynamicMap(1000)
	}
}

// BenchmarkDynamicMapLarge benchmarks dynamic map growth for large sizes.
func BenchmarkDynamicMapLarge(b *testing.B) {
	for b.Loop() {
		_ = topics.DynamicMap(10000)
	}
}

// BenchmarkPreallocatedMapSmall benchmarks preallocated map growth for small sizes.
func BenchmarkPreallocatedMapSmall(b *testing.B) {
	for b.Loop() {
		_ = topics.PreallocatedMap(100)
	}
}

// BenchmarkPreallocatedMapMedium benchmarks preallocated map growth for medium sizes.
func BenchmarkPreallocatedMapMedium(b *testing.B) {
	for b.Loop() {
		_ = topics.PreallocatedMap(1000)
	}
}

// BenchmarkPreallocatedMapLarge benchmarks preallocated map growth for large sizes.
func BenchmarkPreallocatedMapLarge(b *testing.B) {
	for b.Loop() {
		_ = topics.PreallocatedMap(10000)
	}
}
