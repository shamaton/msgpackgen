package enc

import (
	"github.com/shamaton/msgpack/v2/def"
)

// CalcStructHeaderFix returns data size that need.
func (e *Encoder) CalcStructHeaderFix(fieldNum int) int {
	return def.Byte1
}

// CalcStructHeader16 returns data size that need.
func (e *Encoder) CalcStructHeader16(fieldNum int) int {
	return def.Byte1 + def.Byte2
}

// CalcStructHeader32 returns data size that need.
func (e *Encoder) CalcStructHeader32(fieldNum int) int {
	return def.Byte1 + def.Byte4
}

// WriteStructHeaderFixAsArray sets num of fields to the buffer as array type.
func (e *Encoder) WriteStructHeaderFixAsArray(fieldNum, offset int) int {
	offset = e.setByte1Int(def.FixArray+fieldNum, offset)
	return offset
}

// WriteStructHeader16AsArray sets num of fields to the buffer as array type.
func (e *Encoder) WriteStructHeader16AsArray(fieldNum, offset int) int {
	offset = e.setByte1Int(def.Array16, offset)
	offset = e.setByte2Int(fieldNum, offset)
	return offset
}

// WriteStructHeader32AsArray sets num of fields to the buffer as array type.
func (e *Encoder) WriteStructHeader32AsArray(fieldNum, offset int) int {
	offset = e.setByte1Int(def.Array32, offset)
	offset = e.setByte4Int(fieldNum, offset)
	return offset
}

// WriteStructHeaderFixAsMap sets num of fields to the buffer as map type.
func (e *Encoder) WriteStructHeaderFixAsMap(fieldNum, offset int) int {
	offset = e.setByte1Int(def.FixMap+fieldNum, offset)
	return offset
}

// WriteStructHeader16AsMap sets num of fields to the buffer as map type.
func (e *Encoder) WriteStructHeader16AsMap(fieldNum, offset int) int {
	offset = e.setByte1Int(def.Map16, offset)
	offset = e.setByte2Int(fieldNum, offset)
	return offset
}

// WriteStructHeader32AsMap sets num of fields to the buffer as map type.
func (e *Encoder) WriteStructHeader32AsMap(fieldNum, offset int) int {
	offset = e.setByte1Int(def.Map32, offset)
	offset = e.setByte4Int(fieldNum, offset)
	return offset
}
