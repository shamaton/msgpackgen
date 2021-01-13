package enc

func (e *Encoder) CalcRune(v rune) int {
	return e.calcInt(int64(v))
}

func (e *Encoder) WriteRune(v rune, offset int) int {
	return e.writeInt(int64(v), offset)
}
