package enc

import (
	"time"

	"github.com/shamaton/msgpack/v3/def"
)

// CalcTime checks value and returns data size that need.
func CalcTime(t time.Time) int {
	t = t.UTC()
	secs := uint64(t.Unix())
	if secs>>34 == 0 {
		data := uint64(t.Nanosecond())<<34 | secs
		if data&0xffffffff00000000 == 0 {
			return def.Byte1 + def.Byte1 + def.Byte4
		}
		return def.Byte1 + def.Byte1 + def.Byte8
	}

	return def.Byte1 + def.Byte1 + def.Byte1 + def.Byte4 + def.Byte8
}

// CalcTimeMax returns the maximum data size that a time value can need.
func CalcTimeMax(t time.Time) int {
	return def.Byte1 + def.Byte1 + def.Byte1 + def.Byte4 + def.Byte8
}

// WriteTime sets the contents of t to buf at offset.
func WriteTime(buf []byte, t time.Time, offset int) int {
	t = t.UTC()
	secs := uint64(t.Unix())
	if secs>>34 == 0 {
		data := uint64(t.Nanosecond())<<34 | secs
		if data&0xffffffff00000000 == 0 {
			offset = setByte1Int(buf, def.Fixext4, offset)
			offset = setByte1Int(buf, def.TimeStamp, offset)
			offset = setByte4Uint64(buf, data, offset)
			return offset
		}

		offset = setByte1Int(buf, def.Fixext8, offset)
		offset = setByte1Int(buf, def.TimeStamp, offset)
		offset = setByte8Uint64(buf, data, offset)
		return offset
	}

	offset = setByte1Int(buf, def.Ext8, offset)
	offset = setByte1Int(buf, 12, offset)
	offset = setByte1Int(buf, def.TimeStamp, offset)
	offset = setByte4Int(buf, t.Nanosecond(), offset)
	offset = setByte8Uint64(buf, secs, offset)
	return offset
}
