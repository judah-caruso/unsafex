package union

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/judah-caruso/unsafex"
	"github.com/judah-caruso/unsafex/rawptr"
)

var (
	ErrUninitializedAccess = errors.New("access of uninitialized union")
	ErrInvalidType         = errors.New("type does not exist within union")
)

// anystruct represents a struct type with any members.
//
// Note: because Go's type constraint system can't enforce
// this, anystruct is here for documentation purposes.
type anystruct any

// @note(judah): is there a way to declare the type parameters
// to allow 'type Value union.Of[...]' so users can define their
// own methods?

// Of represents a union of different types.
//
// Since members are not accessed by type instead of name,
// T is expected to be a struct of types like so:
//
//	type Value = union.Of[struct {
//		int32
//		uint32
//		float32
//	})
type Of[T anystruct] struct {
	typ reflect.Kind
	mem []byte
}

// String returns the string representation of a union.
func (u Of[T]) String() string {
	t := getInternalType(u)

	var b strings.Builder
	b.WriteString("union[")
	if u.typ == reflect.Invalid {
		b.WriteString("none")
	} else {
		b.WriteString(u.typ.String())
	}
	b.WriteString("] {")

	if t.Kind() == reflect.Struct {
		b.WriteByte(' ')
		fields := reflect.VisibleFields(t)
		for i, field := range fields {
			b.WriteString(field.Type.String())
			if i < len(fields)-1 {
				b.WriteString("; ")
			}
		}
		b.WriteByte(' ')
	}

	b.WriteByte('}')
	return b.String()
}

// Is returns true if the given type is currently stored in the union.
func Is[E any, T anystruct](u Of[T]) bool {
	// Explicit invalid check to make sure invalid types don't result in false-positives.
	if u.typ == reflect.Invalid {
		return false
	}

	return u.typ == reflect.TypeFor[E]().Kind()
}

// Set overwrites the backing memory of a union with the given value; initializing the union if uninitialized.
//
// Set is unsafe and will not verify if the backing memory has enough capacity to store the value.
// Use [SetSafe] for more safety checks.
func Set[V any, T anystruct](u *Of[T], value V) {
	if u.mem == nil {
		u.mem = make([]byte, unsafex.SizeOf[T]())
	}

	*rawptr.To[V](rawptr.From(&u.mem[0])) = value
	u.typ = reflect.TypeFor[V]().Kind()
}

// SetSafe overwrites the backing memory of a union with the given value,
// returning an error if the value cannot be stored in the union.
//
// Use [Set] for fewer safety checks.
func SetSafe[V any, T anystruct](u *Of[T], value V) error {
	if u.mem == nil {
		u.mem = make([]byte, unsafex.SizeOf[T]())
	}

	vt := reflect.TypeFor[V]()
	for _, field := range getInternalFields(*u) {
		if field.Type == vt {
			*rawptr.To[V](rawptr.From(&u.mem[0])) = value
			u.typ = reflect.TypeFor[V]().Kind()
			return nil
		}
	}

	return fmt.Errorf("%s - %w", vt, ErrInvalidType)
}

// Get returns the union's backing memory interpreted as a value of type V, panicking if the union is uninitialized.
//
// Get is unsafe and will not verify if the type exists within the union.
// Use [GetSafe] for more safety checks.
func Get[V any, T anystruct](u Of[T]) V {
	if u.mem == nil {
		panic(ErrUninitializedAccess)
	}

	return rawptr.Cast[V](rawptr.From(&u.mem[0])).Deref()
}

// GetSafe returns the union's backing memory interpreted as a value of type V, returning an error if the type
// does not exist within the union or the union is uninitialized.
//
// Use [Get] for fewer safety checks.
func GetSafe[V any, T anystruct](u Of[T]) (V, error) {
	if u.mem == nil {
		var zero V
		return zero, ErrUninitializedAccess
	}

	ut := getInternalType(u)
	vt := reflect.TypeFor[V]()
	for _, field := range reflect.VisibleFields(ut) {
		if field.Type == vt {
			return rawptr.Cast[V](rawptr.From(&u.mem[0])).Deref(), nil
		}
	}

	var zero V
	return zero, ErrInvalidType
}

// getInternalType returns the internal type for a union.
//
// This is required because calling reflect.TypeFor with a union
// returns the union type itself, so we need some extra destructuring
// of the type.
func getInternalType[U Of[T], T anystruct](_ U) reflect.Type {
	return reflect.TypeFor[T]()
}

// getInternalFields returns an array of reflect.StructField belonging
// to the internal type of a union.
//
// It returns an empty array if the internal type is not a struct.
func getInternalFields[U Of[T], T anystruct](_ U) []reflect.StructField {
	var fields []reflect.StructField

	backing := reflect.TypeFor[T]()
	if backing.Kind() != reflect.Struct {
		return fields
	}

	for i := range backing.NumField() {
		fields = append(fields, backing.Field(i))
	}

	return fields
}
