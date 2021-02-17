package enc

import (
	"fmt"
	"math"

	"github.com/shamaton/msgpack/v2/def"
)

// CalcSliceLength checks values and returns data size that need.
func (e *Encoder) CalcSliceLength(l int, isChildTypeByte bool) (int, error) {

	if isChildTypeByte {
		return e.calcByteSlice(l)
	}

	if l <= 0x0f {
		// format code only
		return def.Byte1, nil
	} else if l <= math.MaxUint16 {
		return def.Byte1 + def.Byte2, nil
	} else if uint(l) <= math.MaxUint32 {
		return def.Byte1 + def.Byte4, nil
	}
	return 0, fmt.Errorf("not support this array length : %d", l)
}

// WriteSliceLength sets the contents of l to the buffer.
func (e *Encoder) WriteSliceLength(l int, offset int, isChildTypeByte bool) int {
	if isChildTypeByte {
		return e.writeByteSliceLength(l, offset)
	}

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
