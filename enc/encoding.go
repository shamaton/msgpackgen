package encoding

import (
	"reflect"

	"github.com/shamaton/msgpack"
)

type Encoder struct {
	d       []byte
	asArray bool
	//common.Common
	mk map[uintptr][]reflect.Value
	mv map[uintptr][]reflect.Value
}

func NewEncoder() *Encoder {
	return &Encoder{asArray: msgpack.StructAsArray}
}

func (e *Encoder) MakeBytes(size int) {
	e.d = make([]byte, size)
}

func (e *Encoder) EncodedBytes() []byte { return e.d }

func (e *Encoder) create(rv reflect.Value, offset int) int {
	return -1
	/*
		switch rv.Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			v := rv.Uint()
			offset = e.WriteUint(v, offset)

		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			v := rv.Int()
			offset = e.WriteInt(v, offset)

		case reflect.Float32:
			offset = e.WriteFloat32(rv.Float(), offset)

		case reflect.Float64:
			offset = e.WriteFloat64(rv.Float(), offset)

		case reflect.Bool:
			offset = e.WriteBool(rv.Bool(), offset)

		case reflect.String:
			offset = e.WriteString(rv.String(), offset)

		case reflect.Slice:
			if rv.IsNil() {
				return e.WriteNil(offset)
			}
			l := rv.Len()
			// bin format
			if e.isByteSlice(rv) {
				offset = e.WriteByteSliceLength(l, offset)
				offset = e.setBytes(rv.Bytes(), offset)
				return offset
			}

			// format
			offset = e.WriteSliceLength(l, offset)

			if offset, find := e.writeFixedSlice(rv, offset); find {
				return offset
			}

			// func
			elem := rv.Type().Elem()
			var f structWriteFunc
			if elem.Kind() == reflect.Struct {
				f = e.getStructWriter(elem)
			} else {
				f = e.create
			}

			// objects
			for i := 0; i < l; i++ {
				offset = f(rv.Index(i), offset)
			}

		case reflect.Array:
			l := rv.Len()
			// bin format
			if e.isByteSlice(rv) {
				offset = e.WriteByteSliceLength(l, offset)
				// objects
				for i := 0; i < l; i++ {
					offset = e.setByte1Uint64(rv.Index(i).Uint(), offset)
				}
				return offset
			}

			// format
			offset = e.WriteSliceLength(l, offset)

			// func
			elem := rv.Type().Elem()
			var f structWriteFunc
			if elem.Kind() == reflect.Struct {
				f = e.getStructWriter(elem)
			} else {
				f = e.create
			}

			// objects
			for i := 0; i < l; i++ {
				offset = f(rv.Index(i), offset)
			}

		case reflect.Map:
			if rv.IsNil() {
				return e.WriteNil(offset)
			}

			l := rv.Len()
			offset = e.WriteMapLength(l, offset)

			if offset, find := e.writeFixedMap(rv, offset); find {
				return offset
			}

			// key-value
			p := rv.Pointer()
			for i := range e.mk[p] {
				offset = e.create(e.mk[p][i], offset)
				offset = e.create(e.mv[p][i], offset)
			}

		case reflect.Struct:
			offset = e.writeStruct(rv, offset)

		case reflect.Ptr:
			if rv.IsNil() {
				return e.WriteNil(offset)
			}

			offset = e.create(rv.Elem(), offset)

		case reflect.Interface:
			offset = e.create(rv.Elem(), offset)

		case reflect.Invalid:
			return e.WriteNil(offset)

		}
		return offset

	*/
}
