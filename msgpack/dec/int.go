package dec

import (
	"encoding/binary"

	"github.com/shamaton/msgpack/v3/def"
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

	start := offset
	code, offset, err := d.readSize1(offset)
	if err != nil {
		return 0, 0, err
	}

	switch {
	case d.isPositiveFixNum(code):
		return int64(code), offset, nil

	case d.isNegativeFixNum(code):
		return int64(int8(code)), offset, nil

	case code == def.Uint8:
		b, offset, err := d.readSize1(offset)
		if err != nil {
			return 0, 0, err
		}
		return int64(uint8(b)), offset, nil

	case code == def.Int8:
		b, offset, err := d.readSize1(offset)
		if err != nil {
			return 0, 0, err
		}
		return int64(int8(b)), offset, nil

	case code == def.Uint16:
		bs, offset, err := d.readSize2(offset)
		if err != nil {
			return 0, 0, err
		}
		v := binary.BigEndian.Uint16(bs)
		return int64(v), offset, nil

	case code == def.Int16:
		bs, offset, err := d.readSize2(offset)
		if err != nil {
			return 0, 0, err
		}
		v := int16(binary.BigEndian.Uint16(bs))
		return int64(v), offset, nil

	case code == def.Uint32:
		bs, offset, err := d.readSize4(offset)
		if err != nil {
			return 0, 0, err
		}
		v := binary.BigEndian.Uint32(bs)
		return int64(v), offset, nil

	case code == def.Int32:
		bs, offset, err := d.readSize4(offset)
		if err != nil {
			return 0, 0, err
		}
		v := int32(binary.BigEndian.Uint32(bs))
		return int64(v), offset, nil

	case code == def.Uint64:
		bs, offset, err := d.readSize8(offset)
		if err != nil {
			return 0, 0, err
		}
		return int64(binary.BigEndian.Uint64(bs)), offset, nil

	case code == def.Int64:
		bs, offset, err := d.readSize8(offset)
		if err != nil {
			return 0, 0, err
		}
		return int64(binary.BigEndian.Uint64(bs)), offset, nil

	case code == def.Float32:
		v, offset, err := d.AsFloat32(start)
		if err != nil {
			return 0, 0, err
		}
		return int64(v), offset, nil

	case code == def.Float64:
		v, offset, err := d.AsFloat64(start)
		if err != nil {
			return 0, 0, err
		}
		return int64(v), offset, nil

	case code == def.Nil:
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
