package msgpack

import "github.com/shamaton/msgpack/v2"

// UnmarshalAsMap decodes data that is encoded as map format.
// This is the same thing that StructAsArray sets false.
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

// UnmarshalAsArray decodes data that is encoded as array format.
// This is the same thing that StructAsArray sets true.
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
