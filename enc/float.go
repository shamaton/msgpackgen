package enc

import (
	"math"

	"github.com/shamaton/msgpack/def"
)

func (e *Encoder) CalcFloat32(v float64) int {
	return def.Byte1 + def.Byte4
}

func (e *Encoder) CalcFloat64(v float64) int {
	return def.Byte1 + def.Byte8
}

func (e *Encoder) WriteFloat32(v float32, offset int) int {
	offset = e.setByte1Int(def.Float32, offset)
	offset = e.setByte4Uint64(uint64(math.Float32bits(v)), offset)
	return offset
}

func (e *Encoder) WriteFloat64(v float64, offset int) int {
	offset = e.setByte1Int(def.Float64, offset)
	offset = e.setByte8Uint64(math.Float64bits(v), offset)
	return offset
}
