package enc

import (
	"fmt"
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

// CalcBool returns data size that need.
func CalcBool(v bool) int {
	return def.Byte1
}

// CalcBoolMax returns the maximum data size that a bool value can need.
func CalcBoolMax(v bool) int {
	return def.Byte1
}

// CalcNil returns data size that need.
func CalcNil() int {
	return def.Byte1
}

// CalcByte returns data size that need.
func CalcByte(b byte) int {
	return def.Byte1
}

// CalcByteMax returns the maximum data size that a byte value can need.
func CalcByteMax(b byte) int {
	return def.Byte1
}

// CalcRune checks value and returns data size that need.
func CalcRune(v rune) int {
	return calcIntSize(int64(v))
}

// CalcRuneMax returns the maximum data size that a rune value can need.
func CalcRuneMax(v rune) int {
	return def.Byte1 + def.Byte4
}

// CalcString checks value and returns data size that need.
func CalcString(v string) int {
	l := len(v)
	if l < 32 {
		return CalcStringFix(l)
	} else if l <= math.MaxUint8 {
		return CalcString8(l)
	} else if l <= math.MaxUint16 {
		return CalcString16(l)
	}
	return CalcString32(l)
}

// CalcStringMax returns the maximum data size that a string value can need.
func CalcStringMax(v string) int {
	return def.Byte1 + def.Byte4 + len(v)
}

// CalcStringFix returns data size that need.
func CalcStringFix(length int) int {
	return def.Byte1 + length
}

// CalcString8 returns data size that need.
func CalcString8(length int) int {
	return def.Byte1 + def.Byte1 + length
}

// CalcString16 returns data size that need.
func CalcString16(length int) int {
	return def.Byte1 + def.Byte2 + length
}

// CalcString32 returns data size that need.
func CalcString32(length int) int {
	return def.Byte1 + def.Byte4 + length
}

// CalcSliceLength checks values and returns data size that need.
func CalcSliceLength(l int, isChildTypeByte bool) (int, error) {
	var e Encoder
	return e.CalcSliceLength(l, isChildTypeByte)
}

// CalcSliceLengthMax returns the maximum data size that a slice header can need.
func CalcSliceLengthMax(l int, isChildTypeByte bool) (int, error) {
	if uint(l) > math.MaxUint32 {
		return 0, fmt.Errorf("not support this array length : %d", l)
	}
	return def.Byte1 + def.Byte4, nil
}

// CalcMapLength checks value and returns data size that need.
func CalcMapLength(l int) (int, error) {
	var e Encoder
	return e.CalcMapLength(l)
}

// CalcMapLengthMax returns the maximum data size that a map header can need.
func CalcMapLengthMax(l int) (int, error) {
	if uint(l) > math.MaxUint32 {
		return 0, fmt.Errorf("not support this map length : %d", l)
	}
	return def.Byte1 + def.Byte4, nil
}

// CalcStructHeaderFix returns data size that need.
func CalcStructHeaderFix(fieldNum int) int {
	return def.Byte1
}

// CalcStructHeader16 returns data size that need.
func CalcStructHeader16(fieldNum int) int {
	return def.Byte1 + def.Byte2
}

// CalcStructHeader32 returns data size that need.
func CalcStructHeader32(fieldNum int) int {
	return def.Byte1 + def.Byte4
}

// CalcTime checks value and returns data size that need.
func CalcTime(t time.Time) int {
	t = t.UTC()
	secs := uint64(t.Unix())
	if secs>>34 == 0 {
		data := uint64(t.Nanosecond())<<34 | secs
		if data&0xffffffff00000000 == 0 {
			return def.Byte1 + def.Byte1 + def.Byte4
		}
		return def.Byte1 + def.Byte1 + def.Byte8
	}

	return def.Byte1 + def.Byte1 + def.Byte1 + def.Byte4 + def.Byte8
}

// CalcTimeMax returns the maximum data size that a time value can need.
func CalcTimeMax(t time.Time) int {
	return def.Byte1 + def.Byte1 + def.Byte1 + def.Byte4 + def.Byte8
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
	offset += copy(buf[offset:offset+length], str)
	return offset
}

// WriteString8To sets the contents of str to buf at offset.
func WriteString8To(buf []byte, str string, length, offset int) int {
	offset = setByte1IntTo(buf, def.Str8, offset)
	offset = setByte1IntTo(buf, length, offset)
	offset += copy(buf[offset:offset+length], str)
	return offset
}

// WriteString16To sets the contents of str to buf at offset.
func WriteString16To(buf []byte, str string, length, offset int) int {
	offset = setByte1IntTo(buf, def.Str16, offset)
	offset = setByte2IntTo(buf, length, offset)
	offset += copy(buf[offset:offset+length], str)
	return offset
}

// WriteString32To sets the contents of str to buf at offset.
func WriteString32To(buf []byte, str string, length, offset int) int {
	offset = setByte1IntTo(buf, def.Str32, offset)
	offset = setByte4IntTo(buf, length, offset)
	offset += copy(buf[offset:offset+length], str)
	return offset
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
	_ = buf[offset+1]
	buf[offset+0] = byte(value >> 8)
	buf[offset+1] = byte(value)
	return offset + 2
}

func setByte4Int64To(buf []byte, value int64, offset int) int {
	_ = buf[offset+3]
	buf[offset+0] = byte(value >> 24)
	buf[offset+1] = byte(value >> 16)
	buf[offset+2] = byte(value >> 8)
	buf[offset+3] = byte(value)
	return offset + 4
}

func setByte8Int64To(buf []byte, value int64, offset int) int {
	_ = buf[offset+7]
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
	_ = buf[offset+1]
	buf[offset] = byte(value >> 8)
	buf[offset+1] = byte(value)
	return offset + 2
}

func setByte4Uint64To(buf []byte, value uint64, offset int) int {
	_ = buf[offset+3]
	buf[offset] = byte(value >> 24)
	buf[offset+1] = byte(value >> 16)
	buf[offset+2] = byte(value >> 8)
	buf[offset+3] = byte(value)
	return offset + 4
}

func setByte8Uint64To(buf []byte, value uint64, offset int) int {
	_ = buf[offset+7]
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
	_ = buf[offset+1]
	buf[offset] = byte(value >> 8)
	buf[offset+1] = byte(value)
	return offset + 2
}

func setByte4IntTo(buf []byte, value int, offset int) int {
	_ = buf[offset+3]
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
