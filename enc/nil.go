package encoding

import "github.com/shamaton/msgpack/def"

func (e *Encoder) WriteNil(offset int) int {
	offset = e.setByte1Int(def.Nil, offset)
	return offset
}
