package unsafex

import (
	"reflect"
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

// AddrOf returns the memory address v.
//
// AddrOf panics if v has a pass-by-value type.
func AddrOf(v any) uintptr {
	return uintptr(reflect.ValueOf(v).UnsafePointer())
}

// ByteString converts a byte-slice to a string without copying.
// Because strings in Go are immutable, the original slice should
// not be modified during the lifetime of the returned string.
//
// However, if the original slice is heap allocated, modifications will carry over.
func ByteString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// StringBytes returns the underlying bytes for a string without copying.
// Because strings in Go are immutable, the returned bytes must not be modified.
//
// However, if the original string is heap alloated, modifications will carry over.
func StringBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// AsBool returns the exact value of i cast to a boolean.
//
// Note: if 'i' is an integer with the first bit set to zero,
// the boolean will be interpreted as false. Because of this,
// values like -2 == false and -1 == true.
func AsBool(i int) bool {
	return *(*bool)(unsafe.Pointer(&i))
}

// AsInt returns the integer value of a boolean.
func AsInt(b bool) int {
	// @note(judah): we cast to int8 instead of uint8 to ensure the sign persists.
	return int(*(*int8)(unsafe.Pointer(&b)))
}
