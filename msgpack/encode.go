package msgpack

import "github.com/shamaton/msgpack/v3"

// MarshalAsMap encodes data as map format.
// This is the same thing that StructAsArray sets false.
func MarshalAsMap(v any) ([]byte, error) {
	return MarshalAsMapTo(v, nil)
}

// MarshalAsMapTo encodes data as map format by appending to buf.
// This is the same thing that StructAsArray sets false.
func MarshalAsMapTo(v any, buf []byte) ([]byte, error) {
	base := buf
	if b, handled, err := encAsMapToResolver(v, buf); err != nil {
		return nil, err
	} else if handled {
		return b, nil
	}
	if b, err := encAsMapResolver(v); err != nil {
		return nil, err
	} else if b != nil {
		return append(base, b...), nil
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
	return MarshalAsArrayTo(v, nil)
}

// MarshalAsArrayTo encodes data as array format by appending to buf.
// This is the same thing that StructAsArray sets true.
func MarshalAsArrayTo(v any, buf []byte) ([]byte, error) {
	base := buf
	if b, handled, err := encAsArrayToResolver(v, buf); err != nil {
		return nil, err
	} else if handled {
		return b, nil
	}
	if b, err := encAsArrayResolver(v); err != nil {
		return nil, err
	} else if b != nil {
		return append(base, b...), nil
	}

	b, err := msgpack.MarshalAsArray(v)
	if err != nil {
		return nil, err
	}
	return append(base, b...), nil
}
