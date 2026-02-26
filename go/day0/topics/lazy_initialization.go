// Package topics provides Go performance optimization demonstrations.
package topics

import (
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// LAZY INITIALIZATION
// =============================================================================
//
// This file demonstrates lazy initialization - deferring expensive operations
// until they are actually needed.
//
// ANALOGY:
// - Eager loading: Pre-packaging every possible lunch option (wasteful!)
// - Lazy loading: Making lunch only when you're hungry (efficient!)
//
// BENEFITS:
// - Reduces startup time
// - Saves memory for unused features
// - Defers expensive operations until needed
//
// =============================================================================
// EXAMPLE 1: Simple Lazy Initialization
// =============================================================================

// ExpensiveConfig represents a configuration that takes time to load.
type ExpensiveConfig struct {
	DatabaseURL string
	APIKey      string
	Timeout     time.Duration
}

// simulateLoad simulates loading config from a slow source.
func simulateLoad() ExpensiveConfig {
	// Simulate slow loading (e.g., reading from disk, network)
	time.Sleep(100 * time.Millisecond)
	return ExpensiveConfig{
		DatabaseURL: "postgres://localhost:5432/db",
		APIKey:      "secret-key-12345",
		Timeout:     30 * time.Second,
	}
}

// LazyConfig demonstrates basic lazy initialization.
type LazyConfig struct {
	config   ExpensiveConfig
	loaded   bool
	loadFunc func() ExpensiveConfig
	mu       sync.Mutex
}

// NewLazyConfig creates a new lazy config loader.
func NewLazyConfig(loadFunc func() ExpensiveConfig) *LazyConfig {
	return &LazyConfig{
		loadFunc: loadFunc,
		loaded:   false,
	}
}

// Get returns the config, loading it on first access.
func (lc *LazyConfig) Get() ExpensiveConfig {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if !lc.loaded {
		fmt.Println("  [Lazy] Loading configuration...")
		lc.config = lc.loadFunc()
		lc.loaded = true
	}
	return lc.config
}

// IsLoaded checks if the config has been loaded.
func (lc *LazyConfig) IsLoaded() bool {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	return lc.loaded
}

// =============================================================================
// EXAMPLE 2: sync.Once for Thread-Safe Lazy Initialization
// =============================================================================

// Service represents a service that requires expensive initialization.
type Service struct {
	Name    string
	Clients int
}

// ServiceRegistry manages services with lazy initialization.
type ServiceRegistry struct {
	services map[string]*Service
	once     sync.Once
	initDone bool
}

// Initialize performs one-time initialization.
func (sr *ServiceRegistry) Initialize() {
	// sync.Once ensures this runs only once, even with concurrent access
	sr.once.Do(func() {
		fmt.Println("  [sync.Once] Initializing services...")
		time.Sleep(50 * time.Millisecond) // Simulate expensive init
		sr.services = map[string]*Service{
			"database": {Name: "Database", Clients: 0},
			"cache":    {Name: "Cache", Clients: 0},
			"queue":    {Name: "Queue", Clients: 0},
		}
		sr.initDone = true
	})
}

// GetService returns a service by name.
func (sr *ServiceRegistry) GetService(name string) *Service {
	// Initialize on first access
	sr.Initialize()
	return sr.services[name]
}

// =============================================================================
// EXAMPLE 3: Lazy Initialization with sync.RWMutex
// =============================================================================

// Cache represents a lazy-loaded cache.
type Cache struct {
	mu     sync.RWMutex
	data   map[string]any
	loader func(string) any
}

// NewCache creates a new lazy cache.
func NewCache(loader func(string) any) *Cache {
	return &Cache{
		data:   make(map[string]any),
		loader: loader,
	}
}

// Get retrieves or loads a value.
func (c *Cache) Get(key string) any {
	// Fast path: check with read lock first
	c.mu.RLock()
	if val, ok := c.data[key]; ok {
		c.mu.RUnlock()
		return val
	}
	c.mu.RUnlock()

	// Slow path: need to load
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if val, ok := c.data[key]; ok {
		return val
	}

	// Load and store
	val := c.loader(key)
	c.data[key] = val
	return val
}

// =============================================================================
// DEMO: Lazy Initialization
// =============================================================================

// demoBasicLazy demonstrates basic lazy initialization.
func demoBasicLazy() {
	fmt.Println("=== BASIC LAZY INITIALIZATION ===")

	lazyConfig := NewLazyConfig(simulateLoad)

	fmt.Println("Config created (not loaded yet)")
	fmt.Printf("Is loaded: %v\n", lazyConfig.IsLoaded())
	fmt.Println()

	// First access - triggers loading
	fmt.Println("First access:")
	config1 := lazyConfig.Get()
	fmt.Printf("  Database URL: %s\n", config1.DatabaseURL)
	fmt.Printf("  Is loaded: %v\n", lazyConfig.IsLoaded())
	fmt.Println()

	// Second access - uses cached value
	fmt.Println("Second access:")
	config2 := lazyConfig.Get()
	fmt.Printf("  API Key: %s\n", config2.APIKey)
	fmt.Printf("  Is loaded: %v\n", lazyConfig.IsLoaded())
	fmt.Println()
}

// demoSyncOnce demonstrates sync.Once pattern.
func demoSyncOnce() {
	fmt.Println("=== SYNC.ONCE PATTERN ===")

	registry := &ServiceRegistry{}

	fmt.Println("Registry created (not initialized)")
	fmt.Println()

	// First access - triggers initialization
	fmt.Println("First access:")
	svc1 := registry.GetService("database")
	fmt.Printf("  Got service: %s\n", svc1.Name)
	fmt.Println()

	// Second access - uses existing
	fmt.Println("Second access:")
	svc2 := registry.GetService("cache")
	fmt.Printf("  Got service: %s\n", svc2.Name)
	fmt.Println()

	// Concurrent access - still only initializes once
	fmt.Println("Concurrent access (5 goroutines):")
	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			svc := registry.GetService("queue")
			fmt.Printf("  Goroutine %d got: %s\n", id, svc.Name)
		}(i)
	}
	wg.Wait()
	fmt.Println()
}

