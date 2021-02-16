package dec

import (
	"encoding/binary"
	"math"

	"github.com/shamaton/msgpack/v2/def"
)

// AsFloat32 checks codes and returns the got bytes as float32
func (d *Decoder) AsFloat32(offset int) (float32, int, error) {
	code := d.data[offset]

	switch {
	case code == def.Float32:
		offset++
		bs, offset := d.readSize4(offset)
		v := math.Float32frombits(binary.BigEndian.Uint32(bs))
		return v, offset, nil

	case d.isPositiveFixNum(code), code == def.Uint8, code == def.Uint16, code == def.Uint32, code == def.Uint64:
		v, offset, _ := d.AsUint(offset)
		return float32(v), offset, nil

	case d.isNegativeFixNum(code), code == def.Int8, code == def.Int16, code == def.Int32, code == def.Int64:
		v, offset, _ := d.AsInt(offset)
		return float32(v), offset, nil

	case code == def.Nil:
		offset++
		return 0, offset, nil
	}
	return 0, 0, d.errorTemplate(code, "AsFloat32")
}

// AsFloat64 checks codes and returns the got bytes as float64
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
		v, offset, _ := d.AsUint(offset)
		return float64(v), offset, nil

	case d.isNegativeFixNum(code), code == def.Int8, code == def.Int16, code == def.Int32, code == def.Int64:
		v, offset, _ := d.AsInt(offset)
		return float64(v), offset, nil

	case code == def.Nil:
		offset++
		return 0, offset, nil
	}
	return 0, 0, d.errorTemplate(code, "AsFloat64")
}
