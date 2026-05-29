package dec

import (
	"encoding/binary"

	"github.com/shamaton/msgpack/v3/def"
)

var emptyString = ""
var emptyBytes = []byte{}

// AsString checks codes and returns the got bytes as string
func (d *Decoder) AsString(offset int) (string, int, error) {
	bs, offset, err := d.AsStringBytes(offset)
	if err != nil {
		return emptyString, 0, err
	}
	return string(bs), offset, nil
}

func (d *Decoder) AsStringBytes(offset int) ([]byte, int, error) {
	code, offset, err := d.readSize1(offset)
	if err != nil {
		return emptyBytes, 0, err
	}

	var l int
	if d.isFixString(code) {
		l = int(code - def.FixStr)
	} else if code == def.Str8 {
		b, next, err := d.readSize1(offset)
		if err != nil {
			return emptyBytes, 0, err
		}
		l = int(b)
		offset = next
	} else if code == def.Str16 {
		b, next, err := d.readSize2(offset)
		if err != nil {
			return emptyBytes, 0, err
		}
		l = int(binary.BigEndian.Uint16(b))
		offset = next
	} else if code == def.Str32 {
		b, next, err := d.readSize4(offset)
		if err != nil {
			return emptyBytes, 0, err
		}
		l = int(binary.BigEndian.Uint32(b))
		offset = next
	} else if code == def.Nil {
		return emptyBytes, offset, nil
	} else {
		return emptyBytes, 0, d.errorTemplate(code, "StringByteLength")
	}

	if l < 1 {
		return emptyBytes, offset, nil
	}

	bs, offset, err := d.readSizeN(offset, l)
	if err != nil {
		return emptyBytes, 0, err
	}
	return bs, offset, nil
}

func (d *Decoder) isFixString(v byte) bool {
	return def.FixStr <= v && v <= def.FixStr+0x1f
}
