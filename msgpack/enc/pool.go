package enc

import (
	"bytes"
	"sync"
)

const fixedSize = 32767

type bufferPool struct {
	pool sync.Pool
}

var empty = make([]byte, fixedSize)

var bufPool = bufferPool{
	pool: sync.Pool{
		New: func() interface{} {
			buf := new(bytes.Buffer)
			buf.Grow(fixedSize)
			buf.Write(empty)
			return buf
		},
	},
}

func (b *bufferPool) Get(size int) *bytes.Buffer {
	buf := b.pool.Get().(*bytes.Buffer)
	if buf.Len() < size {
		buf.Grow(size - buf.Len())
		buf.Write(make([]byte, size-buf.Len()))
	}
	return buf
}

func (b *bufferPool) Put(buf *bytes.Buffer) {
	//buf.Reset()
	b.pool.Put(buf)
}

//type bufferPool struct {
//	pool sync.Pool
//}
//
//var bufPool = bufferPool{
//	pool: sync.Pool{
//		New: func() interface{} {
//			b := make([]byte, fixedSize)
//			return &b
//		},
//	},
//}
//
//var bbb []byte
//var using bool
//
//func init() {
//	fmt.Println("gggggggggggggggggggggget")
//	mp := map[int][]byte{}
//	for i := 0; i < 10; i++ {
//		mp[i] = bufPool.Get(1)
//	}
//	for _, v := range mp {
//		bufPool.Put(v)
//	}
//
//	bbb = make([]byte, fixedSize)
//}
//
//func (b *bufferPool) Get(size int) []byte {
//	if !using {
//		using = true
//		return bbb
//	} else {
//		return make([]byte, size)
//	}
//	return bbb
//	if size < fixedSize+1 {
//		b := b.pool.Get().(*[]byte)
//		return *b
//	}
//	buf := make([]byte, size)
//	return buf
//
//	var ff *bytes.Buffer
//	ff.WriteByte()
//}
//
//func (b *bufferPool) Put(bs []byte) {
//
//	return
//	if len(bs) < fixedSize+1 {
//		b.pool.Put(&bs)
//	}
//}
