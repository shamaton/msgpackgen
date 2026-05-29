// Package msgpack provides shared encoding settings and errors used by generated
// code.
package msgpack

import (
	"errors"

	"github.com/shamaton/msgpack/v3"
)

// ErrUndefinedType is returned by generated code when the value type is not
// handled by generated serializers.
var ErrUndefinedType = errors.New("undefined type")

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
