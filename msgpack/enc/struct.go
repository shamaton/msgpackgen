package enc

import (
	"fmt"
	"math"

	"github.com/shamaton/msgpack/def"
)

// todo : コード生成側ですでにわかる
func (e *Encoder) CalcStructHeader(fieldNum int) (int, error) {
	ret := def.Byte1
	if e.asArray {
		if fieldNum <= 0x0f {
			// format code only
		} else if fieldNum <= math.MaxUint16 {
			ret += def.Byte2
		} else if uint(fieldNum) <= math.MaxUint32 {
			ret += def.Byte4
		} else {
			// not supported error
			return 0, fmt.Errorf("not support this array length : %d", fieldNum)
		}
	} else {

		if fieldNum <= 0x0f {
			// format code only
		} else if fieldNum <= math.MaxUint16 {
			ret += def.Byte2
		} else if uint(fieldNum) <= math.MaxUint32 {
			ret += def.Byte4
		} else {
			// not supported error
			return 0, fmt.Errorf("not support this array length : %d", fieldNum)
		}
	}
	return ret, nil
}

// todo : コード生成側ですでにわかる
func (e *Encoder) WriteStructHeader(fieldNum, offset int) int {
	if e.asArray {
		if fieldNum <= 0x0f {
			offset = e.setByte1Int(def.FixArray+fieldNum, offset)
		} else if fieldNum <= math.MaxUint16 {
			offset = e.setByte1Int(def.Array16, offset)
			offset = e.setByte2Int(fieldNum, offset)
		} else if uint(fieldNum) <= math.MaxUint32 {
			offset = e.setByte1Int(def.Array32, offset)
			offset = e.setByte4Int(fieldNum, offset)
		}
	} else {
		if fieldNum <= 0x0f {
			offset = e.setByte1Int(def.FixMap+fieldNum, offset)
		} else if fieldNum <= math.MaxUint16 {
			offset = e.setByte1Int(def.Map16, offset)
			offset = e.setByte2Int(fieldNum, offset)
		} else if uint(fieldNum) <= math.MaxUint32 {
			offset = e.setByte1Int(def.Map32, offset)
			offset = e.setByte4Int(fieldNum, offset)
		}
	}
	return offset
}
