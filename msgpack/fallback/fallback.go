package fallback

import "github.com/shamaton/msgpack/v3"

func MarshalAsMap(v any) ([]byte, error) {
	return msgpack.MarshalAsMap(v)
}

func MarshalAsArray(v any) ([]byte, error) {
	return msgpack.MarshalAsArray(v)
}

func UnmarshalAsMap(data []byte, v any) error {
	return msgpack.UnmarshalAsMap(data, v)
}

func UnmarshalAsArray(data []byte, v any) error {
	return msgpack.UnmarshalAsArray(data, v)
}
