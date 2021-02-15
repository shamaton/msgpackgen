package msgpack

import "github.com/shamaton/msgpack/v2"

// MarshalAsMap encodes data as map format.
// This is the same thing that StructAsArray sets false.
func MarshalAsMap(v interface{}) ([]byte, error) {
	if b, err := encAsMapResolver(v); err != nil {
		return nil, err
	} else if b != nil {
		return b, nil
	}

	return msgpack.MarshalAsMap(v)
}

// MarshalAsArray encodes data as array format.
// This is the same thing that StructAsArray sets true.
func MarshalAsArray(v interface{}) ([]byte, error) {
	if b, err := encAsArrayResolver(v); err != nil {
		return nil, err
	} else if b != nil {
		return b, nil
	}

	return msgpack.MarshalAsArray(v)
}
