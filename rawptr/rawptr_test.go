package rawptr_test

import (
	"slices"
	"structs"
	"testing"

	"github.com/judah-caruso/unsafex"
	"github.com/judah-caruso/unsafex/rawptr"
)

func TestPointer(t *testing.T) {
	var value uint32 = 0xAAAA_FFFF

	ptr := rawptr.From(&value)
	base := ptr

	ptr.Add(unsafex.SizeOf[uint16]())
	if ptr != base+2 {
		t.Errorf("expected ptr to be %X, was %X", ptr, base+2)
	}

	half := rawptr.To[uint16](ptr)
	if *half != 0xAAAA {
		t.Errorf("expected half to be 0xAAAA, was %X", *half)
	}

	*half = 0xBBBB
	if value != 0xBBBB_FFFF {
		t.Errorf("expected value to be 0xBBBBFFFF, was %X", value)
	}
}

func TestCast(t *testing.T) {
	values := [2]uint32{0xAAAA_AAAA, 0xBBBB_BBBB}

	aptr := rawptr.From(&values)
	bptr := rawptr.Cast[uint16](aptr)

	v0 := rawptr.To[uint16](bptr.Nth(0))
	*v0 = 0xFAFA

	v1 := rawptr.To[uint16](bptr.Nth(1))
	*v1 = 0xAFAF

	if values[0] != 0xAFAF_FAFA {
		t.Errorf("expected value to be 0xAFAFFAFA, was %X", values[0])
	}

	v2 := rawptr.To[uint16](bptr.Nth(2))
	*v2 = 0xBEEF

	v3 := rawptr.To[uint16](bptr.Nth(3))
	*v3 = 0xDEAD

	if values[1] != 0xDEAD_BEEF {
		t.Errorf("expected value to be 0xDEADBEEF, was %X", values[1])
	}

	asU64 := rawptr.To[uint64](bptr)
	if *asU64 != 0xDEAD_BEEF_AFAF_FAFA {
		t.Errorf("expected value to be 0xDEADBEEFAFAFFAFA, was %X", *asU64)
	}
}

func TestNth(t *testing.T) {
	values := []uint32{1, 2, 3, 4}

	base := rawptr.From(&values[0])
	for i := range values {
		ptr := base.Nth(i)
		nval := rawptr.To[uint32](ptr)
		*nval = uint32(i * 2)

		if values[i] != *nval {
			t.Errorf("expected val to be %d, was %d", values[i], *nval)
		}
	}
}

func TestSlice(t *testing.T) {
	values := []uint8{0xAA, 0xBB, 0xCC, 0xDD}

	type sliceHeader struct {
		_   structs.HostLayout
		ptr rawptr.T[uint8]
		len int
		cap int
	}

	base := rawptr.From(&values)
	header := rawptr.To[sliceHeader](base)
	if val := header.ptr.Nth(1).Deref(); val != values[1] {
		t.Errorf("expected value at index to equal %X, was %X", values[1], val)
	}

	if header.len != len(values) {
		t.Errorf("expected len of %d, got %d", len(values), header.len)
	}

	if header.cap != cap(values) {
		t.Errorf("expected cap of %d, got %d", cap(values), header.cap)
	}

	values = append(values, 5)

	if val := header.ptr.Nth(len(values) - 1).Deref(); val != values[len(values)-1] {
		t.Errorf("expected last value to equal %X, was %X", values[len(values)-1], val)
	}

	if header.len != len(values) {
		t.Errorf("expected len of %d after append, got %d", len(values), header.len)
	}

	if header.cap != cap(values) {
		t.Errorf("expected cap of %d after append, got %d", cap(values), header.cap)
	}
}

func TestArray(t *testing.T) {
	vals := [4]uint64{1, 2, 3, 4}

	base := rawptr.From(&vals)
	nvals := base.Deref()

	if slices.Compare(vals[:], nvals[:]) != 0 {
		t.Errorf("dereferencing an array pointer did not yield the same result %v vs. %v", vals, nvals)
	}

	slices.Reverse(nvals[:])
	if slices.Compare(vals[:], nvals[:]) == 0 {
		t.Errorf("dereferenced array and original still shared the same memory %v vs. %v", vals, nvals)
	}
}

func TestSafeWithNonNil(t *testing.T) {
	val := uint16(10)
	ptr, ok := rawptr.FromSafe(&val)
	if !ok {
		t.Error("expected FromSafe of known non-nil value to return true")
	}

	rawptr, ok := rawptr.ToSafe[uint16](ptr)
	if !ok {
		t.Errorf("expected ToSafe of known non-nil pointer to return true")
	}

	*rawptr = 0
	if val != 0 {
		t.Errorf("mutation of invalid address")
	}
}

func TestSafeWithNil(t *testing.T) {
	var value *uint16 = nil

	ptr, ok := rawptr.FromSafe(value)
	if ok {
		t.Errorf("expected FromSafe to fail when given nil pointer")
	}

	if ptr != 0 {
		t.Errorf("expected FromSafe to return 0, instead was %X", ptr)
	}

	invalid := ptr.Deref()
	if invalid != 0 {
		t.Error("expected Deref to return 0 for invalid pointers")
	}

	nval, ok := rawptr.ToSafe[uint16](ptr)
	if ok {
		t.Errorf("expected ToSafe to fail when given invalid pointer")
	}

	if nval != nil {
		t.Errorf("expected ToSafe to return nil pointer")
	}
}

func TestAlignForward(t *testing.T) {
	val := uint32(10)
	ptr := rawptr.From(&val)
	old := ptr

	ptr.Sub(1)
	if ptr.IsAligned() {
		t.Errorf("IsAligned returned true for an unaligned address")
	}

	ptr.AlignForward()
	if !ptr.IsAligned() {
		t.Errorf("IsAligned returned false for an aligned address")
	}

	if ptr != old {
		t.Errorf("AlignForward did not align back to the previous address, new %x, old %x", ptr, old)
	}
}

func TestAlignBackward(t *testing.T) {
	val := uint32(10)
	ptr := rawptr.From(&val)
	old := ptr

	ptr.Add(1)
	if ptr.IsAligned() {
		t.Errorf("IsAligned returned true for an unaligned address")
	}

	ptr.AlignBackward()
	if !ptr.IsAligned() {
		t.Errorf("IsAligned returned false for an aligned address")
	}

	if ptr != old {
		t.Errorf("AlignBackward did not align back to the previous address, new %x, old %x", ptr, old)
	}
}
