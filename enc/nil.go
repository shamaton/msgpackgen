package enc

import "github.com/shamaton/msgpack/def"

func (e *Encoder) CalcNil() int {
	return def.Byte1
}

func (e *Encoder) WriteNil(offset int) int {
	offset = e.setByte1Int(def.Nil, offset)
	return offset
}
