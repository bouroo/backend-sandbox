package main

// =============================================================================
// STACK VS HEAP ALLOCATION
// =============================================================================
//
// This file demonstrates the difference between stack and heap allocation
// and when each is used.
//
// WHAT'S THE STACK?
// - Fast, limited size (~MBs)
// - Automatic cleanup (just move stack pointer)
// - Great for short-lived data
//
// WHAT'S THE HEAP?
// - Slower, larger (~GBs)
// - Requires garbage collection
// - For data that outlives its function

// createLargeStructOnStack creates a large struct and returns it by value.
//
// WHAT'S HAPPENING:
// - LargeStruct lives on stack (fast, no GC)
// - But copy cost is still paid (1KB copied to caller)
// - Compiler CAN optimize this with RVO, but not guaranteed
//
// KEY TAKEAWAY: Even stack-allocated large structs have copy overhead.
func createLargeStructOnStack() LargeStruct {
	s := LargeStruct{
		Field1: 1,
		Field2: 2,
		Field3: 3,
		Field4: 4,
		Field5: 5,
		Field6: 6,
		Field7: 7,
		Field8: 8,
	}
	for i := range len(s.Data) {
		s.Data[i] = int64(i)
	}
	return s
}

// createLargeStructOnHeap creates a large struct on the HEAP.
//
// WHAT'S HAPPENING:
// - Go must allocate 1KB on the heap (slower than stack)
// - This allocation triggers the garbage collector
// - Data survives after function returns (caller gets the pointer)
//
// WHY DO WE NEED HEAP?
// - The returned pointer &s must be valid after createLargeStructOnHeap() ends
// - Stack data is automatically cleaned up when function returns
// - So we NEED heap to persist the data beyond the function call!
//
// KEY TAKEAWAY: Heap = when data must outlive its creating function.
func createLargeStructOnHeap() *LargeStruct {
	s := LargeStruct{
		Field1: 1,
		Field2: 2,
		Field3: 3,
		Field4: 4,
		Field5: 5,
		Field6: 6,
		Field7: 7,
		Field8: 8,
	}
	for i := range len(s.Data) {
		s.Data[i] = int64(i)
	}
	return &s
}
