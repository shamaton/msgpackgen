package enc

import (
	"math"

	"github.com/shamaton/msgpack/v2/def"
)

// CalcFloat32 returns data size that need.
func (e *Encoder) CalcFloat32(v float32) int {
	return def.Byte1 + def.Byte4
}

// CalcFloat64 returns data size that need.
func (e *Encoder) CalcFloat64(v float64) int {
	return def.Byte1 + def.Byte8
}

// WriteFloat32 sets the contents of v to the buffer.
func (e *Encoder) WriteFloat32(v float32, offset int) int {
	offset = e.setByte1Int(def.Float32, offset)
	offset = e.setByte4Uint64(uint64(math.Float32bits(v)), offset)
	return offset
}

// WriteFloat64 sets the contents of v to the buffer.
func (e *Encoder) WriteFloat64(v float64, offset int) int {
	offset = e.setByte1Int(def.Float64, offset)
	offset = e.setByte8Uint64(math.Float64bits(v), offset)
	return offset
}
