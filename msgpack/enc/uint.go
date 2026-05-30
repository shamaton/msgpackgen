package enc

import (
	"math"

	"github.com/shamaton/msgpack/v3/def"
)

// CalcUint checks value and returns data size that need.
func CalcUint(v uint) int {
	return calcUintSize(uint64(v))
}

// CalcUintMax returns the maximum data size that a uint value can need.
func CalcUintMax(v uint) int {
	return def.Byte1 + def.Byte8
}

// CalcUint8 checks value and returns data size that need.
func CalcUint8(v uint8) int {
	return calcUintSize(uint64(v))
}

// CalcUint8Max returns the maximum data size that a uint8 value can need.
func CalcUint8Max(v uint8) int {
	return def.Byte1 + def.Byte1
}

// CalcUint16 checks value and returns data size that need.
func CalcUint16(v uint16) int {
	return calcUintSize(uint64(v))
}

// CalcUint16Max returns the maximum data size that a uint16 value can need.
func CalcUint16Max(v uint16) int {
	return def.Byte1 + def.Byte2
}

// CalcUint32 checks value and returns data size that need.
func CalcUint32(v uint32) int {
	return calcUintSize(uint64(v))
}

// CalcUint32Max returns the maximum data size that a uint32 value can need.
func CalcUint32Max(v uint32) int {
	return def.Byte1 + def.Byte4
}

// CalcUint64 checks value and returns data size that need.
func CalcUint64(v uint64) int {
	return calcUintSize(v)
}

// CalcUint64Max returns the maximum data size that a uint64 value can need.
func CalcUint64Max(v uint64) int {
	return def.Byte1 + def.Byte8
}

// WriteUint sets the contents of v to buf at offset.
func WriteUint(buf []byte, v uint, offset int) int {
	return writeUint(buf, uint64(v), offset)
}

// WriteUint8 sets the contents of v to buf at offset.
func WriteUint8(buf []byte, v uint8, offset int) int {
	return writeUint(buf, uint64(v), offset)
}

// WriteUint16 sets the contents of v to buf at offset.
func WriteUint16(buf []byte, v uint16, offset int) int {
	return writeUint(buf, uint64(v), offset)
}

// WriteUint32 sets the contents of v to buf at offset.
func WriteUint32(buf []byte, v uint32, offset int) int {
	return writeUint(buf, uint64(v), offset)
}

// WriteUint64 sets the contents of v to buf at offset.
func WriteUint64(buf []byte, v uint64, offset int) int {
	return writeUint(buf, v, offset)
}

func calcUintSize(v uint64) int {
	if v <= math.MaxInt8 {
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

func writeUint(buf []byte, v uint64, offset int) int {
	if v <= math.MaxInt8 {
		offset = setByte1Uint64(buf, v, offset)
	} else if v <= math.MaxUint8 {
		offset = setByte1Int(buf, def.Uint8, offset)
		offset = setByte1Uint64(buf, v, offset)
	} else if v <= math.MaxUint16 {
		offset = setByte1Int(buf, def.Uint16, offset)
		offset = setByte2Uint64(buf, v, offset)
	} else if v <= math.MaxUint32 {
		offset = setByte1Int(buf, def.Uint32, offset)
		offset = setByte4Uint64(buf, v, offset)
	} else {
		offset = setByte1Int(buf, def.Uint64, offset)
		offset = setByte8Uint64(buf, v, offset)
	}
	return offset
}
