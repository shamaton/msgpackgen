package enc

import "github.com/shamaton/msgpack/v3/def"

// CalcRune checks value and returns data size that need.
func CalcRune(v rune) int {
	return calcIntSize(int64(v))
}

// CalcRuneMax returns the maximum data size that a rune value can need.
func CalcRuneMax(v rune) int {
	return def.Byte1 + def.Byte4
}

// WriteRune sets the contents of v to buf at offset.
func WriteRune(buf []byte, v rune, offset int) int {
	return writeInt(buf, int64(v), offset)
}
