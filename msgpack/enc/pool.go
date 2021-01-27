package enc

//const fixedSize = 32767
//
//type bufferPool struct {
//	pool sync.Pool
//}
//
//var empty = make([]byte, fixedSize)
//
//var bufPool = bufferPool{
//	pool: sync.Pool{
//		New: func() interface{} {
//			buf := new(bytes.Buffer)
//			buf.Grow(fixedSize)
//			buf.Write(empty)
//			return buf
//		},
//	},
//}
//
//func (b *bufferPool) Get(size int) *bytes.Buffer {
//	buf := b.pool.Get().(*bytes.Buffer)
//	if buf.Len() < size {
//		buf.Grow(size - buf.Len())
//		buf.Write(make([]byte, size-buf.Len()))
//	}
//	return buf
//}
//
//func (b *bufferPool) Put(buf *bytes.Buffer) {
//	//buf.Reset()
//	b.pool.Put(buf)
//}
