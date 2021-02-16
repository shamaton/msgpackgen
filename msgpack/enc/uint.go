package enc

import (
	"math"

	"github.com/shamaton/msgpack/v2/def"
)

//func (e *Encoder) isPositiveFixUint64(v uint64) bool {
//	return def.PositiveFixIntMin <= v && v <= def.PositiveFixIntMax
//}

// CalcUint check value and returns data size that need.
func (e *Encoder) CalcUint(v uint) int {
	return e.calcUint(uint64(v))
}

// CalcUint8 check value and returns data size that need.
func (e *Encoder) CalcUint8(v uint8) int {
	return e.calcUint(uint64(v))
}

// CalcUint16 check value and returns data size that need.
func (e *Encoder) CalcUint16(v uint16) int {
	return e.calcUint(uint64(v))
}

// CalcUint32 check value and returns data size that need.
func (e *Encoder) CalcUint32(v uint32) int {
	return e.calcUint(uint64(v))
}

// CalcUint64 check value and returns data size that need.
func (e *Encoder) CalcUint64(v uint64) int {
	return e.calcUint(v)
}

func (e *Encoder) calcUint(v uint64) int {
	if v <= math.MaxInt8 {
		// format code only
		return def.Byte1
	} else if v <= math.MaxUint8 {
		return def.Byte1 + def.Byte1
	} else if v <= math.MaxUint16 {
		return def.Byte1 + def.Byte2
	} else if v <= math.MaxUint32 {
		return def.Byte1 + def.Byte4
	}
	return def.Byte1 + def.Byte8
}

// WriteUint sets the contents of v to the buffer.
func (e *Encoder) WriteUint(v uint, offset int) int {
	return e.writeUint(uint64(v), offset)
}

// WriteUint8 sets the contents of v to the buffer.
func (e *Encoder) WriteUint8(v uint8, offset int) int {
	return e.writeUint(uint64(v), offset)
}

// WriteUint16 sets the contents of v to the buffer.
func (e *Encoder) WriteUint16(v uint16, offset int) int {
	return e.writeUint(uint64(v), offset)
}

// WriteUint32 sets the contents of v to the buffer.
func (e *Encoder) WriteUint32(v uint32, offset int) int {
	return e.writeUint(uint64(v), offset)
}

// WriteUint64 sets the contents of v to the buffer.
func (e *Encoder) WriteUint64(v uint64, offset int) int {
	return e.writeUint(v, offset)
}

func (e *Encoder) writeUint(v uint64, offset int) int {
	if v <= math.MaxInt8 {
		offset = e.setByte1Uint64(v, offset)
	} else if v <= math.MaxUint8 {
		offset = e.setByte1Int(def.Uint8, offset)
		offset = e.setByte1Uint64(v, offset)
	} else if v <= math.MaxUint16 {
		offset = e.setByte1Int(def.Uint16, offset)
		offset = e.setByte2Uint64(v, offset)
	} else if v <= math.MaxUint32 {
		offset = e.setByte1Int(def.Uint32, offset)
		offset = e.setByte4Uint64(v, offset)
	} else {
		offset = e.setByte1Int(def.Uint64, offset)
		offset = e.setByte8Uint64(v, offset)
	}
	return offset
}
