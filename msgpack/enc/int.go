package enc

import (
	"math"

	"github.com/shamaton/msgpack/v2/def"
)

//func (e *Encoder) isPositiveFixInt64(v int64) bool {
//	return def.PositiveFixIntMin <= v && v <= def.PositiveFixIntMax
//}

func (e *Encoder) isNegativeFixInt64(v int64) bool {
	return def.NegativeFixintMin <= v && v <= def.NegativeFixintMax
}

// CalcInt checks value and returns data size that need.
func (e *Encoder) CalcInt(v int) int {
	return e.calcInt(int64(v))
}

// CalcInt8 checks value and returns data size that need.
func (e *Encoder) CalcInt8(v int8) int {
	return e.calcInt(int64(v))
}

// CalcInt16 checks value and returns data size that need.
func (e *Encoder) CalcInt16(v int16) int {
	return e.calcInt(int64(v))
}

// CalcInt32 checks value and returns data size that need.
func (e *Encoder) CalcInt32(v int32) int {
	return e.calcInt(int64(v))
}

// CalcInt64 checks value and returns data size that need.
func (e *Encoder) CalcInt64(v int64) int {
	return e.calcInt(v)
}

func (e *Encoder) calcInt(v int64) int {
	if v >= 0 {
		return e.calcUint(uint64(v))
	} else if e.isNegativeFixInt64(v) {
		// format code only
		return def.Byte1
	} else if v >= math.MinInt8 {
		return def.Byte1 + def.Byte1
	} else if v >= math.MinInt16 {
		return def.Byte1 + def.Byte2
	} else if v >= math.MinInt32 {
		return def.Byte1 + def.Byte4
	}
	return def.Byte1 + def.Byte8
}

// WriteInt sets the contents of v to the buffer.
func (e *Encoder) WriteInt(v int, offset int) int {
	return e.writeInt(int64(v), offset)
}

// WriteInt8 sets the contents of v to the buffer.
func (e *Encoder) WriteInt8(v int8, offset int) int {
	return e.writeInt(int64(v), offset)
}

// WriteInt16 sets the contents of v to the buffer.
func (e *Encoder) WriteInt16(v int16, offset int) int {
	return e.writeInt(int64(v), offset)
}

// WriteInt32 sets the contents of v to the buffer.
func (e *Encoder) WriteInt32(v int32, offset int) int {
	return e.writeInt(int64(v), offset)
}

// WriteInt64 sets the contents of v to the buffer.
func (e *Encoder) WriteInt64(v int64, offset int) int {
	return e.writeInt(v, offset)
}

func (e *Encoder) writeInt(v int64, offset int) int {
	if v >= 0 {
		offset = e.writeUint(uint64(v), offset)
	} else if e.isNegativeFixInt64(v) {
		offset = e.setByte1Int64(v, offset)
	} else if v >= math.MinInt8 {
		offset = e.setByte1Int(def.Int8, offset)
		offset = e.setByte1Int64(v, offset)
	} else if v >= math.MinInt16 {
		offset = e.setByte1Int(def.Int16, offset)
		offset = e.setByte2Int64(v, offset)
	} else if v >= math.MinInt32 {
		offset = e.setByte1Int(def.Int32, offset)
		offset = e.setByte4Int64(v, offset)
	} else {
		offset = e.setByte1Int(def.Int64, offset)
		offset = e.setByte8Int64(v, offset)
	}
	return offset
}
