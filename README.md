# MessagePack Code Generator for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/shamaton/msgpackgen.svg)](https://pkg.go.dev/github.com/shamaton/msgpackgen)
![test](https://github.com/shamaton/msgpackgen/workflows/test/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/shamaton/msgpackgen)](https://goreportcard.com/report/github.com/shamaton/msgpackgen)
[![codecov](https://codecov.io/gh/shamaton/msgpackgen/branch/main/graph/badge.svg?token=K7M3778X7C)](https://codecov.io/gh/shamaton/msgpackgen)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fshamaton%2Fmsgpackgen.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fshamaton%2Fmsgpackgen?ref=badge_shield)

msgpackgen provides a code generation tool and serialization library for [MessagePack](http://msgpack.org/). 
A notable feature is that it is **easy to maintain** and runs **extremely fast**.

## Quickstart
In a source file(ex. main.go), include the following directive:

```go
//go:generate msgpackgen
or
//go:generate go run github.com/shamaton/msgpackgen
```

And run the following command in your shell:

```shell
go generate
```

It will generate one `.go` file for serialization, default is `resolver.msgpackgen.go`.
You can call one method to use generated code.

```go
func main() {
	// this method is defined in resolver.msgpackgen.go
	RegisterGeneratedResolver()
	
	// ... your code ...
}
```

`Marshal` and `Unmarshal` look like this:
```go
    // import github.com/shamaton/msgpackgen/msgpack
    v := ResolvedStruct{}
    b, err := msgpack.Marshal(v)
    if err != nil {
        panic(err)
    }
    
    var vv ResolvedStruct
    err = msgpack.Unmarshal(b, &vv)
    if err != nil {
        panic(err)
    }
```

## Serializer
### Supported Types
primitive types:  
`int`, `int8`, `int16`, `int32`, `int64`,
`uint`, `uint8`, `uint16`, `uint32`, `uint64`,
`float32`, `float64`, `string`, `bool`, `byte`, `rune`,
`complex64`, `complex128`

slice, array: `[]`, `[cap]`

map: `map[key]value`

struct: `time.Time` and structures you defined

### Tags
Renaming or omitting are available.

* Renaming fields via `msgpack:"field_name"`
* Omitting fields via `msgpack:"-"`


### Switch Default Behaviour
Default serialization behaviour is map type. But the performance of array type is better.
If you want to switch array type as default, use `SetStructAsArray`.

Also, you can use `MarshalAsArray`, `UnmarshalAsArray`.

## Code Generator
### Easy maintenance
This tool generates only one `.go` file.
All you have to delete is one generated `.go` file. 

### Analyzing
`-input-dir` needs to be in $GOPATH.

Resolver is generated by recursively searching directories,
but some directories and files are ignored.
* Prefix `_` and `.` directory.
* `testdata` and `vendor` directory
* `_test.go` file

If you use `-input-file` option, it will work without considering the above conditions.

---

Compatible with various import rules.

```go
import ("
	"example.com/user/a/b"
	d "example.com/user/a/c"
	. "example.com/user/a/e"
)
```

### Not Generated Case
Not generated in the following cases:

```go
// ex. a/example.go
type Example struct {
	// unsupported types
	Interface interface{}
	Uintptr uintptr
	Error error
	Chan chan
	Func func()
	
	// nested struct is also unsupported
	NestedStruct struct {}
	
	// because b.Example is not generated
	B b.Example
	
	// because bytes.Butffer is in outside package
	Buf bytes.Buffer
}

func (e Example) F() {
	// unsupported  struct defined in func
	type InFunction struct {}
}

// ex a/b/example.go
type Example struct {
	Interface interface{}
}
```
If you serialize a struct that wasn't code generated, it will be processed by [shamaton/msgpack](https://github.com/shamaton/msgpack).

### Strict Mode
If you use strict mode(option `-strict`), you will get an error if an unrecognized structure is passed.
In other words,  [shamaton/msgpack](https://github.com/shamaton/msgpack) is not used.

---

See also `msgpackgen -h`
```shell
Usage of msgpackgen:
  -dry-run
        dry run mode
  -input-dir string
        input directory. input-file cannot be used at the same time (default ".")
  -input-file string
        input a specific file. input-dir cannot be used at the same time
  -output-dir string
        output directory (default ".")
  -output-file string
        name of generated file (default "resolver.msgpackgen.go")
  -pointer int
        pointer level to consider (default 1)
  -strict
        strict mode
  -v    verbose diagnostics
```

## Benchmarks

These results are recorded by [msgpack_bench](https://github.com/shamaton/msgpack_bench) at 2021/02.
<img width="866" alt="bench" src="https://user-images.githubusercontent.com/4637556/107843994-23439e00-6e13-11eb-9303-296be7c24282.png">