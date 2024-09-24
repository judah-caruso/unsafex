// Package union emulates a tagged union where values of different types overlap in memory.
//
// This provides the same benefits you'd get from a C-style union, with the addition of
// optional runtime safety checks to ensure valid use. However, this package should still
// be used with caution as it doesn't guarantee the garbage collector won't touch stored values.
package union
