package enc

import (
	"fmt"
	"math"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

// todo : can delete ??
func (e *Encoder) calcFixedMap(rv reflect.Value) (int, bool) {
	size := 0

	switch m := rv.Interface().(type) {
	case map[string]int:
		for k, v := range m {
			size += def.Byte1 + e.CalcString(k)
			size += def.Byte1 + e.CalcInt(int64(v))
		}
		return size, true

	case map[string]uint:
		for k, v := range m {
			size += def.Byte1 + e.CalcString(k)
			size += def.Byte1 + e.CalcUint(uint64(v))
		}
		return size, true

	case map[string]string:
		for k, v := range m {
			size += def.Byte1 + e.CalcString(k)
			size += def.Byte1 + e.CalcString(v)
		}
		return size, true

	case map[string]float32:
		for k := range m {
			size += def.Byte1 + e.CalcString(k)
			size += def.Byte1 + e.CalcFloat32(0)
		}
		return size, true

	case map[string]float64:
		for k := range m {
			size += def.Byte1 + e.CalcString(k)
			size += def.Byte1 + e.CalcFloat64(0)
		}
		return size, true

	case map[string]bool:
		for k := range m {
			size += def.Byte1 + e.CalcString(k)
			size += def.Byte1 /*+ e.CalcBool()*/
		}
		return size, true

	case map[string]int8:
		for k, v := range m {
			size += def.Byte1 + e.CalcString(k)
			size += def.Byte1 + e.CalcInt(int64(v))
		}
		return size, true
	case map[string]int16:
		for k, v := range m {
			size += def.Byte1 + e.CalcString(k)
			size += def.Byte1 + e.CalcInt(int64(v))
		}
		return size, true
	case map[string]int32:
		for k, v := range m {
			size += def.Byte1 + e.CalcString(k)
			size += def.Byte1 + e.CalcInt(int64(v))
		}
		return size, true
	case map[string]int64:
		for k, v := range m {
			size += def.Byte1 + e.CalcString(k)
			size += def.Byte1 + e.CalcInt(v)
		}
		return size, true
	case map[string]uint8:
		for k, v := range m {
			size += def.Byte1 + e.CalcString(k)
			size += def.Byte1 + e.CalcUint(uint64(v))
		}
		return size, true
	case map[string]uint16:
		for k, v := range m {
			size += def.Byte1 + e.CalcString(k)
			size += def.Byte1 + e.CalcUint(uint64(v))
		}
		return size, true
	case map[string]uint32:
		for k, v := range m {
			size += def.Byte1 + e.CalcString(k)
			size += def.Byte1 + e.CalcUint(uint64(v))
		}
		return size, true
	case map[string]uint64:
		for k, v := range m {
			size += def.Byte1 + e.CalcString(k)
			size += def.Byte1 + e.CalcUint(v)
		}
		return size, true

	case map[int]string:
		for k, v := range m {
			size += def.Byte1 + e.CalcInt(int64(k))
			size += def.Byte1 + e.CalcString(v)
		}
		return size, true
	case map[int]bool:
		for k := range m {
			size += def.Byte1 + e.CalcInt(int64(k))
			size += def.Byte1 /* + e.CalcBool()*/
		}
		return size, true

	case map[uint]string:
		for k, v := range m {
			size += def.Byte1 + e.CalcUint(uint64(k))
			size += def.Byte1 + e.CalcString(v)
		}
		return size, true
	case map[uint]bool:
		for k := range m {
			size += def.Byte1 + e.CalcUint(uint64(k))
			size += def.Byte1 /* + e.CalcBool()*/
		}
		return size, true

	case map[float32]string:
		for k, v := range m {
			size += def.Byte1 + e.CalcFloat32(float64(k))
			size += def.Byte1 + e.CalcString(v)
		}
		return size, true
	case map[float32]bool:
		for k := range m {
			size += def.Byte1 + e.CalcFloat32(float64(k))
			size += def.Byte1 /* + e.CalcBool()*/
		}
		return size, true

	case map[float64]string:
		for k, v := range m {
			size += def.Byte1 + e.CalcFloat64(k)
			size += def.Byte1 + e.CalcString(v)
		}
		return size, true
	case map[float64]bool:
		for k := range m {
			size += def.Byte1 + e.CalcFloat64(k)
			size += def.Byte1 /* + e.CalcBool()*/
		}
		return size, true

	case map[int8]string:
		for k, v := range m {
			size += def.Byte1 + e.CalcInt(int64(k))
			size += def.Byte1 + e.CalcString(v)
		}
		return size, true
	case map[int8]bool:
		for k := range m {
			size += def.Byte1 + e.CalcInt(int64(k))
			size += def.Byte1 /* + e.CalcBool()*/
		}
		return size, true
	case map[int16]string:
		for k, v := range m {
			size += def.Byte1 + e.CalcInt(int64(k))
			size += def.Byte1 + e.CalcString(v)
		}
		return size, true
	case map[int16]bool:
		for k := range m {
			size += def.Byte1 + e.CalcInt(int64(k))
			size += def.Byte1 /* + e.CalcBool()*/
		}
		return size, true
	case map[int32]string:
		for k, v := range m {
			size += def.Byte1 + e.CalcInt(int64(k))
			size += def.Byte1 + e.CalcString(v)
		}
		return size, true
	case map[int32]bool:
		for k := range m {
			size += def.Byte1 + e.CalcInt(int64(k))
			size += def.Byte1 /* + e.CalcBool()*/
		}
		return size, true
	case map[int64]string:
		for k, v := range m {
			size += def.Byte1 + e.CalcInt(k)
			size += def.Byte1 + e.CalcString(v)
		}
		return size, true
	case map[int64]bool:
		for k := range m {
			size += def.Byte1 + e.CalcInt(k)
			size += def.Byte1 /* + e.CalcBool()*/
		}
		return size, true

	case map[uint8]string:
		for k, v := range m {
			size += def.Byte1 + e.CalcUint(uint64(k))
			size += def.Byte1 + e.CalcString(v)
		}
		return size, true
	case map[uint8]bool:
		for k := range m {
			size += def.Byte1 + e.CalcUint(uint64(k))
			size += def.Byte1 /* + e.CalcBool()*/
		}
		return size, true
	case map[uint16]string:
		for k, v := range m {
			size += def.Byte1 + e.CalcUint(uint64(k))
			size += def.Byte1 + e.CalcString(v)
		}
		return size, true
	case map[uint16]bool:
		for k := range m {
			size += def.Byte1 + e.CalcUint(uint64(k))
			size += def.Byte1 /* + e.CalcBool()*/
		}
		return size, true
	case map[uint32]string:
		for k, v := range m {
			size += def.Byte1 + e.CalcUint(uint64(k))
			size += def.Byte1 + e.CalcString(v)
		}
		return size, true
	case map[uint32]bool:
		for k := range m {
			size += def.Byte1 + e.CalcUint(uint64(k))
			size += def.Byte1 /* + e.CalcBool()*/
		}
		return size, true
	case map[uint64]string:
		for k, v := range m {
			size += def.Byte1 + e.CalcUint(k)
			size += def.Byte1 + e.CalcString(v)
		}
		return size, true
	case map[uint64]bool:
		for k := range m {
			size += def.Byte1 + e.CalcUint(k)
			size += def.Byte1 /* + e.CalcBool()*/
		}
		return size, true

	}
	return size, false
}

func (e *Encoder) CalcMapLength(l int) (int, error) {
	ret := def.Byte1

	if l <= 0x0f {
		// do nothing
	} else if l <= math.MaxUint16 {
		ret += def.Byte2
	} else if uint(l) <= math.MaxUint32 {
		ret += def.Byte4
	} else {
		// not supported error
		return 0, fmt.Errorf("not support this map length : %d", l)
	}
	return ret, nil
}

func (e *Encoder) WriteMapLength(l int, offset int) int {

	// format
	if l <= 0x0f {
		offset = e.setByte1Int(def.FixMap+l, offset)
	} else if l <= math.MaxUint16 {
		offset = e.setByte1Int(def.Map16, offset)
		offset = e.setByte2Int(l, offset)
	} else if uint(l) <= math.MaxUint32 {
		offset = e.setByte1Int(def.Map32, offset)
		offset = e.setByte4Int(l, offset)
	}
	return offset
}

func (e *Encoder) writeFixedMap(rv reflect.Value, offset int) (int, bool) {
	switch m := rv.Interface().(type) {
	case map[string]int:
		for k, v := range m {
			offset = e.WriteString(k, offset)
			offset = e.WriteInt(int64(v), offset)
		}
		return offset, true

	case map[string]uint:
		for k, v := range m {
			offset = e.WriteString(k, offset)
			offset = e.WriteUint(uint64(v), offset)
		}
		return offset, true

	case map[string]float32:
		for k, v := range m {
			offset = e.WriteString(k, offset)
			offset = e.WriteFloat32(v, offset)
		}
		return offset, true

	case map[string]float64:
		for k, v := range m {
			offset = e.WriteString(k, offset)
			offset = e.WriteFloat64(v, offset)
		}
		return offset, true

	case map[string]bool:
		for k, v := range m {
			offset = e.WriteString(k, offset)
			offset = e.WriteBool(v, offset)
		}
		return offset, true

	case map[string]string:
		for k, v := range m {
			offset = e.WriteString(k, offset)
			offset = e.WriteString(v, offset)
		}
		return offset, true

	case map[string]int8:
		for k, v := range m {
			offset = e.WriteString(k, offset)
			offset = e.WriteInt(int64(v), offset)
		}
		return offset, true
	case map[string]int16:
		for k, v := range m {
			offset = e.WriteString(k, offset)
			offset = e.WriteInt(int64(v), offset)
		}
		return offset, true
	case map[string]int32:
		for k, v := range m {
			offset = e.WriteString(k, offset)
			offset = e.WriteInt(int64(v), offset)
		}
		return offset, true
	case map[string]int64:
		for k, v := range m {
			offset = e.WriteString(k, offset)
			offset = e.WriteInt(int64(v), offset)
		}
		return offset, true

	case map[string]uint8:
		for k, v := range m {
			offset = e.WriteString(k, offset)
			offset = e.WriteUint(uint64(v), offset)
		}
		return offset, true
	case map[string]uint16:
		for k, v := range m {
			offset = e.WriteString(k, offset)
			offset = e.WriteUint(uint64(v), offset)
		}
		return offset, true
	case map[string]uint32:
		for k, v := range m {
			offset = e.WriteString(k, offset)
			offset = e.WriteUint(uint64(v), offset)
		}
		return offset, true
	case map[string]uint64:
		for k, v := range m {
			offset = e.WriteString(k, offset)
			offset = e.WriteUint(uint64(v), offset)
		}
		return offset, true

	case map[int]string:
		for k, v := range m {
			offset = e.WriteInt(int64(k), offset)
			offset = e.WriteString(v, offset)
		}
		return offset, true
	case map[int]bool:
		for k, v := range m {
			offset = e.WriteInt(int64(k), offset)
			offset = e.WriteBool(v, offset)
		}
		return offset, true

	case map[uint]string:
		for k, v := range m {
			offset = e.WriteUint(uint64(k), offset)
			offset = e.WriteString(v, offset)
		}
		return offset, true
	case map[uint]bool:
		for k, v := range m {
			offset = e.WriteUint(uint64(k), offset)
			offset = e.WriteBool(v, offset)
		}
		return offset, true

	case map[float32]string:
		for k, v := range m {
			offset = e.WriteFloat32(k, offset)
			offset = e.WriteString(v, offset)
		}
		return offset, true
	case map[float32]bool:
		for k, v := range m {
			offset = e.WriteFloat32(k, offset)
			offset = e.WriteBool(v, offset)
		}
		return offset, true

	case map[float64]string:
		for k, v := range m {
			offset = e.WriteFloat64(k, offset)
			offset = e.WriteString(v, offset)
		}
		return offset, true
	case map[float64]bool:
		for k, v := range m {
			offset = e.WriteFloat64(k, offset)
			offset = e.WriteBool(v, offset)
		}
		return offset, true

	case map[int8]string:
		for k, v := range m {
			offset = e.WriteInt(int64(k), offset)
			offset = e.WriteString(v, offset)
		}
		return offset, true
	case map[int8]bool:
		for k, v := range m {
			offset = e.WriteInt(int64(k), offset)
			offset = e.WriteBool(v, offset)
		}
		return offset, true
	case map[int16]string:
		for k, v := range m {
			offset = e.WriteInt(int64(k), offset)
			offset = e.WriteString(v, offset)
		}
		return offset, true
	case map[int16]bool:
		for k, v := range m {
			offset = e.WriteInt(int64(k), offset)
			offset = e.WriteBool(v, offset)
		}
		return offset, true
	case map[int32]string:
		for k, v := range m {
			offset = e.WriteInt(int64(k), offset)
			offset = e.WriteString(v, offset)
		}
		return offset, true
	case map[int32]bool:
		for k, v := range m {
			offset = e.WriteInt(int64(k), offset)
			offset = e.WriteBool(v, offset)
		}
		return offset, true
	case map[int64]string:
		for k, v := range m {
			offset = e.WriteInt(k, offset)
			offset = e.WriteString(v, offset)
		}
		return offset, true
	case map[int64]bool:
		for k, v := range m {
			offset = e.WriteInt(k, offset)
			offset = e.WriteBool(v, offset)
		}
		return offset, true

	case map[uint8]string:
		for k, v := range m {
			offset = e.WriteUint(uint64(k), offset)
			offset = e.WriteString(v, offset)
		}
		return offset, true
	case map[uint8]bool:
		for k, v := range m {
			offset = e.WriteUint(uint64(k), offset)
			offset = e.WriteBool(v, offset)
		}
		return offset, true
	case map[uint16]string:
		for k, v := range m {
			offset = e.WriteUint(uint64(k), offset)
			offset = e.WriteString(v, offset)
		}
		return offset, true
	case map[uint16]bool:
		for k, v := range m {
			offset = e.WriteUint(uint64(k), offset)
			offset = e.WriteBool(v, offset)
		}
		return offset, true
	case map[uint32]string:
		for k, v := range m {
			offset = e.WriteUint(uint64(k), offset)
			offset = e.WriteString(v, offset)
		}
		return offset, true
	case map[uint32]bool:
		for k, v := range m {
			offset = e.WriteUint(uint64(k), offset)
			offset = e.WriteBool(v, offset)
		}
		return offset, true
	case map[uint64]string:
		for k, v := range m {
			offset = e.WriteUint(k, offset)
			offset = e.WriteString(v, offset)
		}
		return offset, true
	case map[uint64]bool:
		for k, v := range m {
			offset = e.WriteUint(k, offset)
			offset = e.WriteBool(v, offset)
		}
		return offset, true

	}
	return offset, false
}
