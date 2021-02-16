package enc

// CalcRune check value and returns data size that need.
func (e *Encoder) CalcRune(v rune) int {
	return e.calcInt(int64(v))
}

// WriteRune sets the contents of v to the buffer.
func (e *Encoder) WriteRune(v rune, offset int) int {
	return e.writeInt(int64(v), offset)
}
