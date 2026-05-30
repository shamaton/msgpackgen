package dec

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/shamaton/msgpack/v3/def"
)

// AsDateTime checks codes and returns the got bytes as time.Time
func (d *Decoder) AsDateTime(offset int) (time.Time, int, error) {
	code, offset, err := d.readSize1(offset)
	if err != nil {
		return time.Time{}, 0, err
	}

	switch code {
	case def.Fixext4:
		t, offset, err := d.readSize1(offset)
		if err != nil {
			return time.Time{}, 0, err
		}
		if int8(t) != def.TimeStamp {
			return time.Time{}, 0, fmt.Errorf("fixext4. time type is different %d, %d", t, def.TimeStamp)
		}
		bs, offset, err := d.readSize4(offset)
		if err != nil {
			return time.Time{}, 0, err
		}
		return time.Unix(int64(binary.BigEndian.Uint32(bs)), 0).UTC(), offset, nil

	case def.Fixext8:
		t, offset, err := d.readSize1(offset)
		if err != nil {
			return time.Time{}, 0, err
		}
		if int8(t) != def.TimeStamp {
			return time.Time{}, 0, fmt.Errorf("fixext8. time type is different %d, %d", t, def.TimeStamp)
		}
		bs, offset, err := d.readSize8(offset)
		if err != nil {
			return time.Time{}, 0, err
		}
		data64 := binary.BigEndian.Uint64(bs)
		nano := int64(data64 >> 34)
		if nano > 999999999 {
			return time.Time{}, 0, fmt.Errorf("in timestamp 64 formats, nanoseconds must not be larger than 999999999 : %d", nano)
		}
		return time.Unix(int64(data64&0x00000003ffffffff), nano).UTC(), offset, nil

	case def.Ext8:
		c, offset, err := d.readSize1(offset)
		if err != nil {
			return time.Time{}, 0, err
		}
		if int8(c) != 12 {
			return time.Time{}, 0, fmt.Errorf("ext8. time ext length is different %d, %d", c, 12)
		}
		t, offset, err := d.readSize1(offset)
		if err != nil {
			return time.Time{}, 0, err
		}
		if int8(t) != def.TimeStamp {
			return time.Time{}, 0, fmt.Errorf("ext8. time type is different %d, %d", t, def.TimeStamp)
		}
		nanobs, offset, err := d.readSize4(offset)
		if err != nil {
			return time.Time{}, 0, err
		}
		secbs, offset, err := d.readSize8(offset)
		if err != nil {
			return time.Time{}, 0, err
		}
		nano := binary.BigEndian.Uint32(nanobs)
		if nano > 999999999 {
			return time.Time{}, 0, fmt.Errorf("in timestamp 96 formats, nanoseconds must not be larger than 999999999 : %d", nano)
		}
		sec := binary.BigEndian.Uint64(secbs)
		return time.Unix(int64(sec), int64(nano)).UTC(), offset, nil
	}

	return time.Time{}, 0, d.errorTemplate(code, "AsDateTime")
}
