package msgpack

import "github.com/shamaton/msgpack/v3"

// MarshalAsMap encodes data as map format.
// This is the same thing that StructAsArray sets false.
func MarshalAsMap(v any) ([]byte, error) {
	return marshalAsMapTo(v, nil)
}

func marshalAsMapTo(v any, buf []byte) ([]byte, error) {
	base := buf
	if b, handled, err := encAsMapResolver(v, buf); err != nil {
		return nil, err
	} else if handled {
		return b, nil
	}

	b, err := msgpack.MarshalAsMap(v)
	if err != nil {
		return nil, err
	}
	return append(base, b...), nil
}

// MarshalAsArray encodes data as array format.
// This is the same thing that StructAsArray sets true.
func MarshalAsArray(v any) ([]byte, error) {
	return marshalAsArrayTo(v, nil)
}

func marshalAsArrayTo(v any, buf []byte) ([]byte, error) {
	base := buf
	if b, handled, err := encAsArrayResolver(v, buf); err != nil {
		return nil, err
	} else if handled {
		return b, nil
	}

	b, err := msgpack.MarshalAsArray(v)
	if err != nil {
		return nil, err
	}
	return append(base, b...), nil
}
