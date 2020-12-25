package enc

import (
	"fmt"
	"math"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

var typeByte = reflect.TypeOf(byte(0))

func (e *Encoder) isByteSlice(rv reflect.Value) bool {
	return rv.Type().Elem() == typeByte
}

func (e *Encoder) CalcByteSlice(l int) (int, error) {
	if l <= math.MaxUint8 {
		return def.Byte1 + def.Byte1 + l, nil
	} else if l <= math.MaxUint16 {
		return def.Byte1 + def.Byte2 + l, nil
	} else if uint(l) <= math.MaxUint32 {
		return def.Byte1 + def.Byte4 + l, nil
	}
	// not supported error
	return 0, fmt.Errorf("not support this array length : %d", l)
}

func (e *Encoder) WriteByteSliceLength(l int, offset int) int {
	if l <= math.MaxUint8 {
		offset = e.setByte1Int(def.Bin8, offset)
		offset = e.setByte1Int(l, offset)
	} else if l <= math.MaxUint16 {
		offset = e.setByte1Int(def.Bin16, offset)
		offset = e.setByte2Int(l, offset)
	} else if uint(l) <= math.MaxUint32 {
		offset = e.setByte1Int(def.Bin32, offset)
		offset = e.setByte4Int(l, offset)
	}
	return offset
}

func (e *Encoder) WriteByteSlice(bytes []byte, offset int) int {
	return e.setBytes(bytes, offset)
}
