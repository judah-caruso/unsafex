//go:build !UNSAFEX_DISABLE_ASSERT

package unsafex

import "fmt"

// Assert will panic if the given condition is false. It optionally takes a format string and arguments to control the failure message.
// To disable assertions, use the UNSAFEX_DISABLE_ASSERT build flag.
func Assert(cond bool, message ...any) {
	if cond {
		return
	}

	msg := "assertion failed"
	if len(message) >= 1 {
		format, ok := message[0].(string)
		if ok {
			msg = fmt.Sprint(format, message[1:])
		}
	}

	panic(msg)
}
