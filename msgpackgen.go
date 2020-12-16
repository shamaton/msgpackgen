package msgpackgen

import (
	"github.com/shamaton/msgpack"
)

type (
	EncResolver func(i interface{}) ([]byte, error)
	DecResolver func(data []byte, i interface{}) (bool, error)
)

var (
	encResolver EncResolver
	decResolver DecResolver
)

func SetStructAsArray(on bool) {
	msgpack.StructAsArray = on
}

func StructAsArray() bool {
	return msgpack.StructAsArray
}

func SetResolver(er EncResolver, dr DecResolver) {
	encResolver = er
	decResolver = dr
}

// Encode returns the MessagePack-encoded byte array of v.
func Encode(v interface{}) ([]byte, error) {

	if encResolver != nil {
		if b, err := encResolver(v); err != nil {
			return nil, err
		} else if b != nil {
			return b, nil
		}
	}

	return msgpack.Encode(v)
}

// Decode analyzes the MessagePack-encoded data and stores
// the result into the pointer of v.
func Decode(data []byte, v interface{}) error {

	if decResolver != nil {
		b, err := decResolver(data, v)
		if err != nil {
			return err
		}
		if b {
			return nil
		}
	}

	return msgpack.Decode(data, v)
}
