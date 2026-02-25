// Package topics provides Go performance optimization demonstrations.
package topics

import (
	"fmt"
	"unsafe"
)

// LargeStruct is a 1KB struct to demonstrate heap vs stack allocation.
type LargeStruct struct {
	Field1 int64
	Field2 int64
	Field3 int64
	Field4 int64
	Field5 int64
	Field6 int64
	Field7 int64
	Field8 int64
	Data   [120]int64
}

// GetLargeStructSize returns the size of LargeStruct.
func GetLargeStructSize() int {
	return int(unsafe.Sizeof(LargeStruct{}))
}

// AddByValue demonstrates PASS BY VALUE.
func AddByValue(a, b LargeStruct) int64 {
	return a.Field1 + b.Field2 + b.Field2
}

// AddByPointer demonstrates PASS BY POINTER.
func AddByPointer(a, b *LargeStruct) int64 {
	return a.Field1 + b.Field1 + a.Field2 + b.Field2
}

// RunPassByValueDemo demonstrates pass by value vs pointer.
func RunPassByValueDemo() {
	fmt.Println("Pass by Value vs Pointer Demo")
}
