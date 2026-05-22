package enc

import (
	"math"
	"time"

	"github.com/shamaton/msgpack/v3/def"
)

// EnsureLen extends buf so that len(buf) is at least targetLen.
func EnsureLen(buf []byte, targetLen int) []byte {
	if targetLen < 0 {
		panic("msgpackgen: negative target length")
	}
	if targetLen <= len(buf) {
		return buf
	}
	if targetLen <= cap(buf) {
		return buf[:targetLen]
	}
	return append(buf, make([]byte, targetLen-len(buf))...)
}

// RequireAt extends buf so that writing extra bytes at offset is valid.
func RequireAt(buf []byte, offset, extra int) []byte {
	if offset < 0 || extra < 0 {
		panic("msgpackgen: negative offset or extra length")
	}
	return EnsureLen(buf, offset+extra)
}

// CalcInt checks value and returns data size that need.
func CalcInt(v int) int {
	var e Encoder
	return e.CalcInt(v)
}

// CalcInt8 checks value and returns data size that need.
func CalcInt8(v int8) int {
	var e Encoder
	return e.CalcInt8(v)
}

// CalcInt16 checks value and returns data size that need.
func CalcInt16(v int16) int {
	var e Encoder
	return e.CalcInt16(v)
}

// CalcInt32 checks value and returns data size that need.
func CalcInt32(v int32) int {
	var e Encoder
	return e.CalcInt32(v)
}

// CalcInt64 checks value and returns data size that need.
func CalcInt64(v int64) int {
	var e Encoder
	return e.CalcInt64(v)
}

// CalcUint checks value and returns data size that need.
func CalcUint(v uint) int {
	var e Encoder
	return e.CalcUint(v)
}

// CalcUint8 checks value and returns data size that need.
func CalcUint8(v uint8) int {
	var e Encoder
	return e.CalcUint8(v)
}

// CalcUint16 checks value and returns data size that need.
func CalcUint16(v uint16) int {
	var e Encoder
	return e.CalcUint16(v)
}

// CalcUint32 checks value and returns data size that need.
func CalcUint32(v uint32) int {
	var e Encoder
	return e.CalcUint32(v)
}

// CalcUint64 checks value and returns data size that need.
func CalcUint64(v uint64) int {
	var e Encoder
	return e.CalcUint64(v)
}

// CalcFloat32 returns data size that need.
func CalcFloat32(v float32) int {
	var e Encoder
	return e.CalcFloat32(v)
}

// CalcFloat64 returns data size that need.
func CalcFloat64(v float64) int {
	var e Encoder
	return e.CalcFloat64(v)
}

// CalcComplex64 returns data size that need.
func CalcComplex64(v complex64) int {
	var e Encoder
	return e.CalcComplex64(v)
}

// CalcComplex128 returns data size that need.
func CalcComplex128(v complex128) int {
	var e Encoder
	return e.CalcComplex128(v)
}

// CalcBool returns data size that need.
func CalcBool(v bool) int {
	var e Encoder
	return e.CalcBool(v)
}

// CalcNil returns data size that need.
func CalcNil() int {
	var e Encoder
	return e.CalcNil()
}

// CalcByte returns data size that need.
func CalcByte(b byte) int {
	var e Encoder
	return e.CalcByte(b)
}

// CalcRune checks value and returns data size that need.
func CalcRune(v rune) int {
	var e Encoder
	return e.CalcRune(v)
}

// CalcString checks value and returns data size that need.
func CalcString(v string) int {
	var e Encoder
	return e.CalcString(v)
}

// CalcStringFix returns data size that need.
func CalcStringFix(length int) int {
	var e Encoder
	return e.CalcStringFix(length)
}

// CalcString8 returns data size that need.
func CalcString8(length int) int {
	var e Encoder
	return e.CalcString8(length)
}

// CalcString16 returns data size that need.
func CalcString16(length int) int {
	var e Encoder
	return e.CalcString16(length)
}

// CalcString32 returns data size that need.
func CalcString32(length int) int {
	var e Encoder
	return e.CalcString32(length)
}

