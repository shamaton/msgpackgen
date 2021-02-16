package enc

import (
	"fmt"
	"math"

	"github.com/shamaton/msgpack/v2/def"
)

func (e *Encoder) calcByteSlice(l int) (int, error) {
	if l <= math.MaxUint8 {
		return def.Byte1 + def.Byte1, nil
	} else if l <= math.MaxUint16 {
		return def.Byte1 + def.Byte2, nil
	} else if uint(l) <= math.MaxUint32 {
		return def.Byte1 + def.Byte4, nil
	}
	// not supported error
	return 0, fmt.Errorf("not support this array length : %d", l)
}

func (e *Encoder) writeByteSliceLength(l int, offset int) int {
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

// CalcByte returns data size that need.
func (e *Encoder) CalcByte(b byte) int {
	return def.Byte1
}

// WriteByte sets the contents of v to the buffer.
func (e *Encoder) WriteByte(b byte, offset int) int {
	return e.setByte(b, offset)
}
