package enc

func (e *Encoder) setByte1Int64(value int64, offset int) int {
	e.d[offset] = byte(value)
	return offset + 1
}

func (e *Encoder) setByte2Int64(value int64, offset int) int {
	e.d[offset+0] = byte(value >> 8)
	e.d[offset+1] = byte(value)
	return offset + 2
}

func (e *Encoder) setByte4Int64(value int64, offset int) int {
	e.d[offset+0] = byte(value >> 24)
	e.d[offset+1] = byte(value >> 16)
	e.d[offset+2] = byte(value >> 8)
	e.d[offset+3] = byte(value)
	return offset + 4
}

func (e *Encoder) setByte8Int64(value int64, offset int) int {
	e.d[offset] = byte(value >> 56)
	e.d[offset+1] = byte(value >> 48)
	e.d[offset+2] = byte(value >> 40)
	e.d[offset+3] = byte(value >> 32)
	e.d[offset+4] = byte(value >> 24)
	e.d[offset+5] = byte(value >> 16)
	e.d[offset+6] = byte(value >> 8)
	e.d[offset+7] = byte(value)
	return offset + 8
}

func (e *Encoder) setByte1Uint64(value uint64, offset int) int {
	e.d[offset] = byte(value)
	return offset + 1
}

func (e *Encoder) setByte2Uint64(value uint64, offset int) int {
	e.d[offset] = byte(value >> 8)
	e.d[offset+1] = byte(value)
	return offset + 2
}

func (e *Encoder) setByte4Uint64(value uint64, offset int) int {
	e.d[offset] = byte(value >> 24)
	e.d[offset+1] = byte(value >> 16)
	e.d[offset+2] = byte(value >> 8)
	e.d[offset+3] = byte(value)
	return offset + 4
}

func (e *Encoder) setByte8Uint64(value uint64, offset int) int {
	e.d[offset] = byte(value >> 56)
	e.d[offset+1] = byte(value >> 48)
	e.d[offset+2] = byte(value >> 40)
	e.d[offset+3] = byte(value >> 32)
	e.d[offset+4] = byte(value >> 24)
	e.d[offset+5] = byte(value >> 16)
	e.d[offset+6] = byte(value >> 8)
	e.d[offset+7] = byte(value)
	return offset + 8
}

func (e *Encoder) setByte1Int(code, offset int) int {
	e.d[offset] = byte(code)
	return offset + 1
}

func (e *Encoder) setByte2Int(value int, offset int) int {
	e.d[offset] = byte(value >> 8)
	e.d[offset+1] = byte(value)
	return offset + 2
}

func (e *Encoder) setByte4Int(value int, offset int) int {
	e.d[offset] = byte(value >> 24)
	e.d[offset+1] = byte(value >> 16)
	e.d[offset+2] = byte(value >> 8)
	e.d[offset+3] = byte(value)
	return offset + 4
}

func (e *Encoder) setByte(b byte, offset int) int {
	e.d[offset] = b
	return offset + 1
}
