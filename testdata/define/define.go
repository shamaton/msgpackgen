package define

import (
	define2 "github.com/shamaton/msgpackgen/testdata/define/define"
)

type A struct {
	Int int
	B   define2.B
}

type NotGeneratedChild struct {
	Interface interface{}
}
