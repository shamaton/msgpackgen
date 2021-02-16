package enc

import "github.com/shamaton/msgpack/v2/def"

// CalcBool returns data size that need.
func (e *Encoder) CalcBool(v bool) int {
	return def.Byte1
}

// WriteBool sets the contents of v to the buffer.
func (e *Encoder) WriteBool(v bool, offset int) int {
	if v {
		offset = e.setByte1Int(def.True, offset)
	} else {
		offset = e.setByte1Int(def.False, offset)
	}
	return offset
}
