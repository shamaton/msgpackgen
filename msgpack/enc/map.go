package enc

import (
	"fmt"
	"math"

	"github.com/shamaton/msgpack/v2/def"
)

// CalcMapLength checks value and returns data size that need.
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

// WriteMapLength sets the contents of l to the buffer.
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