// demoLazyCache demonstrates lazy cache pattern.
func demoLazyCache() {
	fmt.Println("=== LAZY CACHE PATTERN ===")

	cache := NewCache(func(key string) any {
		// Simulate expensive load
		time.Sleep(10 * time.Millisecond)
		return fmt.Sprintf("value-%s", key)
	})

	// First access - loads
	fmt.Println("First access to 'user:1':")
	val1 := cache.Get("user:1")
	fmt.Printf("  Value: %v\n", val1)
	fmt.Println()

	// Second access - cached
	fmt.Println("Second access to 'user:1':")
	val2 := cache.Get("user:1")
	fmt.Printf("  Value: %v\n", val2)
	fmt.Println()

	// Different key - loads
	fmt.Println("First access to 'user:2':")
	val3 := cache.Get("user:2")
	fmt.Printf("  Value: %v\n", val3)
	fmt.Println()
}

// RunLazyInitDemo demonstrates all lazy initialization patterns.
func RunLazyInitDemo() {
	fmt.Println("================================================================================")
	fmt.Println("                     LAZY INITIALIZATION DEMONSTRATION                         ")
	fmt.Println("================================================================================")
	fmt.Println()

	demoBasicLazy()
	demoSyncOnce()
	demoLazyCache()

	// Run micro-benchmarks for lazy initialization
	const benchIterations = 100000

	// Lazy config benchmarks - use faster load function for benchmarking
	fastLoad := func() ExpensiveConfig {
		time.Sleep(1 * time.Millisecond) // Simulate faster load
		return ExpensiveConfig{
			DatabaseURL: "postgres://localhost:5432/db",
			APIKey:      "secret-key-12345",
			Timeout:     30 * time.Second,
		}
	}

	// First access benchmark (includes load)
	lazyConfigBench := NewLazyConfig(fastLoad)
	// Don't print during benchmark
	configGet := func() ExpensiveConfig {
		lazyConfigBench.mu.Lock()
		defer lazyConfigBench.mu.Unlock()
		if !lazyConfigBench.loaded {
			lazyConfigBench.config = fastLoad()
			lazyConfigBench.loaded = true
		}
		return lazyConfigBench.config
	}

	configGet() // Load once
	configFirstStart := time.Now()
	for range benchIterations {
		_ = configGet()
	}
	configFirstTime := time.Since(configFirstStart)
	configFirstNsOp := float64(configFirstTime.Nanoseconds()) / float64(benchIterations)

	// IsLoaded check benchmark
	isLoadedStart := time.Now()
	for range benchIterations {
		_ = lazyConfigBench.IsLoaded()
	}
	isLoadedTime := time.Since(isLoadedStart)
	isLoadedNsOp := float64(isLoadedTime.Nanoseconds()) / float64(benchIterations)

	// Lazy cache benchmarks
	cache := NewCache(func(key string) any {
		time.Sleep(1 * time.Millisecond) // Simulate fast load
		return fmt.Sprintf("value-%s", key)
	})

	// First access (cache miss)
	cache.Get("benchkey") // Load once
	cacheMissStart := time.Now()
	for i := range benchIterations {
		_ = cache.Get(fmt.Sprintf("benchkey-%d", i))
	}
	cacheMissTime := time.Since(cacheMissStart)
	cacheMissNsOp := float64(cacheMissTime.Nanoseconds()) / float64(benchIterations)

	// Cached access (cache hit)
	cache.Get("hitkey") // Load once
	cacheHitStart := time.Now()
	for range benchIterations {
		_ = cache.Get("hitkey")
	}
	cacheHitTime := time.Since(cacheHitStart)
	cacheHitNsOp := float64(cacheHitTime.Nanoseconds()) / float64(benchIterations)

	// Multiple keys benchmark
	cacheMulti := NewCache(func(key string) any {
		time.Sleep(1 * time.Millisecond)
		return fmt.Sprintf("value-%s", key)
	})
	for i := range 10 {
		cacheMulti.Get(fmt.Sprintf("key%d", i))
	}
	multiStart := time.Now()
	for range 10000 {
		for i := range 10 {
			_ = cacheMulti.Get(fmt.Sprintf("key%d", i))
		}
	}
	multiTime := time.Since(multiStart)
	multiNsOp := float64(multiTime.Nanoseconds()) / 100000

	// Print benchmark results with actual measurements
	fmt.Println("=== BENCHMARK RESULTS ===")
	fmt.Println("Lazy Config:")
	fmt.Printf("  - First access: ~%.0f ns/op\n", configFirstNsOp)
	fmt.Printf("  - Cached access: ~%.0f ns/op\n", configFirstNsOp)
	fmt.Printf("  - IsLoaded check: ~%.0f ns/op\n", isLoadedNsOp)
	fmt.Println()
	fmt.Println("Lazy Cache:")
	fmt.Printf("  - First access (cache miss): ~%.0f ns/op\n", cacheMissNsOp)
	fmt.Printf("  - Cached access (cache hit): ~%.0f ns/op\n", cacheHitNsOp)
	fmt.Printf("  - Multiple keys (10): ~%.0f ns/op\n", multiNsOp)
	fmt.Println()
	fmt.Println("Key Insight:")
	cachedSpeedup := cacheMissNsOp / cacheHitNsOp
	fmt.Printf("  - Lazy initialization is ~%.0fx faster on cached access\n", cachedSpeedup)
	fmt.Println("  - Trade-off: first access slower, subsequent access much faster")
	fmt.Println()

	// Explain when to use lazy initialization
	fmt.Println("=== WHEN TO USE LAZY INITIALIZATION ===")
	fmt.Println("✓ Expensive initialization (database, network, file I/O)")
	fmt.Println("✓ Features that may not be used")
	fmt.Println("✓ Reducing startup time")
	fmt.Println("✓ Resource conservation")
	fmt.Println()
	fmt.Println("✗ Don't use lazy initialization when:")
	fmt.Println("  - Required at startup anyway")
	fmt.Println("  - Multiple threads need it early (adds complexity)")
	fmt.Println("  - Error handling is critical (errors deferred)")
	fmt.Println()

	// Patterns comparison
	fmt.Println("=== PATTERN COMPARISON ===")
	fmt.Println("1. Basic (with mutex): Simple but has lock overhead")
	fmt.Println("2. sync.Once: Thread-safe, only runs once, no lock on reads")
	fmt.Println("3. RWMutex: Read-heavy workloads, double-check locking")
	fmt.Println()

	fmt.Println("================================================================================")
}
