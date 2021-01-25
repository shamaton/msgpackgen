package tst

import (
	"time"

	other "github.com/shamaton/msgpackgen/internal/tst/tst"
	. "github.com/shamaton/msgpackgen/internal/tst/tst/tst"
)

//go:generate go run github.com/shamaton/msgpackgen -s -p 2 -v -g resolver.msgpackgen_test.go

// point
// ドットインポートできる
// 別名インポートも出力できる
// can embedded
// ワンファイル
// msgp以上にこうそく

type TestingValue struct {
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

	Slice  []int
	Array1 [8]float32
	Array2 [31280]string
	Array3 [1031280]bool
	Array4 [0b11]int
	Array5 [0o22]int
	Array6 [0x33]int

	MapIntInt map[string]int

	Pint      *int
	P2string  **string
	P3float32 ***float32

	IntPointers []*int
	MapPointers map[*uint]**string
}

func (v TestingValue) Function() int {
	return v.Int + v.Int
}

type TestingTime struct {
	Time        time.Time
	TimePointer *time.Time
}

type TestingTag struct {
	Tag    int `msgpack:"tag_tag_tag_tag_tag"`
	Ignore int `msgpack:"ignore"`
	Omit   int `msgpack:"-"`
}

type TestingStruct struct {
	Int int
	// embedded
	Embedded
	Emb embedded
	// todo : same file struct
	// todo : other file struct

	// package name
	other.A

	// dot import
	BB DotImport
	Time
}

type Private struct {
	i int
}

type NotGenerated1 struct {
	Int       int
	Interface interface{}
}

type NotGenerated2 struct {
	Int int
	Ptr uintptr
}

type NotGenerated3 struct {
	Error error
	Int   int
}

type NotGenerated4 struct {
	Chan chan int
	Int  int
}

type NotGenerated5 struct {
	InnerStruct struct {
		Int int
	}
	Int int
}

type NotGenerated6 struct {
	Func func() int
	Int  int
}

type NotGenerated7 struct {
	Child other.NotGeneratedChild
	Int   int
}

type NotGenerated8 struct {
	Child other.NotGeneratedChild // todo : outside struct
	Int   int
}
