package enc

import (
	"github.com/shamaton/msgpack/v2/def"
)

func (e *Encoder) CalcStructHeaderFix(fieldNum int) int {
	return def.Byte1
}

func (e *Encoder) CalcStructHeader16(fieldNum int) int {
	return def.Byte1 + def.Byte2
}

func (e *Encoder) CalcStructHeader32(fieldNum int) int {
	return def.Byte1 + def.Byte4
}

func (e *Encoder) WriteStructHeaderFixAsArray(fieldNum, offset int) int {
	offset = e.setByte1Int(def.FixArray+fieldNum, offset)
	return offset
}

func (e *Encoder) WriteStructHeader16AsArray(fieldNum, offset int) int {
	offset = e.setByte1Int(def.Array16, offset)
	offset = e.setByte2Int(fieldNum, offset)
	return offset
}

func (e *Encoder) WriteStructHeader32AsArray(fieldNum, offset int) int {
	offset = e.setByte1Int(def.Array32, offset)
	offset = e.setByte4Int(fieldNum, offset)
	return offset
}

func (e *Encoder) WriteStructHeaderFixAsMap(fieldNum, offset int) int {
	offset = e.setByte1Int(def.FixMap+fieldNum, offset)
	return offset
}

func (e *Encoder) WriteStructHeader16AsMap(fieldNum, offset int) int {
	offset = e.setByte1Int(def.Map16, offset)
	offset = e.setByte2Int(fieldNum, offset)
	return offset
}

func (e *Encoder) WriteStructHeader32AsMap(fieldNum, offset int) int {
	offset = e.setByte1Int(def.Map32, offset)
	offset = e.setByte4Int(fieldNum, offset)
	return offset
}
