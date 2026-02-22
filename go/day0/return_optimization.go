package main

// =============================================================================
// RETURN VALUE OPTIMIZATION (RVO)
// =============================================================================
//
// This file demonstrates how Go optimizes return by value through
// Return Value Optimization (RVO).
//
// ANALOGY: Instead of packing a box and shipping it, the caller prepares the box
//          location first, and we just fill it in place. No moving required!

// globalResult is used to prevent compiler from optimizing away allocation.
// Without this, Go might realize the pointer isn't actually used and skip the heap alloc.
var globalResult *LargeStruct

// returnAddByValue demonstrates RETURN BY VALUE with RVO (Return Value Optimization).
//
// HOW IT WORKS:
// - Without RVO: Create temp box → copy to caller's box → destroy temp
// - With RVO: Caller allocates space → we fill it directly = ZERO copies!
//
// BENEFIT: No heap allocation, no copying, no GC pressure.
// COST: None! This is the ideal case.
//
// KEY TAKEAWAY: Let the compiler help you! Return by value when possible.
func returnAddByValue(a, b LargeStruct) LargeStruct {
	c := LargeStruct{
		Field1: a.Field1 + b.Field1,
		Field2: a.Field2 + b.Field2,
	}
	return c
}

// returnAddByPointer demonstrates HEAP ESCAPE - returning a pointer to local data.
//
// ANALOGY: We wrote our return address on the box and mailed it to the caller.
//          Now the caller has the box, so we can't throw it away!
//          This forces Go to put the box in the "warehouse" (heap).
//
// WHY ESCAPE HAPPENS:
// - We return &c (address of local variable c)
// - This proves c "escapes" the function - it outlives the function call
// - Go's escape analysis says: "Can't keep this on stack, must go to heap!"
//
// BENEFIT: Caller gets direct access to the data (no copy)
// COST: Heap allocation + garbage collector work + potential cache misses
//
// KEY TAKEAWAY: Returning pointers to locals = heap allocation. Use sparingly!
func returnAddByPointer(a, b LargeStruct) *LargeStruct {
	c := LargeStruct{
		Field1: a.Field1 + b.Field1,
		Field2: a.Field2 + b.Field2,
	}
	// Store to global to prevent compiler from optimizing away allocation
	globalResult = &c
	return &c // This &c is the "escape hatch" - forces heap allocation!
}
