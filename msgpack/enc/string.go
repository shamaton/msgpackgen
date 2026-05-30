package enc

import (
	"math"

	"github.com/shamaton/msgpack/v3/def"
)

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

// WriteString sets the contents of str to buf at offset.
func WriteString(buf []byte, str string, offset int) int {
	l := len(str)
	if l < 32 {
		return WriteStringFix(buf, str, l, offset)
	} else if l <= math.MaxUint8 {
		return WriteString8(buf, str, l, offset)
	} else if l <= math.MaxUint16 {
		return WriteString16(buf, str, l, offset)
	}
	return WriteString32(buf, str, l, offset)
}

// WriteStringFix sets the contents of str to buf at offset.
func WriteStringFix(buf []byte, str string, length, offset int) int {
	offset = setByte1Int(buf, def.FixStr+length, offset)
	offset += copy(buf[offset:offset+length], str)
	return offset
}

// WriteString8 sets the contents of str to buf at offset.
func WriteString8(buf []byte, str string, length, offset int) int {
	offset = setByte1Int(buf, def.Str8, offset)
	offset = setByte1Int(buf, length, offset)
	offset += copy(buf[offset:offset+length], str)
	return offset
}

// WriteString16 sets the contents of str to buf at offset.
func WriteString16(buf []byte, str string, length, offset int) int {
	offset = setByte1Int(buf, def.Str16, offset)
	offset = setByte2Int(buf, length, offset)
	offset += copy(buf[offset:offset+length], str)
	return offset
}

// WriteString32 sets the contents of str to buf at offset.
func WriteString32(buf []byte, str string, length, offset int) int {
	offset = setByte1Int(buf, def.Str32, offset)
	offset = setByte4Int(buf, length, offset)
	offset += copy(buf[offset:offset+length], str)
	return offset
}
