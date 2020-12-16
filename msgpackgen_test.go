package msgpackgen_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/shamaton/msgpack"

	"github.com/shamaton/msgpack/def"

	"github.com/shamaton/msgpackgen"
	"github.com/shamaton/msgpackgen/dec"
	encoding "github.com/shamaton/msgpackgen/enc"
)

// todo : extまわりの対応、time.Timeとか

type structTest struct {
	A        int
	B        float32
	String   string
	Bool     bool
	Uint64   uint64
	Slice    []uint
	ItemData Item
	Items    []Item
}

var v = structTest{
	A:        123456789,
	B:        12.34,
	String:   "my name is msgpack gen.",
	Bool:     true,
	Uint64:   math.MaxUint32 * 2,
	Slice:    []uint{1, 100, 10000, 1000000},
	ItemData: Item{ID: 1, Name: "abc", Effect: 7.89, Num: 999},
}
var num = 8

type Item struct {
	ID     int
	Name   string
	Effect float32
	Num    uint
}

func _TestA(t *testing.T) {
	e := func(interface{}) ([]byte, error) { return nil, nil }
	d := func([]byte, interface{}) (bool, error) { return false, nil }
	add()
	msgpackgen.SetEncodingOption(true)

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
	if v.Items != nil {
		return
	}

	n := 1
	v.Items = make([]Item, n)
	for i := 0; i < n; i++ {
		name := "item" + fmt.Sprint(i)
		item := Item{
			ID:     i,
			Name:   name,
			Effect: float32(i*i) / 3.0,
			Num:    uint(i * i * i * i),
		}
		v.Items[i] = item
	}
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

	switch v := i.(type) {
	case *structTest:
		_, err := decodestructTest(v, dec.NewDecoder(data), 0)
		return true, err

	case *Item:
		_, err := decodeItem(v, dec.NewDecoder(data), 0)
		return true, err

	}

	return false, nil
}

// todo : auto generate
func encode(i interface{}) ([]byte, error) {

	switch v := i.(type) {
	case structTest:
		e := encoding.NewEncoder()
		size, err := calcSizestructTest(v, e)
		if err != nil {
			return nil, err
		}
		e.MakeBytes(size)
		b, _, err := encodestructTest(v, e, 0)
		return b, err

	case Item:
		e := encoding.NewEncoder()
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

func calcSizestructTest(v structTest, encoder *encoding.Encoder) (int, error) {
	size := def.Byte1
	s, err := encoder.CalcStruct(num)
	if err != nil {
		return 0, err
	}
	size += s

	// todo : タグが設定されているパターン
	if !msgpack.StructAsArray {
		size += encoder.CalcString("A")
		size += encoder.CalcString("B")
		size += encoder.CalcString("String")
		size += encoder.CalcString("Bool")
		size += encoder.CalcString("Uint64")
		size += encoder.CalcString("Slice")
		size += encoder.CalcString("ItemData")
		size += encoder.CalcString("Items")
	}

	size += def.Byte1
	size += encoder.CalcInt(int64(v.A))

	size += def.Byte1
	size += encoder.CalcFloat32(float64(v.B))

	size += def.Byte1
	size += encoder.CalcString(v.String)

	size += def.Byte1
	size += encoder.CalcBool()

	size += def.Byte1
	size += encoder.CalcUint(v.Uint64)

	// todo : nilのパターン
	size += def.Byte1
	s, err = encoder.CalcSliceLength(len(v.Slice))
	if err != nil {
		return 0, err
	}
	size += s
	for _, v := range v.Slice {
		size += def.Byte1
		size += encoder.CalcUint(uint64(v))
	}

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

	return size, nil
}

func encodestructTest(v structTest, encoder *encoding.Encoder, offset int) ([]byte, int, error) {
	var err error
	offset = encoder.WriteStruct(num, offset)

	offset = encoder.WriteInt(int64(v.A), offset)
	offset = encoder.WriteFloat32(v.B, offset)
	offset = encoder.WriteString(v.String, offset)
	offset = encoder.WriteBool(v.Bool, offset)
	offset = encoder.WriteUint(v.Uint64, offset)

	// todo : nilのパターン
	offset = encoder.WriteSliceLength(len(v.Slice), offset)
	for _, vv := range v.Slice {
		offset = encoder.WriteUint(uint64(vv), offset)
	}

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
	return encoder.EncodedBytes(), offset, nil
}

func decodestructTest(v *structTest, decoder *dec.Decoder, offset int) (int, error) {

	// todo : mapの場合はここでstringをみてswitchする

	a := "abc"
	switch a {
	case "abc":
		fmt.Println("correct")

	case "b":
		fmt.Println("wrong")
	}

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
	}
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
	return offset, nil
}

///////////////////////////////

func calcSizeItem(v Item, encoder *encoding.Encoder) (int, error) {
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

	return size, nil
}

func encodeItem(v Item, encoder *encoding.Encoder, offset int) ([]byte, int, error) {
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
