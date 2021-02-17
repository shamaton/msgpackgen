package dec

import (
	"encoding/binary"

	"github.com/shamaton/msgpack/v2/def"
)

func (d *Decoder) isFixMap(v byte) bool {
	return def.FixMap <= v && v <= def.FixMap+0x0f
}

// MapLength reads the need bytes and convert to length value.
func (d *Decoder) MapLength(offset int) (int, int, error) {
	code, offset := d.readSize1(offset)

	switch {
	case d.isFixMap(code):
		return int(code - def.FixMap), offset, nil
	case code == def.Map16:
		bs, offset := d.readSize2(offset)
		return int(binary.BigEndian.Uint16(bs)), offset, nil
	case code == def.Map32:
		bs, offset := d.readSize4(offset)
		return int(binary.BigEndian.Uint32(bs)), offset, nil
	}
	return 0, 0, d.errorTemplate(code, "MapLength")
}
