package benchmarks

import (
	"testing"

	"day0/topics"
)

// =============================================================================
// OBJECT POOLING BENCHMARKS
// =============================================================================
//
// This file benchmarks the object pooling topic to demonstrate the performance
// benefits of reusing objects instead of creating new ones.
//
// KEY INSIGHTS:
// - Object pooling reduces GC pressure
// - Improves performance in high-frequency scenarios
// - Trade-off: memory usage vs allocation overhead

// =============================================================================
// WITH/WITHOUT POOL COMPARISON BENCHMARKS
// =============================================================================

// BenchmarkWithoutPoolSmall benchmarks allocations without pooling (small).
func BenchmarkWithoutPoolSmall(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		buf := &topics.Buffer{Data: make([]byte, 1024)}
		buf.Write([]byte("hello"))
		_ = buf.Length
	}
}

// BenchmarkWithPoolSmall benchmarks allocations with pooling (small).
func BenchmarkWithPoolSmall(b *testing.B) {
	// Warm up the pool
	for range 10 {
		buf := topics.GetBuffer()
		topics.PutBuffer(buf)
	}

	b.ResetTimer()
	for b.Loop() {
		buf := topics.GetBuffer()
		buf.Write([]byte("hello"))
		_ = buf.Length
		topics.PutBuffer(buf)
	}
}

// BenchmarkWithoutPoolMedium benchmarks allocations without pooling (medium).
func BenchmarkWithoutPoolMedium(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		buf := &topics.Buffer{Data: make([]byte, 10240)}
		buf.Write([]byte("hello world this is a longer string"))
		_ = buf.Length
	}
}

// BenchmarkWithPoolMedium benchmarks allocations with pooling (medium).
func BenchmarkWithPoolMedium(b *testing.B) {
	// Warm up the pool
	for range 10 {
		buf := topics.GetBuffer()
		topics.PutBuffer(buf)
	}

	b.ResetTimer()
	for b.Loop() {
		buf := topics.GetBuffer()
		buf.Write([]byte("hello world this is a longer string"))
		_ = buf.Length
		topics.PutBuffer(buf)
	}
}

// BenchmarkWithoutPoolLarge benchmarks allocations without pooling (large).
func BenchmarkWithoutPoolLarge(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		buf := &topics.Buffer{Data: make([]byte, 102400)}
		buf.Write([]byte("hello world this is a much longer string for benchmarking"))
		_ = buf.Length
	}
}

// BenchmarkWithPoolLarge benchmarks allocations with pooling (large).
func BenchmarkWithPoolLarge(b *testing.B) {
	// Warm up the pool
	for range 10 {
		buf := topics.GetBuffer()
		topics.PutBuffer(buf)
	}

	b.ResetTimer()
	for b.Loop() {
		buf := topics.GetBuffer()
		buf.Write([]byte("hello world this is a much longer string for benchmarking"))
		_ = buf.Length
		topics.PutBuffer(buf)
	}
}

// =============================================================================
// POOL SIZE BENCHMARKS
// =============================================================================

// BenchmarkPoolSmallIterations benchmarks pool with small iterations.
func BenchmarkPoolSmallIterations(b *testing.B) {
	// Warm up
	buf := topics.GetBuffer()
	topics.PutBuffer(buf)

	const iterations = 100

	b.ResetTimer()
	for b.Loop() {
		for range iterations {
			buf := topics.GetBuffer()
			buf.Write([]byte("test"))
			_ = buf.Length
			topics.PutBuffer(buf)
		}
	}
}

// BenchmarkPoolMediumIterations benchmarks pool with medium iterations.
func BenchmarkPoolMediumIterations(b *testing.B) {
	// Warm up
	buf := topics.GetBuffer()
	topics.PutBuffer(buf)

	const iterations = 1000

	b.ResetTimer()
	for b.Loop() {
		for range iterations {
			buf := topics.GetBuffer()
			buf.Write([]byte("test"))
			_ = buf.Length
			topics.PutBuffer(buf)
		}
	}
}

// BenchmarkPoolLargeIterations benchmarks pool with large iterations.
func BenchmarkPoolLargeIterations(b *testing.B) {
	// Warm up
	buf := topics.GetBuffer()
	topics.PutBuffer(buf)

	const iterations = 10000

	b.ResetTimer()
	for b.Loop() {
		for range iterations {
			buf := topics.GetBuffer()
			buf.Write([]byte("test"))
			_ = buf.Length
			topics.PutBuffer(buf)
		}
	}
}

// =============================================================================
// CONCURRENT POOL ACCESS BENCHMARKS
// =============================================================================

// BenchmarkPoolConcurrentSmall benchmarks concurrent pool access (small).
func BenchmarkPoolConcurrentSmall(b *testing.B) {
	// Warm up
	buf := topics.GetBuffer()
	topics.PutBuffer(buf)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := topics.GetBuffer()
			buf.Write([]byte("test"))
			_ = buf.Length
			topics.PutBuffer(buf)
		}
	})
}

// BenchmarkPoolConcurrentMedium benchmarks concurrent pool access (medium).
func BenchmarkPoolConcurrentMedium(b *testing.B) {
	// Warm up
	for range 100 {
		buf := topics.GetBuffer()
		topics.PutBuffer(buf)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := topics.GetBuffer()
			buf.Write([]byte("test data for concurrent access"))
			_ = buf.Length
			topics.PutBuffer(buf)
		}
	})
}

// BenchmarkPoolConcurrentLarge benchmarks concurrent pool access (large).
func BenchmarkPoolConcurrentLarge(b *testing.B) {
	// Warm up with more buffers
	for range 500 {
		buf := topics.GetBuffer()
		topics.PutBuffer(buf)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			buf := topics.GetBuffer()
			buf.Write([]byte("test data for concurrent access benchmark"))
			_ = buf.Length
			_ = counter
			counter++
			topics.PutBuffer(buf)
		}
	})
}

// =============================================================================
// BUFFER REUSE PATTERN BENCHMARKS
// =============================================================================

// BenchmarkBufferReuseSequential benchmarks sequential buffer reuse.
func BenchmarkBufferReuseSequential(b *testing.B) {
	buf := topics.GetBuffer()
	defer topics.PutBuffer(buf)

	data := []byte("sequential test data")

	b.ResetTimer()
	for b.Loop() {
		buf.Reset()
		buf.Write(data)
		_ = buf.Length
	}
}

// BenchmarkBufferReuseMultipleSizes benchmarks reusing buffer with multiple sizes.
func BenchmarkBufferReuseMultipleSizes(b *testing.B) {
	buf := topics.GetBuffer()
	defer topics.PutBuffer(buf)

	smallData := []byte("small")
	mediumData := []byte("medium size data")
	largeData := []byte("this is a much larger data set for testing")

	b.ResetTimer()
	for b.Loop() {
		buf.Reset()
		buf.Write(smallData)
		_ = buf.Length

		buf.Reset()
		buf.Write(mediumData)
		_ = buf.Length

		buf.Reset()
		buf.Write(largeData)
		_ = buf.Length
	}
}

// BenchmarkBufferWithoutReset benchmarks buffer without proper reset.
func BenchmarkBufferWithoutReset(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		buf := topics.GetBuffer()
		buf.Write([]byte("test without reset"))
		_ = buf.Length
		topics.PutBuffer(buf)
	}
}
