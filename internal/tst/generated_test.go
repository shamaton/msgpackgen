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

func TestInt(t *testing.T) {
	f := func(t *testing.T, v tst.A) {
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
		var v1, v2 tst.A
		err = msgpack.Decode(b1, &v1)
		if err != nil {
			t.Error(err)
		}
		err = msgpack.DecodeAsArray(b2, &v2)
		if err != nil {
			t.Error(err)
		}

		// encode single pointer
		b3, err := msgpack.Encode(&v)
		if err != nil {
			t.Error(err)
		}
		b4, err := msgpack.EncodeAsArray(&v)
		if err != nil {
			t.Error(err)
		}

		// decode double pointer
		v3, v4 := new(tst.A), new(tst.A)
		err = msgpack.Decode(b3, &v3)
		if err != nil {
			t.Error(err)
		}
		err = msgpack.DecodeAsArray(b4, &v4)
		if err != nil {
			t.Error(err)
		}

		if v.Int != v1.Int || v.Uint != v1.Uint {
			t.Error("not equal v1", v, v1)
		}
		if v.Int != v2.Int || v.Uint != v2.Uint {
			t.Error("not equal v2", v, v2)
		}
		if v.Int != v3.Int || v.Uint != v3.Uint {
			t.Error("not equal v3", v, v3)
		}
		if v.Int != v4.Int || v.Uint != v4.Uint {
			t.Error("not equal v4", v, v4)
		}

		//// NG

		// encode double pointer
		b5, err := msgpack.Encode(&v3)
		if err != nil && !strings.Contains(err.Error(), "strict") {
			t.Error(err)
		}
		b6, err := msgpack.EncodeAsArray(&v4)
		if err != nil && !strings.Contains(err.Error(), "strict") {
			t.Error(err)
		}

		// decode triple pointer
		_v5, _v6 := new(tst.A), new(tst.A)
		v5, v6 := &_v5, &_v6
		err = msgpack.Decode(b5, &v5)
		if err != nil && !strings.Contains(err.Error(), "strict") {
			t.Error(err)
		}
		err = msgpack.DecodeAsArray(b6, &v6)
		if err != nil && !strings.Contains(err.Error(), "strict") {
			t.Error(err)
		}
	}

	v := tst.A{Int: -8, Uint: 8}
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
