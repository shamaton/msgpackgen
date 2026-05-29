package fallback

import "github.com/shamaton/msgpack/v3"

// MarshalAsMap is a fallback runtime API for generated code.
// Users should call generated MarshalAsMap instead of this function directly.
func MarshalAsMap(v any) ([]byte, error) {
	return msgpack.MarshalAsMap(v)
}

// MarshalAsArray is a fallback runtime API for generated code.
// Users should call generated MarshalAsArray instead of this function directly.
func MarshalAsArray(v any) ([]byte, error) {
	return msgpack.MarshalAsArray(v)
}

// UnmarshalAsMap is a fallback runtime API for generated code.
// Users should call generated UnmarshalAsMap instead of this function directly.
func UnmarshalAsMap(data []byte, v any) error {
	return msgpack.UnmarshalAsMap(data, v)
}

// UnmarshalAsArray is a fallback runtime API for generated code.
// Users should call generated UnmarshalAsArray instead of this function directly.
func UnmarshalAsArray(data []byte, v any) error {
	return msgpack.UnmarshalAsArray(data, v)
}
