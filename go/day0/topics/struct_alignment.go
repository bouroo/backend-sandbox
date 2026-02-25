// Package topics provides Go performance optimization demonstrations.
package topics

import (
	"fmt"
	"math/rand"
	"time"
	"unsafe"
)

// =============================================================================
// STRUCT ALIGNMENT AND DATA PADDING
// =============================================================================

// UnalignedStruct demonstrates BAD field ordering causing PADDING WASTE.
type UnalignedStruct struct {
	Field1 int8
	Field2 int64
	Field3 int8
	Field4 int64
	Field5 int8
	Field6 int64
}

// AlignedStruct demonstrates GOOD field ordering to MINIMIZE PADDING.
type AlignedStruct struct {
	Field2 int64
	Field4 int64
	Field6 int64
	Field1 int8
	Field3 int8
	Field5 int8
}

// PoorlyPaddedStruct demonstrates WORST CASE with many padding bytes.
type PoorlyPaddedStruct struct {
	Field1 int8
	Field3 int8
	Field5 int8
	Field2 int64
	Field4 int64
	Field6 int64
}

// MixedTypesAligned shows alignment with various field types.
type MixedTypesAligned struct {
	Pointer *int64
	Float   float64
	Counter int64
	Count   int32
	Flag    float32
	Short   int16
	Char    int16
	Byte    byte
	Bool    bool
}

// MixedTypesUnaligned shows the same fields in poor order.
type MixedTypesUnaligned struct {
	Bool    bool
	Byte    byte
	Char    int16
	Short   int16
	Flag    float32
	Count   int32
	Counter int64
	Float   float64
	Pointer *int64
}

// GetStructSizes demonstrates how to check struct sizes at runtime.
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

const (
	BenchSliceSize   = 1000000
	BenchCacheLine   = 64
	BenchL1CacheSize = 32 * 1024
)

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

	fmt.Println("=== STRUCT SIZES ===")
	sizes := GetStructSizes()
	for name, size := range sizes {
		fmt.Printf("%-24s: %3d bytes\n", name, size)
	}
	fmt.Println()

	unalignedSize := sizes["UnalignedStruct"]
	alignedSize := sizes["AlignedStruct"]
	savings := unalignedSize - alignedSize
	savingsPercent := float64(savings) / float64(unalignedSize) * 100

	fmt.Println("=== MEMORY SAVINGS ===")
	fmt.Printf("UnalignedStruct: %d bytes\n", unalignedSize)
	fmt.Printf("AlignedStruct:   %d bytes\n", alignedSize)
	fmt.Printf("Savings:         %d bytes (%.1f%% reduction)\n", savings, savingsPercent)
	fmt.Println()

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

	fmt.Println("=== CACHE EFFECTS ===")
	elementsPerCacheLineUnaligned := BenchCacheLine / unalignedSize
	elementsPerCacheLineAligned := BenchCacheLine / alignedSize

	fmt.Printf("Cache line size: %d bytes\n", BenchCacheLine)
	fmt.Printf("Unaligned: %d elements per cache line\n", elementsPerCacheLineUnaligned)
	fmt.Printf("Aligned:   %d elements per cache line\n", elementsPerCacheLineAligned)
	fmt.Printf("Efficiency improvement: %.1fx\n",
		float64(elementsPerCacheLineAligned)/float64(elementsPerCacheLineUnaligned))
	fmt.Println()

	elementsInL1Unaligned := BenchL1CacheSize / unalignedSize
	elementsInL1Aligned := BenchL1CacheSize / alignedSize

	fmt.Printf("L1 Cache: %d KB\n", BenchL1CacheSize/1024)
	fmt.Printf("Unaligned: ~%d elements fit in L1\n", elementsInL1Unaligned)
	fmt.Printf("Aligned:   ~%d elements fit in L1\n", elementsInL1Aligned)
	fmt.Printf("Cache capacity improvement: %.1fx\n",
		float64(elementsInL1Aligned)/float64(elementsInL1Unaligned))
	fmt.Println()

	fmt.Println("=== RUN BENCHMARKS ===")
	fmt.Println("To run benchmarks, execute:")
	fmt.Println("  go test -bench=. -benchmem -run=^$ .")
	fmt.Println()

	fmt.Println("================================================================================")

	fmt.Println()
	fmt.Println("=== QUICK TIMING TEST ===")
	
	unalignedData := createUnalignedSliceForDemo(100000)
	alignedData := createAlignedSliceForDemo(100000)

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
	}
}
