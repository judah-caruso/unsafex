# UnsafeX

An **unsafe** utility package for [Go](https://go.dev).

[Documentation](https://pkg.go.dev/github.com/judah-caruso/unsafex)

## Installation

```
go get github.com/judah-caruso/unsafex
```

## Unsafex

The root package `unsafex` contains general-purpose utilites for:

- Retreiving memory size using a type rather than a zero-value.
- Retreiving memory alignment using a type rather than a zero-value.
- Retreiving the address of any value.
- Converting `byte`<->`string` without an allocation.
- Converting `int`<->`bool` without branching.

## Rawptr

The `rawptr` makes working with unsafe pointers a little bit safer and a little bit easier.

- Casting to/from different pointer types.
- Typed memory addresses (prevents invalid alignment and makes casting explicit)
- Aligning addresses (forward/backward) based on their underlying type's requirement.

See `rawptr/rawptr_test.go` for usage examples.

## License

Public Domain. See [LICENSE](./LICENSE) for details.
