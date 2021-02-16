package enc

// A Encoder calculates and writes the required byte size.
type Encoder struct {
	d []byte
}

// NewEncoder creates a new Encoder for serialization.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// MakeBytes reserves the required byte array
func (e *Encoder) MakeBytes(size int) {
	e.d = make([]byte, size)
}

// EncodedBytes gets encoded bytes.
func (e *Encoder) EncodedBytes() []byte {
	return e.d
}
