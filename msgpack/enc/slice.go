package enc

import (
	"fmt"
	"math"

	"github.com/shamaton/msgpack/v3/def"
)

// CalcSliceLength checks values and returns data size that need.
func CalcSliceLength(l int, isChildTypeByte bool) (int, error) {
	if isChildTypeByte {
		return calcByteSliceLength(l)
	}

	if l <= 0x0f {
		return def.Byte1, nil
	} else if l <= math.MaxUint16 {
		return def.Byte1 + def.Byte2, nil
	} else if uint(l) <= math.MaxUint32 {
		return def.Byte1 + def.Byte4, nil
	}
	return 0, fmt.Errorf("not support this array length : %d", l)
}

// CalcSliceLengthMax returns the maximum data size that a slice header can need.
func CalcSliceLengthMax(l int, isChildTypeByte bool) (int, error) {
	if uint(l) > math.MaxUint32 {
		return 0, fmt.Errorf("not support this array length : %d", l)
	}
	return def.Byte1 + def.Byte4, nil
}

// WriteSliceLength sets the contents of l to buf at offset.
func WriteSliceLength(buf []byte, l int, offset int, isChildTypeByte bool) int {
	if isChildTypeByte {
		return writeByteSliceLength(buf, l, offset)
	}

	if l <= 0x0f {
		offset = setByte1Int(buf, def.FixArray+l, offset)
	} else if l <= math.MaxUint16 {
		offset = setByte1Int(buf, def.Array16, offset)
		offset = setByte2Int(buf, l, offset)
	} else if uint(l) <= math.MaxUint32 {
		offset = setByte1Int(buf, def.Array32, offset)
		offset = setByte4Int(buf, l, offset)
	}
	return offset
}
