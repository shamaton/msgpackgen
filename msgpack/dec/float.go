package dec

import (
	"encoding/binary"
	"github.com/shamaton/msgpack/def"
	"math"
)

func (d *Decoder) AsFloat32(offset int) (float32, int, error) {
	code := d.data[offset]

	switch {
	case code == def.Float32:
		offset++
		bs, offset := d.readSize4(offset)
		v := math.Float32frombits(binary.BigEndian.Uint32(bs))
		return v, offset, nil

	case d.isPositiveFixNum(code), code == def.Uint8, code == def.Uint16, code == def.Uint32, code == def.Uint64:
		v, offset, err := d.AsUint(offset)
		if err != nil {
			break
		}
		return float32(v), offset, nil

	case d.isNegativeFixNum(code), code == def.Int8, code == def.Int16, code == def.Int32, code == def.Int64:
		v, offset, err := d.AsInt(offset)
		if err != nil {
			break
		}
		return float32(v), offset, nil

	case code == def.Nil:
		offset++
		return 0, offset, nil
	}
	return 0, 0, d.errorTemplate(code, "AsFloat32")
}

func (d *Decoder) AsFloat64(offset int) (float64, int, error) {
	code := d.data[offset]

	switch {
	case code == def.Float64:
		offset++
		bs, offset := d.readSize8(offset)
		v := math.Float64frombits(binary.BigEndian.Uint64(bs))
		return v, offset, nil

	case code == def.Float32:
		offset++
		bs, offset := d.readSize4(offset)
		v := math.Float32frombits(binary.BigEndian.Uint32(bs))
		return float64(v), offset, nil

	case d.isPositiveFixNum(code), code == def.Uint8, code == def.Uint16, code == def.Uint32, code == def.Uint64:
		v, offset, err := d.AsUint(offset)
		if err != nil {
			break
		}
		return float64(v), offset, nil

	case d.isNegativeFixNum(code), code == def.Int8, code == def.Int16, code == def.Int32, code == def.Int64:
		v, offset, err := d.AsInt(offset)
		if err != nil {
			break
		}
		return float64(v), offset, nil

	case code == def.Nil:
		offset++
		return 0, offset, nil
	}
	return 0, 0, d.errorTemplate(code, "AsFloat64")
}
