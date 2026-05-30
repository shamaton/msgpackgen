package enc

import (
	"math"

	"github.com/shamaton/msgpack/v3/def"
)

// CalcComplex64 returns data size that need.
func CalcComplex64(v complex64) int {
	return def.Byte1 + def.Byte1 + def.Byte8
}

// CalcComplex64Max returns the maximum data size that a complex64 value can need.
func CalcComplex64Max(v complex64) int {
	return def.Byte1 + def.Byte1 + def.Byte8
}

// CalcComplex128 returns data size that need.
func CalcComplex128(v complex128) int {
	return def.Byte1 + def.Byte1 + def.Byte16
}

// CalcComplex128Max returns the maximum data size that a complex128 value can need.
func CalcComplex128Max(v complex128) int {
	return def.Byte1 + def.Byte1 + def.Byte16
}

// WriteComplex64 sets the contents of v to buf at offset.
func WriteComplex64(buf []byte, v complex64, offset int) int {
	offset = setByte1Int(buf, def.Fixext8, offset)
	offset = setByte1Int(buf, int(def.ComplexTypeCode()), offset)
	offset = setByte4Uint64(buf, uint64(math.Float32bits(real(v))), offset)
	offset = setByte4Uint64(buf, uint64(math.Float32bits(imag(v))), offset)
	return offset
}

// WriteComplex128 sets the contents of v to buf at offset.
func WriteComplex128(buf []byte, v complex128, offset int) int {
	offset = setByte1Int(buf, def.Fixext16, offset)
	offset = setByte1Int(buf, int(def.ComplexTypeCode()), offset)
	offset = setByte8Uint64(buf, math.Float64bits(real(v)), offset)
	offset = setByte8Uint64(buf, math.Float64bits(imag(v)), offset)
	return offset
}
