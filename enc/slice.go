package encoding

import (
	"fmt"
	"math"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

// todo : can delete ??
func (e *Encoder) calcFixedSlice(rv reflect.Value) (int, bool) {
	size := 0

	switch sli := rv.Interface().(type) {
	case []int:
		for _, v := range sli {
			size += def.Byte1 + e.CalcInt(int64(v))
		}
		return size, true

	case []uint:
		for _, v := range sli {
			size += def.Byte1 + e.CalcUint(uint64(v))
		}
		return size, true

	case []string:
		for _, v := range sli {
			size += def.Byte1 + e.CalcString(v)
		}
		return size, true

	case []float32:
		for _, v := range sli {
			size += def.Byte1 + e.CalcFloat32(float64(v))
		}
		return size, true

	case []float64:
		for _, v := range sli {
			size += def.Byte1 + e.CalcFloat64(v)
		}
		return size, true

	case []bool:
		size += def.Byte1 * len(sli)
		return size, true

	case []int8:
		for _, v := range sli {
			size += def.Byte1 + e.CalcInt(int64(v))
		}
		return size, true

	case []int16:
		for _, v := range sli {
			size += def.Byte1 + e.CalcInt(int64(v))
		}
		return size, true

	case []int32:
		for _, v := range sli {
			size += def.Byte1 + e.CalcInt(int64(v))
		}
		return size, true

	case []int64:
		for _, v := range sli {
			size += def.Byte1 + e.CalcInt(v)
		}
		return size, true

	case []uint8:
		for _, v := range sli {
			size += def.Byte1 + e.CalcUint(uint64(v))
		}
		return size, true

	case []uint16:
		for _, v := range sli {
			size += def.Byte1 + e.CalcUint(uint64(v))
		}
		return size, true

	case []uint32:
		for _, v := range sli {
			size += def.Byte1 + e.CalcUint(uint64(v))
		}
		return size, true

	case []uint64:
		for _, v := range sli {
			size += def.Byte1 + e.CalcUint(v)
		}
		return size, true
	}

	return size, false
}

func (e *Encoder) CalcSliceLength(l int) (int, error) {

	if l <= 0x0f {
		// format code only
		return 0, nil
	} else if l <= math.MaxUint16 {
		return def.Byte2, nil
	} else if uint(l) <= math.MaxUint32 {
		return def.Byte4, nil
	}
	return 0, fmt.Errorf("not support this array length : %d", l)
}

func (e *Encoder) WriteSliceLength(l int, offset int) int {
	// format size
	if l <= 0x0f {
		offset = e.setByte1Int(def.FixArray+l, offset)
	} else if l <= math.MaxUint16 {
		offset = e.setByte1Int(def.Array16, offset)
		offset = e.setByte2Int(l, offset)
	} else if uint(l) <= math.MaxUint32 {
		offset = e.setByte1Int(def.Array32, offset)
		offset = e.setByte4Int(l, offset)
	}
	return offset
}

func (e *Encoder) writeFixedSlice(rv reflect.Value, offset int) (int, bool) {

	switch sli := rv.Interface().(type) {
	case []int:
		for _, v := range sli {
			offset = e.WriteInt(int64(v), offset)
		}
		return offset, true

	case []uint:
		for _, v := range sli {
			offset = e.WriteUint(uint64(v), offset)
		}
		return offset, true

	case []string:
		for _, v := range sli {
			offset = e.WriteString(v, offset)
		}
		return offset, true

	case []float32:
		for _, v := range sli {
			offset = e.WriteFloat32(v, offset)
		}
		return offset, true

	case []float64:
		for _, v := range sli {
			offset = e.WriteFloat64(v, offset)
		}
		return offset, true

	case []bool:
		for _, v := range sli {
			offset = e.WriteBool(v, offset)
		}
		return offset, true

	case []int8:
		for _, v := range sli {
			offset = e.WriteInt(int64(v), offset)
		}
		return offset, true

	case []int16:
		for _, v := range sli {
			offset = e.WriteInt(int64(v), offset)
		}
		return offset, true

	case []int32:
		for _, v := range sli {
			offset = e.WriteInt(int64(v), offset)
		}
		return offset, true

	case []int64:
		for _, v := range sli {
			offset = e.WriteInt(v, offset)
		}
		return offset, true

	case []uint8:
		for _, v := range sli {
			offset = e.WriteUint(uint64(v), offset)
		}
		return offset, true

	case []uint16:
		for _, v := range sli {
			offset = e.WriteUint(uint64(v), offset)
		}
		return offset, true

	case []uint32:
		for _, v := range sli {
			offset = e.WriteUint(uint64(v), offset)
		}
		return offset, true

	case []uint64:
		for _, v := range sli {
			offset = e.WriteUint(v, offset)
		}
		return offset, true
	}

	return offset, false
}
