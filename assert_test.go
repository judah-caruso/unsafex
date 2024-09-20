//go:build !UNSAFEX_DISABLE_ASSERT

package unsafex_test

import (
	"testing"

	"github.com/judah-caruso/unsafex"
)

func TestAssertPanics(t *testing.T) {
	defer func() {
		v := recover()
		str, ok := v.(string)
		if !ok {
			t.Error("Assert did not return a string!")
		}

		if str != "value:10" {
			t.Errorf("unexpected string returned from Assert %q", str)
		}
	}()

	unsafex.Assert(false, "value:%d", 10)
}

func TestAssertDoesNotPanic(t *testing.T) {
	defer func() {
		v := recover()
		if v != nil {
			t.Errorf("expected Assert(true) to not panic: %v", v)
		}
	}()

	unsafex.Assert(true)
}
