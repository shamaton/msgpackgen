package enc

func setByte1Int64(buf []byte, value int64, offset int) int {
	buf[offset] = byte(value)
	return offset + 1
}

func setByte2Int64(buf []byte, value int64, offset int) int {
	_ = buf[offset+1]
	buf[offset+0] = byte(value >> 8)
	buf[offset+1] = byte(value)
	return offset + 2
}

func setByte4Int64(buf []byte, value int64, offset int) int {
	_ = buf[offset+3]
	buf[offset+0] = byte(value >> 24)
	buf[offset+1] = byte(value >> 16)
	buf[offset+2] = byte(value >> 8)
	buf[offset+3] = byte(value)
	return offset + 4
}

func setByte8Int64(buf []byte, value int64, offset int) int {
	_ = buf[offset+7]
	buf[offset] = byte(value >> 56)
	buf[offset+1] = byte(value >> 48)
	buf[offset+2] = byte(value >> 40)
	buf[offset+3] = byte(value >> 32)
	buf[offset+4] = byte(value >> 24)
	buf[offset+5] = byte(value >> 16)
	buf[offset+6] = byte(value >> 8)
	buf[offset+7] = byte(value)
	return offset + 8
}

func setByte1Uint64(buf []byte, value uint64, offset int) int {
	buf[offset] = byte(value)
	return offset + 1
}

func setByte2Uint64(buf []byte, value uint64, offset int) int {
	_ = buf[offset+1]
	buf[offset] = byte(value >> 8)
	buf[offset+1] = byte(value)
	return offset + 2
}

func setByte4Uint64(buf []byte, value uint64, offset int) int {
	_ = buf[offset+3]
	buf[offset] = byte(value >> 24)
	buf[offset+1] = byte(value >> 16)
	buf[offset+2] = byte(value >> 8)
	buf[offset+3] = byte(value)
	return offset + 4
}

func setByte8Uint64(buf []byte, value uint64, offset int) int {
	_ = buf[offset+7]
	buf[offset] = byte(value >> 56)
	buf[offset+1] = byte(value >> 48)
	buf[offset+2] = byte(value >> 40)
	buf[offset+3] = byte(value >> 32)
	buf[offset+4] = byte(value >> 24)
	buf[offset+5] = byte(value >> 16)
	buf[offset+6] = byte(value >> 8)
	buf[offset+7] = byte(value)
	return offset + 8
}

func setByte1Int(buf []byte, code, offset int) int {
	buf[offset] = byte(code)
	return offset + 1
}

func setByte2Int(buf []byte, value int, offset int) int {
	_ = buf[offset+1]
	buf[offset] = byte(value >> 8)
	buf[offset+1] = byte(value)
	return offset + 2
}

func setByte4Int(buf []byte, value int, offset int) int {
	_ = buf[offset+3]
	buf[offset] = byte(value >> 24)
	buf[offset+1] = byte(value >> 16)
	buf[offset+2] = byte(value >> 8)
	buf[offset+3] = byte(value)
	return offset + 4
}

func setByte(buf []byte, b byte, offset int) int {
	buf[offset] = b
	return offset + 1
}
