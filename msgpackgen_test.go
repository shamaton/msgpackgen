package msgpackgen_test

import (
	"fmt"
	"testing"

	"github.com/shamaton/msgpack"

	"github.com/shamaton/msgpack/def"

	"github.com/shamaton/msgpackgen"
	"github.com/shamaton/msgpackgen/dec"
	encoding "github.com/shamaton/msgpackgen/enc"
)

func TestA(t *testing.T) {
	e := func(interface{}) ([]byte, error) { return nil, nil }
	d := func([]byte, interface{}) (bool, error) { return false, nil }
	msgpackgen.SetEncodingOption(true)

	msgpackgen.SetResolver(e, d)
	check(t)
	msgpackgen.SetResolver(encode, decode)
	check(t)
}

func check(t *testing.T) {

	v := structTest{A: 123456789}
	b, err := msgpackgen.Encode(v)
	if err != nil {
		t.Error(err)
	}

	var vv structTest
	err = msgpackgen.Decode(b, &vv)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(vv, v, b)
	fmt.Printf("% x \n", b)
}

func BenchmarkMsgGenEncShamaton(b *testing.B) {

	msgpackgen.SetResolver(encode, decode)

	v := structTest{A: 123456789}
	for i := 0; i < b.N; i++ {
		_, err := msgpackgen.Encode(v)
		if err != nil {
			b.Error(err)
		}
	}

}

func BenchmarkMsgEncShamaton(b *testing.B) {

	v := structTest{A: 123456789}
	for i := 0; i < b.N; i++ {
		_, err := msgpack.Encode(v)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkMsgGenDecShamaton(b *testing.B) {

	msgpackgen.SetResolver(encode, decode)

	v := structTest{A: 123456789}
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

	v := structTest{A: 123456789}
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
		return true, decodestructTest(v, dec.NewDecoder(data))

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
		return encodestructTest(v, e)

	}
	return nil, nil
}

type structTest struct {
	A int
}

func calcSizestructTest(v structTest, encoder *encoding.Encoder) (int, error) {
	size := def.Byte1
	s, err := encoder.CalcStruct(1)
	if err != nil {
		return 0, err
	}
	size += s

	size += def.Byte1
	size += encoder.CalcInt(int64(v.A))
	return size, nil
}

func encodestructTest(v structTest, encoder *encoding.Encoder) ([]byte, error) {
	offset := 0
	offset = encoder.WriteStruct(1, offset)
	offset = encoder.WriteInt(int64(v.A), offset)
	return encoder.EncodedBytes(), nil
}

func decodestructTest(v *structTest, decoder *dec.Decoder) error {

	offset, err := decoder.CheckStruct(1, 0)
	if err != nil {
		return err
	}
	{
		var vv int64
		vv, offset, err = decoder.AsInt(offset)
		if err != nil {
			return err
		}
		v.A = int(vv)
	}
	return nil
}
