package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"unsafe"
)

// =============================================================================
// COMPREHENSIVE GO OPTIMIZATION DEMO
// =============================================================================
// This program demonstrates 6 key Go optimization topics:
// 1. Struct Alignment - How field ordering affects memory usage and cache efficiency
// 2. Pass by Value vs Pointer - Copy cost vs indirection overhead
// 3. Receiver Types - Value vs pointer receivers for methods
// 4. Return Value Optimization (RVO) - How Go optimizes return-by-value
// 5. Slice Escape Analysis - When slices escape to heap vs stay on stack
// 6. Stack vs Heap - Where Go allocates data and performance implications
// =============================================================================

func main() {
	printHeader("GO PERFORMANCE OPTIMIZATION DEMONSTRATION")
	fmt.Println()
	fmt.Println("This demo covers 6 key optimization topics in Go:")
	fmt.Println("  1. Struct Alignment & Memory Padding")
	fmt.Println("  2. Pass by Value vs Pointer")
	fmt.Println("  3. Receiver Types (Value vs Pointer)")
	fmt.Println("  4. Return Value Optimization (RVO)")
	fmt.Println("  5. Slice Escape Analysis")
	fmt.Println("  6. Stack vs Heap Allocation")
	fmt.Println()

	// Run all demos
	runAllDemos()
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func printHeader(title string) {
	fmt.Println("================================================================================")
	fmt.Printf("%73s\n", title)
	fmt.Println("================================================================================")
}

func printSection(title string) {
	fmt.Println()
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("%73s\n", title)
	fmt.Println("--------------------------------------------------------------------------------")
}

func printSubsection(title string) {
	fmt.Println()
	fmt.Println("### " + title)
}

func formatBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.2f MB", float64(bytes)/1024/1024)
}

func runAllDemos() {
	// Demo 1: Struct Alignment
	demoStructAlignment()

	// Demo 2: Pass by Value vs Pointer
	demoPassByValue()

	// Demo 3: Receiver Types
	demoReceiverTypes()

	// Demo 4: Return Value Optimization
	demoReturnOptimization()

	// Demo 5: Slice Escape Analysis
	demoSliceEscape()

	// Demo 6: Stack vs Heap
	demoStackVsHeap()

	printHeader("DEMONSTRATION COMPLETE")
	fmt.Println()
	fmt.Println("Key Takeaways:")
	fmt.Println("  1. Order struct fields largest to smallest to minimize padding")
	fmt.Println("  2. Pass small structs by value, large structs by pointer")
	fmt.Println("  3. Use pointer receivers for large types or when mutation is needed")
	fmt.Println("  4. Return by value when possible - let RVO handle optimization")
	fmt.Println("  5. Keep data local to avoid heap escape and GC pressure")
	fmt.Println("  6. Stack allocation is faster but data must not outlive function")
}

// =============================================================================
// DEMO 1: STRUCT ALIGNMENT
// =============================================================================

