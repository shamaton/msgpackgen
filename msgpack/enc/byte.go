package enc

import (
	"fmt"
	"math"

	"github.com/shamaton/msgpack/v3/def"
)

// CalcByte returns data size that need.
func CalcByte(b byte) int {
	return def.Byte1
}

// CalcByteMax returns the maximum data size that a byte value can need.
func CalcByteMax(b byte) int {
	return def.Byte1
}

// WriteByte sets the contents of b to buf at offset.
func WriteByte(buf []byte, b byte, offset int) int {
	return setByte(buf, b, offset)
}

func calcByteSliceLength(l int) (int, error) {
	if l <= math.MaxUint8 {
		return def.Byte1 + def.Byte1, nil
	} else if l <= math.MaxUint16 {
		return def.Byte1 + def.Byte2, nil
	} else if uint(l) <= math.MaxUint32 {
		return def.Byte1 + def.Byte4, nil
	}
	return 0, fmt.Errorf("not support this array length : %d", l)
}

func writeByteSliceLength(buf []byte, l int, offset int) int {
	if l <= math.MaxUint8 {
		offset = setByte1Int(buf, def.Bin8, offset)
		offset = setByte1Int(buf, l, offset)
	} else if l <= math.MaxUint16 {
		offset = setByte1Int(buf, def.Bin16, offset)
		offset = setByte2Int(buf, l, offset)
	} else if uint(l) <= math.MaxUint32 {
		offset = setByte1Int(buf, def.Bin32, offset)
		offset = setByte4Int(buf, l, offset)
	}
	return offset
}
