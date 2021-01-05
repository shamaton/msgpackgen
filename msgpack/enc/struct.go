package enc

import (
	"fmt"
	"math"

	"github.com/shamaton/msgpack/def"
)

// todo : delete
func (e *Encoder) CalcStructHeader(fieldNum int) (int, error) {
	if fieldNum <= 0x0f {
		return e.CalcStructHeaderFix(fieldNum), nil
	} else if fieldNum <= math.MaxUint16 {
		return e.CalcStructHeader16(fieldNum), nil
	} else if uint(fieldNum) <= math.MaxUint32 {
		return e.CalcStructHeader32(fieldNum), nil
	}
	return 0, fmt.Errorf("not support this array length : %d", fieldNum)
}

func (e *Encoder) CalcStructHeaderFix(fieldNum int) int {
	return def.Byte1
}

func (e *Encoder) CalcStructHeader16(fieldNum int) int {
	return def.Byte1 + def.Byte2
}

func (e *Encoder) CalcStructHeader32(fieldNum int) int {
	return def.Byte1 + def.Byte4
}

// todo : delete
func (e *Encoder) WriteStructHeaderAsArray(fieldNum, offset int) int {
	if fieldNum <= 0x0f {
		return e.WriteStructHeaderFixAsArray(fieldNum, offset)
	} else if fieldNum <= math.MaxUint16 {
		return e.WriteStructHeader16AsArray(fieldNum, offset)
	} else {
		return e.WriteStructHeader32AsArray(fieldNum, offset)
	}
}

// todo : delete
func (e *Encoder) WriteStructHeaderAsMap(fieldNum, offset int) int {
	if fieldNum <= 0x0f {
		return e.WriteStructHeaderFixAsMap(fieldNum, offset)
	} else if fieldNum <= math.MaxUint16 {
		return e.WriteStructHeader16AsMap(fieldNum, offset)
	} else {
		return e.WriteStructHeader32AsMap(fieldNum, offset)
	}
}

func (e *Encoder) WriteStructHeaderFixAsArray(fieldNum, offset int) int {
	offset = e.setByte1Int(def.FixArray+fieldNum, offset)
	return offset
}

func (e *Encoder) WriteStructHeader16AsArray(fieldNum, offset int) int {
	offset = e.setByte1Int(def.Array16, offset)
	offset = e.setByte2Int(fieldNum, offset)
	return offset
}

func (e *Encoder) WriteStructHeader32AsArray(fieldNum, offset int) int {
	offset = e.setByte1Int(def.Array32, offset)
	offset = e.setByte4Int(fieldNum, offset)
	return offset
}

func (e *Encoder) WriteStructHeaderFixAsMap(fieldNum, offset int) int {
	offset = e.setByte1Int(def.FixMap+fieldNum, offset)
	return offset
}

func (e *Encoder) WriteStructHeader16AsMap(fieldNum, offset int) int {
	offset = e.setByte1Int(def.Map16, offset)
	offset = e.setByte2Int(fieldNum, offset)
	return offset
}

func (e *Encoder) WriteStructHeader32AsMap(fieldNum, offset int) int {
	offset = e.setByte1Int(def.Map32, offset)
	offset = e.setByte4Int(fieldNum, offset)
	return offset
}
