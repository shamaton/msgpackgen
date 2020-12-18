package msgpackgen_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/shamaton/msgpack"

	"github.com/shamaton/msgpack/def"

	"github.com/shamaton/msgpackgen"
	"github.com/shamaton/msgpackgen/dec"
	"github.com/shamaton/msgpackgen/enc"
)

// todo : extまわりの対応、time.Timeとか

type structTest struct {
	A      int
	B      float32
	String string
	Bool   bool
	Uint64 uint64
	Now    time.Time
	Slice  []uint
	Map    map[string]float64
	//ItemData Item
	//Items    []Item
	//Interface interface{}
}

var v = structTest{
	A:      123456789,
	B:      12.34,
	String: "my name is msgpack gen.",
	Bool:   true,
	Uint64: math.MaxUint32 * 2,
	Now:    time.Now(),
	Slice:  nil, // []uint{1, 100, 10000, 1000000},
	Map:    map[string]float64{"a": 1.23, "b": 2.34, "c": 3.45},
	//ItemData: Item{ID: 1, Name: "abc", Effect: 7.89, Num: 999},
	//Interface: "interface is not supported",
}
var num = 8

type Item struct {
	ID     int
	Name   string
	Effect float32
	Num    uint
}

func TestA(t *testing.T) {
	e := func(interface{}) ([]byte, error) { return nil, nil }
	d := func([]byte, interface{}) (bool, error) { return false, nil }
	add()
	msgpackgen.SetStructAsArray(true)

	msgpackgen.SetResolver(e, d)
	check(t)
	msgpackgen.SetResolver(encode, decode)
	check(t)
}

func check(t *testing.T) {

	b, err := msgpackgen.Encode(v)
	if err != nil {
		t.Error(err)
	}

	var vv structTest
	err = msgpackgen.Decode(b, &vv)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(v, vv)
	fmt.Printf("% x \n", b)
}

func add() {
	//if v.Items != nil {
	//	return
	//}
	//
	//n := 1
	//v.Items = make([]Item, n)
	//for i := 0; i < n; i++ {
	//	name := "item" + fmt.Sprint(i)
	//	item := Item{
	//		ID:     i,
	//		Name:   name,
	//		Effect: float32(i*i) / 3.0,
	//		Num:    uint(i * i * i * i),
	//	}
	//	v.Items[i] = item
	//}
}

func BenchmarkMsgGenEncShamaton(b *testing.B) {

	add()
	msgpackgen.SetResolver(encode, decode)

	for i := 0; i < b.N; i++ {
		_, err := msgpackgen.Encode(v)
		if err != nil {
			b.Error(err)
		}
	}

}

