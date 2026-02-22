package main

// LargeStruct is a 1KB struct to demonstrate heap vs stack allocation.
// Think of it like a big box - whether we keep it on the stack (quick desk)
// or heap (warehouse) affects performance dramatically.
// Size: ~1KB (960 bytes for array + 64 bytes for fields)
type LargeStruct struct {
	Field1 int64
	Field2 int64
	Field3 int64
	Field4 int64
	Field5 int64
	Field6 int64
	Field7 int64
	Field8 int64
	// Array of 120 int64s = 960 bytes + 8*8 = 1024 bytes total
	// This large size makes copy overhead very visible in benchmarks
	Data [120]int64
}
