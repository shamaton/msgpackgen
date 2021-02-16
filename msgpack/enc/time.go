package enc

import (
	"time"

	"github.com/shamaton/msgpack/v2/def"
)

// CalcTime check value and returns data size that need.
func (e *Encoder) CalcTime(t time.Time) int {
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

// WriteTime sets the contents of t to the buffer.
func (e *Encoder) WriteTime(t time.Time, offset int) int {
	secs := uint64(t.Unix())
	if secs>>34 == 0 {
		data := uint64(t.Nanosecond())<<34 | secs
		if data&0xffffffff00000000 == 0 {
			offset = e.setByte1Int(def.Fixext4, offset)
			offset = e.setByte1Int(def.TimeStamp, offset)
			offset = e.setByte4Uint64(data, offset)
			return offset
		}

		offset = e.setByte1Int(def.Fixext8, offset)
		offset = e.setByte1Int(def.TimeStamp, offset)
		offset = e.setByte8Uint64(data, offset)
		return offset
	}

	offset = e.setByte1Int(def.Ext8, offset)
	offset = e.setByte1Int(12, offset)
	offset = e.setByte1Int(def.TimeStamp, offset)
	offset = e.setByte4Int(t.Nanosecond(), offset)
	offset = e.setByte8Uint64(secs, offset)
	return offset
}
