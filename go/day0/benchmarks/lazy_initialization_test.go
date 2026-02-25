package benchmarks

import (
	"testing"
	"time"

	"day0/topics"
)

// =============================================================================
// LAZY INITIALIZATION BENCHMARKS
// =============================================================================
//
// This file benchmarks the lazy initialization topic to demonstrate the
// performance benefits of deferring expensive operations until needed.
//
// KEY INSIGHTS:
// - Lazy initialization reduces startup time
// - Saves memory for unused features
// - Trade-off: first-access latency vs memory usage

// =============================================================================
// LAZY CONFIG BENCHMARKS
// =============================================================================

// BenchmarkLazyConfigFirstAccess benchmarks first access to lazy config.
func BenchmarkLazyConfigFirstAccess(b *testing.B) {
	// Each iteration creates a new lazy config (not cached)
	b.ResetTimer()
	for b.Loop() {
		lazyConfig := topics.NewLazyConfig(func() topics.ExpensiveConfig {
			return topics.ExpensiveConfig{
				DatabaseURL: "postgres://localhost:5432/db",
				APIKey:      "secret-key-12345",
				Timeout:     30 * time.Second,
			}
		})
		_ = lazyConfig.Get()
	}
}

// BenchmarkLazyConfigCachedAccess benchmarks subsequent access to cached config.
func BenchmarkLazyConfigCachedAccess(b *testing.B) {
	lazyConfig := topics.NewLazyConfig(func() topics.ExpensiveConfig {
		return topics.ExpensiveConfig{
			DatabaseURL: "postgres://localhost:5432/db",
			APIKey:      "secret-key-12345",
			Timeout:     30 * time.Second,
		}
	})
	// First access loads it
	_ = lazyConfig.Get()

	b.ResetTimer()
	for b.Loop() {
		_ = lazyConfig.Get()
	}
}

// BenchmarkLazyConfigIsLoaded benchmarks checking if config is loaded.
func BenchmarkLazyConfigIsLoaded(b *testing.B) {
	lazyConfig := topics.NewLazyConfig(func() topics.ExpensiveConfig {
		return topics.ExpensiveConfig{}
	})
	_ = lazyConfig.Get() // Load it

	b.ResetTimer()
	for b.Loop() {
		_ = lazyConfig.IsLoaded()
	}
}

// =============================================================================
// SYNC.ONCE BENCHMARKS
// =============================================================================

// BenchmarkServiceRegistryInitialize benchmarks service registry initialization.
func BenchmarkServiceRegistryInitialize(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		registry := &topics.ServiceRegistry{}
		_ = registry.GetService("database")
	}
}

// BenchmarkServiceRegistryMultipleAccess benchmarks multiple service accesses.
func BenchmarkServiceRegistryMultipleAccess(b *testing.B) {
	registry := &topics.ServiceRegistry{}
	// Initialize once
	_ = registry.GetService("database")

	b.ResetTimer()
	for b.Loop() {
		_ = registry.GetService("cache")
		_ = registry.GetService("queue")
	}
}

// =============================================================================
// LAZY CACHE BENCHMARKS
// =============================================================================

// BenchmarkLazyCacheFirstAccess benchmarks first access to lazy cache.
func BenchmarkLazyCacheFirstAccess(b *testing.B) {
	loader := func(key string) any {
		time.Sleep(time.Microsecond) // Simulate slow load
		return "value-" + key
	}

	b.ResetTimer()
	for b.Loop() {
		cache := topics.NewCache(loader)
		_ = cache.Get("key1")
	}
}

// BenchmarkLazyCacheCacheHit benchmarks cache hit scenario.
func BenchmarkLazyCacheCacheHit(b *testing.B) {
	cache := topics.NewCache(func(key string) any {
		return "value-" + key
	})
	// Populate cache
	_ = cache.Get("key1")

	b.ResetTimer()
	for b.Loop() {
		_ = cache.Get("key1")
	}
}

// BenchmarkLazyCacheMultipleKeys benchmarks accessing multiple different keys.
func BenchmarkLazyCacheMultipleKeys(b *testing.B) {
	cache := topics.NewCache(func(key string) any {
		return "value-" + key
	})

	b.ResetTimer()
	for b.Loop() {
		for i := range 10 {
			_ = cache.Get(string(rune('a' + i)))
		}
	}
}

// =============================================================================
// EAGER VS LAZY COMPARISON BENCHMARKS
// =============================================================================

// BenchmarkEagerLoadAll benchmarks eager loading all configs at startup.
func BenchmarkEagerLoadAll(b *testing.B) {
	// Simulate eager loading all configs at startup
	loadConfig := func() topics.ExpensiveConfig {
		return topics.ExpensiveConfig{
			DatabaseURL: "postgres://localhost:5432/db",
			APIKey:      "secret-key-12345",
			Timeout:     30 * time.Second,
		}
	}

	configs := make([]topics.ExpensiveConfig, 10)
	for i := range configs {
		configs[i] = loadConfig()
	}

	b.ResetTimer()
	for b.Loop() {
		// Access first config
		_ = configs[0]
	}
}

// BenchmarkLazyLoadOnDemand benchmarks lazy loading configs on demand.
func BenchmarkLazyLoadOnDemand(b *testing.B) {
	loadConfig := func() topics.ExpensiveConfig {
		return topics.ExpensiveConfig{
			DatabaseURL: "postgres://localhost:5432/db",
			APIKey:      "secret-key-12345",
			Timeout:     30 * time.Second,
		}
	}

	b.ResetTimer()
	for b.Loop() {
		// Only load what we need
		lazyConfig := topics.NewLazyConfig(loadConfig)
		_ = lazyConfig.Get()
	}
}

// BenchmarkLazyCacheWriteHeavy benchmarks write-heavy workload on lazy cache.
func BenchmarkLazyCacheWriteHeavy(b *testing.B) {
	cache := topics.NewCache(func(key string) any {
		return "value-" + key
	})

	// First populate with some data
	for i := range 10 {
		_ = cache.Get(string(rune('a' + i)))
	}

	b.ResetTimer()
	for b.Loop() {
		// Mix of reads and new writes
		for i := range 100 {
			key := string(rune('a' + i%10))
			if i%3 == 0 {
				// New key - will load
				_ = cache.Get(string(rune('z' - i%26)))
			} else {
				// Existing key - cache hit
				_ = cache.Get(key)
			}
		}
	}
}
