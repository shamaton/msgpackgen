package msgpack

import (
	"github.com/shamaton/msgpack"
)

// Encode returns the MessagePack-encoded byte array of v.
func Encode(v interface{}) ([]byte, error) {
	if StructAsArray() {
		return EncodeAsArray(v)
	}
	return EncodeAsMap(v)
}

func EncodeAsMap(v interface{}) ([]byte, error) {
	if b, err := encAsMapResolver(v); err != nil {
		return nil, err
	} else if b != nil {
		return b, nil
	}

	return msgpack.MarshalAsMap(v)
}

func EncodeAsArray(v interface{}) ([]byte, error) {
	if b, err := encAsArrayResolver(v); err != nil {
		return nil, err
	} else if b != nil {
		return b, nil
	}

	return msgpack.MarshalAsArray(v)
}

// Decode analyzes the MessagePack-encoded data and stores
// the result into the pointer of v.
func Decode(data []byte, v interface{}) error {
	if StructAsArray() {
		return DecodeAsArray(data, v)
	}
	return DecodeAsMap(data, v)
}

func DecodeAsMap(data []byte, v interface{}) error {
	b, err := decAsMapResolver(data, v)
	if err != nil {
		return err
	}
	if b {
		return nil
	}
	return msgpack.UnmarshalAsMap(data, v)
}

func DecodeAsArray(data []byte, v interface{}) error {
	b, err := decAsArrayResolver(data, v)
	if err != nil {
		return err
	}
	if b {
		return nil
	}
	return msgpack.UnmarshalAsArray(data, v)
}

func SetStructAsArray(on bool) {
	msgpack.StructAsArray = on
}

func StructAsArray() bool {
	return msgpack.StructAsArray
}

// SetComplexTypeCode sets def.complexTypeCode in github.com/shamaton/msgpack
func SetComplexTypeCode(code int8) {
	msgpack.SetComplexTypeCode(code)
}
