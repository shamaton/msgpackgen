package tst

import "github.com/shamaton/msgpackgen/internal/tst/tst/tst"

type A struct {
	Int int
	B   tst.B
}

type NotGeneratedChild struct {
	Interface interface{}
}
