// Package msgpack provides fallback MessagePack encode/decode APIs and shared
// encoding settings used by generated code.
package msgpack

import (
	"errors"

	"github.com/shamaton/msgpack/v3"
)

// ErrUndefinedType is returned by generated code when the value type is not
// handled by generated serializers.
var ErrUndefinedType = errors.New("undefined type")

// Marshal returns the MessagePack-encoded byte array of v.
func Marshal(v any) ([]byte, error) {
	if StructAsArray() {
		return MarshalAsArray(v)
	}
	return MarshalAsMap(v)
}

func marshalWithBuffer(v any, buf []byte) ([]byte, error) {
	if StructAsArray() {
		return marshalAsArrayTo(v, buf)
	}
	return marshalAsMapTo(v, buf)
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
