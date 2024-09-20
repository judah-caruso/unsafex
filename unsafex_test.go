package unsafex_test

import (
	"strings"
	"testing"
	"unsafe"

	"github.com/judah-caruso/unsafex"
)

func TestSizeOf(t *testing.T) {
	// Ensure SizeOf matches unsafe.Sizeof
	var ch chan int
	cases := [][2]uintptr{
		{unsafex.SizeOf[bool](), unsafe.Sizeof(false)},

		{unsafex.SizeOf[uint](), unsafe.Sizeof(uint(0))},
		{unsafex.SizeOf[uint8](), unsafe.Sizeof(uint8(0))},
		{unsafex.SizeOf[uint16](), unsafe.Sizeof(uint16(0))},
		{unsafex.SizeOf[uint32](), unsafe.Sizeof(uint32(0))},
		{unsafex.SizeOf[uint64](), unsafe.Sizeof(uint64(0))},

		{unsafex.SizeOf[int](), unsafe.Sizeof(int(0))},
		{unsafex.SizeOf[int8](), unsafe.Sizeof(int8(0))},
		{unsafex.SizeOf[int16](), unsafe.Sizeof(int16(0))},
		{unsafex.SizeOf[int32](), unsafe.Sizeof(int32(0))},
		{unsafex.SizeOf[int64](), unsafe.Sizeof(int64(0))},

		{unsafex.SizeOf[float32](), unsafe.Sizeof(float32(0))},
		{unsafex.SizeOf[float64](), unsafe.Sizeof(float64(0))},

		{unsafex.SizeOf[complex64](), unsafe.Sizeof(complex64(0))},
		{unsafex.SizeOf[complex128](), unsafe.Sizeof(complex128(0))},

		{unsafex.SizeOf[string](), unsafe.Sizeof("")},
		{unsafex.SizeOf[[]byte](), unsafe.Sizeof([]byte{})},

		{unsafex.SizeOf[map[int]int](), unsafe.Sizeof(map[int]int{})},
		{unsafex.SizeOf[map[string]int](), unsafe.Sizeof(map[string]int{})},

		{unsafex.SizeOf[chan int](), unsafe.Sizeof(ch)},
		{unsafex.SizeOf[struct{}](), unsafe.Sizeof(struct{}{})},
		{unsafex.SizeOf[interface{}](), unsafe.Sizeof(interface{}(0))},
	}

	for i, c := range cases {
		if c[0] != c[1] {
			t.Errorf("expected case %d to be %d, was %d", i, c[1], c[0])
		}
	}
}

func TestAlignOf(t *testing.T) {
	// Ensure AlignOf matches unsafe.Alignof
	var ch chan int
	cases := [][2]uintptr{
		{unsafex.AlignOf[bool](), unsafe.Alignof(false)},

		{unsafex.AlignOf[uint](), unsafe.Alignof(uint(0))},
		{unsafex.AlignOf[uint8](), unsafe.Alignof(uint8(0))},
		{unsafex.AlignOf[uint16](), unsafe.Alignof(uint16(0))},
		{unsafex.AlignOf[uint32](), unsafe.Alignof(uint32(0))},
		{unsafex.AlignOf[uint64](), unsafe.Alignof(uint64(0))},

		{unsafex.AlignOf[int](), unsafe.Alignof(int(0))},
		{unsafex.AlignOf[int8](), unsafe.Alignof(int8(0))},
		{unsafex.AlignOf[int16](), unsafe.Alignof(int16(0))},
		{unsafex.AlignOf[int32](), unsafe.Alignof(int32(0))},
		{unsafex.AlignOf[int64](), unsafe.Alignof(int64(0))},

		{unsafex.AlignOf[float32](), unsafe.Alignof(float32(0))},
		{unsafex.AlignOf[float64](), unsafe.Alignof(float64(0))},

		{unsafex.AlignOf[complex64](), unsafe.Alignof(complex64(0))},
		{unsafex.AlignOf[complex128](), unsafe.Alignof(complex128(0))},

		{unsafex.AlignOf[string](), unsafe.Alignof("")},
		{unsafex.AlignOf[[]byte](), unsafe.Alignof([]byte{})},

		{unsafex.AlignOf[map[int]int](), unsafe.Alignof(map[int]int{})},
		{unsafex.AlignOf[map[string]int](), unsafe.Alignof(map[string]int{})},

		{unsafex.AlignOf[chan int](), unsafe.Alignof(ch)},
		{unsafex.AlignOf[struct{}](), unsafe.Alignof(struct{}{})},
		{unsafex.AlignOf[interface{}](), unsafe.Alignof(interface{}(0))},
	}

	for i, c := range cases {
		if c[0] != c[1] {
			t.Errorf("expected case %d to be %d, was %d", i, c[1], c[0])
		}
	}
}

