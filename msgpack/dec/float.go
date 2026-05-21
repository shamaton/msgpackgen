package dec

import (
	"encoding/binary"
	"math"

	"github.com/shamaton/msgpack/v3/def"
)

// AsFloat32 checks codes and returns the got bytes as float32
func (d *Decoder) AsFloat32(offset int) (float32, int, error) {
	start := offset
	code, offset, err := d.readSize1Checked(offset)
	if err != nil {
		return 0, 0, err
	}

	switch {
	case code == def.Float32:
		bs, offset, err := d.readSize4Checked(offset)
		if err != nil {
			return 0, 0, err
		}
		v := math.Float32frombits(binary.BigEndian.Uint32(bs))
		return v, offset, nil

	case d.isPositiveFixNum(code), code == def.Uint8, code == def.Uint16, code == def.Uint32, code == def.Uint64:
		v, offset, err := d.AsUint(start)
		if err != nil {
			return 0, 0, err
		}
		return float32(v), offset, nil

	case d.isNegativeFixNum(code), code == def.Int8, code == def.Int16, code == def.Int32, code == def.Int64:
		v, offset, err := d.AsInt(start)
		if err != nil {
			return 0, 0, err
		}
		return float32(v), offset, nil

	case code == def.Nil:
		return 0, offset, nil
	}
	return 0, 0, d.errorTemplate(code, "AsFloat32")
}

// AsFloat64 checks codes and returns the got bytes as float64
func (d *Decoder) AsFloat64(offset int) (float64, int, error) {
	start := offset
	code, offset, err := d.readSize1Checked(offset)
	if err != nil {
		return 0, 0, err
	}

	switch {
	case code == def.Float64:
		bs, offset, err := d.readSize8Checked(offset)
		if err != nil {
			return 0, 0, err
		}
		v := math.Float64frombits(binary.BigEndian.Uint64(bs))
		return v, offset, nil

	case code == def.Float32:
		bs, offset, err := d.readSize4Checked(offset)
		if err != nil {
			return 0, 0, err
		}
		v := math.Float32frombits(binary.BigEndian.Uint32(bs))
		return float64(v), offset, nil

	case d.isPositiveFixNum(code), code == def.Uint8, code == def.Uint16, code == def.Uint32, code == def.Uint64:
		v, offset, err := d.AsUint(start)
		if err != nil {
			return 0, 0, err
		}
		return float64(v), offset, nil

	case d.isNegativeFixNum(code), code == def.Int8, code == def.Int16, code == def.Int32, code == def.Int64:
		v, offset, err := d.AsInt(start)
		if err != nil {
			return 0, 0, err
		}
		return float64(v), offset, nil

	case code == def.Nil:
		return 0, offset, nil
	}
	return 0, 0, d.errorTemplate(code, "AsFloat64")
}
