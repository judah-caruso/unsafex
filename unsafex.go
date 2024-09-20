package unsafex

import (
	"unsafe"
)

// SizeOf returns the memory size in bytes required to store a value of type T.
func SizeOf[T any]() uintptr {
	var zero T
	return unsafe.Sizeof(zero)
}

// AlignOf returns the required alignment size in bytes for type T.
func AlignOf[T any]() uintptr {
	var zero T
	return unsafe.Alignof(zero)
}