func demoStructAlignment() {
	printHeader("DEMO 1: STRUCT ALIGNMENT & MEMORY PADDING")

	fmt.Println()
	fmt.Println("WHAT IS STRUCT ALIGNMENT?")
	fmt.Println("-------------------------")
	fmt.Println("Go aligns fields to their natural boundaries for CPU efficiency.")
	fmt.Println("When fields are ordered poorly, Go adds PADDING bytes between fields.")
	fmt.Println("This can waste 50%+ of memory!")
	fmt.Println()

	// Show struct sizes
	printSubsection("Struct Memory Footprint")
	sizes := GetStructSizes()
	fmt.Println("+-------------------+--------+------------------+------------+")
	fmt.Println("| Struct Name       |  Size  | Theoretical Data  |  Padding   |")
	fmt.Println("+-------------------+--------+------------------+------------+")

	for name, size := range sizes {
		var theoretical, padding int
		switch name {
		case "UnalignedStruct", "PoorlyPaddedStruct":
			theoretical = 24
			padding = size - theoretical
		case "AlignedStruct":
			theoretical = 24
			padding = size - theoretical
		case "MixedTypesAligned":
			theoretical = 38
			padding = size - theoretical
		case "MixedTypesUnaligned":
			theoretical = 38
			padding = size - theoretical
		}
		paddingPct := float64(padding) / float64(size) * 100
		fmt.Printf("| %-17s | %6d | %16d | %6d (%d%%) |\n", name, size, theoretical, padding, int(paddingPct))
	}
	fmt.Println("+-------------------+--------+------------------+------------+")

	// Calculate savings
	printSubsection("Memory Savings with Proper Alignment")
	unalignedSize := sizes["UnalignedStruct"]
	alignedSize := sizes["AlignedStruct"]
	savings := unalignedSize - alignedSize
	savingsPercent := float64(savings) / float64(unalignedSize) * 100

	fmt.Printf("Unaligned: %d bytes\n", unalignedSize)
	fmt.Printf("Aligned:   %d bytes\n", alignedSize)
	fmt.Printf("Savings:   %d bytes (%.1f%% reduction)\n", savings, savingsPercent)
	fmt.Printf("\nWith 1 million elements:\n")
	fmt.Printf("  Unaligned: %s\n", formatBytes(int64(unalignedSize*1000000)))
	fmt.Printf("  Aligned:   %s\n", formatBytes(int64(alignedSize*1000000)))
	fmt.Printf("  Saved:     %s\n", formatBytes(int64(savings*1000000)))

	// Cache effects
	printSubsection("Cache Line Efficiency")
	cacheLine := 64
	elemsPerLineUnaligned := cacheLine / unalignedSize
	elemsPerLineAligned := cacheLine / alignedSize

	fmt.Printf("Cache line size: %d bytes\n", cacheLine)
	fmt.Printf("Unaligned: %d elements per cache line\n", elemsPerLineUnaligned)
	fmt.Printf("Aligned:   %d elements per cache line\n", elemsPerLineAligned)
	fmt.Printf("Efficiency improvement: %.1fx\n", float64(elemsPerLineAligned)/float64(elemsPerLineUnaligned))

	// Run benchmarks
	printSubsection("Performance Benchmarks")
	runBenchmarks("Alignment")

	// Best practices
	printSubsection("Best Practices")
	fmt.Println("✓ Order struct fields by size: largest first, smallest last")
	fmt.Println("  - int64, float64, pointers (8 bytes)")
	fmt.Println("  - int32, float32 (4 bytes)")
	fmt.Println("  - int16 (2 bytes)")
	fmt.Println("  - int8, bool (1 byte)")
	fmt.Println("✓ Use unsafe.Sizeof() to check struct sizes")
	fmt.Println("✓ Consider 'go vet -structtag' for tagged field ordering")
}

// =============================================================================
// DEMO 2: PASS BY VALUE VS POINTER
// =============================================================================

func demoPassByValue() {
	printHeader("DEMO 2: PASS BY VALUE VS POINTER")

	fmt.Println()
	fmt.Println("PASS BY VALUE:")
	fmt.Println("  - Copies the entire struct onto the stack")
	fmt.Println("  - No heap allocation needed (stack-to-stack copy is fast)")
	fmt.Println("  - Good for small structs (< 2 words)")
	fmt.Println()
	fmt.Println("PASS BY POINTER:")
	fmt.Println("  - Only copies the 8-byte pointer")
	fmt.Println("  - Avoids copy overhead for large structs")
	fmt.Println("  - Adds slight indirection cost")
	fmt.Println()

	// Show LargeStruct size
	structSize := int(unsafe.Sizeof(LargeStruct{}))
	printSubsection("LargeStruct Size")
	fmt.Printf("LargeStruct: %d bytes (%s)\n", structSize, formatBytes(int64(structSize)))
	fmt.Println("This is large enough to show significant copy overhead!")

	// Run benchmarks
	printSubsection("Performance Benchmarks")
	runBenchmarks("PassByValue")

	// When to use each
	printSubsection("When to Use Each")
	fmt.Println("PASS BY VALUE when:")
	fmt.Println("  - Struct is small (< 16 bytes / 2 words)")
	fmt.Println("  - You need thread-safety (no aliasing)")
	fmt.Println("  - Data is read-only")
	fmt.Println()
	fmt.Println("PASS BY POINTER when:")
	fmt.Println("  - Struct is large (> 100 bytes)")
	fmt.Println("  - You need to modify the original")
	fmt.Println("  - Performance is critical in hot paths")
}

// =============================================================================
// DEMO 3: RECEIVER TYPES
// =============================================================================

