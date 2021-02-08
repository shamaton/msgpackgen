package define

import (
	define2 "github.com/shamaton/msgpackgen/internal/fortest/define/define"
	. "time"
)

type A struct {
	Int int
	B   define2.B
}

type AA struct {
	Time
}

type NotGeneratedChild struct {
	Interface interface{}
}
