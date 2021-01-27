package enc

type Encoder struct {
	//buf  *bytes.Buffer
	d    []byte
	size int
}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) MakeBytes(size int) {
	//e.size = size
	//e.buf = bufPool.Get(size)
	//e.d = e.buf.Bytes()
	e.d = make([]byte, size)
}

func (e *Encoder) EncodedBytes() []byte {
	//return e.d[:e.size]
	return e.d
}

func (e *Encoder) ReleaseBytes() {
	// bufPool.Put(e.buf)
}
