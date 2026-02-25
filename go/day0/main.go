package main

import (
	"fmt"
	"os/exec"
	"strings"
	"unsafe"

	"day0/topics"
)

// =============================================================================
// COMPREHENSIVE GO OPTIMIZATION DEMO
// =============================================================================
// This program demonstrates 11 key Go optimization topics:
// (From goperf.dev/01-common-patterns/)
//
// ORIGINAL 6 TOPICS:
// 1. Struct Alignment - How field ordering affects memory usage and cache efficiency
// 2. Pass by Value vs Pointer - Copy cost vs indirection overhead
// 3. Receiver Types - Value vs pointer receivers for methods
// 4. Return Value Optimization (RVO) - How Go optimizes return-by-value
// 5. Slice Escape Analysis - When slices escape to heap vs stay on stack
// 6. Stack vs Heap - Where Go allocates data and performance implications
// 7. Object Pooling - Reusing objects to reduce GC pressure
// 8. Batching Operations - Reducing overhead by grouping operations
// 9. Immutable Data Sharing - Safe concurrent access without locks
// 10. Lazy Initialization - Deferring expensive operations until needed
// 11. Memory Preallocation - Preallocating slices and maps for performance
// =============================================================================

func main() {
	printHeader("GO PERFORMANCE OPTIMIZATION DEMONSTRATION")
	fmt.Println()
	fmt.Println("This demo covers 11 key optimization topics in Go:")
	fmt.Println()
	fmt.Println("ORIGINAL TOPICS:")
	fmt.Println("  1. Struct Alignment & Memory Padding")
	fmt.Println("  2. Pass by Value vs Pointer")
	fmt.Println("  3. Receiver Types (Value vs Pointer)")
	fmt.Println("  4. Return Value Optimization (RVO)")
	fmt.Println("  5. Slice Escape Analysis")
	fmt.Println("  6. Stack vs Heap Allocation")
	fmt.Println("  7. Object Pooling")
	fmt.Println("  8. Batching Operations")
	fmt.Println("  9. Immutable Data Sharing")
	fmt.Println(" 10. Lazy Initialization")
	fmt.Println(" 11. Memory Preallocation")
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

	// Demo 7: Object Pooling
	demoObjectPooling()

	// Demo 8: Batching Operations
	demoBatchingOperations()

	// Demo 9: Immutable Data Sharing
	demoImmutableDataSharing()

	// Demo 10: Lazy Initialization
	demoLazyInitialization()

	// Demo 11: Memory Preallocation
	demoMemoryPreallocation()

	printHeader("DEMONSTRATION COMPLETE")
	fmt.Println()
	fmt.Println("Key Takeaways:")
	fmt.Println("  1. Order struct fields largest to smallest to minimize padding")
	fmt.Println("  2. Pass small structs by value, large structs by pointer")
	fmt.Println("  3. Use pointer receivers for large types or when mutation is needed")
	fmt.Println("  4. Return by value when possible - let RVO handle optimization")
	fmt.Println("  5. Keep data local to avoid heap escape and GC pressure")
	fmt.Println("  6. Stack allocation is faster but data must not outlive function")
	fmt.Println("  7. Object pooling reduces GC pressure for high-frequency allocations")
	fmt.Println("  8. Batching reduces overhead for I/O-bound operations")
	fmt.Println("  9. Immutable data enables safe concurrent access without locks")
	fmt.Println(" 10. Lazy initialization defers expensive operations until needed")
	fmt.Println(" 11. Preallocate slices and maps when size is known")
}

// =============================================================================
// DEMO 1: STRUCT ALIGNMENT
// =============================================================================

func demoStructAlignment() {
	topics.RunAlignmentDemo()
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
	structSize := int(unsafe.Sizeof(topics.LargeStruct{}))
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
	counterSize := int(unsafe.Sizeof(topics.Counter{}))
	processorSize := int(unsafe.Sizeof(topics.DataProcessor{}))
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
	topics.RunReturnOptimizationDemo()
}

// =============================================================================
// DEMO 5: SLICE ESCAPE ANALYSIS
// =============================================================================

func demoSliceEscape() {
	topics.RunSliceEscapeDemo()
}

// =============================================================================
// DEMO 6: STACK VS HEAP
// =============================================================================

func demoStackVsHeap() {
	topics.RunStackVsHeapDemo()
}

// =============================================================================
// DEMO 7: OBJECT POOLING
// =============================================================================

func demoObjectPooling() {
	topics.RunPoolingDemo()
}

// =============================================================================
// DEMO 8: BATCHING OPERATIONS
// =============================================================================

func demoBatchingOperations() {
	topics.RunBatchingDemo()
}

// =============================================================================
// DEMO 9: IMMUTABLE DATA SHARING
// =============================================================================

func demoImmutableDataSharing() {
	topics.RunImmutableDemo()
}

// =============================================================================
// DEMO 10: LAZY INITIALIZATION
// =============================================================================

func demoLazyInitialization() {
	topics.RunLazyInitDemo()
}

// =============================================================================
// DEMO 11: MEMORY PREALLOCATION
// =============================================================================

func demoMemoryPreallocation() {
	topics.RunMemoryPreallocationDemo()
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
	// Run the benchmark using go test from the benchmarks directory
	cmd := exec.Command("go", "test", "-bench="+benchmarkName, "-benchmem", "-run=^$", "-count=1", "./benchmarks")
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Try running with different approach
		cmd := exec.Command("go", "test", "-bench="+benchmarkName, "-benchmem", "-run=^$", "./benchmarks")
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
