package union_test

import (
	"testing"

	"github.com/judah-caruso/unsafex/union"
)

func TestBasicUnion(t *testing.T) {
	type Numbers = union.Of[struct {
		uint8
		bool
	}]

	var num Numbers
	union.Set[uint8](&num, 1)

	b := union.Get[bool](num)
	if !b {
		t.Errorf("expected bool value to be true, was %v", b)
	}

	union.Set(&num, false)

	i := union.Get[uint8](num)
	if i != 0 {
		t.Errorf("expected uint8 value to be 0, was %v", i)
	}
}

type (
	expr = union.Of[struct {
		binaryExpr
		intExpr
		floatExpr
	}]
	binaryExpr struct {
		Op  string
		Lhs expr
		Rhs expr
	}
	intExpr   int64
	floatExpr float64
)

func TestUnionOfStructs(t *testing.T) {
	makeInt := func(value int64) (e expr) {
		union.Set(&e, intExpr(value))
		return
	}
	makeFloat := func(value float64) (e expr) {
		union.Set(&e, floatExpr(value))
		return
	}
	makeBinop := func(op string, lhs, rhs expr) (e expr) {
		union.Set(&e, binaryExpr{
			Op:  op,
			Lhs: lhs,
			Rhs: rhs,
		})
		return
	}

	expr1 := makeBinop("+", makeInt(10), makeInt(20))
	bin1 := union.Get[binaryExpr](expr1)
	if bin1.Op != "+" {
		t.Errorf("incorrect op returned from union: %s", bin1.Op)
	}
	if lhs := union.Get[intExpr](bin1.Lhs); lhs != 10 {
		t.Errorf("incorrect lhs returned from union: %v", lhs)
	}
	if rhs := union.Get[intExpr](bin1.Rhs); rhs != 20 {
		t.Errorf("incorrect rhs returned from union: %v", rhs)
	}

	expr2 := makeBinop("-", expr1, makeFloat(3.14))
	bin2 := union.Get[binaryExpr](expr2)
	if bin2.Op != "-" {
		t.Errorf("incorrect op returned from union of union: %s", bin2.Op)
	}
	if lhs := union.Get[binaryExpr](bin2.Lhs); lhs.Op != "+" {
		t.Errorf("incorrect lhs returned from union of union: %v", lhs)
	}
	if rhs := union.Get[floatExpr](bin2.Rhs); rhs != 3.14 {
		t.Errorf("incorrect rhs returned from union of union: %v", rhs)
	}
}

func TestUnionOfPointers(t *testing.T) {
	type Value = union.Of[struct {
		*float64
		*uint64
	}]

	var (
		original uint64 = 100
		value    Value
	)

	if union.Is[*uint64](value) || union.Is[*float64](value) {
		t.Error("union internal type was incorrect before usage")
	}

	union.Set(&value, &original)

	if !union.Is[*uint64](value) {
		t.Error("union internal type was incorrect after Set")
	}

	fptr := union.Get[*float64](value)
	*fptr = 3.14

	if original == 100 {
		t.Error("original value did not change")
	}

	uptr := union.Get[*uint64](value)
	*uptr = 200

	if *fptr == 3.14 {
		t.Error("float pointer value did not change after modification")
	}

	if original != 200 {
		t.Errorf("original value was incorrect: %v", original)
	}
}

func TestUnionString(t *testing.T) {
	type (
		Struct = union.Of[struct {
			int32
			uint32
		}]
		Interface = union.Of[interface {
			Int()
			Bool()
		}]
		Bool = union.Of[bool]
	)

	var (
		s Struct
		i Interface
		b Bool
	)

	if s.String() != "union[none] { int32; uint32 }" {
		t.Errorf("valid union had invalid stringification: %s", s.String())
	}

	if i.String() != b.String() {
		t.Errorf("invalid union had invalid stringification: %s, %s", i.String(), b.String())
	}

	union.Set[int32](&s, 10)

	if s.String() != "union[int32] { int32; uint32 }" {
		t.Errorf("valid union had invalid stringification after Set: %s", s.String())
	}
}

func TestUnionSafe(t *testing.T) {
	type Value = union.Of[struct {
		int32
		uint32
		float32
	}]

	var v Value
	if _, err := union.GetSafe[int32](v); err == nil {
		t.Errorf("GetSafe did not error for an uninitialized union")
	}

	if err := union.SetSafe(&v, false); err == nil {
		t.Error("SetSafe allowed invalid type")
	}

	if err := union.SetSafe[int32](&v, 10); err != nil {
		t.Errorf("SetSafe failed with valid type: %s", err)
	}

	if _, err := union.GetSafe[bool](v); err == nil {
		t.Errorf("GetSafe allowed invalid type")
	}

	if v, err := union.GetSafe[int32](v); err != nil {
		t.Errorf("GetSafe failed with valid type: %s", err)
	} else if v != 10 {
		t.Errorf("GetSafe returned invalid value: %v", v)
	}
}
