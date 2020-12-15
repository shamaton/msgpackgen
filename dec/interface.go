package dec

import (
	"github.com/shamaton/msgpack/def"
)

func (d *Decoder) AsInterface(offset int) (interface{}, int, error) {
	code := d.data[offset]

	switch {
	case code == def.Nil:
		offset++
		return nil, offset, nil

	case code == def.True, code == def.False:
		v, offset, err := d.AsBool(offset)
		if err != nil {
			return nil, 0, err
		}
		return v, offset, nil

	case d.isPositiveFixNum(code), code == def.Uint8:
		v, offset, err := d.AsUint(offset)
		if err != nil {
			return nil, 0, err
		}
		return uint8(v), offset, err
	case code == def.Uint16:
		v, offset, err := d.AsUint(offset)
		if err != nil {
			return nil, 0, err
		}
		return uint16(v), offset, err
	case code == def.Uint32:
		v, offset, err := d.AsUint(offset)
		if err != nil {
			return nil, 0, err
		}
		return uint32(v), offset, err
	case code == def.Uint64:
		v, offset, err := d.AsUint(offset)
		if err != nil {
			return nil, 0, err
		}
		return v, offset, err

	case d.isNegativeFixNum(code), code == def.Int8:
		v, offset, err := d.AsInt(offset)
		if err != nil {
			return nil, 0, err
		}
		return int8(v), offset, err
	case code == def.Int16:
		v, offset, err := d.AsInt(offset)
		if err != nil {
			return nil, 0, err
		}
		return int16(v), offset, err
	case code == def.Int32:
		v, offset, err := d.AsInt(offset)
		if err != nil {
			return nil, 0, err
		}
		return int32(v), offset, err
	case code == def.Int64:
		v, offset, err := d.AsInt(offset)
		if err != nil {
			return nil, 0, err
		}
		return v, offset, err

	case code == def.Float32:
		v, offset, err := d.AsFloat32(offset)
		if err != nil {
			return nil, 0, err
		}
		return v, offset, err
	case code == def.Float64:
		v, offset, err := d.AsFloat64(offset)
		if err != nil {
			return nil, 0, err
		}
		return v, offset, err

	case d.isFixString(code), code == def.Str8, code == def.Str16, code == def.Str32:
		v, offset, err := d.AsString(offset)
		if err != nil {
			return nil, 0, err
		}
		return v, offset, err

	case code == def.Bin8, code == def.Bin16, code == def.Bin32:
		v, offset, err := d.AsBin(offset)
		if err != nil {
			return nil, 0, err
		}
		return v, offset, err

	case d.isFixSlice(code), code == def.Array16, code == def.Array32:
		l, o, err := d.SliceLength(offset)
		if err != nil {
			return nil, 0, err
		}

		v := make([]interface{}, l)
		for i := 0; i < l; i++ {
			vv, o2, err := d.AsInterface(o)
			if err != nil {
				return nil, 0, err
			}
			v[i] = vv
			o = o2
		}
		offset = o
		return v, offset, nil

	case d.isFixMap(code), code == def.Map16, code == def.Map32:
		l, o, err := d.MapLength(offset)
		if err != nil {
			return nil, 0, err
		}
		v := make(map[interface{}]interface{}, l)
		for i := 0; i < l; i++ {
			key, o2, err := d.AsInterface(o)
			if err != nil {
				return nil, 0, err
			}
			value, o2, err := d.AsInterface(o2)
			if err != nil {
				return nil, 0, err
			}
			v[key] = value
			o = o2
		}
		offset = o
		return v, offset, nil
	}

	/* use ext
	if d.isDateTime(offset) {
		v, offset, err := d.asDateTime(offset, k)
		if err != nil {
			return nil, 0, err
		}
		return v, offset, nil
	}
	*/

	// ext
	for i := range extCoders {
		if extCoders[i].IsType(offset, &d.data) {
			v, offset, err := extCoders[i].AsValue(offset, k, &d.data)
			if err != nil {
				return nil, 0, err
			}
			return v, offset, nil
		}
	}

	return nil, 0, d.errorTemplate(code, "AsInterface")
}
