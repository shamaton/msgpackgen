package enc

import "github.com/shamaton/msgpack/v2/def"

func (e *Encoder) CalcBool(v bool) int {
	return def.Byte1
}

func (e *Encoder) WriteBool(v bool, offset int) int {
	if v {
		offset = e.setByte1Int(def.True, offset)
	} else {
		offset = e.setByte1Int(def.False, offset)
	}
	return offset
}
