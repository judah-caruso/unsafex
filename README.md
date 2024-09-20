# UnsafeX

An **unsafe** utility package for [Go](https://go.dev).

## Installation

```
go get github.com/judah-caruso/unsafex
```

## Unsafex

The root package `unsafex` contains general-purpose utilites:

```go
import "github.com/judah-caruso/unsafex"

func main() {
   // Memory size of types
   uint8Size := unsafex.SizeOf[uint8]()

   // Assertions (can be disabled with the UNSAFEX_DISABLE_ASSERT build tag)
   unsafex.Assert(uint8Size == 1, "uint8 actually had a size of %d", uint8Size)
}
```

## Rawptr

The `rawptr` makes working with unsafe pointers a little bit safer and a little bit easier.

See `rawptr/rawptr_test.go` for usage examples.

## License

Public Domain. See [LICENSE](./LICENSE) for details.
