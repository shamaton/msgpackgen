package enc

import (
	"math"

	"github.com/shamaton/msgpack/v2/def"
)

// CalcComplex64 returns data size that need.
func (e *Encoder) CalcComplex64(v complex64) int {
	return def.Byte1 + def.Byte1 + def.Byte8
}

// CalcComplex128 returns data size that need.
func (e *Encoder) CalcComplex128(v complex128) int {
	return def.Byte1 + def.Byte1 + def.Byte16
}

// WriteComplex64 sets the contents of v to the buffer.
func (e *Encoder) WriteComplex64(v complex64, offset int) int {
	offset = e.setByte1Int(def.Fixext8, offset)
	offset = e.setByte1Int(int(def.ComplexTypeCode()), offset)
	offset = e.setByte4Uint64(uint64(math.Float32bits(real(v))), offset)
	offset = e.setByte4Uint64(uint64(math.Float32bits(imag(v))), offset)
	return offset
}

// WriteComplex128 sets the contents of v to the buffer.
func (e *Encoder) WriteComplex128(v complex128, offset int) int {
	offset = e.setByte1Int(def.Fixext16, offset)
	offset = e.setByte1Int(int(def.ComplexTypeCode()), offset)
	offset = e.setByte8Uint64(math.Float64bits(real(v)), offset)
	offset = e.setByte8Uint64(math.Float64bits(imag(v)), offset)
	return offset
}
