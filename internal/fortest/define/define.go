package define

import (
	define2 "github.com/shamaton/msgpackgen/internal/fortest/define/define"
	. "time"
)

// A is a definition for test
type A struct {
	Int int
	B   define2.B
}

// AA is a definition for test
type AA struct {
	Time
}

// NotGeneratedChild is a definition for test
type NotGeneratedChild struct {
	Interface interface{}
}
