package rawptr

import (
	"unsafe"

	"github.com/judah-caruso/unsafex"
)

// T represents an arbitrary address in memory with an associated type.
// Note: T follows the same rules and patterns unsafe.Pointer.
type T[Underlying any] uintptr

// From converts a pointer to a raw pointer.
func From[Underlying any](value *Underlying) T[Underlying] {
	return T[Underlying](unsafe.Pointer(value))
}

// To converts a raw pointer into a pointer of the given type.
func To[To, From any](p T[From]) *To {
	return (*To)(unsafe.Pointer(p))
}

// Cast converts a raw pointer of one type to another.
func Cast[To, From any](p T[From]) T[To] {
	return T[To](uintptr(p))
}

// FromSafe converts a non-nil pointer to a raw pointer.
// Returns the raw pointer and a boolean indicating if the pointer is valid.
func FromSafe[Underlying any](value *Underlying) (T[Underlying], bool) {
	if value == nil {
		return 0, false
	}

	return T[Underlying](unsafe.Pointer(value)), true
}

// To converts a raw pointer into a pointer of the given type.
// Returns the pointer and a boolean indicating if the raw pointer was nil.
func ToSafe[To, From any](p T[From]) (*To, bool) {
	ptr := (*To)(unsafe.Pointer(p))
	if ptr == nil {
		return nil, false
	}

	return ptr, true
}

// Size returns the size in bytes of the type associated with this raw pointer.
func (p T[Underlying]) Size() uintptr {
	return unsafex.SizeOf[Underlying]()
}

// Alignment returns the required alignment in bytes of the type associated with this raw pointer.
func (p T[Underlying]) Alignment() uintptr {
	return unsafex.AlignOf[Underlying]()
}

// IsAligned returns if a raw pointer is aligned to a power of two address.
func (p T[Underlying]) IsAligned() bool {
	return uintptr(p)&(p.Alignment()-1) == 0
}

// AlignForward aligns a raw pointer to the next aligned address following the alignment rules of its associated type.
// Note: AlignForward does nothing if the address is already aligned.
func (p *T[Underlying]) AlignForward() {
	align := p.Alignment()
	*(*uintptr)(p) = (*(*uintptr)(p) + align - 1) & ^(align - 1)
}

// AlignBackward aligns a raw pointer to the previous aligned address following the alignment rules of its associated type.
// Note: AlignBackward does nothing if the address is already aligned.
func (p *T[Underlying]) AlignBackward() {
	*(*uintptr)(p) = *(*uintptr)(p) & ^(p.Alignment() - 1)
}

// Deref safely dereferences a raw pointer and returns its value or its zero value if the pointer was invalid.
func (p T[Underlying]) Deref() Underlying {
	v, ok := ToSafe[Underlying](p)
	if !ok {
		var zero Underlying
		return zero
	}

	return *v
}

// Add modifies a raw pointer by incrementing its address by the given amount.
// Note: Add *does not* align the new address.
func (p *T[Underlying]) Add(amt uintptr) {
	*(*uintptr)(p) += amt
}

// Sub modifies a raw pointer by decrementing its address by the given amount.
// Note: Sub *does not* align the new address.
func (p *T[Underlying]) Sub(amt uintptr) {
	*(*uintptr)(p) -= amt
}

// Nth indexes a raw pointer by its associated type and returns a new raw pointer of the same type.
func (p T[Underlying]) Nth(index int) T[Underlying] {
	nptr := uintptr(p) + uintptr(index)*p.Size()
	return T[Underlying](nptr)
}
