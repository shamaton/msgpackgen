package enc

import (
	"math"

	"github.com/shamaton/msgpack/v3/def"
)

// CalcInt checks value and returns data size that need.
func CalcInt(v int) int {
	return calcIntSize(int64(v))
}

// CalcIntMax returns the maximum data size that an int value can need.
func CalcIntMax(v int) int {
	return def.Byte1 + def.Byte8
}

// CalcInt8 checks value and returns data size that need.
func CalcInt8(v int8) int {
	return calcIntSize(int64(v))
}

// CalcInt8Max returns the maximum data size that an int8 value can need.
func CalcInt8Max(v int8) int {
	return def.Byte1 + def.Byte1
}

// CalcInt16 checks value and returns data size that need.
func CalcInt16(v int16) int {
	return calcIntSize(int64(v))
}

// CalcInt16Max returns the maximum data size that an int16 value can need.
func CalcInt16Max(v int16) int {
	return def.Byte1 + def.Byte2
}

// CalcInt32 checks value and returns data size that need.
func CalcInt32(v int32) int {
	return calcIntSize(int64(v))
}

// CalcInt32Max returns the maximum data size that an int32 value can need.
func CalcInt32Max(v int32) int {
	return def.Byte1 + def.Byte4
}

// CalcInt64 checks value and returns data size that need.
func CalcInt64(v int64) int {
	return calcIntSize(v)
}

// CalcInt64Max returns the maximum data size that an int64 value can need.
func CalcInt64Max(v int64) int {
	return def.Byte1 + def.Byte8
}

// WriteInt sets the contents of v to buf at offset.
func WriteInt(buf []byte, v int, offset int) int {
	return writeInt(buf, int64(v), offset)
}

// WriteInt8 sets the contents of v to buf at offset.
func WriteInt8(buf []byte, v int8, offset int) int {
	return writeInt(buf, int64(v), offset)
}

// WriteInt16 sets the contents of v to buf at offset.
func WriteInt16(buf []byte, v int16, offset int) int {
	return writeInt(buf, int64(v), offset)
}

// WriteInt32 sets the contents of v to buf at offset.
func WriteInt32(buf []byte, v int32, offset int) int {
	return writeInt(buf, int64(v), offset)
}

// WriteInt64 sets the contents of v to buf at offset.
func WriteInt64(buf []byte, v int64, offset int) int {
	return writeInt(buf, v, offset)
}

func calcIntSize(v int64) int {
	if v >= 0 {
		return calcUintSize(uint64(v))
	} else if def.NegativeFixintMin <= v && v <= def.NegativeFixintMax {
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

func writeInt(buf []byte, v int64, offset int) int {
	if v >= 0 {
		offset = writeUint(buf, uint64(v), offset)
	} else if def.NegativeFixintMin <= v && v <= def.NegativeFixintMax {
		offset = setByte1Int64(buf, v, offset)
	} else if v >= math.MinInt8 {
		offset = setByte1Int(buf, def.Int8, offset)
		offset = setByte1Int64(buf, v, offset)
	} else if v >= math.MinInt16 {
		offset = setByte1Int(buf, def.Int16, offset)
		offset = setByte2Int64(buf, v, offset)
	} else if v >= math.MinInt32 {
		offset = setByte1Int(buf, def.Int32, offset)
		offset = setByte4Int64(buf, v, offset)
	} else {
		offset = setByte1Int(buf, def.Int64, offset)
		offset = setByte8Int64(buf, v, offset)
	}
	return offset
}
