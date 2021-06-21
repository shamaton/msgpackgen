package dec

import (
	"encoding/binary"

	"github.com/shamaton/msgpack/v2/def"
)

// AsInt checks codes and returns the got bytes as int
func (d *Decoder) AsInt(offset int) (int, int, error) {
	v, offset, err := d.asInt(offset)
	return int(v), offset, err
}

// AsInt8 checks codes and returns the got bytes as int8
func (d *Decoder) AsInt8(offset int) (int8, int, error) {
	v, offset, err := d.asInt(offset)
	return int8(v), offset, err
}

// AsInt16 checks codes and returns the got bytes as int16
func (d *Decoder) AsInt16(offset int) (int16, int, error) {
	v, offset, err := d.asInt(offset)
	return int16(v), offset, err
}

// AsInt32 checks codes and returns the got bytes as int32
func (d *Decoder) AsInt32(offset int) (int32, int, error) {
	v, offset, err := d.asInt(offset)
	return int32(v), offset, err
}

// AsInt64 checks codes and returns the got bytes as int64
func (d *Decoder) AsInt64(offset int) (int64, int, error) {
	return d.asInt(offset)
}

func (d *Decoder) asInt(offset int) (int64, int, error) {

	code := d.data[offset]

	switch {
	case d.isPositiveFixNum(code):
		b, offset := d.readSize1(offset)
		return int64(b), offset, nil

	case d.isNegativeFixNum(code):
		b, offset := d.readSize1(offset)
		return int64(int8(b)), offset, nil

	case code == def.Uint8:
		offset++
		b, offset := d.readSize1(offset)
		return int64(uint8(b)), offset, nil

	case code == def.Int8:
		offset++
		b, offset := d.readSize1(offset)
		return int64(int8(b)), offset, nil

	case code == def.Uint16:
		offset++
		bs, offset := d.readSize2(offset)
		v := binary.BigEndian.Uint16(bs)
		return int64(v), offset, nil

	case code == def.Int16:
		offset++
		bs, offset := d.readSize2(offset)
		v := int16(binary.BigEndian.Uint16(bs))
		return int64(v), offset, nil

	case code == def.Uint32:
		offset++
		bs, offset := d.readSize4(offset)
		v := binary.BigEndian.Uint32(bs)
		return int64(v), offset, nil

	case code == def.Int32:
		offset++
		bs, offset := d.readSize4(offset)
		v := int32(binary.BigEndian.Uint32(bs))
		return int64(v), offset, nil

	case code == def.Uint64:
		offset++
		bs, offset := d.readSize8(offset)
		return int64(binary.BigEndian.Uint64(bs)), offset, nil

	case code == def.Int64:
		offset++
		bs, offset := d.readSize8(offset)
		return int64(binary.BigEndian.Uint64(bs)), offset, nil

	case code == def.Float32:
		v, offset, err := d.AsFloat32(offset)
		if err != nil {
			return 0, 0, err
		}
		return int64(v), offset, nil

	case code == def.Float64:
		v, offset, err := d.AsFloat64(offset)
		if err != nil {
			return 0, 0, err
		}
		return int64(v), offset, nil

	case code == def.Nil:
		offset++
		return 0, offset, nil
	}

	return 0, 0, d.errorTemplate(code, "AsInt")
}

func (d *Decoder) isPositiveFixNum(v byte) bool {
	return def.PositiveFixIntMin <= v && v <= def.PositiveFixIntMax
}

func (d *Decoder) isNegativeFixNum(v byte) bool {
	return def.NegativeFixintMin <= int8(v) && int8(v) <= def.NegativeFixintMax
}
