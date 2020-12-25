package dec

import (
	"encoding/binary"
	"unsafe"

	"github.com/shamaton/msgpack/def"
)

var emptyString = ""
var emptyBytes = []byte{}

func (d *Decoder) isCodeString(code byte) bool {
	return d.isFixString(code) || code == def.Str8 || code == def.Str16 || code == def.Str32
}

func (d *Decoder) isFixString(v byte) bool {
	return def.FixStr <= v && v <= def.FixStr+0x1f
}

func (d *Decoder) StringByteLength(offset int) (int, int, error) {
	code := d.data[offset]
	offset++

	if def.FixStr <= code && code <= def.FixStr+0x1f {
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

func (d *Decoder) AsString(offset int) (string, int, error) {
	l, offset, err := d.StringByteLength(offset)
	if err != nil {
		return emptyString, 0, err
	}
	bs, offset := d.asStringByte(offset, l)
	return *(*string)(unsafe.Pointer(&bs)), offset, nil
}

func (d *Decoder) asStringByte(offset int, l int) ([]byte, int) {
	if l < 1 {
		return emptyBytes, offset
	}

	return d.readSizeN(offset, l)
}
