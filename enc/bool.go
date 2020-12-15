package encoding

import "github.com/shamaton/msgpack/def"

func (e *Encoder) CalcBool() int {
	return 0
}

func (e *Encoder) WriteBool(v bool, offset int) int {
	if v {
		offset = e.setByte1Int(def.True, offset)
	} else {
		offset = e.setByte1Int(def.False, offset)
	}
	return offset
}
