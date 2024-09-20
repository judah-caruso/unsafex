//go:build UNSAFEX_DISABLE_ASSERT

package unsafex

// Assert does nothing due to UNSAFEX_DISABLE_ASSERT being set.
// To enable assertions, remove the UNSAFEX_DISABLE_ASSERT build flag.
func Assert(_ bool, _ ...any) {}
