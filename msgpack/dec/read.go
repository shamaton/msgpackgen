package dec

import (
	"github.com/shamaton/msgpack/v3/def"
)

func (d *Decoder) readSize1(index int) (byte, int, error) {
	rb := def.Byte1
	if !d.canRead(index, rb) {
		return 0, 0, def.ErrTooShortBytes
	}
	return d.data[index], index + rb, nil
}

func (d *Decoder) readSize2(index int) ([]byte, int, error) {
	return d.readSizeN(index, def.Byte2)
}

func (d *Decoder) readSize4(index int) ([]byte, int, error) {
	return d.readSizeN(index, def.Byte4)
}

func (d *Decoder) readSize8(index int) ([]byte, int, error) {
	return d.readSizeN(index, def.Byte8)
}

func (d *Decoder) readSizeN(index, n int) ([]byte, int, error) {
	if !d.canRead(index, n) {
		return nil, 0, def.ErrTooShortBytes
	}
	return d.data[index : index+n], index + n, nil
}

func (d *Decoder) canRead(index, n int) bool {
	return index >= 0 && n >= 0 && index <= len(d.data) && n <= len(d.data)-index
}
