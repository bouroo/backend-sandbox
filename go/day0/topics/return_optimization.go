// Package topics provides Go performance optimization demonstrations.
package topics

import (
	"fmt"
)

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

// ReturnAddByValue demonstrates RETURN BY VALUE with RVO (Return Value Optimization).
//
// HOW IT WORKS:
// - Without RVO: Create temp box → copy to caller's box → destroy temp
// - With RVO: Caller allocates space → we fill it directly = ZERO copies!
//
// BENEFIT: No heap allocation, no copying, no GC pressure.
// COST: None! This is the ideal case.
//
// KEY TAKEAWAY: Let the compiler help you! Return by value when possible.
func ReturnAddByValue(a, b LargeStruct) LargeStruct {
	c := LargeStruct{
		Field1: a.Field1 + b.Field1,
		Field2: a.Field2 + b.Field2,
	}
	return c
}

// ReturnAddByPointer demonstrates HEAP ESCAPE - returning a pointer to local data.
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
func ReturnAddByPointer(a, b LargeStruct) *LargeStruct {
	c := LargeStruct{
		Field1: a.Field1 + b.Field1,
		Field2: a.Field2 + b.Field2,
	}
	// Store to global to prevent compiler from optimizing away allocation
	globalResult = &c
	return &c // This &c is the "escape hatch" - forces heap allocation!
}

// =============================================================================
// DEMONSTRATION
// =============================================================================

// RunReturnOptimizationDemo demonstrates RVO and heap escape.
func RunReturnOptimizationDemo() {
	fmt.Println("================================================================================")
	fmt.Println("                RETURN VALUE OPTIMIZATION (RVO) DEMONSTRATION                   ")
	fmt.Println("================================================================================")
	fmt.Println()

	fmt.Println("WHAT IS RVO?")
	fmt.Println("Return Value Optimization allows Go to:")
	fmt.Println("  1. Allocate return space in the CALLER (not in the function)")
	fmt.Println("  2. Build the return value directly where it's needed")
	fmt.Println("  3. Zero copies - zero allocations!")
	fmt.Println()
	fmt.Println("RETURN BY VALUE (with RVO):")
	fmt.Println("  ✓ No heap allocation")
	fmt.Println("  ✓ No copy overhead")
	fmt.Println("  ✓ Clean code, great performance")
	fmt.Println()
	fmt.Println("RETURN POINTER (heap escape):")
	fmt.Println("  ✗ Heap allocation required")
	fmt.Println("  ✗ Garbage collector pressure")
	fmt.Println("  ✗ Potential cache misses")

	// Key takeaway
	fmt.Println()
	fmt.Println("=== Key Takeaway ===")
	fmt.Println("✓ Let the compiler help you - return by value when possible!")
	fmt.Println("✓ RVO is one of Go's best optimizations")
	fmt.Println("✓ Avoid returning pointers to local variables unless necessary")
	fmt.Println()

	fmt.Println("================================================================================")
}
