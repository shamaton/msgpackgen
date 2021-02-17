package dec

import (
	"encoding/binary"

	"github.com/shamaton/msgpack/v2/def"
)

func (d *Decoder) isFixSlice(v byte) bool {
	return def.FixArray <= v && v <= def.FixArray+0x0f
}

// SliceLength reads the need bytes and convert to length value.
func (d *Decoder) SliceLength(offset int) (int, int, error) {
	code, offset := d.readSize1(offset)

	switch {
	case d.isFixSlice(code):
		return int(code - def.FixArray), offset, nil
	case code == def.Array16:
		bs, offset := d.readSize2(offset)
		return int(binary.BigEndian.Uint16(bs)), offset, nil
	case code == def.Array32:
		bs, offset := d.readSize4(offset)
		return int(binary.BigEndian.Uint32(bs)), offset, nil

	case code == def.Bin8:
		l, offset := d.readSize1(offset)
		return int(uint8(l)), offset, nil
	case code == def.Bin16:
		bs, offset := d.readSize2(offset)
		return int(binary.BigEndian.Uint16(bs)), offset, nil
	case code == def.Bin32:
		bs, offset := d.readSize4(offset)
		return int(binary.BigEndian.Uint32(bs)), offset, nil
	}
	return 0, 0, d.errorTemplate(code, "SliceLength")
}
