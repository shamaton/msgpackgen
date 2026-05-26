package enc

import "github.com/shamaton/msgpack/v3/def"

// CalcNil returns data size that need.
func CalcNil() int {
	return def.Byte1
}

// WriteNil sets nil to buf at offset.
func WriteNil(buf []byte, offset int) int {
	return setByte1Int(buf, def.Nil, offset)
}