func TestAddrOf(t *testing.T) {
	a := false
	aadr := unsafex.AddrOf(&a) // NOTICE: we take the address of a here because bools are passed by value.

	b := &a
	badr := unsafex.AddrOf(b)

	if aadr != badr {
		t.Errorf("addresses did not match %d vs %d", aadr, badr)
	}
}

func TestByteString(t *testing.T) {
	bytes := []byte{'h', 'e', 'l', 'l', 'o'}
	str := unsafex.ByteString(bytes)

	if str != "hello" {
		t.Errorf("expected ByteString to return \"hello\", was %q", str)
	}

	byteAddr := unsafex.AddrOf(bytes)
	strAddr := unsafex.AddrOf(str)
	if strAddr != byteAddr {
		t.Errorf("expected returned string's address to be %d, was %d", byteAddr, strAddr)
	}

	bytes[0] = 'J'
	if str[0] != bytes[0] {
		t.Errorf("expected byte to be %c, was %c", bytes[0], str[0])
	}

	slice := []byte{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'}

	head := unsafex.ByteString(slice[:5])
	if head != "hello" {
		t.Errorf("expected head to be \"hello\", was %q", head)
	}

	tail := unsafex.ByteString(slice[6:])
	if tail != "world" {
		t.Errorf("expected tail to be \"world\", was %q", tail)
	}

	new := head + " " + tail
	if new != "hello world" {
		t.Errorf("expected head + tail to be \"hello world\", was %q", new)
	}

	slice[0] = 'H'
	slice[7] = 'W'

	if new != "hello world" {
		t.Error("string still shared reference to original slice after concatenation", new)
	}
}

func TestStringBytes(t *testing.T) {
	str := strings.Clone("Hello World") // Ensure we have a heap allocated string
	bytes := unsafex.StringBytes(str)
	for i, b := range bytes {
		if b != str[i] {
			t.Errorf("byte #%d %c was incorrect %c", i, str[i], b)
		}
	}

	strAddr := unsafex.AddrOf(str)
	byteAddr := unsafex.AddrOf(bytes)
	if strAddr != byteAddr {
		t.Errorf("expected returned slice's address to be %d, was %d", strAddr, byteAddr)
	}

	bytes[0] = 'h'
	bytes[6] = 'W'

	if string(bytes) != str {
		t.Errorf("modifying the slice had no affect the original string %q vs. %q", string(bytes), str)
	}
}

func TestAsInt(t *testing.T) {
	if unsafex.AsInt(false) != 0 {
		t.Errorf("expected false to equal 0, was %d", unsafex.AsInt(false))
	}

	if unsafex.AsInt(true) != 1 {
		t.Errorf("expected true to equal 1, was %d", unsafex.AsInt(true))
	}

	x, y := 0, 0
	z := 10 * unsafex.AsInt(x == y) * 2
	if z != 20 {
		t.Errorf("expected %d to be 20", z)
	}
}

func TestAsBool(t *testing.T) {
	falsey := 0
	if !(unsafex.AsBool(falsey) == false) {
		t.Errorf("expected %d to equal false", falsey)
	}

	falsey = -2
	if !(unsafex.AsBool(falsey) == false) {
		t.Errorf("expected %d to equal false", falsey)
	}

	truthy := 1
	if !(unsafex.AsBool(truthy) == true) {
		t.Errorf("expected %d to equal true", truthy)
	}

	a := (truthy + falsey) * 2
	b := unsafex.AsInt(unsafex.AsBool(a))

	if a != b {
		t.Errorf("conversion from int to bool and back was incorrect %d vs %d", a, b)
	}
}