func demoReceiverTypes() {
	printHeader("DEMO 3: RECEIVER TYPES (VALUE VS POINTER)")

	fmt.Println()
	fmt.Println("VALUE RECEIVER:")
	fmt.Println("  - Go makes a COPY of the struct")
	fmt.Println("  - Changes don't affect the original")
	fmt.Println("  - Good for small, read-only operations")
	fmt.Println()
	fmt.Println("POINTER RECEIVER:")
	fmt.Println("  - Go passes a pointer to the original")
	fmt.Println("  - Changes persist")
	fmt.Println("  - No copy overhead - efficient for large types")
	fmt.Println()

	// Show sizes
	printSubsection("Receiver Type Impact")
	counterSize := int(unsafe.Sizeof(Counter{}))
	processorSize := int(unsafe.Sizeof(DataProcessor{}))
	fmt.Printf("Counter: %d bytes (copy cost negligible)\n", counterSize)
	fmt.Printf("DataProcessor: %d bytes (copy cost significant!)\n", processorSize)

	// Run benchmarks - small struct
	printSubsection("Performance Benchmarks - Small Struct (Counter)")
	runBenchmarks("ReceiverSmall")

	// Run benchmarks - large struct
	printSubsection("Performance Benchmarks - Large Struct (DataProcessor)")
	runBenchmarks("ReceiverLarge")

	// Guidelines
	printSubsection("Guidelines")
	fmt.Println("Use VALUE RECEIVER when:")
	fmt.Println("  - Struct is small (< 16 bytes)")
	fmt.Println("  - You don't need to modify the original")
	fmt.Println("  - Thread-safety is important (no aliasing)")
	fmt.Println()
	fmt.Println("Use POINTER RECEIVER when:")
	fmt.Println("  - Struct is large (> 100 bytes)")
	fmt.Println("  - You need to modify the original")
	fmt.Println("  - Method must satisfy an interface that requires pointers")
}

// =============================================================================
// DEMO 4: RETURN VALUE OPTIMIZATION
// =============================================================================

func demoReturnOptimization() {
	printHeader("DEMO 4: RETURN VALUE OPTIMIZATION (RVO)")

	fmt.Println()
	fmt.Println("WHAT IS RVO?")
	fmt.Println("Return Value Optimization allows Go to:")
	fmt.Println("  1. Allocate return space in the CALLER (not in the function)")
	fmt.Println("  2. Build the return value directly where it's needed")
	fmt.Println("  3. Zero copies - zero allocations!")
	fmt.Println()
	fmt.Println("RETURN BY VALUE (with RVO):")
	fmt.Println("  ✓ No heap allocation")
	fmt.Println("  ✓ No copy overhead")
	fmt.Println("  ✓ Clean code, great performance")
	fmt.Println()
	fmt.Println("RETURN POINTER (heap escape):")
	fmt.Println("  ✗ Heap allocation required")
	fmt.Println("  ✗ Garbage collector pressure")
	fmt.Println("  ✗ Potential cache misses")

	// Run benchmarks
	printSubsection("Performance Benchmarks")
	runBenchmarks("Return")

	// Key takeaway
	printSubsection("Key Takeaway")
	fmt.Println("✓ Let the compiler help you - return by value when possible!")
	fmt.Println("✓ RVO is one of Go's best optimizations")
	fmt.Println("✓ Avoid returning pointers to local variables unless necessary")
}

// =============================================================================
// DEMO 5: SLICE ESCAPE ANALYSIS
// =============================================================================

func demoSliceEscape() {
	printHeader("DEMO 5: SLICE ESCAPE ANALYSIS")

	fmt.Println()
	fmt.Println("WHAT IS ESCAPE ANALYSIS?")
	fmt.Println("Go determines at compile time whether data can stay on the stack")
	fmt.Println("or must be moved to the heap (where it 'escapes' the function).")
	fmt.Println()
	fmt.Println("STAY ON STACK (no escape):")
	fmt.Println("  ✓ Fast allocation (just move stack pointer)")
	fmt.Println("  ✓ No GC pressure")
	fmt.Println("  ✓ Automatic cleanup")
	fmt.Println()
	fmt.Println("ESCAPE TO HEAP:")
	fmt.Println("  ✗ Requires allocation")
	fmt.Println("  ✗ GC must track and collect")
	fmt.Println("  ✗ Slower than stack")

	// Demonstrate escape with timing
	printSubsection("Escape Demonstration (Timing Test)")

	// Run timing test to show difference
	for _, size := range []int{100, 1000, 10000} {
		start := time.Now()
		var escapeSum, noEscapeSum int

		// Warm up
		for i := range 1000 {
			if i%2 == 0 {
				escapeSum += processSliceWithEscape(size)
			} else {
				noEscapeSum += processSliceNoEscape(size)
			}
		}
		elapsed := time.Since(start)

		_ = escapeSum
		_ = noEscapeSum

		fmt.Printf("Slice size %5d: %v for 1000 iterations\n", size, elapsed)
	}

	// Run benchmarks
	printSubsection("Performance Benchmarks")
	runBenchmarks("SliceEscape")

	// Guidelines
	printSubsection("Guidelines")
	fmt.Println("✓ Keep data local to avoid escape")
	fmt.Println("✓ Avoid assigning to global variables")
	fmt.Println("✓ Don't return slices unnecessarily")
	fmt.Println("✓ Use //go:noinline to prevent optimization if testing")
}

