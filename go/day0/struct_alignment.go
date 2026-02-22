package main

import (
	"fmt"
	"math/rand"
	"time"
	"unsafe"
)

// =============================================================================
// STRUCT ALIGNMENT AND DATA PADDING
// =============================================================================
//
// This file demonstrates how Go aligns struct fields and how padding
// affects memory usage.
//
// PRINCIPLE: Go aligns fields to their natural boundaries for performance.
// When fields are ordered poorly, Go adds PADDING bytes between fields.

// UnalignedStruct demonstrates BAD field ordering causing PADDING WASTE.
//
// MEMORY LAYOUT (64-bit system):
// Field1 (int8):  1 byte + 7 bytes padding
// Field2 (int64): 8 bytes (aligned to 8-byte boundary)
// Field3 (int8):  1 byte + 7 bytes padding
// Field4 (int64): 8 bytes
// Field5 (int8):  1 byte + 7 bytes padding
// Field6 (int64): 8 bytes
//
// TOTAL: 48 bytes (24 bytes data + 24 bytes padding!)
// WASTED: 66% of memory is padding!
type UnalignedStruct struct {
	Field1 int8
	Field2 int64
	Field3 int8
	Field4 int64
	Field5 int8
	Field6 int64
}

// AlignedStruct demonstrates GOOD field ordering to MINIMIZE PADDING.
//
// PRINCIPLE: Place large fields first, then smaller fields.
// This allows Go to pack smaller fields into the padding gaps.
//
// MEMORY LAYOUT (64-bit system):
// Field2 (int64): 8 bytes (offset 0)
// Field4 (int64): 8 bytes (offset 8)
// Field6 (int64): 8 bytes (offset 16)
// Field1 (int8):  1 byte (offset 24)
// Field3 (int8):  1 byte (offset 25)
// Field5 (int8):  1 byte (offset 26)
//                6 bytes padding to align struct to 8 bytes
//
// TOTAL: 32 bytes (24 bytes data + 8 bytes padding)
// SAVINGS: 16 bytes less than UnalignedStruct! (33% reduction)
type AlignedStruct struct {
	Field2 int64
	Field4 int64
	Field6 int64
	Field1 int8
	Field3 int8
	Field5 int8
}

// PoorlyPaddedStruct demonstrates WORST CASE with many padding bytes.
//
// This has the same fields but ordered worst possible way:
// All 1-byte fields first, then all 8-byte fields.
//
// MEMORY LAYOUT (64-bit system):
// Field1 (int8):  1 byte + 7 bytes padding
// Field3 (int8):  1 byte + 7 bytes padding
// Field5 (int8):  1 byte + 7 bytes padding
// Field2 (int64): 8 bytes
// Field4 (int64): 8 bytes
// Field6 (int64): 8 bytes
//
// TOTAL: 48 bytes (same as UnalignedStruct - both are bad!)
type PoorlyPaddedStruct struct {
	Field1 int8
	Field3 int8
	Field5 int8
	Field2 int64
	Field4 int64
	Field6 int64
}

// MixedTypesAligned shows alignment with various field types.
//
// PRINCIPLE: Order fields by size (largest to smallest):
// 1. int64, float64, pointers (8 bytes on 64-bit)
// 2. int32, float32 (4 bytes)
// 3. int16 (2 bytes)
// 4. int8, bool (1 byte)
type MixedTypesAligned struct {
	// 8-byte fields first (3 fields = 24 bytes)
	Pointer *int64
	Float   float64
	Counter int64

	// 4-byte fields (2 fields = 8 bytes)
	Count  int32
	Flag   float32

	// 2-byte fields (2 fields = 4 bytes)
	Short  int16
	Char   int16

	// 1-byte fields (2 fields = 2 bytes)
	Byte   byte
	Bool   bool

	// 6 bytes padding to align struct
}

// MixedTypesUnaligned shows the same fields in poor order.
type MixedTypesUnaligned struct {
	// 1-byte fields first (causes padding!)
	Bool   bool
	Byte   byte

	// 2-byte fields (more padding)
	Char   int16
	Short  int16

	// 4-byte fields (more padding)
	Flag   float32
	Count  int32

	// 8-byte fields (finally!)
	Counter int64
	Float   float64
	Pointer *int64
}

