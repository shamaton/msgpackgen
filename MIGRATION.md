# Migration Guide

## v1

v1 changes generated code from the v0 resolver registration model to package-level APIs generated in the same package as your target structs.

The generated file defines:

* `Marshal`
* `MarshalAsMap`
* `MarshalAsArray`
* `Unmarshal`
* `UnmarshalAsMap`
* `UnmarshalAsArray`

After updating from v0, regenerate `resolver.msgpackgen.go` and call these generated package-level functions directly.

```go
v := ResolvedStruct{}
b, err := Marshal(v)
if err != nil {
	panic(err)
}

var vv ResolvedStruct
err = Unmarshal(b, &vv)
if err != nil {
	panic(err)
}
```

The generated `RegisterGeneratedResolver` function and the runtime resolver interfaces are no longer provided.

The module path remains `github.com/shamaton/msgpackgen` because v1 does not use a semantic import version suffix.

### Unrecognized Types

Generated package-level functions handle generated types directly. Without `-strict`, unrecognized types fall back to the runtime implementation. With `-strict`, generated APIs return `msgpack.ErrUndefinedType` for unrecognized types.
