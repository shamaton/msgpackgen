package enc

type Encoder struct {
	d []byte
}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) MakeBytes(size int) {
	e.d = make([]byte, size)
}

func (e *Encoder) EncodedBytes() []byte {
	return e.d
}
