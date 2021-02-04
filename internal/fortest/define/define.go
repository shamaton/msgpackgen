package define

import (
	define2 "github.com/shamaton/msgpackgen/internal/fortest/define/define"
)

type A struct {
	Int int
	B   define2.B
}

type NotGeneratedChild struct {
	Interface interface{}
}
