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

type Int struct {
	Int  int
	Uint uint
	i    int
}

type Float struct {
	Float32 float32
	Float64 float64
}

type String struct {
	String string
}

//func (a Int) F() { a.Emb = Emb{Val: 1} }

type AA struct {
	BB B
	R  rune
	C  complex128
	Emb
	tst2.B
	G Time
}

type Emb struct {
	Val int
}
