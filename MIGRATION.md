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

After updating from v0, call these generated package-level functions directly.

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

### Generated Filename

The default generated filename changed from `resolver.msgpackgen.go` to `msgpack.msgpackgen.go`.

If you already generated code with the old default filename, remove `resolver.msgpackgen.go` before regenerating. Keeping both `resolver.msgpackgen.go` and `msgpack.msgpackgen.go` in the same package causes duplicate generated `Marshal` and `Unmarshal` definitions.

### Unrecognized Types

Generated package-level functions handle generated types directly. Without `-strict`, unrecognized types fall back to the runtime implementation. With `-strict`, generated APIs return `msgpack.ErrUndefinedType` for unrecognized types.
