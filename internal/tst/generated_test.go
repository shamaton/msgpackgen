package tst_test

import (
	"html/template"
	"log"
	"math"
	"os"
	"testing"

	"github.com/shamaton/msgpack"

	"github.com/shamaton/msgpackgen/internal/tst"
)

func TestMain(m *testing.M) {
	tst.RegisterGeneratedResolver()

	code := m.Run()

	resetGeneratedCode()

	os.Exit(code)
}

func resetGeneratedCode() {
	tpl := template.Must(template.New("").Parse(`package tst

import "fmt"

func RegisterGeneratedResolver() {
	fmt.Println("this is dummy.")
}
`))

	file, err := os.OpenFile("./resolver.msgpackgen.go", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = tpl.Execute(file, tpl)
	if err != nil {
		log.Fatal(err)
	}
}

func TestInt(t *testing.T) {
	f := func(t *testing.T, v tst.A) {
		b1, err := msgpack.Encode(v)
		if err != nil {
			t.Error(err)
		}
		b2, err := msgpack.EncodeStructAsArray(v)
		if err != nil {
			t.Error(err)
		}
		var v1, v2 tst.A
		err = msgpack.Decode(b1, &v1)
		if err != nil {
			t.Error(err)
		}
		err = msgpack.DecodeStructAsArray(b2, &v2)
		if err != nil {
			t.Error(v2)
		}

		if v.Int != v1.Int {
			t.Error("not equal v1", v.Int, v1.Int)
		}
		if v.Int != v2.Int {
			t.Error("not equal v2", v.Int, v2.Int)
		}

		// todo : ダブルポインタ、トリプル
	}

	v := tst.A{Int: -8}
	f(t, v)
	v.Int = -30108
	f(t, v)
	v.Int = -1030108
	f(t, v)
	v.Int = math.MinInt64 + 12345
	f(t, v)

}
