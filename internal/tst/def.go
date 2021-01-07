package tst

import "github.com/shamaton/msgpackgen/internal/tst/tst"

//go:generate go run github.com/shamaton/msgpackgen

type A struct {
	Int int
	B   tst.B
}
