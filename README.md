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

## Features
### Serializer
* Supported Types
* Unsupported Case
* Tags
* Switch Default Behaviour

### Code Generator
* Generate only one file fo easy maintenance
* Recursively search directories to generate resolver
* Support various import methods
  * dot import
  * other name import
* Can specify the pointer level
* Can use strict mode


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

tdb