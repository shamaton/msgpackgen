package msgpackgen

import (
	"github.com/shamaton/msgpack"
	decoding "github.com/shamaton/msgpackgen/dec"
	encoding "github.com/shamaton/msgpackgen/enc"
)

func SetEncodingOption(asArray bool) {
	msgpack.StructAsArray = asArray
}

func SetResolver(er encoding.EncoderResolver, dr decoding.DecodeResolver) {
	encoding.Resolver = er
	decoding.Resolver = dr
}

// Encode returns the MessagePack-encoded byte array of v.
func Encode(v interface{}) ([]byte, error) {
	return encoding.Encode(v, msgpack.StructAsArray)
}

// Decode analyzes the MessagePack-encoded data and stores
// the result into the pointer of v.
func Decode(data []byte, v interface{}) error {
	return decoding.Decode(data, v, msgpack.StructAsArray)
}
