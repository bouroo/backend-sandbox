package benchmarks

import (
	"sync"
	"testing"

	"day0/topics"
)

// =============================================================================
// IMMUTABLE DATA BENCHMARKS
// =============================================================================
//
// This file benchmarks the immutable data topic to demonstrate the performance
// benefits of immutable data structures for concurrent access.
//
// KEY INSIGHTS:
// - Immutable data eliminates race conditions
// - No locks needed for reading immutable data
// - Trade-off: memory usage vs thread safety

// =============================================================================
// IMMUTABLE STRUCT BENCHMARKS
// =============================================================================

// BenchmarkImmutableUserCreate benchmarks creating immutable users.
func BenchmarkImmutableUserCreate(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		user := topics.NewImmutableUser(1, "Alice", 30, "alice@example.com")
		_ = user
	}
}

// BenchmarkImmutableUserWithAge benchmarks updating immutable user with new age.
func BenchmarkImmutableUserWithAge(b *testing.B) {
	user := topics.NewImmutableUser(1, "Alice", 30, "alice@example.com")

	b.ResetTimer()
	for b.Loop() {
		olderUser := user.WithAge(31)
		_ = olderUser
	}
}

// BenchmarkImmutableUserWithName benchmarks updating immutable user with new name.
func BenchmarkImmutableUserWithName(b *testing.B) {
	user := topics.NewImmutableUser(1, "Alice", 30, "alice@example.com")

	b.ResetTimer()
	for b.Loop() {
		namedUser := user.WithName("Bob")
		_ = namedUser
	}
}

// =============================================================================
// IMMUTABLE MAP BENCHMARKS
// =============================================================================

// BenchmarkImmutableMapGet benchmarks reading from immutable map.
func BenchmarkImmutableMapGet(b *testing.B) {
	m := topics.NewImmutableMap()
	m.Set("key1", 100)
	m.Set("key2", 200)
	m.Set("key3", 300)

	b.ResetTimer()
	for b.Loop() {
		_, _ = m.Get("key2")
	}
}

// BenchmarkImmutableMapSet benchmarks writing to immutable map.
func BenchmarkImmutableMapSet(b *testing.B) {
	m := topics.NewImmutableMap()

	b.ResetTimer()
	for b.Loop() {
		m.Set("key", 100)
	}
}

// BenchmarkImmutableMapConcurrentReads benchmarks concurrent reads on immutable map.
func BenchmarkImmutableMapConcurrentReads(b *testing.B) {
	m := topics.NewImmutableMap()
	for i := range 100 {
		m.Set(string(rune('a'+i)), i*10)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = m.Get("key50")
		}
	})
}

// =============================================================================
// IMMUTABLE SLICE BENCHMARKS
// =============================================================================

// BenchmarkImmutableSliceAppend benchmarks appending to immutable slice.
func BenchmarkImmutableSliceAppend(b *testing.B) {
	s := topics.NewImmutableSlice()

	b.ResetTimer()
	for b.Loop() {
		s.Append(100)
	}
}

// BenchmarkImmutableSliceGet benchmarks reading from immutable slice.
func BenchmarkImmutableSliceGet(b *testing.B) {
	s := topics.NewImmutableSlice()
	for i := range 100 {
		s.Append(i)
	}

	b.ResetTimer()
	for b.Loop() {
		_ = s.Get()
	}
}

// BenchmarkImmutableSliceLen benchmarks getting length of immutable slice.
func BenchmarkImmutableSliceLen(b *testing.B) {
	s := topics.NewImmutableSlice()
	for i := range 100 {
		s.Append(i)
	}

	b.ResetTimer()
	for b.Loop() {
		_ = s.Len()
	}
}

// =============================================================================
// MUTABLE VS IMMUTABLE COMPARISON BENCHMARKS
// =============================================================================

// BenchmarkMutableCounterWithLock benchmarks mutable counter with mutex.
func BenchmarkMutableCounterWithLock(b *testing.B) {
	type MutexCounter struct {
		mu    sync.Mutex
		value int
	}

	counter := MutexCounter{}
	const iterations = 1000

	b.ResetTimer()
	for b.Loop() {
		for range iterations {
			counter.mu.Lock()
			counter.value++
			counter.mu.Unlock()
		}
	}
}

// BenchmarkAtomicCounter benchmarks atomic counter (lock-free).
func BenchmarkAtomicCounter(b *testing.B) {
	var counter int64
	const iterations = 1000

	b.ResetTimer()
	for b.Loop() {
		for range iterations {
			counter++
		}
	}
}

// BenchmarkConcurrentImmutableMapReadWrite benchmarks mixed read/write on immutable map.
func BenchmarkConcurrentImmutableMapReadWrite(b *testing.B) {
	m := topics.NewImmutableMap()
	for i := range 50 {
		m.Set(string(rune('a'+i)), i*10)
	}

	var wg sync.WaitGroup
	readCh := make(chan struct{})
	writeCh := make(chan struct{})

	b.ResetTimer()

	// Start readers
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-readCh:
					return
				default:
					_, _ = m.Get("key25")
				}
			}
		}()
	}

	// Start writers
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			i := 0
			for {
				select {
				case <-writeCh:
					return
				default:
					m.Set("keynew", i)
					i++
				}
			}
		}()
	}

	for b.Loop() {
		// Benchmark the concurrent access pattern
	}

	close(readCh)
	close(writeCh)
	wg.Wait()
}