func BenchmarkMsgEncShamaton(b *testing.B) {

	add()

	for i := 0; i < b.N; i++ {
		_, err := msgpack.Encode(v)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkMsgGenDecShamaton(b *testing.B) {

	add()
	msgpackgen.SetResolver(encode, decode)

	d, _ := msgpack.Encode(v)
	var vv structTest

	for i := 0; i < b.N; i++ {
		err := msgpackgen.Decode(d, &vv)
		if err != nil {
			b.Error(err)
		}
	}

}

func BenchmarkMsgDecShamaton(b *testing.B) {

	add()

	d, _ := msgpack.Encode(v)
	var vv structTest

	for i := 0; i < b.N; i++ {
		err := msgpack.Decode(d, &vv)
		if err != nil {
			b.Error(err)
		}
	}
}

// todo : auto generate
func decode(data []byte, i interface{}) (bool, error) {

	if msgpackgen.StructAsArray() {
		return decodeAsArray(data, i)
	} else {
		return decodeAsMap(data, i)
	}
}

func decodeAsArray(data []byte, i interface{}) (bool, error) {

	switch v := i.(type) {
	case *structTest:
		_, err := decodeArraystructTest(v, dec.NewDecoder(data), 0)
		return true, err

	case *Item:
		_, err := decodeItem(v, dec.NewDecoder(data), 0)
		return true, err

	case **structTest:
		_, err := decodeArraystructTest(*v, dec.NewDecoder(data), 0)
		return true, err

	case **Item:
		_, err := decodeItem(*v, dec.NewDecoder(data), 0)
		return true, err

	}

	return false, nil
}

func decodeAsMap(data []byte, i interface{}) (bool, error) {

	switch v := i.(type) {
	case *structTest:
		_, err := decodeMapstructTest(v, dec.NewDecoder(data), 0)
		return true, err

	case *Item:
		_, err := decodeItem(v, dec.NewDecoder(data), 0)
		return true, err

	case **structTest:
		_, err := decodeMapstructTest(*v, dec.NewDecoder(data), 0)
		return true, err

	case **Item:
		_, err := decodeItem(*v, dec.NewDecoder(data), 0)
		return true, err

	}
	return false, nil
}

// todo : auto generate
func encode(i interface{}) ([]byte, error) {
	if msgpackgen.StructAsArray() {
		return encodeAsArray(i)
	} else {
		return encodeAsMap(i)
	}

}
func encodeAsArray(i interface{}) ([]byte, error) {

	switch v := i.(type) {
	case structTest:
		e := enc.NewEncoder()
		size, err := calcArraySizestructTest(v, e)
		if err != nil {
			return nil, err
		}
		e.MakeBytes(size)
		b, _, err := encodeArraystructTest(v, e, 0)
		return b, err

	case Item:
		e := enc.NewEncoder()
		size, err := calcSizeItem(v, e)
		if err != nil {
			return nil, err
		}
		e.MakeBytes(size)
		b, _, err := encodeItem(v, e, 0)
		return b, err

	}
	return nil, nil
}

func encodeAsMap(i interface{}) ([]byte, error) {

	switch v := i.(type) {
	case structTest:
		e := enc.NewEncoder()
		size, err := calcMapSizestructTest(v, e)
		if err != nil {
			return nil, err
		}
		e.MakeBytes(size)
		b, _, err := encodeMapstructTest(v, e, 0)
		return b, err

	case Item:
		e := enc.NewEncoder()
		size, err := calcSizeItem(v, e)
		if err != nil {
			return nil, err
		}
		e.MakeBytes(size)
		b, _, err := encodeItem(v, e, 0)
		return b, err

	}
	return nil, nil
}

func calcArraySizestructTest(v structTest, encoder *enc.Encoder) (int, error) {
	size := 0
	{
		s, err := encoder.CalcStruct(num)
		if err != nil {
			return 0, err
		}
		size += s
	}

	size += encoder.CalcInt(int64(v.A))

	size += encoder.CalcFloat32(float64(v.B))

	size += encoder.CalcString(v.String)

	size += encoder.CalcBool()

	size += encoder.CalcUint(v.Uint64)

	size += encoder.CalcTime(v.Now)

	// todo : nilのパターン
	if v.Slice != nil {
		s, err := encoder.CalcSliceLength(len(v.Slice))
		if err != nil {
			return 0, err
		}
		size += s
		for _, v := range v.Slice {
			size += encoder.CalcUint(uint64(v))
		}
	} else {
		size += encoder.CalcNil()
	}

	// todo : nilのパターン
	if v.Map != nil {
		s, err := encoder.CalcMapLength(len(v.Map))
		if err != nil {
			return 0, err
		}
		size += s
		for k, v := range v.Map {
			size += encoder.CalcString(k)
			size += encoder.CalcFloat64(v)
		}
	} else {
		size += encoder.CalcNil()
	}

	/*
		s, err = calcSizeItem(v.ItemData, encoder)
		if err != nil {
			return 0, err
		}
		size += s

		// todo : nilのパターン
		size += def.Byte1
		s, err = encoder.CalcSliceLength(len(v.Items))
		if err != nil {
			return 0, err
		}
		size += s
		for _, v := range v.Items {
			s, err = calcSizeItem(v, encoder)
			if err != nil {
				return 0, err
			}
			size += s
		}

	*/

	return size, nil
}

func calcMapSizestructTest(v structTest, encoder *enc.Encoder) (int, error) {
	size := 0
	{
		s, err := encoder.CalcStruct(num)
		if err != nil {
			return 0, err
		}
		size += s
	}

	size += encoder.CalcString("A")
	size += encoder.CalcInt(int64(v.A))

	size += encoder.CalcString("B")
	size += encoder.CalcFloat32(float64(v.B))

	size += encoder.CalcString("String")
	size += encoder.CalcString(v.String)

	size += encoder.CalcString("Bool")
	size += encoder.CalcBool()

	size += encoder.CalcString("Uint64")
	size += encoder.CalcUint(v.Uint64)

	size += encoder.CalcString("Now")
	size += encoder.CalcTime(v.Now)

	// todo : nilのパターン
	size += encoder.CalcString("Slice")
	if v.Slice != nil {
		s, err := encoder.CalcSliceLength(len(v.Slice))
		if err != nil {
			return 0, err
		}
		size += s
		for _, v := range v.Slice {
			size += encoder.CalcUint(uint64(v))
		}
	} else {
		size += encoder.CalcNil()
	}

	// todo : nilのパターン
	size += encoder.CalcString("Map")
	if v.Map != nil {
		s, err := encoder.CalcMapLength(len(v.Map))
		if err != nil {
			return 0, err
		}
		size += s
		for k, v := range v.Map {
			size += encoder.CalcString(k)
			size += encoder.CalcFloat64(v)
		}
	} else {
		size += encoder.CalcNil()
	}

	/*
		s, err = calcSizeItem(v.ItemData, encoder)
		if err != nil {
			return 0, err
		}
		size += s

		// todo : nilのパターン
		size += def.Byte1
		s, err = encoder.CalcSliceLength(len(v.Items))
		if err != nil {
			return 0, err
		}
		size += s
		for _, v := range v.Items {
			s, err = calcSizeItem(v, encoder)
			if err != nil {
				return 0, err
			}
			size += s
		}

	*/

	return size, nil
}

func encodeArraystructTest(v structTest, encoder *enc.Encoder, offset int) ([]byte, int, error) {
	var err error
	offset = encoder.WriteStruct(num, offset)

	offset = encoder.WriteInt(int64(v.A), offset)
	offset = encoder.WriteFloat32(v.B, offset)
	offset = encoder.WriteString(v.String, offset)
	offset = encoder.WriteBool(v.Bool, offset)
	offset = encoder.WriteUint(v.Uint64, offset)
	offset = encoder.WriteTime(v.Now, offset)

	// todo : nilのパターン
	if v.Slice != nil {
		offset = encoder.WriteSliceLength(len(v.Slice), offset)
		for _, vv := range v.Slice {
			offset = encoder.WriteUint(uint64(vv), offset)
		}
	} else {
		offset = encoder.WriteNil(offset)
	}

	// todo : nilのパターン
	if v.Map != nil {
		offset = encoder.WriteMapLength(len(v.Map), offset)
		for kk, vv := range v.Map {
			offset = encoder.WriteString(kk, offset)
			offset = encoder.WriteFloat64(vv, offset)
		}
	} else {
		offset = encoder.WriteNil(offset)
	}

	/*
		_, offset, err = encodeItem(v.ItemData, encoder, offset)
		if err != nil {
			return nil, 0, err
		}

		// todo : nilのパターン
		offset = encoder.WriteSliceLength(len(v.Items), offset)
		for _, vv := range v.Items {
			_, offset, err = encodeItem(vv, encoder, offset)
			if err != nil {
				return nil, 0, err
			}
		}

	*/
	return encoder.EncodedBytes(), offset, err
}

func encodeMapstructTest(v structTest, encoder *enc.Encoder, offset int) ([]byte, int, error) {
	var err error
	offset = encoder.WriteStruct(num, offset)

	offset = encoder.WriteString("A", offset)
	offset = encoder.WriteInt(int64(v.A), offset)

	offset = encoder.WriteString("B", offset)
	offset = encoder.WriteFloat32(v.B, offset)

	offset = encoder.WriteString("String", offset)
	offset = encoder.WriteString(v.String, offset)

	offset = encoder.WriteString("Bool", offset)
	offset = encoder.WriteBool(v.Bool, offset)

	offset = encoder.WriteString("Uint64", offset)
	offset = encoder.WriteUint(v.Uint64, offset)

	offset = encoder.WriteString("Now", offset)
	offset = encoder.WriteTime(v.Now, offset)

	// todo : nilのパターン
	offset = encoder.WriteString("Slice", offset)
	if v.Slice != nil {
		offset = encoder.WriteSliceLength(len(v.Slice), offset)
		for _, vv := range v.Slice {
			offset = encoder.WriteUint(uint64(vv), offset)
		}
	} else {
		offset = encoder.WriteNil(offset)
	}

	// todo : nilのパターン
	offset = encoder.WriteString("Map", offset)
	if v.Map != nil {
		offset = encoder.WriteMapLength(len(v.Map), offset)
		for kk, vv := range v.Map {
			offset = encoder.WriteString(kk, offset)
			offset = encoder.WriteFloat64(vv, offset)
		}
	} else {
		offset = encoder.WriteNil(offset)
	}

	/*
		_, offset, err = encodeItem(v.ItemData, encoder, offset)
		if err != nil {
			return nil, 0, err
		}

		// todo : nilのパターン
		offset = encoder.WriteSliceLength(len(v.Items), offset)
		for _, vv := range v.Items {
			_, offset, err = encodeItem(vv, encoder, offset)
			if err != nil {
				return nil, 0, err
			}
		}

	*/
	return encoder.EncodedBytes(), offset, err
}

func decodeArraystructTest(v *structTest, decoder *dec.Decoder, offset int) (int, error) {

	offset, err := decoder.CheckStruct(num, 0)
	if err != nil {
		return 0, err
	}
	{
		var vv int64
		vv, offset, err = decoder.AsInt(offset)
		if err != nil {
			return 0, err
		}
		v.A = int(vv)
	}
	{
		var vv float32
		vv, offset, err = decoder.AsFloat32(offset)
		if err != nil {
			return 0, err
		}
		v.B = vv
	}
	{
		var vv string
		vv, offset, err = decoder.AsString(offset)
		if err != nil {
			return 0, err
		}
		v.String = vv
	}
	{
		var vv bool
		vv, offset, err = decoder.AsBool(offset)
		if err != nil {
			return 0, err
		}
		v.Bool = vv
	}
	{
		var vv uint64
		vv, offset, err = decoder.AsUint(offset)
		if err != nil {
			return 0, err
		}
		v.Uint64 = vv
	}
	{
		var vv time.Time
		vv, offset, err = decoder.AsDateTime(offset)
		if err != nil {
			return 0, err
		}
		v.Now = vv
	}
	if !decoder.IsCodeNil(offset) {

		var vv []uint
		l, o, err := decoder.SliceLength(offset)
		if err != nil {
			return 0, err
		}

		vv = make([]uint, l)
		for i := range vv {
			vvv, oo, err := decoder.AsUint(o)
			if err != nil {
				return 0, err
			}
			vv[i] = uint(vvv)
			o = oo
		}
		offset = o
		v.Slice = vv
	} else {
		offset++
	}

	if !decoder.IsCodeNil(offset) {
		// todo : nilのパターン
		var vv map[string]float64
		l, o, err := decoder.MapLength(offset)
		if err != nil {
			return 0, err
		}

		vv = make(map[string]float64, l)
		for i := 0; i < l; i++ {
			vvv, oo, err := decoder.AsString(o)
			if err != nil {
				return 0, err
			}
			o = oo

			vvvv, oo, err := decoder.AsFloat64(o)
			if err != nil {
				return 0, err
			}
			o = oo
			vv[vvv] = vvvv
		}
		offset = o
		v.Map = vv
	} else {
		offset++
	}
	/*
		{
			var vv Item
			offset, err = decodeItem(&vv, decoder, offset)
			if err != nil {
				return 0, err
			}
			v.ItemData = vv
		}
		{
			// todo : nilのパターン
			var vv []Item
			l, o, err := decoder.SliceLength(offset)
			if err != nil {
				return 0, err
			}

			vv = make([]Item, l)
			for i := range vv {
				var vvv Item
				oo, err := decodeItem(&vvv, decoder, o)
				if err != nil {
					return 0, err
				}
				vv[i] = vvv
				o = oo
			}
			offset = o
			v.Items = vv
		}

	*/
	return offset, err
}

func decodeMapstructTest(v *structTest, decoder *dec.Decoder, offset int) (int, error) {

	// todo : mapの場合はここでstringをみてswitchする
	offset, err := decoder.CheckStruct(num, 0)
	if err != nil {
		return 0, err
	}

	dataLen := decoder.Len()
	for offset < dataLen {
		s, o, err := decoder.AsString(offset)
		if err != nil {
			return 0, err
		}
		offset = o

		switch s {
		case "A":
			{
				var vv int64
				vv, offset, err = decoder.AsInt(offset)
				if err != nil {
					return 0, err
				}
				v.A = int(vv)
			}
		case "B":
			{
				var vv float32
				vv, offset, err = decoder.AsFloat32(offset)
				if err != nil {
					return 0, err
				}
				v.B = vv
			}
		case "String":
			{
				var vv string
				vv, offset, err = decoder.AsString(offset)
				if err != nil {
					return 0, err
				}
				v.String = vv
			}
		case "Bool":
			{
				var vv bool
				vv, offset, err = decoder.AsBool(offset)
				if err != nil {
					return 0, err
				}
				v.Bool = vv
			}
		case "Uint64":
			{
				var vv uint64
				vv, offset, err = decoder.AsUint(offset)
				if err != nil {
					return 0, err
				}
				v.Uint64 = vv
			}
		case "Now":
			{
				var vv time.Time
				vv, offset, err = decoder.AsDateTime(offset)
				if err != nil {
					return 0, err
				}
				v.Now = vv
			}
		case "Slice":
			if !decoder.IsCodeNil(offset) {
				// todo : nilのパターン
				var vv []uint
				l, o, err := decoder.SliceLength(offset)
				if err != nil {
					return 0, err
				}

				vv = make([]uint, l)
				for i := range vv {
					vvv, oo, err := decoder.AsUint(o)
					if err != nil {
						return 0, err
					}
					vv[i] = uint(vvv)
					o = oo
				}
				offset = o
				v.Slice = vv
			} else {
				offset++
			}

		case "Map":
			if !decoder.IsCodeNil(offset) {
				// todo : nilのパターン
				var vv map[string]float64
				l, o, err := decoder.MapLength(offset)
				if err != nil {
					return 0, err
				}

				vv = make(map[string]float64, l)
				for i := 0; i < l; i++ {
					vvv, oo, err := decoder.AsString(o)
					if err != nil {
						return 0, err
					}
					o = oo

					vvvv, oo, err := decoder.AsFloat64(o)
					if err != nil {
						return 0, err
					}
					o = oo
					vv[vvv] = vvvv
				}
				offset = o
				v.Map = vv
			} else {
				offset++
			}

		default:
			offset = decoder.JumpOffset(offset)
		}
	}

	if offset != dataLen {
		return 0, fmt.Errorf("structure check failed %d : %d", offset, dataLen)
	}

	/*
		{
			var vv Item
			offset, err = decodeItem(&vv, decoder, offset)
			if err != nil {
				return 0, err
			}
			v.ItemData = vv
		}
		{
			// todo : nilのパターン
			var vv []Item
			l, o, err := decoder.SliceLength(offset)
			if err != nil {
				return 0, err
			}

			vv = make([]Item, l)
			for i := range vv {
				var vvv Item
				oo, err := decodeItem(&vvv, decoder, o)
				if err != nil {
					return 0, err
				}
				vv[i] = vvv
				o = oo
			}
			offset = o
			v.Items = vv
		}

	*/
	return offset, err
}

func _decodeMapstructTest(v *structTest, decoder *dec.Decoder, offset int) (int, error) {

	// todo : mapの場合はここでstringをみてswitchする
	offset, err := decoder.CheckStruct(num, 0)
	if err != nil {
		return 0, err
	}

	// 最初に文字列情報を受取る
	fields := make(map[string]int, num)

	fieldOffset := offset
	dataLen := decoder.Len()
	index := 0

	// けつまでいったら終わり
	for fieldOffset < dataLen && index < num {
		s, o, err := decoder.AsString(fieldOffset)
		if err != nil {
			return 0, err
		}
		fields[s] = o
		fieldOffset = decoder.JumpOffset(o)
		index++
	}
	if fieldOffset != dataLen {
		return 0, fmt.Errorf("structure check failed %d : %d", fieldOffset, dataLen)
	}

	for field, off := range fields {
		switch field {
		case "A":
			{
				var vv int64
				vv, offset, err = decoder.AsInt(off)
				if err != nil {
					return 0, err
				}
				v.A = int(vv)
			}
		case "B":
			{
				var vv float32
				vv, offset, err = decoder.AsFloat32(off)
				if err != nil {
					return 0, err
				}
				v.B = vv
			}
		case "String":
			{
				var vv string
				vv, offset, err = decoder.AsString(off)
				if err != nil {
					return 0, err
				}
				v.String = vv
			}
		case "Bool":
			{
				var vv bool
				vv, offset, err = decoder.AsBool(off)
				if err != nil {
					return 0, err
				}
				v.Bool = vv
			}
		case "Uint64":
			{
				var vv uint64
				vv, offset, err = decoder.AsUint(off)
				if err != nil {
					return 0, err
				}
				v.Uint64 = vv
			}
		case "Slice":
			{
				// todo : nilのパターン
				var vv []uint
				l, o, err := decoder.SliceLength(off)
				if err != nil {
					return 0, err
				}

				vv = make([]uint, l)
				for i := range vv {
					vvv, oo, err := decoder.AsUint(o)
					if err != nil {
						return 0, err
					}
					vv[i] = uint(vvv)
					o = oo
				}
				offset = o
				v.Slice = vv
			}
		}
	}

	/*
		{
			var vv Item
			offset, err = decodeItem(&vv, decoder, offset)
			if err != nil {
				return 0, err
			}
			v.ItemData = vv
		}
		{
			// todo : nilのパターン
			var vv []Item
			l, o, err := decoder.SliceLength(offset)
			if err != nil {
				return 0, err
			}

			vv = make([]Item, l)
			for i := range vv {
				var vvv Item
				oo, err := decodeItem(&vvv, decoder, o)
				if err != nil {
					return 0, err
				}
				vv[i] = vvv
				o = oo
			}
			offset = o
			v.Items = vv
		}

	*/
	return offset, err
}

///////////////////////////////

func calcSizeItem(v Item, encoder *enc.Encoder) (int, error) {
	size := def.Byte1
	s, err := encoder.CalcStruct(4)
	if err != nil {
		return 0, err
	}
	size += s

	size += def.Byte1
	size += encoder.CalcInt(int64(v.ID))

	size += def.Byte1
	size += encoder.CalcString(v.Name)

	size += def.Byte1
	size += encoder.CalcFloat32(float64(v.Effect))

	size += def.Byte1
	size += encoder.CalcUint(uint64(v.Num))

	return size, err
}

func encodeItem(v Item, encoder *enc.Encoder, offset int) ([]byte, int, error) {
	offset = encoder.WriteStruct(4, offset)
	offset = encoder.WriteInt(int64(v.ID), offset)
	offset = encoder.WriteString(v.Name, offset)
	offset = encoder.WriteFloat32(v.Effect, offset)
	offset = encoder.WriteUint(uint64(v.Num), offset)
	return encoder.EncodedBytes(), offset, nil
}

func decodeItem(v *Item, decoder *dec.Decoder, offset int) (int, error) {

	offset, err := decoder.CheckStruct(4, offset)
	if err != nil {
		return 0, err
	}
	{
		var vv int64
		vv, offset, err = decoder.AsInt(offset)
		if err != nil {
			return 0, err
		}
		v.ID = int(vv)
	}
	{
		var vv string
		vv, offset, err = decoder.AsString(offset)
		if err != nil {
			return 0, err
		}
		v.Name = vv
	}
	{
		var vv float32
		vv, offset, err = decoder.AsFloat32(offset)
		if err != nil {
			return 0, err
		}
		v.Effect = vv
	}
	{
		var vv uint64
		vv, offset, err = decoder.AsUint(offset)
		if err != nil {
			return 0, err
		}
		v.Num = uint(vv)
	}
	return offset, nil
}
