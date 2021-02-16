package dec

import (
	"encoding/binary"

	"github.com/shamaton/msgpack/v2/def"
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

	code := d.data[offset]

	switch {
	case d.isPositiveFixNum(code):
		b, offset := d.readSize1(offset)
		return uint64(b), offset, nil

	case d.isNegativeFixNum(code):
		b, offset := d.readSize1(offset)
		return uint64(int8(b)), offset, nil

	case code == def.Uint8:
		offset++
		b, offset := d.readSize1(offset)
		return uint64(uint8(b)), offset, nil

	case code == def.Int8:
		offset++
		b, offset := d.readSize1(offset)
		return uint64(int8(b)), offset, nil

	case code == def.Uint16:
		offset++
		bs, offset := d.readSize2(offset)
		v := binary.BigEndian.Uint16(bs)
		return uint64(v), offset, nil

	case code == def.Int16:
		offset++
		bs, offset := d.readSize2(offset)
		v := int16(binary.BigEndian.Uint16(bs))
		return uint64(v), offset, nil

	case code == def.Uint32:
		offset++
		bs, offset := d.readSize4(offset)
		v := binary.BigEndian.Uint32(bs)
		return uint64(v), offset, nil

	case code == def.Int32:
		offset++
		bs, offset := d.readSize4(offset)
		v := int32(binary.BigEndian.Uint32(bs))
		return uint64(v), offset, nil

	case code == def.Uint64:
		offset++
		bs, offset := d.readSize8(offset)
		return binary.BigEndian.Uint64(bs), offset, nil

	case code == def.Int64:
		offset++
		bs, offset := d.readSize8(offset)
		return binary.BigEndian.Uint64(bs), offset, nil

	case code == def.Nil:
		offset++
		return 0, offset, nil
	}

	return 0, 0, d.errorTemplate(code, "AsUint")
}
