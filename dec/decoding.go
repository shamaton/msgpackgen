package dec

import (
	"fmt"

	"github.com/shamaton/msgpack"
)

type Decoder struct {
	data    []byte
	asArray bool
	//common.Common
}

func NewDecoder(data []byte) *Decoder {
	return &Decoder{data: data, asArray: msgpack.StructAsArray}
}

func (d *Decoder) Len() int { return len(d.data) }

func (d *Decoder) errorTemplate(code byte, str string) error {
	return fmt.Errorf("msgpack : invalid code %x decoding %s", code, str)
}
