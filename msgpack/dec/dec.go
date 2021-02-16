package dec

import (
	"fmt"
)

// A Decoder holds an encoded bytes array and has methods for encoding.
type Decoder struct {
	data []byte
}

// NewDecoder creates a new Decoder for deserialization.
func NewDecoder(data []byte) *Decoder {
	return &Decoder{data: data}
}

// Len get encoded data length.
func (d *Decoder) Len() int { return len(d.data) }

func (d *Decoder) errorTemplate(code byte, str string) error {
	return fmt.Errorf("msgpackgen : invalid code %x decoding %s", code, str)
}
