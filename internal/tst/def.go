package tst

import "github.com/shamaton/msgpackgen/internal/tst/tst"

//go:generate go run github.com/shamaton/msgpackgen -strict

type A struct {
	Int  int
	Uint uint
	B    tst.B
}
