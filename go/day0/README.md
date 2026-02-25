# Go Performance Optimization Demo

A comprehensive demonstration of Go performance optimization techniques based on [goperf.dev](https://goperf.dev/01-common-patterns/) including struct alignment, memory management, and efficient data handling patterns.

## 📋 Overview

This project demonstrates **10 key Go optimization topics** from goperf.dev through interactive demonstrations and benchmarks:

### Original Topics
1. **Struct Alignment & Memory Padding** - How field ordering affects memory usage and cache efficiency
2. **Pass by Value vs Pointer** - Copy cost vs indirection overhead trade-offs
3. **Receiver Types (Value vs Pointer)** - Method performance implications
4. **Return Value Optimization (RVO)** - How Go optimizes return-by-value
5. **Slice Escape Analysis** - When slices escape to heap vs stay on stack
6. **Stack vs Heap Allocation** - Performance implications of allocation strategies
7. **Object Pooling** - Reusing objects to reduce GC pressure
8. **Batching Operations** - Reducing overhead by grouping I/O operations
9. **Immutable Data Sharing** - Safe concurrent access without locks
10. **Lazy Initialization** - Deferring expensive operations until needed

## 🚀 Quick Start

### Prerequisites
- Go 1.26 or later
- Standard Go toolchain

### Running the Demo
```bash
# Run the interactive demonstration
go run main.go

# Run all benchmarks
go test -bench=. -benchmem -run=^$ .

# Run specific benchmark categories
go test -bench=BenchmarkProcessAligned -benchmem
go test -bench=BenchmarkAddByPointer -benchmem
```

## 📁 Project Structure

```
├── main.go                     # Interactive demo runner
├── types.go                    # Common type definitions
├── topics/                     # Topic implementations
│   ├── struct_alignment.go         # Struct alignment demonstrations
│   ├── pass_by_value.go            # Pass by value vs pointer examples
│   ├── receiver_types.go           # Value vs pointer receiver methods
│   ├── return_optimization.go      # Return value optimization examples
│   ├── slice_escape.go             # Slice escape analysis
│   ├── stack_vs_heap.go            # Stack vs heap allocation
│   ├── object_pooling.go           # Object pooling pattern
│   ├── batching_operations.go      # Batching operations
│   ├── immutable_data.go           # Immutable data sharing
│   ├── lazy_initialization.go      # Lazy initialization
│   └── memory_preallocation.go     # Memory preallocation
├── benchmarks/                 # Benchmark tests
│   └── *_test.go
└── go.mod                      # Go module definition
```

## 🎯 Key Topics

### 1. Struct Alignment & Memory Padding

**Problem**: Poor field ordering causes memory padding waste
- **UnalignedStruct**: 48 bytes (66% padding!)
- **AlignedStruct**: 32 bytes (25% padding)
- **Savings**: 16 bytes per struct (33% reduction)

**Best Practices**:
- Order fields by size: largest first, smallest last
- Use `unsafe.Sizeof()` to check struct sizes
- Consider cache line efficiency (typically 64 bytes)

**Example**:
```go
// Bad - lots of padding
type BadStruct struct {
    Small int8   // 1 byte + 7 padding
    Big   int64  // 8 bytes
}

// Good - minimal padding
type GoodStruct struct {
    Big   int64  // 8 bytes
    Small int8   // 1 byte + 7 padding
}
```

### 2. Pass by Value vs Pointer

**Trade-offs**:
- **Value**: No heap allocation, but copy overhead for large structs
- **Pointer**: No copy cost, but indirection and potential heap allocation

**Guidelines**:
- Pass small structs (< 16 bytes) by value
- Pass large structs (> 100 bytes) by pointer
- Consider mutation needs

### 3. Receiver Types (Value vs Pointer)

**Value Receivers**:
- Copy the entire struct
- Good for small, read-only operations
- Thread-safe (no aliasing)

**Pointer Receivers**:
- Pass only 8-byte pointer
- Essential for large structs or mutation
- Required for interface methods that modify data

### 4. Return Value Optimization (RVO)

**How it works**: Go allocates return space in the caller, avoiding copies
- **Return by value**: Zero copies, zero allocations (with RVO)
- **Return pointer**: Heap allocation required, GC pressure

**Best Practice**: Return by value when possible - let RVO help you!

### 5. Slice Escape Analysis

**Escape happens when**:
- Slices are assigned to global variables
- Slices are returned from functions
- Slices are stored in long-lived data structures

**No escape when**:
- Slices are used locally only
- Data stays within function scope

**Impact**: Non-escaping slices stay on stack (fast), escaping slices go to heap (slow + GC)

### 6. Stack vs Heap Allocation

**Stack**:
- Fast allocation (just move pointer)
- Automatic cleanup
- Limited size (~MBs)
- Great for short-lived data

**Heap**:
- Slower allocation (requires finding space)
- Requires garbage collection
- Much larger (~GBs)
- For data that outlives its function

### 7. Object Pooling

**Problem**: Creating and destroying objects repeatedly causes GC pressure

**Solution**: Reuse objects from a pool instead of allocating new ones

**When to use**:
- High-frequency allocations in tight loops
- Objects with expensive initialization
- Burstable workloads

**Example**:
```go
var pool = sync.Pool{
    New: func() interface{} {
        return &Buffer{data: make([]byte, 1024)}
    },
}

// Use buffer from pool
buf := pool.Get().(*Buffer)
// ... use buffer ...
pool.Put(buf)  // Return to pool for reuse
```

### 8. Batching Operations

**Problem**: Individual operations have high overhead (syscalls, network round-trips)

**Solution**: Group multiple operations together to reduce per-operation cost

**When to use**:
- Database writes (batch inserts)
- Network requests (HTTP batching)
- File I/O (buffered writes)

**Example**:
```go
// Without batching - 1000 individual writes
for _, entry := range entries {
    db.Write(entry.Key, entry.Value)
}

// With batching - 1 batch write
db.BatchWrite(entries)
```

### 9. Immutable Data Sharing

**Problem**: Mutable shared data requires locks, causing contention

**Solution**: Use immutable data structures - no locks needed for reading!

**When to use**:
- Concurrent access without locking
- Functional programming patterns
- Event sourcing architectures

**Example**:
```go
type ImmutableUser struct {
    ID   int64
    Name string
}

// "Modification" creates new instance
func (u ImmutableUser) WithName(name string) ImmutableUser {
    return ImmutableUser{ID: u.ID, Name: name}
}
```

### 10. Lazy Initialization

**Problem**: Expensive initialization at startup slows down application

**Solution**: Defer initialization until the value is actually needed

**When to use**:
- Expensive operations not always needed
- Reducing startup time
- Resource conservation

**Example**:
```go
type LazyConfig struct {
    config   ExpensiveConfig
    loaded   bool
    once     sync.Once
}

func (lc *LazyConfig) Get() ExpensiveConfig {
    lc.once.Do(func() {
        lc.config = loadConfig()
    })
    return lc.config
}
```

## 📊 Benchmarks

### Running Benchmarks

```bash
# All benchmarks with memory allocation info
go test -bench=. -benchmem -run=^$ .

# Specific benchmarks
go test -bench=BenchmarkProcessAligned -benchmem
go test -bench=BenchmarkAddByPointer -benchmem
```

### Key Benchmark Categories

1. **Alignment Benchmarks**:
   - `BenchmarkProcessUnaligned` vs `BenchmarkProcessAligned`
   - `BenchmarkSequentialUnaligned` vs `BenchmarkSequentialAligned`

2. **Passing Benchmarks**:
   - `BenchmarkAddByValue` vs `BenchmarkAddByPointer`

3. **Receiver Benchmarks**:
   - `BenchmarkIncrementByValue` vs `BenchmarkIncrementByPointer`
   - `BenchmarkProcessByValue` vs `BenchmarkProcessByPointer`

4. **Return Optimization**:
   - `BenchmarkReturnAddByValue` vs `BenchmarkReturnAddByPointer`

5. **Slice Escape**:
   - `BenchmarkProcessSliceWithEscape` vs `BenchmarkProcessSliceNoEscape`

6. **Stack vs Heap**:
   - `BenchmarkCreateLargeStructOnStack` vs `BenchmarkCreateLargeStructOnHeap`

## 🔍 Expected Results

### Struct Alignment
- Aligned structs should be 25-50% faster than unaligned
- Memory usage should be 25-50% less

### Pass by Value vs Pointer
- For small structs: Value may be faster (no indirection)
- For large structs: Pointer should be significantly faster

### Receiver Types
- Small structs: Value receiver fine
- Large structs: Pointer receiver much faster

### Return Optimization
- Return by value should show 0 allocations
- Return pointer should show 1 allocation per call

### Slice Escape
- Non-escaping slices should be much faster
- Escaping slices should show heap allocations

### Stack vs Heap
- Stack allocation should be faster
- Heap allocation should show allocation counts

## 💡 Key Takeaways

### Original Topics
1. **Order struct fields largest to smallest** to minimize padding
2. **Pass small structs by value, large structs by pointer**
3. **Use pointer receivers for large types or when mutation is needed**
4. **Return by value when possible** - let RVO handle optimization
5. **Keep data local to avoid heap escape** and GC pressure
6. **Stack allocation is faster** but data must not outlive function
7. **Object pooling reduces GC pressure** for high-frequency allocations
8. **Batching reduces overhead** for I/O-bound operations
9. **Immutable data enables safe concurrent access** without locks
10. **Lazy initialization defers expensive operations** until needed

## 🛠️ Development

### Adding New Benchmarks

1. Create new benchmark functions in the appropriate `*_test.go` file
2. Follow the naming convention: `Benchmark[Topic][Technique]`
3. Include comments explaining what the benchmark demonstrates
4. Compare against existing benchmarks to show performance differences

### Running Tests

```bash
# Run all tests
go test

# Run with verbose output
go test -v

# Run tests with coverage
go test -cover

# Run benchmarks with detailed output
go test -bench=. -benchmem -count=5
```

## 📚 Further Reading

- [Go Memory Model](https://golang.org/ref/mem)
- [Go Optimizations](https://goperf.dev/)
- [Struct Padding](https://kawin.dev/posts/go/struct_memory/)
- [Escape Analysis](https://goperf.dev/01-common-patterns/stack-alloc/)

## 🤝 Contributing

This is an educational demo project. Contributions are welcome to:
- Add more optimization examples
- Improve benchmark coverage
- Enhance documentation
- Add visualizations of memory layouts

## 📄 License

This project is for educational purposes. Feel free to use the examples in your own projects.
