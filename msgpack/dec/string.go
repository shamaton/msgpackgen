package dec

import (
	"encoding/binary"

	"github.com/shamaton/msgpack/v2/def"
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
	bs, offset := d.asStringByte(offset, l)
	return bs, offset, nil
}

func (d *Decoder) stringByteLength(offset int) (int, int, error) {
	code := d.data[offset]
	offset++

	if d.isFixString(code) {
		l := int(code - def.FixStr)
		return l, offset, nil
	} else if code == def.Str8 {
		b, offset := d.readSize1(offset)
		return int(b), offset, nil
	} else if code == def.Str16 {
		b, offset := d.readSize2(offset)
		return int(binary.BigEndian.Uint16(b)), offset, nil
	} else if code == def.Str32 {
		b, offset := d.readSize4(offset)
		return int(binary.BigEndian.Uint32(b)), offset, nil
	} else if code == def.Nil {
		return 0, offset, nil
	}
	return 0, 0, d.errorTemplate(code, "StringByteLength")
}

func (d *Decoder) isFixString(v byte) bool {
	return def.FixStr <= v && v <= def.FixStr+0x1f
}

func (d *Decoder) asStringByte(offset int, l int) ([]byte, int) {
	if l < 1 {
		return emptyBytes, offset
	}

	return d.readSizeN(offset, l)
}
