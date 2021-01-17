package tst

import (
	tst2 "github.com/shamaton/msgpackgen/internal/tst/tst"
	. "github.com/shamaton/msgpackgen/internal/tst/tst/tst"
)

//go:generate go run github.com/shamaton/msgpackgen -s -p 2

// point
// ドットインポートできる
// 別名インポートも出力できる
// can embedded
// ワンファイル
// msgp以上にこうそく

type ValueChecking struct {
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
	Slice      []int
	// chan
	// func
}

func (v ValueChecking) Function() int {
	return v.Int + v.Int
}

type SliceArray struct {
	Slice  []int
	Array1 [8]float32
	Array2 [31280]string
	Array3 [1031280]bool
	Array4 [0b11]int
	Array5 [0o22]int
	Array6 [0x33]int
	Ints   **[]**int
}

type Complexity struct {
	// array / map / pointer
	// double array pointer
	BB B
	Emb
	tst2.B
	G *Time
}

type embedded struct {
}

type Emb struct {
	Val int
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
	Func func() int
	Int  int
}

type NotGenerated6 struct {
	Child tst2.NotGeneratedChild
	Int   int
}
