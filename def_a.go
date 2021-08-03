package main

import (
	"bytes"
	"time"

	define2 "github.com/shamaton/msgpackgen/internal/fortest/define"
	. "github.com/shamaton/msgpackgen/internal/fortest/define/define"
)

//go:generate go run github.com/shamaton/msgpackgen -output-file resolver_test.go -pointer 2 -strict -v

type testingValue struct {
	Int        int
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Uint       uint
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Float32    float32
	Float64    float64
	String     string
	Bool       bool
	Byte       byte
	Rune       rune
	Complex64  complex64
	Complex128 complex128

	Slice []int8
	Bytes []byte

	DoubleSlice [][]int16
	DoubleArray [3][4]int16

	TripleBytes [][][]byte

	MapIntInt map[string]int

	Pint      *int
	P2string  **string
	P3float32 ***float32

	IntPointers []*int
	MapPointers map[*uint]**string

	P2IntSlice     **[]int
	P2MapStringInt **map[string]int
	P2IntArray     **[1]int

	DoubleSlicePointerMap    [][]**map[string]int
	MapDoubleSlicePointerInt map[string][][]**int

	Abcdefghijabcdefghijabcdefghijabcdefghij int

	AbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghij int
}

func (v testingValue) Function() int {
	type notGenerated9 struct {
		Int int
	}
	ng := notGenerated9{}
	return v.Int + v.Int + ng.Int
}

type testingInt struct {
	I int
}

type testingUint struct {
	U uint
}

type testingFloat32 struct {
	F float32
}

type testingFloat64 struct {
	F float64
}

type testingString struct {
	S string
}

type testingBool struct {
	B bool
}

type testingComplex64 struct {
	C complex64
}

type testingComplex128 struct {
	C complex128
}

type testingSlice struct {
	Slice []int8
}

type testingMap struct {
	Map map[string]int
}

type testingTime struct {
	Time time.Time
}

type testingTimePointer struct {
	Time  *time.Time
	Times []*time.Time
}

type testingArrays struct {
	Array1 [8]float32
	Array2 [31280]string
	Array3 [1031280]bool
	Array4 [0b11]int
	Array5 [0o22]int
	Array6 [0x33]int
}

type testingTag struct {
	Tag    int `msgpack:"tag_tag_tag_tag_tag"`
	Ignore int `msgpack:"ignore"`
	Omit   int `msgpack:"-"`
}

type testingStruct struct {
	Int int

	Inside  inside
	Outside outside

	// package name
	define2.A

	// dot import
	BB DotImport
	Time

	// recursive
	R *recursive

	TmpSlice   [][]inside
	TmpArray   [1]inside
	TmpMap     map[inside]inside
	TmpPointer *inside
}

type inside struct {
	Int int
}

type recursive struct {
	Int int
	R   *recursive
}

type private struct {
	i int
}

func (p *private) SetInt() {
	p.i = 1
}

type notGenerated1 struct {
	Int       int
	Interface interface{}
}

type notGenerated2 struct {
	Int int
	Ptr uintptr
}

type notGenerated3 struct {
	Error error
	Int   int
}

type notGenerated4 struct {
	Chan chan int
	Int  int
}

type notGenerated5 struct {
	InnerStruct struct {
		Int int
	}
	Int int
}

type notGenerated6 struct {
	Func func() int
	Int  int
}

type notGenerated7 struct {
	Child define2.NotGeneratedChild
	Int   int
}

type notGenerated8 struct {
	Child bytes.Buffer
	Int   int
}

type notGeneratedInt int
type notGenerated10 struct {
	Int notGeneratedInt
}
