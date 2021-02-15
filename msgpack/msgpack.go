package msgpack

import (
	"github.com/shamaton/msgpack/v2"
)

// Marshal returns the MessagePack-encoded byte array of v.
func Marshal(v interface{}) ([]byte, error) {
	if StructAsArray() {
		return MarshalAsArray(v)
	}
	return MarshalAsMap(v)
}

// Unmarshal analyzes the MessagePack-encoded data and stores
// the result into the pointer of v.
func Unmarshal(data []byte, v interface{}) error {
	if StructAsArray() {
		return UnmarshalAsArray(data, v)
	}
	return UnmarshalAsMap(data, v)
}

// SetStructAsArray sets default encoding option.
// If this option sets true, default encoding sets to array-format.
func SetStructAsArray(on bool) {
	msgpack.StructAsArray = on
}

// StructAsArray gets default encoding option.
// If this option sets true, default encoding sets to array-format.
func StructAsArray() bool {
	return msgpack.StructAsArray
}

// SetComplexTypeCode sets def.complexTypeCode in github.com/shamaton/msgpack
func SetComplexTypeCode(code int8) {
	msgpack.SetComplexTypeCode(code)
}
