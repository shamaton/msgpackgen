package enc

import (
	"math"

	"github.com/shamaton/msgpack/v3/def"
)

// CalcFloat32 returns data size that need.
func CalcFloat32(v float32) int {
	return def.Byte1 + def.Byte4
}

// CalcFloat32Max returns the maximum data size that a float32 value can need.
func CalcFloat32Max(v float32) int {
	return def.Byte1 + def.Byte4
}

// CalcFloat64 returns data size that need.
func CalcFloat64(v float64) int {
	return def.Byte1 + def.Byte8
}

// CalcFloat64Max returns the maximum data size that a float64 value can need.
func CalcFloat64Max(v float64) int {
	return def.Byte1 + def.Byte8
}

// WriteFloat32 sets the contents of v to buf at offset.
func WriteFloat32(buf []byte, v float32, offset int) int {
	offset = setByte1Int(buf, def.Float32, offset)
	offset = setByte4Uint64(buf, uint64(math.Float32bits(v)), offset)
	return offset
}

// WriteFloat64 sets the contents of v to buf at offset.
func WriteFloat64(buf []byte, v float64, offset int) int {
	offset = setByte1Int(buf, def.Float64, offset)
	offset = setByte8Uint64(buf, math.Float64bits(v), offset)
	return offset
}
