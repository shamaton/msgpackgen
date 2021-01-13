package tst

import (
	. "github.com/shamaton/msgpackgen/internal/gen"
	tst2 "github.com/shamaton/msgpackgen/internal/tst/tst"
	. "github.com/shamaton/msgpackgen/internal/tst/tst/tst"
)

//go:generate go run github.com/shamaton/msgpackgen -strict

// point
// ドットインポートできる
// 別名インポートも出力できる
// can embedded
// ワンファイル
// msgp以上にこうそく

type A struct {
	Int  int
	Uint uint
	BB   B
	R    rune
	C    complex128
	Emb
	tst2.B
	G Generator
}

//func (a A) F() { a.Emb = Emb{Val: 1} }

type Emb struct {
	Val int
}
