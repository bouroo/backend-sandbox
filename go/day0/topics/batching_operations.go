// Package topics provides Go performance optimization demonstrations.
package topics

import (
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// BATCHING OPERATIONS
// =============================================================================
//
// This file demonstrates the batching pattern - grouping multiple operations
// together to reduce overhead.
//
// ANALOGY:
// - Without batching: Mailing 100 letters one at a time (100 trips to post office)
// - With batching: Mailing 100 letters in one envelope (1 trip!)
//
// BENEFITS:
// - Reduces syscall overhead
// - Improves I/O efficiency
// - Better utilization of network/disk bandwidth
//
// =============================================================================
// EXAMPLE 1: Batch Database Writes
// =============================================================================

// SimulatedDB represents a database connection for demonstration.
type SimulatedDB struct {
	writeCount int
	mu         sync.Mutex
}

// Write simulates a single database write operation.
func (db *SimulatedDB) Write(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.writeCount++
	// Simulate some work
	time.Sleep(time.Microsecond)
}

// BatchWrite simulates a batch write operation.
func (db *SimulatedDB) BatchWrite(entries map[string]string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	for range entries {
		db.writeCount++
	}
	// Simulate batch work - much faster than individual writes
	time.Sleep(time.Microsecond * 10)
}

// =============================================================================
// EXAMPLE 2: Batch HTTP Requests
// =============================================================================

// HTTPRequest represents an HTTP request.
type HTTPRequest struct {
	URL     string
	Method  string
	Payload []byte
}

// HTTPResponse represents an HTTP response.
type HTTPResponse struct {
	StatusCode int
	Body       []byte
}

// BatchHTTPClient demonstrates batching HTTP requests.
type BatchHTTPClient struct {
	mu         sync.Mutex
	pending    []HTTPRequest
	batchSize  int
	flushDelay time.Duration
}

// NewBatchHTTPClient creates a new batch HTTP client.
func NewBatchHTTPClient(batchSize int, flushDelay time.Duration) *BatchHTTPClient {
	return &BatchHTTPClient{
		batchSize:  batchSize,
		flushDelay: flushDelay,
	}
}

// Send adds a request to the batch and flushes if batch is full.
func (c *BatchHTTPClient) Send(req HTTPRequest) HTTPResponse {
	c.mu.Lock()
	c.pending = append(c.pending, req)

	// Flush if batch is full
	if len(c.pending) >= c.batchSize {
		c.mu.Unlock()
		return c.flush()
	}
	c.mu.Unlock()

	// In real implementation, would also flush after flushDelay
	return HTTPResponse{StatusCode: 200}
}

// flush sends all pending requests as a batch.
func (c *BatchHTTPClient) flush() HTTPResponse {
	requests := c.pending
	c.pending = nil

	// Simulate batch request processing
	_ = len(requests)

	return HTTPResponse{StatusCode: 200, Body: []byte("batch response")}
}

// =============================================================================
// EXAMPLE 3: Batch Processing with Worker Pool
// =============================================================================

// Task represents a work task.
type Task struct {
	ID   int
	Data string
}

// Result represents a task result.
type Result struct {
	TaskID  int
	Success bool
}

// BatchProcessor processes tasks in batches for efficiency.
type BatchProcessor struct {
	taskChan    chan Task
	resultChan  chan Result
	workerCount int
	batchSize   int
}

// NewBatchProcessor creates a new batch processor.
func NewBatchProcessor(workerCount, batchSize int) *BatchProcessor {
	return &BatchProcessor{
		taskChan:    make(chan Task, batchSize*2),
		resultChan:  make(chan Result, batchSize*2),
		workerCount: workerCount,
		batchSize:   batchSize,
	}
}

// ProcessBatch processes a batch of tasks together.
func (bp *BatchProcessor) ProcessBatch(tasks []Task) []Result {
	results := make([]Result, len(tasks))

	// Process all tasks in the batch
	for i, task := range tasks {
		// Simulate work
		_ = task.Data
		results[i] = Result{TaskID: task.ID, Success: true}
	}

	return results
}

// =============================================================================
// DEMO: Batching Operations
// =============================================================================

// demoDatabaseBatching demonstrates batching database writes.
func demoDatabaseBatching() {
	fmt.Println("=== DATABASE WRITE BATCHING ===")

	db := &SimulatedDB{}
	entries := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	// Without batching - individual writes
	start := time.Now()
	for range 10 {
		for k, v := range entries {
			db.Write(k, v)
		}
	}
	individualTime := time.Since(start)
	fmt.Printf("Individual writes (30 total): %v\n", individualTime)
	fmt.Printf("Write count: %d\n", db.writeCount)

	// Reset
	db.writeCount = 0

	// With batching
	start = time.Now()
	for range 10 {
		db.BatchWrite(entries)
	}
	batchTime := time.Since(start)
	fmt.Printf("Batch writes (10 batches): %v\n", batchTime)
	fmt.Printf("Write count: %d\n", db.writeCount)

	improvement := float64(individualTime.Nanoseconds()) / float64(batchTime.Nanoseconds())
	fmt.Printf("Speedup: %.2fx\n", improvement)
	fmt.Println()
}

// demoHTTPBatching demonstrates batching HTTP requests.
func demoHTTPBatching() {
	fmt.Println("=== HTTP REQUEST BATCHING ===")

	client := NewBatchHTTPClient(10, time.Millisecond)

	// Simulate individual requests
	start := time.Now()
	for i := range 100 {
		client.Send(HTTPRequest{
			URL:    fmt.Sprintf("/api/item/%d", i),
			Method: "POST",
		})
	}
	individualTime := time.Since(start)
	fmt.Printf("Individual requests (100): %v\n", individualTime)

	// Simulate batched requests (would be actual batching in production)
	client2 := NewBatchHTTPClient(100, time.Millisecond)
	start = time.Now()
	for i := range 100 {
		client2.Send(HTTPRequest{
			URL:    fmt.Sprintf("/api/item/%d", i),
			Method: "POST",
		})
	}
	batchTime := time.Since(start)
	fmt.Printf("Batched requests (1 batch of 100): %v\n", batchTime)

	improvement := float64(individualTime.Nanoseconds()) / float64(batchTime.Nanoseconds())
	fmt.Printf("Speedup: %.2fx\n", improvement)
	fmt.Println()
}

// RunBatchingDemo demonstrates all batching patterns.
func RunBatchingDemo() {
	fmt.Println("================================================================================")
	fmt.Println("                        BATCHING OPERATIONS DEMONSTRATION                     ")
	fmt.Println("================================================================================")
	fmt.Println()

	demoDatabaseBatching()
	demoHTTPBatching()

	// Run micro-benchmarks for database operations
	db := &SimulatedDB{}
	entries := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	dbIterations := 1000
	// Individual writes benchmark
	dbStart := time.Now()
	for range dbIterations {
		for k, v := range entries {
			db.Write(k, v)
		}
	}
	dbIndividualTime := time.Since(dbStart)
	dbIndividualNsOp := float64(dbIndividualTime.Nanoseconds()) / float64(dbIterations*len(entries))

	// Reset and test batch writes
	db.writeCount = 0
	dbBatchStart := time.Now()
	for range dbIterations {
		db.BatchWrite(entries)
	}
	dbBatchTime := time.Since(dbBatchStart)
	dbBatchNsOp := float64(dbBatchTime.Nanoseconds()) / float64(dbIterations)

	// HTTP benchmarks
	httpIterations := 10000
	// Single requests benchmark
	singleClient := NewBatchHTTPClient(1, time.Millisecond)
	httpSingleStart := time.Now()
	for i := range httpIterations {
		_ = i
		singleClient.Send(HTTPRequest{
			URL:    "/api/item",
			Method: "POST",
		})
	}
	httpSingleTime := time.Since(httpSingleStart)
	httpSingleNsOp := float64(httpSingleTime.Nanoseconds()) / float64(httpIterations)

	// Small batch (10) benchmark
	smallBatchClient := NewBatchHTTPClient(10, time.Millisecond)
	httpSmallStart := time.Now()
	for i := range httpIterations {
		smallBatchClient.Send(HTTPRequest{
			URL:    fmt.Sprintf("/api/item/%d", i),
			Method: "POST",
		})
	}
	httpSmallTime := time.Since(httpSmallStart)
	httpSmallNsOp := float64(httpSmallTime.Nanoseconds()) / float64(httpIterations)

	// Large batch (100) benchmark
	largeBatchClient := NewBatchHTTPClient(100, time.Millisecond)
	httpLargeStart := time.Now()
	for i := range httpIterations {
		largeBatchClient.Send(HTTPRequest{
			URL:    fmt.Sprintf("/api/item/%d", i),
			Method: "POST",
		})
	}
	httpLargeTime := time.Since(httpLargeStart)
	httpLargeNsOp := float64(httpLargeTime.Nanoseconds()) / float64(httpIterations)

	// Print benchmark results with actual measurements
	fmt.Println("=== BENCHMARK RESULTS ===")
	fmt.Println("Database Write Comparison (3 items):")
	fmt.Printf("  - Individual writes: ~%.0f ns/op\n", dbIndividualNsOp)
	fmt.Printf("  - Batch write (3 items): ~%.0f ns/op\n", dbBatchNsOp)
	fmt.Println("  -> Batch amortizes overhead across items")
	fmt.Println()
	fmt.Println("HTTP Request Comparison:")
	fmt.Printf("  - Single request: ~%.0f ns/op\n", httpSingleNsOp)
	fmt.Printf("  - Small batch (10): ~%.0f ns/op (~%.0f ns/item)\n", httpSmallNsOp, httpSmallNsOp/10)
	fmt.Printf("  - Large batch (100): ~%.0f ns/op (~%.0f ns/item)\n", httpLargeNsOp, httpLargeNsOp/100)
	fmt.Println("  -> Larger batches reduce per-item overhead")
	fmt.Println()

	// Explain when to use batching
	fmt.Println("=== WHEN TO USE BATCHING ===")
	fmt.Println("✓ Database writes - group inserts/updates")
	fmt.Println("✓ Network requests - combine multiple API calls")
	fmt.Println("✓ File I/O - buffer writes before flushing")
	fmt.Println("✓ Message queues - batch messages for throughput")
	fmt.Println()
	fmt.Println("✗ Don't batch when:")
	fmt.Println("  - Latency is critical (batching adds delay)")
	fmt.Println("  - Operations are unrelated (complexity not worth it)")
	fmt.Println("  - Single-item latency matters more than throughput")
	fmt.Println()

	// Key insight
	fmt.Println("=== KEY INSIGHT ===")
	fmt.Println("Batching trades latency for throughput:")
	fmt.Println("  - Individual ops: Low latency, high overhead")
	fmt.Println("  - Batched ops: Higher latency, lower overhead")
	fmt.Println()
	fmt.Println("Choose based on your use case:")
	fmt.Println("  - User-facing: Lower latency (fewer batches)")
	fmt.Println("  - Background processing: Higher throughput (larger batches)")
	fmt.Println()

	fmt.Println("================================================================================")
}
