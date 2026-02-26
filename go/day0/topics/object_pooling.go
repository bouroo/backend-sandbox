// Package main provides Go performance optimization demonstrations.
package topics

import (
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// OBJECT POOLING
// =============================================================================
//
// This file demonstrates object pooling - a pattern for reusing memory
// instead of creating new objects repeatedly.
//
// ANALOGY:
// - Without pooling: Like buying a new coffee cup every time you want coffee
// - With pooling: Buy one cup, wash it, reuse it for the next coffee
//
// BENEFITS:
// - Reduces GC pressure (fewer allocations)
// - Improves performance in high-frequency scenarios
// - Reuses pre-allocated memory
//
// Pool is a generic object pool implementation.
// Using sync.Pool for thread-safe object reuse.
var pool = sync.Pool{
	New: func() any {
		// Create a new buffer when pool is empty
		return &Buffer{Data: make([]byte, 1024)}
	},
}

// Buffer represents a reusable data buffer.
// In real scenarios, this could be a connection, a parser, etc.
type Buffer struct {
	Data   []byte
	Length int
}

// Reset clears the buffer for reuse.
func (b *Buffer) Reset() {
	b.Length = 0
}

// Write adds data to the buffer.
func (b *Buffer) Write(data []byte) {
	copy(b.Data[b.Length:], data)
	b.Length += len(data)
}

// GetBuffer retrieves a buffer from the pool.
func GetBuffer() *Buffer {
	return pool.Get().(*Buffer)
}

// PutBuffer returns a buffer to the pool.
func PutBuffer(b *Buffer) {
	b.Reset()
	pool.Put(b)
}

// =============================================================================
// DEMO: Object Pooling
// =============================================================================

// simulateWorkWithoutPool demonstrates creating new objects each time.
// This causes GC pressure and slower performance.
func simulateWorkWithoutPool(iterations int) time.Duration {
	start := time.Now()

	for range iterations {
		// Create new buffer each time - causes allocation!
		buf := &Buffer{Data: make([]byte, 1024)}
		buf.Write([]byte("hello"))
		_ = buf.Length
		// Buffer is abandoned and GC will collect it
	}

	return time.Since(start)
}

// simulateWorkWithPool demonstrates reusing objects from the pool.
// This reduces GC pressure and improves performance.
func simulateWorkWithPool(iterations int) time.Duration {
	start := time.Now()

	for range iterations {
		// Get buffer from pool - reuse instead of allocate!
		buf := pool.Get().(*Buffer)
		buf.Write([]byte("hello"))
		_ = buf.Length
		// Return buffer to pool for reuse
		buf.Reset()
		pool.Put(buf)
	}

	return time.Since(start)
}

// RunPoolingDemo demonstrates the performance difference.
func RunPoolingDemo() {
	fmt.Println("================================================================================")
	fmt.Println("                         OBJECT POOLING DEMONSTRATION                         ")
	fmt.Println("================================================================================")
	fmt.Println()

	const iterations = 100000

	// Warm up the pool
	for range 10 {
		buf := pool.Get()
		pool.Put(buf)
	}

	// Test without pooling
	fmt.Println("=== WITHOUT OBJECT POOL ===")
	timeWithoutPool := simulateWorkWithoutPool(iterations)
	fmt.Printf("Iterations: %d\n", iterations)
	fmt.Printf("Time taken: %v\n", timeWithoutPool)
	fmt.Println()

	// Test with pooling
	fmt.Println("=== WITH OBJECT POOL ===")
	timeWithPool := simulateWorkWithPool(iterations)
	fmt.Printf("Iterations: %d\n", iterations)
	fmt.Printf("Time taken: %v\n", timeWithPool)
	fmt.Println()

	// Calculate improvement
	improvement := float64(timeWithoutPool.Nanoseconds()) / float64(timeWithPool.Nanoseconds())
	fmt.Printf("=== PERFORMANCE IMPROVEMENT ===\n")
	fmt.Printf("Speedup: %.2fx\n", improvement)
	fmt.Printf("Time saved: %v\n", timeWithoutPool-timeWithPool)
	fmt.Println()

	// Run micro-benchmarks for different buffer sizes
	const benchIterations = 100000

	// Small buffer (1KB) benchmark
	smallWithoutStart := time.Now()
	for range benchIterations {
		buf := &Buffer{Data: make([]byte, 1024)}
		buf.Write([]byte("hello"))
		_ = buf.Length
	}
	smallWithoutTime := time.Since(smallWithoutStart)
	smallWithoutNsOp := float64(smallWithoutTime.Nanoseconds()) / float64(benchIterations)

	smallWithStart := time.Now()
	for range benchIterations {
		buf := pool.Get().(*Buffer)
		buf.Write([]byte("hello"))
		_ = buf.Length
		buf.Reset()
		pool.Put(buf)
	}
	smallWithTime := time.Since(smallWithStart)
	smallWithNsOp := float64(smallWithTime.Nanoseconds()) / float64(benchIterations)

	// Medium buffer (10KB) benchmark
	mediumWithoutStart := time.Now()
	for range benchIterations {
		buf := &Buffer{Data: make([]byte, 10240)}
		buf.Write([]byte("hello"))
		_ = buf.Length
	}
	mediumWithoutTime := time.Since(mediumWithoutStart)
	mediumWithoutNsOp := float64(mediumWithoutTime.Nanoseconds()) / float64(benchIterations)

	mediumPool := sync.Pool{
		New: func() any {
			return &Buffer{Data: make([]byte, 10240)}
		},
	}
	// Warm up medium pool
	for range 10 {
		buf := mediumPool.Get()
		mediumPool.Put(buf)
	}

	mediumWithStart := time.Now()
	for range benchIterations {
		buf := mediumPool.Get().(*Buffer)
		buf.Write([]byte("hello"))
		_ = buf.Length
		buf.Reset()
		mediumPool.Put(buf)
	}
	mediumWithTime := time.Since(mediumWithStart)
	mediumWithNsOp := float64(mediumWithTime.Nanoseconds()) / float64(benchIterations)

	// Large buffer (100KB) benchmark
	largeWithoutStart := time.Now()
	for range benchIterations {
		buf := &Buffer{Data: make([]byte, 102400)}
		buf.Write([]byte("hello"))
		_ = buf.Length
	}
	largeWithoutTime := time.Since(largeWithoutStart)
	largeWithoutNsOp := float64(largeWithoutTime.Nanoseconds()) / float64(benchIterations)

	largePool := sync.Pool{
		New: func() any {
			return &Buffer{Data: make([]byte, 102400)}
		},
	}
	// Warm up large pool
	for range 10 {
		buf := largePool.Get()
		largePool.Put(buf)
	}

	largeWithStart := time.Now()
	for range benchIterations {
		buf := largePool.Get().(*Buffer)
		buf.Write([]byte("hello"))
		_ = buf.Length
		buf.Reset()
		largePool.Put(buf)
	}
	largeWithTime := time.Since(largeWithStart)
	largeWithNsOp := float64(largeWithTime.Nanoseconds()) / float64(benchIterations)

	// Print benchmark results with actual measurements
	fmt.Println("=== BENCHMARK RESULTS ===")
	fmt.Println("Size Comparison (1KB buffer):")
	fmt.Printf("  - Without pool: ~%.0f ns/op\n", smallWithoutNsOp)
	fmt.Printf("  - With pool: ~%.0f ns/op\n", smallWithNsOp)
	smallSpeedup := smallWithoutNsOp / smallWithNsOp
	fmt.Printf("  -> Speedup: %.1fx\n", smallSpeedup)
	fmt.Println()
	fmt.Println("Size Comparison (10KB buffer):")
	fmt.Printf("  - Without pool: ~%.0f ns/op\n", mediumWithoutNsOp)
	fmt.Printf("  - With pool: ~%.0f ns/op\n", mediumWithNsOp)
	mediumSpeedup := mediumWithoutNsOp / mediumWithNsOp
	fmt.Printf("  -> Speedup: %.0fx (massive improvement!)\n", mediumSpeedup)
	fmt.Println()
	fmt.Println("Size Comparison (100KB buffer):")
	fmt.Printf("  - Without pool: ~%.0f ns/op\n", largeWithoutNsOp)
	fmt.Printf("  - With pool: ~%.0f ns/op\n", largeWithNsOp)
	largeSpeedup := largeWithoutNsOp / largeWithNsOp
	fmt.Printf("  -> Speedup: %.0fx (AMAZING!)\n", largeSpeedup)
	fmt.Println()
	fmt.Println("Key Insight:")
	fmt.Println("  - Pooling is MORE effective for larger objects")
	fmt.Println("  - Larger allocations benefit more from reuse")
	fmt.Println("  - GC pressure reduction is significant")
	fmt.Println()

	// Explain when to use pooling
	fmt.Println("=== WHEN TO USE OBJECT POOLING ===")
	fmt.Println("✓ High-frequency allocations (loops, request handlers)")
	fmt.Println("✓ Objects with expensive initialization")
	fmt.Println("✓ Burstable workloads with many short-lived objects")
	fmt.Println()
	fmt.Println("✗ Don't pool objects that are:")
	fmt.Println("  - Rarely used (pool overhead not worth it)")
	fmt.Println("  - Very small (allocation cost negligible)")
	fmt.Println("  - Held for long periods (defeats pooling purpose)")
	fmt.Println()

	fmt.Println("================================================================================")
}
