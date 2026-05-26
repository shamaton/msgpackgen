package enc

import "github.com/shamaton/msgpack/v3/def"

// CalcStructHeaderFix returns data size that need.
func CalcStructHeaderFix(fieldNum int) int {
	return def.Byte1
}

// CalcStructHeader16 returns data size that need.
func CalcStructHeader16(fieldNum int) int {
	return def.Byte1 + def.Byte2
}

// CalcStructHeader32 returns data size that need.
func CalcStructHeader32(fieldNum int) int {
	return def.Byte1 + def.Byte4
}

// WriteStructHeaderFixAsArray sets num of fields to buf as array type.
func WriteStructHeaderFixAsArray(buf []byte, fieldNum, offset int) int {
	return setByte1Int(buf, def.FixArray+fieldNum, offset)
}

// WriteStructHeader16AsArray sets num of fields to buf as array type.
func WriteStructHeader16AsArray(buf []byte, fieldNum, offset int) int {
	offset = setByte1Int(buf, def.Array16, offset)
	return setByte2Int(buf, fieldNum, offset)
}

// WriteStructHeader32AsArray sets num of fields to buf as array type.
func WriteStructHeader32AsArray(buf []byte, fieldNum, offset int) int {
	offset = setByte1Int(buf, def.Array32, offset)
	return setByte4Int(buf, fieldNum, offset)
}

// WriteStructHeaderFixAsMap sets num of fields to buf as map type.
func WriteStructHeaderFixAsMap(buf []byte, fieldNum, offset int) int {
	return setByte1Int(buf, def.FixMap+fieldNum, offset)
}

// WriteStructHeader16AsMap sets num of fields to buf as map type.
func WriteStructHeader16AsMap(buf []byte, fieldNum, offset int) int {
	offset = setByte1Int(buf, def.Map16, offset)
	return setByte2Int(buf, fieldNum, offset)
}

// WriteStructHeader32AsMap sets num of fields to buf as map type.
func WriteStructHeader32AsMap(buf []byte, fieldNum, offset int) int {
	offset = setByte1Int(buf, def.Map32, offset)
	return setByte4Int(buf, fieldNum, offset)
}
