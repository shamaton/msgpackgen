package tst_test

import (
	"html/template"
	"log"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/shamaton/msgpackgen/internal/tst"
	"github.com/shamaton/msgpackgen/msgpack"
)

func TestMain(m *testing.M) {
	tst.RegisterGeneratedResolver()

	code := m.Run()

	// resetGeneratedCode()

	os.Exit(code)
}

func resetGeneratedCode() {
	tpl := template.Must(template.New("").Parse(`package tst

import "fmt"

func RegisterGeneratedResolver() {
	fmt.Println("this is dummy.")
}
`))

	file, err := os.Create("./resolver.msgpackgen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = tpl.Execute(file, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func marshal(v1, v2 interface{}) ([]byte, []byte, error, error) {
	b1, e1 := msgpack.Encode(v1)
	b2, e2 := msgpack.EncodeAsArray(v2)
	return b1, b2, e1, e2
}

func unmarshal(b1, b2 []byte, v1, v2 interface{}) (error, error) {
	return msgpack.Decode(b1, v1), msgpack.DecodeAsArray(b2, v2)
}

func TestPointer(t *testing.T) {
	f := func(t *testing.T, v tst.Int) {
		//// OK
		// encode single pointer
		b1, b2, err1, err2 := marshal(&v, &v)
		if err1 != nil {
			t.Error(err1)
		}
		if err2 != nil {
			t.Error(err2)
		}

		// decode double pointer
		v1, v2 := new(tst.Int), new(tst.Int)
		err1, err2 = unmarshal(b1, b2, &v1, &v2)
		if err1 != nil {
			t.Error(err1)
		}
		if err2 != nil {
			t.Error(err2)
		}

		// encode double pointer
		b3, b4, err1, err2 := marshal(&v1, &v2)
		if err1 != nil {
			t.Error(err1)
		}
		if err2 != nil {
			t.Error(err2)
		}

		// decode triple pointer
		_v3, _v4 := new(tst.Int), new(tst.Int)
		v3, v4 := &_v3, &_v4
		err1, err2 = unmarshal(b3, b4, &v3, &v4)
		if err1 != nil {
			t.Error(err1)
		}
		if err2 != nil {
			t.Error(err2)
		}

		if v.Int != v1.Int || v.Uint != v1.Uint {
			t.Error("not equal v1", v, v1)
		}
		if v.Int != v2.Int || v.Uint != v2.Uint {
			t.Error("not equal v2", v, v2)
		}
		if _v := *v3; v.Int != _v.Int || v.Uint != _v.Uint {
			t.Error("not equal v3", v, _v)
		}
		if _v := *v4; v.Int != _v.Int || v.Uint != _v.Uint {
			t.Error("not equal v4", v, _v)
		}

		//// NG

		// encode triple pointer
		b5, b6, err1, err2 := marshal(&v3, &v4)
		if err1 != nil && !strings.Contains(err1.Error(), "strict") {
			t.Error(err1)
		}
		if err2 != nil && !strings.Contains(err2.Error(), "strict") {
			t.Error(err2)
		}

		// decode quad pointer
		__v5, __v6 := new(tst.Int), new(tst.Int)
		_v5, _v6 := &__v5, &__v6
		v5, v6 := &_v5, &_v6
		err1, err2 = unmarshal(b5, b6, &v5, &v6)
		if err1 != nil && !strings.Contains(err1.Error(), "strict") {
			t.Error(err1)
		}
		if err2 != nil && !strings.Contains(err2.Error(), "strict") {
			t.Error(err2)
		}
	}

	v := tst.Int{Int: -1, Uint: 1}
	f(t, v)

}

func TestInt(t *testing.T) {
	f := func(t *testing.T, v tst.Int) {
		//// OK
		// encode value
		b1, err := msgpack.Encode(v)
		if err != nil {
			t.Error(err)
		}
		b2, err := msgpack.EncodeAsArray(v)
		if err != nil {
			t.Error(err)
		}

		// decode single pointer
		var v1, v2 tst.Int
		err = msgpack.Decode(b1, &v1)
		if err != nil {
			t.Error(err)
		}
		err = msgpack.DecodeAsArray(b2, &v2)
		if err != nil {
			t.Error(err)
		}

		if v.Int != v1.Int || v.Uint != v1.Uint {
			t.Error("not equal v1", v, v1)
		}
		if v.Int != v2.Int || v.Uint != v2.Uint {
			t.Error("not equal v2", v, v2)
		}

	}

	v := tst.Int{Int: -8, Uint: 8}
	f(t, v)
	v.Int = -108
	v.Uint = 130
	f(t, v)
	v.Int = -30108
	v.Uint = 30130
	f(t, v)
	v.Int = -1030108
	v.Uint = 1030130
	f(t, v)
	v.Int = math.MinInt64 + 12345
	v.Uint = math.MaxUint64 - 12345
	f(t, v)

}

func TestFloat(t *testing.T) {
	f := func(t *testing.T, v tst.Float) {
		//// OK
		// encode value
		b1, err := msgpack.Encode(v)
		if err != nil {
			t.Error(err)
		}
		b2, err := msgpack.EncodeAsArray(v)
		if err != nil {
			t.Error(err)
		}

		// decode single pointer
		var v1, v2 tst.Float
		err = msgpack.Decode(b1, &v1)
		if err != nil {
			t.Error(err)
		}
		err = msgpack.DecodeAsArray(b2, &v2)
		if err != nil {
			t.Error(err)
		}

		if v.Float32 != v1.Float32 || v.Float64 != v1.Float64 {
			t.Error("not equal v1", v, v1)
		}
		if v.Float32 != v2.Float32 || v.Float64 != v2.Float64 {
			t.Error("not equal v2", v, v2)
		}
	}

	v := tst.Float{Float32: 0, Float64: 0}
	f(t, v)
	v = tst.Float{Float32: -1, Float64: -1}
	f(t, v)
	v = tst.Float{Float32: math.SmallestNonzeroFloat32, Float64: math.SmallestNonzeroFloat64}
	f(t, v)
	v = tst.Float{Float32: math.MaxFloat32, Float64: math.MaxFloat64}
	f(t, v)

}
