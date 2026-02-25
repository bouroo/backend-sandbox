package benchmarks

import (
	"testing"

	"day0/topics"
)

// =============================================================================
// BATCHING OPERATIONS BENCHMARKS
// =============================================================================
//
// This file benchmarks the batching operations topic to demonstrate the
// performance difference between individual operations vs batched operations.
//
// KEY INSIGHTS:
// - Batching reduces syscall overhead
// - Batching improves I/O efficiency
// - Trade-off: latency vs throughput

// =============================================================================
// DATABASE BATCHING BENCHMARKS
// =============================================================================

// BenchmarkDBWriteIndividual benchmarks individual database writes.
func BenchmarkDBWriteIndividual(b *testing.B) {
	db := &topics.SimulatedDB{}
	entries := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	b.ResetTimer()
	for b.Loop() {
		for k, v := range entries {
			db.Write(k, v)
		}
	}
}

// BenchmarkDBWriteBatch benchmarks batched database writes.
func BenchmarkDBWriteBatch(b *testing.B) {
	db := &topics.SimulatedDB{}
	entries := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	b.ResetTimer()
	for b.Loop() {
		db.BatchWrite(entries)
	}
}

// BenchmarkDBWriteBatchSize10 benchmarks batch writes with 10 items.
func BenchmarkDBWriteBatchSize10(b *testing.B) {
	db := &topics.SimulatedDB{}
	entries := make(map[string]string, 10)
	for i := range 10 {
		entries[string(rune('a'+i))] = string(rune('0' + i))
	}

	b.ResetTimer()
	for b.Loop() {
		db.BatchWrite(entries)
	}
}

// BenchmarkDBWriteBatchSize100 benchmarks batch writes with 100 items.
func BenchmarkDBWriteBatchSize100(b *testing.B) {
	db := &topics.SimulatedDB{}
	entries := make(map[string]string, 100)
	for i := range 100 {
		entries[string(rune('a'+i%26))+string(rune('a'+(i/26)%26))] = string(rune('0' + i%10))
	}

	b.ResetTimer()
	for b.Loop() {
		db.BatchWrite(entries)
	}
}

// =============================================================================
// HTTP BATCHING BENCHMARKS
// =============================================================================

// BenchmarkHTTPSingleRequest benchmarks sending single HTTP requests individually.
func BenchmarkHTTPSingleRequest(b *testing.B) {
	client := topics.NewBatchHTTPClient(1, 0) // Flush immediately

	b.ResetTimer()
	for b.Loop() {
		client.Send(topics.HTTPRequest{
			URL:    "/api/item/1",
			Method: "POST",
		})
	}
}

// BenchmarkHTTPSmallBatch benchmarks sending requests in small batches.
func BenchmarkHTTPSmallBatch(b *testing.B) {
	client := topics.NewBatchHTTPClient(10, 0)

	b.ResetTimer()
	for b.Loop() {
		for i := range 10 {
			client.Send(topics.HTTPRequest{
				URL:    "/api/item/1",
				Method: "POST",
			})
			_ = i // Avoid unused variable
		}
	}
}

// BenchmarkHTTPLargeBatch benchmarks sending requests in large batches.
func BenchmarkHTTPLargeBatch(b *testing.B) {
	client := topics.NewBatchHTTPClient(100, 0)

	b.ResetTimer()
	for b.Loop() {
		for i := range 100 {
			client.Send(topics.HTTPRequest{
				URL:    "/api/item/1",
				Method: "POST",
			})
			_ = i // Avoid unused variable
		}
	}
}

// =============================================================================
// BATCH PROCESSING BENCHMARKS
// =============================================================================

// BenchmarkBatchProcessorSmall benchmarks batch processor with small batches.
func BenchmarkBatchProcessorSmall(b *testing.B) {
	processor := topics.NewBatchProcessor(4, 10)
	tasks := make([]topics.Task, 10)
	for i := range tasks {
		tasks[i] = topics.Task{ID: i, Data: "test"}
	}

	b.ResetTimer()
	for b.Loop() {
		_ = processor.ProcessBatch(tasks)
	}
}

// BenchmarkBatchProcessorMedium benchmarks batch processor with medium batches.
func BenchmarkBatchProcessorMedium(b *testing.B) {
	processor := topics.NewBatchProcessor(4, 100)
	tasks := make([]topics.Task, 100)
	for i := range tasks {
		tasks[i] = topics.Task{ID: i, Data: "test"}
	}

	b.ResetTimer()
	for b.Loop() {
		_ = processor.ProcessBatch(tasks)
	}
}

// BenchmarkBatchProcessorLarge benchmarks batch processor with large batches.
func BenchmarkBatchProcessorLarge(b *testing.B) {
	processor := topics.NewBatchProcessor(4, 1000)
	tasks := make([]topics.Task, 1000)
	for i := range tasks {
		tasks[i] = topics.Task{ID: i, Data: "test"}
	}

	b.ResetTimer()
	for b.Loop() {
		_ = processor.ProcessBatch(tasks)
	}
}
