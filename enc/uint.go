package encoding

import (
	"math"

	"github.com/shamaton/msgpack/def"
)

func (e *Encoder) isPositiveFixUint64(v uint64) bool {
	return def.PositiveFixIntMin <= v && v <= def.PositiveFixIntMax
}

func (e *Encoder) CalcUint(v uint64) int {
	if v <= math.MaxInt8 {
		// format code only
		return def.Byte1
	} else if v <= math.MaxUint8 {
		return def.Byte1 + def.Byte1
	} else if v <= math.MaxUint16 {
		return def.Byte1 + def.Byte2
	} else if v <= math.MaxUint32 {
		return def.Byte1 + def.Byte4
	}
	return def.Byte1 + def.Byte8
}

func (e *Encoder) WriteUint(v uint64, offset int) int {
	if v <= math.MaxInt8 {
		offset = e.setByte1Uint64(v, offset)
	} else if v <= math.MaxUint8 {
		offset = e.setByte1Int(def.Uint8, offset)
		offset = e.setByte1Uint64(v, offset)
	} else if v <= math.MaxUint16 {
		offset = e.setByte1Int(def.Uint16, offset)
		offset = e.setByte2Uint64(v, offset)
	} else if v <= math.MaxUint32 {
		offset = e.setByte1Int(def.Uint32, offset)
		offset = e.setByte4Uint64(v, offset)
	} else {
		offset = e.setByte1Int(def.Uint64, offset)
		offset = e.setByte8Uint64(v, offset)
	}
	return offset
}
