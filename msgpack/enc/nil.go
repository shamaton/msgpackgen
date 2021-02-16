package enc

import "github.com/shamaton/msgpack/v2/def"

// CalcNil returns data size that need.
func (e *Encoder) CalcNil() int {
	return def.Byte1
}

// WriteNil sets the contents of v to the buffer.
func (e *Encoder) WriteNil(offset int) int {
	offset = e.setByte1Int(def.Nil, offset)
	return offset
}
