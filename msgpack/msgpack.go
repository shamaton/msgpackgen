// Package msgpack provides MessagePack encode/decode APIs backed by generated
// resolvers and a fallback runtime.
//
// Generated resolver registration is intended for init or startup time. Call
// RegisterGeneratedResolver from the generated file before goroutines start
// calling Marshal, MarshalTo, or Unmarshal. SetResolver and SetToResolver are
// also registration-time APIs; they are not synchronized with concurrent
// encoding or decoding.
package msgpack

import (
	"github.com/shamaton/msgpack/v3"
)

// Marshal returns the MessagePack-encoded byte array of v.
func Marshal(v any) ([]byte, error) {
	if StructAsArray() {
		return MarshalAsArray(v)
	}
	return MarshalAsMap(v)
}

// MarshalTo returns the MessagePack-encoded byte array of v by appending to buf.
func MarshalTo(v any, buf []byte) ([]byte, error) {
	if StructAsArray() {
		return MarshalAsArrayTo(v, buf)
	}
	return MarshalAsMapTo(v, buf)
}

// Unmarshal analyzes the MessagePack-encoded data and stores
// the result into the pointer of v.
func Unmarshal(data []byte, v any) error {
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
