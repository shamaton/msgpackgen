package define

import (
	define2 "github.com/shamaton/msgpackgen/internal/fortest/define/define"
	. "time"
)

// DefinedInt is a primitive-compatible defined type for test.
type DefinedInt int

// ChainedDefinedInt is a primitive-compatible defined type chain for test.
type ChainedDefinedInt DefinedInt

// AliasInt is a primitive-compatible alias type for test.
type AliasInt = int

type hiddenInt int

// A is a definition for test
type A struct {
	Int     int
	B       define2.B
	Defined DefinedInt
	Chained ChainedDefinedInt
	Alias   AliasInt
	Slice   []DefinedInt
	Map     map[DefinedInt]AliasInt
}

// AA is a definition for test
type AA struct {
	Time
}

// NotGeneratedChild is a definition for test
type NotGeneratedChild struct {
	Interface any
}

// NotGeneratedNamedPrimitive uses an inaccessible named primitive for test.
type NotGeneratedNamedPrimitive struct {
	Hidden hiddenInt
}
