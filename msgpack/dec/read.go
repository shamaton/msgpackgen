package dec

import (
	"github.com/shamaton/msgpack/v2/def"
)

func (d *Decoder) readSize1(index int) (byte, int) {
	rb := def.Byte1
	return d.data[index], index + rb
}

func (d *Decoder) readSize2(index int) ([]byte, int) {
	rb := def.Byte2
	return d.data[index : index+rb], index + rb
}

func (d *Decoder) readSize4(index int) ([]byte, int) {
	rb := def.Byte4
	return d.data[index : index+rb], index + rb
}

func (d *Decoder) readSize8(index int) ([]byte, int) {
	rb := def.Byte8
	return d.data[index : index+rb], index + rb
}

func (d *Decoder) readSizeN(index, n int) ([]byte, int) {
	return d.data[index : index+n], index + n
}