// GetStructSizes demonstrates how to check struct sizes at runtime.
// Use unsafe.Sizeof() to see the actual memory footprint.
func GetStructSizes() map[string]int {
	return map[string]int{
		"UnalignedStruct":      int(unsafe.Sizeof(UnalignedStruct{})),
		"AlignedStruct":        int(unsafe.Sizeof(AlignedStruct{})),
		"PoorlyPaddedStruct":   int(unsafe.Sizeof(PoorlyPaddedStruct{})),
		"MixedTypesAligned":    int(unsafe.Sizeof(MixedTypesAligned{})),
		"MixedTypesUnaligned":  int(unsafe.Sizeof(MixedTypesUnaligned{})),
	}
}

// ProcessUnaligned demonstrates processing with poor alignment.
func ProcessUnaligned(s UnalignedStruct) int64 {
	return s.Field2 + s.Field4 + s.Field6
}

// ProcessAligned demonstrates processing with good alignment.
func ProcessAligned(s AlignedStruct) int64 {
	return s.Field2 + s.Field4 + s.Field6
}

// ProcessUnalignedPtr demonstrates processing with pointer (no copy).
func ProcessUnalignedPtr(s *UnalignedStruct) int64 {
	return s.Field2 + s.Field4 + s.Field6
}

// ProcessAlignedPtr demonstrates processing with pointer (no copy).
func ProcessAlignedPtr(s *AlignedStruct) int64 {
	return s.Field2 + s.Field4 + s.Field6
}

// UnalignedInts shows worst case: alternating small and large ints.
type UnalignedInts struct {
	A int8  // 1 byte + 7 padding
	B int64 // 8 bytes
	C int8  // 1 byte + 7 padding
	D int64 // 8 bytes
	E int8  // 1 byte + 7 padding
	F int64 // 8 bytes
	// Total: 48 bytes (24 data + 24 padding)
}

// AlignedInts shows best case: large ints first, then small.
type AlignedInts struct {
	B int64 // 8 bytes
	D int64 // 8 bytes
	F int64 // 8 bytes
	A int8  // 1 byte
	C int8  // 1 byte
	E int8  // 1 byte
	// 5 bytes padding at end
	// Total: 32 bytes (24 data + 8 padding)
}

// =============================================================================
// ALIGNMENT DEMONSTRATION
// =============================================================================

// Benchmark constants - large enough to exceed L1/L2 cache
const (
	BenchSliceSize   = 1000000 // 1M elements to exceed L1/L2 cache
	BenchCacheLine   = 64      // Typical cache line size
	BenchL1CacheSize = 32 * 1024
)

