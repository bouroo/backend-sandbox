// Package topics provides Go performance optimization demonstrations.
package topics

import (
	"fmt"
	"unsafe"
)

// =============================================================================
// RECEIVER TYPES: VALUE VS POINTER
// =============================================================================
//
// This file demonstrates the difference between value and pointer receivers
// for methods, and when to use each.
//
// ANALOGY:
// - Value receiver: Photocopy the document, write on the copy (original unchanged)
// - Pointer receiver: Get the original document, write on it (original modified)

// Counter demonstrates VALUE vs POINTER receiver methods.
type Counter struct {
	value int
}

// IncrementByValue uses VALUE RECEIVER - operates on a COPY.
//
// WHAT'S HAPPENING:
// - Go makes a copy of the Counter struct
// - We modify the copy, not the original
// - Original remains unchanged!
//
// USE WHEN: You want read-only access or thread-local temporary work
func (c Counter) IncrementByValue() int {
	c.value++
	return c.value
}

// IncrementByPointer uses POINTER RECEIVER - operates on the ORIGINAL.
//
// WHAT'S HAPPENING:
// - Go passes a pointer to the original struct
// - We modify the actual original data
// - Changes persist!
//
// USE WHEN: You need to modify the original or avoid copying large structs
func (c *Counter) IncrementByPointer() int {
	c.value++
	return c.value
}

// Increment with pointer receiver - satisfies interface and mutations persist.
// This exists to demonstrate interface implementation with pointer receivers.
func (c *Counter) Increment() int {
	c.value++
	return c.value
}

// DataProcessor embeds LargeStruct to demonstrate receiver performance.
// Used to show the cost difference between value and pointer receivers.
type DataProcessor struct {
	LargeStruct
}

// ProcessByValue uses VALUE RECEIVER - COPIES the entire 1KB struct!
//
// PERFORMANCE IMPACT:
// - Every method call copies 1KB of data
// - In a tight loop, this adds up quickly
// - Cache locality might help, but copy cost dominates
//
// ANALOGY: Handing someone a 1kg box vs. telling them where the box is
//
// WHEN TO USE: Small structs (<~2KB) where copy cost is negligible
func (dp DataProcessor) ProcessByValue() int {
	sum := int64(dp.Field1 + dp.Field2 + dp.Field3 + dp.Field4)
	sum += int64(dp.Field5 + dp.Field6 + dp.Field7 + dp.Field8)
	for i := range dp.Data {
		sum += dp.Data[i]
	}
	return int(sum)
}

// ProcessByPointer uses POINTER RECEIVER - only passes 8-byte pointer.
//
// PERFORMANCE IMPACT:
// - No copy of LargeStruct - just the pointer
// - Much faster for large structs
// - Slight indirection cost (following the pointer)
//
// ANALOGY: Giving directions to the warehouse vs. carrying the whole box
//
// WHEN TO USE: Large structs (>~2KB) or when mutation is needed
func (dp *DataProcessor) ProcessByPointer() int {
	sum := int64(dp.Field1 + dp.Field2 + dp.Field3 + dp.Field4)
	sum += int64(dp.Field5 + dp.Field6 + dp.Field7 + dp.Field8)
	for i := range dp.Data {
		sum += dp.Data[i]
	}
	return int(sum)
}

// Incrementer interface requires mutation capability.
// Only pointer receivers can properly implement this interface.
//
// WHY?
// - Value receivers receive a COPY - changes don't persist
// - Pointer receivers get the ORIGINAL - changes DO persist
// - Interface methods that need mutation MUST use pointer receivers!
type Incrementer interface {
	Increment() int
}

// GetCounterSize returns the size of Counter for demonstration.
func GetCounterSize() int {
	return int(unsafe.Sizeof(Counter{}))
}

// GetDataProcessorSize returns the size of DataProcessor for demonstration.
func GetDataProcessorSize() int {
	return int(unsafe.Sizeof(DataProcessor{}))
}

// =============================================================================
// DEMONSTRATION
// =============================================================================

// RunReceiverTypesDemo demonstrates value vs pointer receiver performance.
func RunReceiverTypesDemo() {
	fmt.Println("================================================================================")
	fmt.Println("                    RECEIVER TYPES DEMONSTRATION                               ")
	fmt.Println("================================================================================")
	fmt.Println()

	fmt.Println("VALUE RECEIVER:")
	fmt.Println("  - Go makes a COPY of the struct")
	fmt.Println("  - Changes don't affect the original")
	fmt.Println("  - Good for small, read-only operations")
	fmt.Println()
	fmt.Println("POINTER RECEIVER:")
	fmt.Println("  - Go passes a pointer to the original")
	fmt.Println("  - Changes persist")
	fmt.Println("  - No copy overhead - efficient for large types")
	fmt.Println()

	// Show sizes
	fmt.Println("=== Receiver Type Impact ===")
	counterSize := GetCounterSize()
	processorSize := GetDataProcessorSize()
	fmt.Printf("Counter: %d bytes (copy cost negligible)\n", counterSize)
	fmt.Printf("DataProcessor: %d bytes (copy cost significant!)\n", processorSize)
	fmt.Println()

	// Guidelines
	fmt.Println("=== Guidelines ===")
	fmt.Println("Use VALUE RECEIVER when:")
	fmt.Println("  - Struct is small (< 16 bytes)")
	fmt.Println("  - You don't need to modify the original")
	fmt.Println("  - Thread-safety is important (no aliasing)")
	fmt.Println()
	fmt.Println("Use POINTER RECEIVER when:")
	fmt.Println("  - Struct is large (> 100 bytes)")
	fmt.Println("  - You need to modify the original")
	fmt.Println("  - Method must satisfy an interface that requires pointers")
	fmt.Println()

	fmt.Println("================================================================================")
}
