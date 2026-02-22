package main

// =============================================================================
// PASS BY VALUE VS POINTER
// =============================================================================
//
// This file demonstrates the difference between passing structs by value
// versus by pointer, and how it affects performance.
//
// ANALOGY:
// - Stack = your scratch paper (fast, temporary)
// - Heap = filing cabinet (slower, persistent)

// addByValue demonstrates PASS BY VALUE - the entire 1KB struct is COPIED.
//
// ANALOGY: Stack = your scratch paper (fast, temporary)
//           Heap = filing cabinet (slower, persistent)
// 
// When we pass by value, Go copies the entire struct onto the stack.
// BENEFIT: No heap allocation needed = no garbage collector (GC) work.
// COST: Copying 1KB takes time, especially in tight loops.
//
// KEY TAKEAWAY: For small structs (< 2 words), pass by value is usually faster.
//               For large structs, consider passing by pointer instead.
func addByValue(a, b LargeStruct) int64 {
	return a.Field1 + b.Field2 + b.Field2
}

// addByPointer demonstrates PASS BY POINTER - only the pointer (8 bytes) is copied.
//
// ANALOGY: Instead of copying a big box, we just write down its location (address).
//          The pointer is like a Post-it note with a warehouse location.
//
// WHAT'S HAPPENING:
// - The pointers (a, b) live on the stack
// - The data they point to MIGHT escape to heap depending on usage
// - Here, data is only read, so it likely stays wherever the caller allocated it
//
// BENEFIT: No copying cost - just 8 bytes per pointer
// COST: Slight indirection (need to follow the pointer to get data)
//
// KEY TAKEAWAY: Pointers avoid copy overhead but add a small dereference cost.
func addByPointer(a, b *LargeStruct) int64 {
	return a.Field1 + b.Field1 + a.Field2 + b.Field2
}
