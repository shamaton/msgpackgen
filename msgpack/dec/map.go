package dec

import (
	"encoding/binary"

	"github.com/shamaton/msgpack/v3/def"
)

func (d *Decoder) isFixMap(v byte) bool {
	return def.FixMap <= v && v <= def.FixMap+0x0f
}

// MapLength reads the need bytes and convert to length value.
func (d *Decoder) MapLength(offset int) (int, int, error) {
	code, offset, err := d.readSize1Checked(offset)
	if err != nil {
		return 0, 0, err
	}

	switch {
	case d.isFixMap(code):
		return int(code - def.FixMap), offset, nil
	case code == def.Map16:
		bs, offset, err := d.readSize2Checked(offset)
		if err != nil {
			return 0, 0, err
		}
		return int(binary.BigEndian.Uint16(bs)), offset, nil
	case code == def.Map32:
		bs, offset, err := d.readSize4Checked(offset)
		if err != nil {
			return 0, 0, err
		}
		return int(binary.BigEndian.Uint32(bs)), offset, nil
	}
	return 0, 0, d.errorTemplate(code, "MapLength")
}
