package enc

import (
	"math"

	"github.com/shamaton/msgpack/def"
)

func (e *Encoder) CalcString(v string) int {
	l := len(v)
	if l < 32 {
		return e.CalcStringFix(l)
	} else if l <= math.MaxUint8 {
		return e.CalcString8(l)
	} else if l <= math.MaxUint16 {
		return e.CalcString16(l)
	}
	return e.CalcString32(l)
	// NOTE : length over uint32
}

func (e *Encoder) CalcStringFix(length int) int {
	return def.Byte1 + length
}

func (e *Encoder) CalcString8(length int) int {
	return def.Byte1 + def.Byte1 + length
}

func (e *Encoder) CalcString16(length int) int {
	return def.Byte1 + def.Byte2 + length
}

func (e *Encoder) CalcString32(length int) int {
	return def.Byte1 + def.Byte4 + length
}

func (e *Encoder) WriteString(str string, offset int) int {
	l := len(str)
	if l < 32 {
		return e.WriteStringFix(str, l, offset)
	} else if l <= math.MaxUint8 {
		return e.WriteString8(str, l, offset)
	} else if l <= math.MaxUint16 {
		return e.WriteString16(str, l, offset)
	} else {
		return e.WriteString32(str, l, offset)
	}
}

func (e *Encoder) WriteStringFix(str string, length, offset int) int {
	offset = e.setByte1Int(def.FixStr+length, offset)
	offset += copy(e.d[offset:], str)
	return offset
}

func (e *Encoder) WriteString8(str string, length, offset int) int {
	offset = e.setByte1Int(def.Str8, offset)
	offset = e.setByte1Int(length, offset)
	offset += copy(e.d[offset:], str)
	return offset
}

func (e *Encoder) WriteString16(str string, length, offset int) int {
	offset = e.setByte1Int(def.Str16, offset)
	offset = e.setByte2Int(length, offset)
	offset += copy(e.d[offset:], str)
	return offset
}

func (e *Encoder) WriteString32(str string, length, offset int) int {
	offset = e.setByte1Int(def.Str32, offset)
	offset = e.setByte4Int(length, offset)
	offset += copy(e.d[offset:], str)
	return offset
}
