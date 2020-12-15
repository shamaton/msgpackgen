package encoding

import "github.com/shamaton/msgpack/def"

func (e *Encoder) calcBool() int {
	return 0
}

func (e *Encoder) writeBool(v bool, offset int) int {
	if v {
		offset = e.setByte1Int(def.True, offset)
	} else {
		offset = e.setByte1Int(def.False, offset)
	}
	return offset
}