// createUnalignedSliceForDemo creates a large slice of UnalignedStruct
func createUnalignedSliceForDemo(size int) []UnalignedStruct {
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

// createAlignedSliceForDemo creates a large slice of AlignedStruct
func createAlignedSliceForDemo(size int) []AlignedStruct {
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

func init() {
	rand.Seed(42)
}

// RunAlignmentDemo demonstrates the performance impact of struct alignment
func RunAlignmentDemo() {
	fmt.Println("================================================================================")
	fmt.Println("                    STRUCT ALIGNMENT PERFORMANCE ANALYSIS                     ")
	fmt.Println("================================================================================")
	fmt.Println()

	// Print struct sizes
	fmt.Println("=== STRUCT SIZES ===")
	sizes := GetStructSizes()
	for name, size := range sizes {
		fmt.Printf("%-24s: %3d bytes\n", name, size)
	}
	fmt.Println()

	// Calculate memory savings
	unalignedSize := sizes["UnalignedStruct"]
	alignedSize := sizes["AlignedStruct"]
	savings := unalignedSize - alignedSize
	savingsPercent := float64(savings) / float64(unalignedSize) * 100

	fmt.Println("=== MEMORY SAVINGS ===")
	fmt.Printf("UnalignedStruct: %d bytes\n", unalignedSize)
	fmt.Printf("AlignedStruct:   %d bytes\n", alignedSize)
	fmt.Printf("Savings:         %d bytes (%.1f%% reduction)\n", savings, savingsPercent)
	fmt.Println()

	// Demonstrate with actual data
	fmt.Println("=== PERFORMANCE DEMONSTRATION ===")
	fmt.Printf("Slice size: %d elements\n", BenchSliceSize)
	fmt.Printf("Unaligned memory: %d bytes (%.2f MB)\n",
		BenchSliceSize*unalignedSize,
		float64(BenchSliceSize*unalignedSize)/1024/1024)
	fmt.Printf("Aligned memory:   %d bytes (%.2f MB)\n",
		BenchSliceSize*alignedSize,
		float64(BenchSliceSize*alignedSize)/1024/1024)
	fmt.Printf("Memory saved:     %d bytes (%.2f MB)\n",
		BenchSliceSize*savings,
		float64(BenchSliceSize*savings)/1024/1024)
	fmt.Println()

	// Estimate cache effects
	fmt.Println("=== CACHE EFFECTS ===")
	elementsPerCacheLineUnaligned := BenchCacheLine / unalignedSize
	elementsPerCacheLineAligned := BenchCacheLine / alignedSize

	fmt.Printf("Cache line size: %d bytes\n", BenchCacheLine)
	fmt.Printf("Unaligned: %d elements per cache line\n", elementsPerCacheLineUnaligned)
	fmt.Printf("Aligned:   %d elements per cache line\n", elementsPerCacheLineAligned)
	fmt.Printf("Efficiency improvement: %.1fx\n",
		float64(elementsPerCacheLineAligned)/float64(elementsPerCacheLineUnaligned))
	fmt.Println()

	// Estimate L1 cache capacity
	elementsInL1Unaligned := BenchL1CacheSize / unalignedSize
	elementsInL1Aligned := BenchL1CacheSize / alignedSize

	fmt.Printf("L1 Cache: %d KB\n", BenchL1CacheSize/1024)
	fmt.Printf("Unaligned: ~%d elements fit in L1\n", elementsInL1Unaligned)
	fmt.Printf("Aligned:   ~%d elements fit in L1\n", elementsInL1Aligned)
	fmt.Printf("Cache capacity improvement: %.1fx\n",
		float64(elementsInL1Aligned)/float64(elementsInL1Unaligned))
	fmt.Println()

	// Print benchmark hints
	fmt.Println("=== RUN BENCHMARKS ===")
	fmt.Println("To run benchmarks, execute:")
	fmt.Println("  go test -bench=. -benchmem -run=^$ .")
	fmt.Println()
	fmt.Println("Benchmark categories:")
	fmt.Println("  - Sequential*   : Sequential access patterns")
	fmt.Println("  - Random*       : Random access patterns")
	fmt.Println("  - Strided*      : Strided access patterns")
	fmt.Println("  - Value*        : Pass by value")
	fmt.Println("  - Pointer*      : Pass by pointer")
	fmt.Println("  - MixedTypes*   : Various field types")
	fmt.Println()
	fmt.Println("================================================================================")
	fmt.Println("Note: Actual performance impact depends on:")
	fmt.Println("  - CPU cache size and architecture")
	fmt.Println("  - Memory bandwidth")
	fmt.Println("  - Access patterns (sequential vs random)")
	fmt.Println("  - Workload characteristics")
	fmt.Println("================================================================================")

	// Run quick timing demonstration
	fmt.Println()
	fmt.Println("=== QUICK TIMING TEST ===")
	runQuickTimingTestDemo()
}

func runQuickTimingTestDemo() {
	// Create test data
	unalignedData := createUnalignedSliceForDemo(100000)
	alignedData := createAlignedSliceForDemo(100000)

	// Test sequential access
	start := time.Now()
	var sum1 int64
	for _, v := range unalignedData {
		sum1 += v.Field2 + v.Field4 + v.Field6
	}
	unalignedTime := time.Since(start)

	start = time.Now()
	var sum2 int64
	for _, v := range alignedData {
		sum2 += v.Field2 + v.Field4 + v.Field6
	}
	alignedTime := time.Since(start)

	_ = sum1
	_ = sum2

	fmt.Printf("Unaligned sequential: %v\n", unalignedTime)
	fmt.Printf("Aligned sequential:   %v\n", alignedTime)
	if alignedTime < unalignedTime {
		fmt.Printf("Speedup:              %.2fx\n",
			float64(unalignedTime.Nanoseconds())/float64(alignedTime.Nanoseconds()))
	} else {
		fmt.Printf("Ratio:                %.2fx\n",
			float64(alignedTime.Nanoseconds())/float64(unalignedTime.Nanoseconds()))
	}
}
