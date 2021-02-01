package dec

import (
	"fmt"
)

type Decoder struct {
	data []byte
}

func NewDecoder(data []byte) *Decoder {
	return &Decoder{data: data}
}

func (d *Decoder) Len() int { return len(d.data) }

func (d *Decoder) errorTemplate(code byte, str string) error {
	return fmt.Errorf("msgpackgen : invalid code %x decoding %s", code, str)
}