// =============================================================================
// DEMO 6: STACK VS HEAP
// =============================================================================

func demoStackVsHeap() {
	printHeader("DEMO 6: STACK VS HEAP ALLOCATION")

	fmt.Println()
	fmt.Println("THE STACK:")
	fmt.Println("  - Fast allocation (just move a pointer)")
	fmt.Println("  - Automatic cleanup (no GC needed)")
	fmt.Println("  - Limited size (~MB)")
	fmt.Println("  - Great for short-lived data")
	fmt.Println()
	fmt.Println("THE HEAP:")
	fmt.Println("  - Slower allocation (requires finding free space)")
	fmt.Println("  - Requires garbage collection")
	fmt.Println("  - Much larger (~GB)")
	fmt.Println("  - For data that outlives its function")

	// Get current runtime info
	printSubsection("Runtime Information")
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("NumCPU: %d\n", runtime.NumCPU())
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("GC cycles: %d\n", m.NumGC)

	// Run benchmarks
	printSubsection("Performance Benchmarks")
	runBenchmarks("StackHeap")

	// Key insights
	printSubsection("Key Insights")
	fmt.Println("✓ Stack allocation is ~10-100x faster than heap")
	fmt.Println("✓ Small allocations may not trigger GC at all")
	fmt.Println("✓ Escape analysis happens at compile time")
	fmt.Println("✓ Large objects (> 64KB) go directly to heap")
	fmt.Println("✓ Use pprof to identify heap allocations: go tool pprof")
}

// =============================================================================
// BENCHMARK RUNNER
// =============================================================================

func runBenchmarks(category string) {
	fmt.Println()

	// Define benchmark patterns for each category
	benchmarks := getBenchmarksForCategory(category)

	if len(benchmarks) == 0 {
		fmt.Println("No benchmarks available for this category")
		return
	}

	fmt.Printf("%-45s | %12s | %10s\n", "Benchmark", "Time/op", "Allocations")
	fmt.Println(strings.Repeat("-", 75))

	// Run benchmarks using go test
	for _, bm := range benchmarks {
		result := runGoTestBenchmark(bm)
		if result != "" {
			fmt.Println(result)
		}
	}
}

func getBenchmarksForCategory(category string) []string {
	switch category {
	case "Alignment":
		return []string{
			"BenchmarkProcessUnaligned",
			"BenchmarkProcessAligned",
			"BenchmarkProcessUnalignedPtr",
			"BenchmarkProcessAlignedPtr",
			"BenchmarkMixedTypesAligned",
			"BenchmarkMixedTypesUnaligned",
		}
	case "PassByValue":
		return []string{
			"BenchmarkAddByValue",
			"BenchmarkAddByPointer",
		}
	case "ReceiverSmall":
		return []string{
			"BenchmarkIncrementByValue",
			"BenchmarkIncrementByPointer",
		}
	case "ReceiverLarge":
		return []string{
			"BenchmarkProcessByValue",
			"BenchmarkProcessByPointer",
		}
	case "Return":
		return []string{
			"BenchmarkReturnAddByValue",
			"BenchmarkReturnAddByPointer",
		}
	case "SliceEscape":
		return []string{
			"BenchmarkProcessSliceWithEscape",
			"BenchmarkProcessSliceNoEscape",
		}
	case "StackHeap":
		return []string{
			"BenchmarkCreateLargeStructOnStack",
			"BenchmarkCreateLargeStructOnHeap",
		}
	}
	return []string{}
}

func runGoTestBenchmark(benchmarkName string) string {
	// Run the benchmark using go test
	cmd := exec.Command("go", "test", "-bench="+benchmarkName, "-benchmem", "-run=^$", "-count=1", ".")
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Try running with different approach
		cmd := exec.Command("go", "test", "-bench="+benchmarkName, "-benchmem", "-run=^$", ".")
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Sprintf("%-45s | Error: %v", benchmarkName, err)
		}
	}

	// Parse the output
	lines := strings.SplitSeq(string(output), "\n")
	for line := range lines {
		if strings.HasPrefix(line, benchmarkName) {
			// Parse benchmark output: BenchmarkName	N	ns/op	A	bytes/op	B	allocs/op
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				timeOp := "N/A"
				allocs := "N/A"

				for i, part := range parts {
					if part == "ns/op" && i > 0 {
						timeOp = parts[i-1]
					}
					if part == "allocs/op" && i > 0 {
						allocs = parts[i-1]
					}
				}

				return fmt.Sprintf("%-45s | %12s | %10s", benchmarkName, timeOp, allocs)
			}
		}
	}

	return fmt.Sprintf("%-45s | (no output)", benchmarkName)
}
