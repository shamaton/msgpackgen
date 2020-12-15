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

type DecodeResolver func(data []byte, i interface{}) (bool, error)

var Resolver DecodeResolver

// Decode analyzes the MessagePack-encoded data and stores
// the result into the pointer of v.
func Decode(data []byte, v interface{}, asArray bool) error {
	if Resolver == nil {
		return fmt.Errorf("error")
	}

	b, err := Resolver(data, v)
	if err != nil {
		return err
	}
	if b {
		return nil
	}

	return msgpack.Decode(data, v)
}

func (d *Decoder) errorTemplate(code byte, str string) error {
	return fmt.Errorf("msgpack : invalid code %x decoding %s", code, str)
}
