// Package unsafex is utility package for unsafe things.
//
// # Assertions
//
// Unsafex exposes a function for assertions that panics if the given condition was false.
// Because assertions are not always wanted (for instance in release builds), a build tag
// can be given to disable them: UNSAFEX_DISABLE_ASSERT
package unsafex