// CalcSliceLength checks values and returns data size that need.
func CalcSliceLength(l int, isChildTypeByte bool) (int, error) {
	var e Encoder
	return e.CalcSliceLength(l, isChildTypeByte)
}

// CalcMapLength checks value and returns data size that need.
func CalcMapLength(l int) (int, error) {
	var e Encoder
	return e.CalcMapLength(l)
}

// CalcStructHeaderFix returns data size that need.
func CalcStructHeaderFix(fieldNum int) int {
	var e Encoder
	return e.CalcStructHeaderFix(fieldNum)
}

// CalcStructHeader16 returns data size that need.
func CalcStructHeader16(fieldNum int) int {
	var e Encoder
	return e.CalcStructHeader16(fieldNum)
}

// CalcStructHeader32 returns data size that need.
func CalcStructHeader32(fieldNum int) int {
	var e Encoder
	return e.CalcStructHeader32(fieldNum)
}

// CalcTime checks value and returns data size that need.
func CalcTime(t time.Time) int {
	t = t.UTC()
	var e Encoder
	return e.CalcTime(t)
}

// WriteIntTo sets the contents of v to buf at offset.
func WriteIntTo(buf []byte, v int, offset int) int {
	return writeIntTo(buf, int64(v), offset)
}

// WriteInt8To sets the contents of v to buf at offset.
func WriteInt8To(buf []byte, v int8, offset int) int {
	return writeIntTo(buf, int64(v), offset)
}

// WriteInt16To sets the contents of v to buf at offset.
func WriteInt16To(buf []byte, v int16, offset int) int {
	return writeIntTo(buf, int64(v), offset)
}

// WriteInt32To sets the contents of v to buf at offset.
func WriteInt32To(buf []byte, v int32, offset int) int {
	return writeIntTo(buf, int64(v), offset)
}

// WriteInt64To sets the contents of v to buf at offset.
func WriteInt64To(buf []byte, v int64, offset int) int {
	return writeIntTo(buf, v, offset)
}

// WriteUintTo sets the contents of v to buf at offset.
func WriteUintTo(buf []byte, v uint, offset int) int {
	return writeUintTo(buf, uint64(v), offset)
}

// WriteUint8To sets the contents of v to buf at offset.
func WriteUint8To(buf []byte, v uint8, offset int) int {
	return writeUintTo(buf, uint64(v), offset)
}

// WriteUint16To sets the contents of v to buf at offset.
func WriteUint16To(buf []byte, v uint16, offset int) int {
	return writeUintTo(buf, uint64(v), offset)
}

// WriteUint32To sets the contents of v to buf at offset.
func WriteUint32To(buf []byte, v uint32, offset int) int {
	return writeUintTo(buf, uint64(v), offset)
}

// WriteUint64To sets the contents of v to buf at offset.
func WriteUint64To(buf []byte, v uint64, offset int) int {
	return writeUintTo(buf, v, offset)
}

// WriteFloat32To sets the contents of v to buf at offset.
func WriteFloat32To(buf []byte, v float32, offset int) int {
	offset = setByte1IntTo(buf, def.Float32, offset)
	offset = setByte4Uint64To(buf, uint64(math.Float32bits(v)), offset)
	return offset
}

// WriteFloat64To sets the contents of v to buf at offset.
func WriteFloat64To(buf []byte, v float64, offset int) int {
	offset = setByte1IntTo(buf, def.Float64, offset)
	offset = setByte8Uint64To(buf, math.Float64bits(v), offset)
	return offset
}

// WriteComplex64To sets the contents of v to buf at offset.
func WriteComplex64To(buf []byte, v complex64, offset int) int {
	offset = setByte1IntTo(buf, def.Fixext8, offset)
	offset = setByte1IntTo(buf, int(def.ComplexTypeCode()), offset)
	offset = setByte4Uint64To(buf, uint64(math.Float32bits(real(v))), offset)
	offset = setByte4Uint64To(buf, uint64(math.Float32bits(imag(v))), offset)
	return offset
}

