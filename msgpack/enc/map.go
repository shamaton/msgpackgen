package enc

import (
	"fmt"
	"math"

	"github.com/shamaton/msgpack/v3/def"
)

// CalcMapLength checks value and returns data size that need.
func CalcMapLength(l int) (int, error) {
	if l <= 0x0f {
		return def.Byte1, nil
	} else if l <= math.MaxUint16 {
		return def.Byte1 + def.Byte2, nil
	} else if uint(l) <= math.MaxUint32 {
		return def.Byte1 + def.Byte4, nil
	}
	return 0, fmt.Errorf("not support this map length : %d", l)
}

// CalcMapLengthMax returns the maximum data size that a map header can need.
func CalcMapLengthMax(l int) (int, error) {
	if uint(l) > math.MaxUint32 {
		return 0, fmt.Errorf("not support this map length : %d", l)
	}
	return def.Byte1 + def.Byte4, nil
}

// WriteMapLength sets the contents of l to buf at offset.
func WriteMapLength(buf []byte, l int, offset int) int {
	if l <= 0x0f {
		offset = setByte1Int(buf, def.FixMap+l, offset)
	} else if l <= math.MaxUint16 {
		offset = setByte1Int(buf, def.Map16, offset)
		offset = setByte2Int(buf, l, offset)
	} else if uint(l) <= math.MaxUint32 {
		offset = setByte1Int(buf, def.Map32, offset)
		offset = setByte4Int(buf, l, offset)
	}
	return offset
}
