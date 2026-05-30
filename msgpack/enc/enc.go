package enc

// EnsureLen extends buf so that len(buf) is at least targetLen.
func EnsureLen(buf []byte, targetLen int) []byte {
	if targetLen < 0 {
		panic("msgpackgen: negative target length")
	}
	if targetLen <= len(buf) {
		return buf
	}
	if targetLen <= cap(buf) {
		return buf[:targetLen]
	}
	return append(buf, make([]byte, targetLen-len(buf))...)
}

// RequireAt extends buf so that writing extra bytes at offset is valid.
func RequireAt(buf []byte, offset, extra int) []byte {
	if offset < 0 || extra < 0 {
		panic("msgpackgen: negative offset or extra length")
	}
	return EnsureLen(buf, offset+extra)
}
