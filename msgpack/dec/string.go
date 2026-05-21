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
	l, offset, err := d.stringByteLength(offset)
	if err != nil {
		return emptyBytes, 0, err
	}
	bs, offset, err := d.asStringByte(offset, l)
	if err != nil {
		return emptyBytes, 0, err
	}
	return bs, offset, nil
}

func (d *Decoder) stringByteLength(offset int) (int, int, error) {
	code, offset, err := d.readSize1Checked(offset)
	if err != nil {
		return 0, 0, err
	}

	if d.isFixString(code) {
		l := int(code - def.FixStr)
		return l, offset, nil
	} else if code == def.Str8 {
		b, offset, err := d.readSize1Checked(offset)
		if err != nil {
			return 0, 0, err
		}
		return int(b), offset, nil
	} else if code == def.Str16 {
		b, offset, err := d.readSize2Checked(offset)
		if err != nil {
			return 0, 0, err
		}
		return int(binary.BigEndian.Uint16(b)), offset, nil
	} else if code == def.Str32 {
		b, offset, err := d.readSize4Checked(offset)
		if err != nil {
			return 0, 0, err
		}
		return int(binary.BigEndian.Uint32(b)), offset, nil
	} else if code == def.Nil {
		return 0, offset, nil
	}
	return 0, 0, d.errorTemplate(code, "StringByteLength")
}

func (d *Decoder) isFixString(v byte) bool {
	return def.FixStr <= v && v <= def.FixStr+0x1f
}

func (d *Decoder) asStringByte(offset int, l int) ([]byte, int, error) {
	if l < 1 {
		return emptyBytes, offset, nil
	}

	return d.readSizeNChecked(offset, l)
}
