package enc

import "github.com/shamaton/msgpack/v3/def"

// CalcBool returns data size that need.
func CalcBool(v bool) int {
	return def.Byte1
}

// CalcBoolMax returns the maximum data size that a bool value can need.
func CalcBoolMax(v bool) int {
	return def.Byte1
}

// WriteBool sets the contents of v to buf at offset.
func WriteBool(buf []byte, v bool, offset int) int {
	if v {
		return setByte1Int(buf, def.True, offset)
	}
	return setByte1Int(buf, def.False, offset)
}