// WriteComplex128To sets the contents of v to buf at offset.
func WriteComplex128To(buf []byte, v complex128, offset int) int {
	offset = setByte1IntTo(buf, def.Fixext16, offset)
	offset = setByte1IntTo(buf, int(def.ComplexTypeCode()), offset)
	offset = setByte8Uint64To(buf, math.Float64bits(real(v)), offset)
	offset = setByte8Uint64To(buf, math.Float64bits(imag(v)), offset)
	return offset
}

// WriteBoolTo sets the contents of v to buf at offset.
func WriteBoolTo(buf []byte, v bool, offset int) int {
	if v {
		return setByte1IntTo(buf, def.True, offset)
	}
	return setByte1IntTo(buf, def.False, offset)
}

// WriteNilTo sets nil to buf at offset.
func WriteNilTo(buf []byte, offset int) int {
	return setByte1IntTo(buf, def.Nil, offset)
}

// WriteByteTo sets the contents of b to buf at offset.
func WriteByteTo(buf []byte, b byte, offset int) int {
	return setByteTo(buf, b, offset)
}

// WriteRuneTo sets the contents of v to buf at offset.
func WriteRuneTo(buf []byte, v rune, offset int) int {
	return writeIntTo(buf, int64(v), offset)
}

// WriteStringTo sets the contents of str to buf at offset.
func WriteStringTo(buf []byte, str string, offset int) int {
	l := len(str)
	if l < 32 {
		return WriteStringFixTo(buf, str, l, offset)
	} else if l <= math.MaxUint8 {
		return WriteString8To(buf, str, l, offset)
	} else if l <= math.MaxUint16 {
		return WriteString16To(buf, str, l, offset)
	}
	return WriteString32To(buf, str, l, offset)
}

// WriteStringFixTo sets the contents of str to buf at offset.
func WriteStringFixTo(buf []byte, str string, length, offset int) int {
	offset = setByte1IntTo(buf, def.FixStr+length, offset)
	offset += copy(buf[offset:], str)
	return offset
}

// WriteString8To sets the contents of str to buf at offset.
func WriteString8To(buf []byte, str string, length, offset int) int {
	offset = setByte1IntTo(buf, def.Str8, offset)
	offset = setByte1IntTo(buf, length, offset)
	offset += copy(buf[offset:], str)
	return offset
}

// WriteString16To sets the contents of str to buf at offset.
func WriteString16To(buf []byte, str string, length, offset int) int {
	offset = setByte1IntTo(buf, def.Str16, offset)
	offset = setByte2IntTo(buf, length, offset)
	offset += copy(buf[offset:], str)
	return offset
}

// WriteString32To sets the contents of str to buf at offset.
func WriteString32To(buf []byte, str string, length, offset int) int {
	offset = setByte1IntTo(buf, def.Str32, offset)
	offset = setByte4IntTo(buf, length, offset)
	offset += copy(buf[offset:], str)
	return offset
}

// WriteSliceLengthTo sets the contents of l to buf at offset.
func WriteSliceLengthTo(buf []byte, l int, offset int, isChildTypeByte bool) int {
	if isChildTypeByte {
		return writeByteSliceLengthTo(buf, l, offset)
	}

	if l <= 0x0f {
		offset = setByte1IntTo(buf, def.FixArray+l, offset)
	} else if l <= math.MaxUint16 {
		offset = setByte1IntTo(buf, def.Array16, offset)
		offset = setByte2IntTo(buf, l, offset)
	} else if uint(l) <= math.MaxUint32 {
		offset = setByte1IntTo(buf, def.Array32, offset)
		offset = setByte4IntTo(buf, l, offset)
	}
	return offset
}

// WriteMapLengthTo sets the contents of l to buf at offset.
func WriteMapLengthTo(buf []byte, l int, offset int) int {
	if l <= 0x0f {
		offset = setByte1IntTo(buf, def.FixMap+l, offset)
	} else if l <= math.MaxUint16 {
		offset = setByte1IntTo(buf, def.Map16, offset)
		offset = setByte2IntTo(buf, l, offset)
	} else if uint(l) <= math.MaxUint32 {
		offset = setByte1IntTo(buf, def.Map32, offset)
		offset = setByte4IntTo(buf, l, offset)
	}
	return offset
}

