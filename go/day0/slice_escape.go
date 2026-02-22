package main

// =============================================================================
// SLICE ESCAPE ANALYSIS
// =============================================================================
//
// This file demonstrates when slices escape to heap and when they stay
// on the stack.
//
// WHAT'S A SLICE?
// - A slice has 3 parts: pointer (data location), length, capacity
// - The pointer points to an underlying array (which CAN be on heap)

// globalSlice is a global variable used to demonstrate escape analysis.
// Global variables force anything assigned to them to escape to heap.
var globalSlice []int

// processSliceWithEscape demonstrates when slices ESCAPE to heap.
//
// WHY ESCAPE HERE?
// - We assign s to globalSlice (global variable)
// - Global variables live for the entire program run
// - Therefore s's data MUST survive function scope → HEAP allocation!
//
// ANALOGY: Writing something in a notebook vs. publishing a book.
//          Global = published book (can't be taken back!)
//
// KEY TAKEAWAY: Assigning to globals/returning/storing = escape to heap.
func processSliceWithEscape(n int) int {
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}
	// Assignment to global escapes to heap
	// The slice header AND its data both escape
	globalSlice = s
	return len(s)
}

// processSliceNoEscape demonstrates when slices STAY ON STACK.
//
// WHY NO ESCAPE?
// - Slice s is only used locally within this function
// - It's never returned, stored, or shared with anything that escapes
// - Go's escape analysis: "This slice dies with the function" → stack allocation!
//
// BENEFIT: No heap allocation, no GC pressure
// COST: None! This is ideal.
//
// ANALOGY: Working on scratch paper - throw it away when done, no filing needed.
//
// KEY TAKEAWAY: Keep data local = stack = fast! Don't return/store unless needed.
func processSliceNoEscape(n int) int {
	s := make([]int, n)
	sum := 0
	for i := range s {
		s[i] = i
		sum += s[i]
	}
	return sum
}
