package dec

import (
	"encoding/binary"

	"github.com/shamaton/msgpack/v3/def"
)

// AsUint checks codes and returns the got bytes as uint
func (d *Decoder) AsUint(offset int) (uint, int, error) {
	v, offset, err := d.asUint(offset)
	return uint(v), offset, err
}

// AsUint8 checks codes and returns the got bytes as uint8
func (d *Decoder) AsUint8(offset int) (uint8, int, error) {
	v, offset, err := d.asUint(offset)
	return uint8(v), offset, err
}

// AsUint16 checks codes and returns the got bytes as uint16
func (d *Decoder) AsUint16(offset int) (uint16, int, error) {
	v, offset, err := d.asUint(offset)
	return uint16(v), offset, err
}

// AsUint32 checks codes and returns the got bytes as uint32
func (d *Decoder) AsUint32(offset int) (uint32, int, error) {
	v, offset, err := d.asUint(offset)
	return uint32(v), offset, err
}

// AsUint64 checks codes and returns the got bytes as uint64
func (d *Decoder) AsUint64(offset int) (uint64, int, error) {
	return d.asUint(offset)
}

func (d *Decoder) asUint(offset int) (uint64, int, error) {

	code, offset, err := d.readSize1(offset)
	if err != nil {
		return 0, 0, err
	}

	switch {
	case d.isPositiveFixNum(code):
		return uint64(code), offset, nil

	case d.isNegativeFixNum(code):
		return uint64(int8(code)), offset, nil

	case code == def.Uint8:
		b, offset, err := d.readSize1(offset)
		if err != nil {
			return 0, 0, err
		}
		return uint64(uint8(b)), offset, nil

	case code == def.Int8:
		b, offset, err := d.readSize1(offset)
		if err != nil {
			return 0, 0, err
		}
		return uint64(int8(b)), offset, nil

	case code == def.Uint16:
		bs, offset, err := d.readSize2(offset)
		if err != nil {
			return 0, 0, err
		}
		v := binary.BigEndian.Uint16(bs)
		return uint64(v), offset, nil

	case code == def.Int16:
		bs, offset, err := d.readSize2(offset)
		if err != nil {
			return 0, 0, err
		}
		v := int16(binary.BigEndian.Uint16(bs))
		return uint64(v), offset, nil

	case code == def.Uint32:
		bs, offset, err := d.readSize4(offset)
		if err != nil {
			return 0, 0, err
		}
		v := binary.BigEndian.Uint32(bs)
		return uint64(v), offset, nil

	case code == def.Int32:
		bs, offset, err := d.readSize4(offset)
		if err != nil {
			return 0, 0, err
		}
		v := int32(binary.BigEndian.Uint32(bs))
		return uint64(v), offset, nil

	case code == def.Uint64:
		bs, offset, err := d.readSize8(offset)
		if err != nil {
			return 0, 0, err
		}
		return binary.BigEndian.Uint64(bs), offset, nil

	case code == def.Int64:
		bs, offset, err := d.readSize8(offset)
		if err != nil {
			return 0, 0, err
		}
		return binary.BigEndian.Uint64(bs), offset, nil

	case code == def.Nil:
		return 0, offset, nil
	}

	return 0, 0, d.errorTemplate(code, "AsUint")
}