// WriteStructHeaderFixAsArrayTo sets num of fields to buf as array type.
func WriteStructHeaderFixAsArrayTo(buf []byte, fieldNum, offset int) int {
	return setByte1IntTo(buf, def.FixArray+fieldNum, offset)
}

// WriteStructHeader16AsArrayTo sets num of fields to buf as array type.
func WriteStructHeader16AsArrayTo(buf []byte, fieldNum, offset int) int {
	offset = setByte1IntTo(buf, def.Array16, offset)
	return setByte2IntTo(buf, fieldNum, offset)
}

// WriteStructHeader32AsArrayTo sets num of fields to buf as array type.
func WriteStructHeader32AsArrayTo(buf []byte, fieldNum, offset int) int {
	offset = setByte1IntTo(buf, def.Array32, offset)
	return setByte4IntTo(buf, fieldNum, offset)
}

// WriteStructHeaderFixAsMapTo sets num of fields to buf as map type.
func WriteStructHeaderFixAsMapTo(buf []byte, fieldNum, offset int) int {
	return setByte1IntTo(buf, def.FixMap+fieldNum, offset)
}

// WriteStructHeader16AsMapTo sets num of fields to buf as map type.
func WriteStructHeader16AsMapTo(buf []byte, fieldNum, offset int) int {
	offset = setByte1IntTo(buf, def.Map16, offset)
	return setByte2IntTo(buf, fieldNum, offset)
}

// WriteStructHeader32AsMapTo sets num of fields to buf as map type.
func WriteStructHeader32AsMapTo(buf []byte, fieldNum, offset int) int {
	offset = setByte1IntTo(buf, def.Map32, offset)
	return setByte4IntTo(buf, fieldNum, offset)
}

// WriteTimeTo sets the contents of t to buf at offset.
func WriteTimeTo(buf []byte, t time.Time, offset int) int {
	t = t.UTC()
	secs := uint64(t.Unix())
	if secs>>34 == 0 {
		data := uint64(t.Nanosecond())<<34 | secs
		if data&0xffffffff00000000 == 0 {
			offset = setByte1IntTo(buf, def.Fixext4, offset)
			offset = setByte1IntTo(buf, def.TimeStamp, offset)
			offset = setByte4Uint64To(buf, data, offset)
			return offset
		}

		offset = setByte1IntTo(buf, def.Fixext8, offset)
		offset = setByte1IntTo(buf, def.TimeStamp, offset)
		offset = setByte8Uint64To(buf, data, offset)
		return offset
	}

	offset = setByte1IntTo(buf, def.Ext8, offset)
	offset = setByte1IntTo(buf, 12, offset)
	offset = setByte1IntTo(buf, def.TimeStamp, offset)
	offset = setByte4IntTo(buf, t.Nanosecond(), offset)
	offset = setByte8Uint64To(buf, secs, offset)
	return offset
}

func writeIntTo(buf []byte, v int64, offset int) int {
	if v >= 0 {
		offset = writeUintTo(buf, uint64(v), offset)
	} else if def.NegativeFixintMin <= v && v <= def.NegativeFixintMax {
		offset = setByte1Int64To(buf, v, offset)
	} else if v >= math.MinInt8 {
		offset = setByte1IntTo(buf, def.Int8, offset)
		offset = setByte1Int64To(buf, v, offset)
	} else if v >= math.MinInt16 {
		offset = setByte1IntTo(buf, def.Int16, offset)
		offset = setByte2Int64To(buf, v, offset)
	} else if v >= math.MinInt32 {
		offset = setByte1IntTo(buf, def.Int32, offset)
		offset = setByte4Int64To(buf, v, offset)
	} else {
		offset = setByte1IntTo(buf, def.Int64, offset)
		offset = setByte8Int64To(buf, v, offset)
	}
	return offset
}

