package msgpack

import (
	"github.com/shamaton/msgpack/v2"
)

// Marshal returns the MessagePack-encoded byte array of v.
func Marshal(v interface{}) ([]byte, error) {
	if StructAsArray() {
		return MarshalAsArray(v)
	}
	return MarshalAsMap(v)
}

func MarshalAsMap(v interface{}) ([]byte, error) {
	if b, err := encAsMapResolver(v); err != nil {
		return nil, err
	} else if b != nil {
		return b, nil
	}

	return msgpack.MarshalAsMap(v)
}

func MarshalAsArray(v interface{}) ([]byte, error) {
	if b, err := encAsArrayResolver(v); err != nil {
		return nil, err
	} else if b != nil {
		return b, nil
	}

	return msgpack.MarshalAsArray(v)
}

// Unmarshal analyzes the MessagePack-encoded data and stores
// the result into the pointer of v.
func Unmarshal(data []byte, v interface{}) error {
	if StructAsArray() {
		return UnmarshalAsArray(data, v)
	}
	return UnmarshalAsMap(data, v)
}

func UnmarshalAsMap(data []byte, v interface{}) error {
	b, err := decAsMapResolver(data, v)
	if err != nil {
		return err
	}
	if b {
		return nil
	}
	return msgpack.UnmarshalAsMap(data, v)
}

func UnmarshalAsArray(data []byte, v interface{}) error {
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