func writeUintTo(buf []byte, v uint64, offset int) int {
	if v <= math.MaxInt8 {
		offset = setByte1Uint64To(buf, v, offset)
	} else if v <= math.MaxUint8 {
		offset = setByte1IntTo(buf, def.Uint8, offset)
		offset = setByte1Uint64To(buf, v, offset)
	} else if v <= math.MaxUint16 {
		offset = setByte1IntTo(buf, def.Uint16, offset)
		offset = setByte2Uint64To(buf, v, offset)
	} else if v <= math.MaxUint32 {
		offset = setByte1IntTo(buf, def.Uint32, offset)
		offset = setByte4Uint64To(buf, v, offset)
	} else {
		offset = setByte1IntTo(buf, def.Uint64, offset)
		offset = setByte8Uint64To(buf, v, offset)
	}
	return offset
}

func writeByteSliceLengthTo(buf []byte, l int, offset int) int {
	if l <= math.MaxUint8 {
		offset = setByte1IntTo(buf, def.Bin8, offset)
		offset = setByte1IntTo(buf, l, offset)
	} else if l <= math.MaxUint16 {
		offset = setByte1IntTo(buf, def.Bin16, offset)
		offset = setByte2IntTo(buf, l, offset)
	} else if uint(l) <= math.MaxUint32 {
		offset = setByte1IntTo(buf, def.Bin32, offset)
		offset = setByte4IntTo(buf, l, offset)
	}
	return offset
}

func setByte1Int64To(buf []byte, value int64, offset int) int {
	buf[offset] = byte(value)
	return offset + 1
}

func setByte2Int64To(buf []byte, value int64, offset int) int {
	buf[offset+0] = byte(value >> 8)
	buf[offset+1] = byte(value)
	return offset + 2
}

func setByte4Int64To(buf []byte, value int64, offset int) int {
	buf[offset+0] = byte(value >> 24)
	buf[offset+1] = byte(value >> 16)
	buf[offset+2] = byte(value >> 8)
	buf[offset+3] = byte(value)
	return offset + 4
}

func setByte8Int64To(buf []byte, value int64, offset int) int {
	buf[offset] = byte(value >> 56)
	buf[offset+1] = byte(value >> 48)
	buf[offset+2] = byte(value >> 40)
	buf[offset+3] = byte(value >> 32)
	buf[offset+4] = byte(value >> 24)
	buf[offset+5] = byte(value >> 16)
	buf[offset+6] = byte(value >> 8)
	buf[offset+7] = byte(value)
	return offset + 8
}

func setByte1Uint64To(buf []byte, value uint64, offset int) int {
	buf[offset] = byte(value)
	return offset + 1
}

func setByte2Uint64To(buf []byte, value uint64, offset int) int {
	buf[offset] = byte(value >> 8)
	buf[offset+1] = byte(value)
	return offset + 2
}

func setByte4Uint64To(buf []byte, value uint64, offset int) int {
	buf[offset] = byte(value >> 24)
	buf[offset+1] = byte(value >> 16)
	buf[offset+2] = byte(value >> 8)
	buf[offset+3] = byte(value)
	return offset + 4
}

func setByte8Uint64To(buf []byte, value uint64, offset int) int {
	buf[offset] = byte(value >> 56)
	buf[offset+1] = byte(value >> 48)
	buf[offset+2] = byte(value >> 40)
	buf[offset+3] = byte(value >> 32)
	buf[offset+4] = byte(value >> 24)
	buf[offset+5] = byte(value >> 16)
	buf[offset+6] = byte(value >> 8)
	buf[offset+7] = byte(value)
	return offset + 8
}

func setByte1IntTo(buf []byte, code, offset int) int {
	buf[offset] = byte(code)
	return offset + 1
}

func setByte2IntTo(buf []byte, value int, offset int) int {
	buf[offset] = byte(value >> 8)
	buf[offset+1] = byte(value)
	return offset + 2
}

func setByte4IntTo(buf []byte, value int, offset int) int {
	buf[offset] = byte(value >> 24)
	buf[offset+1] = byte(value >> 16)
	buf[offset+2] = byte(value >> 8)
	buf[offset+3] = byte(value)
	return offset + 4
}

func setByteTo(buf []byte, b byte, offset int) int {
	buf[offset] = b
	return offset + 1
}
