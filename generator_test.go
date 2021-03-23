package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/shamaton/msgpackgen/msgpack"
)

var (
	iDir  = "."
	iFile = ""
	oDir  = "."
	oFile = defaultFileName
	ptr   = defaultPointerLevel
)

func TestMain(m *testing.M) {
	testBeforeRegister()
	RegisterGeneratedResolver()

	code := m.Run()

	os.Exit(code)
}

func testBeforeRegister() {
	{
		v := rand.Int()
		b, err := msgpack.Marshal(v)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var vv int
		err = msgpack.Unmarshal(b, &vv)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if v != vv {
			fmt.Println(v, vv, "different")
			os.Exit(1)
		}
	}
	msgpack.SetStructAsArray(true)
	{
		v := rand.Int()
		b, err := msgpack.Marshal(v)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var vv int
		err = msgpack.Unmarshal(b, &vv)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if v != vv {
			fmt.Println(v, vv, "different")
			os.Exit(1)
		}
	}
	msgpack.SetStructAsArray(false)
}

func TestGenerateCodeErrorInput(t *testing.T) {
	{
		d := "./noname"

		err := generate(d, iFile, oDir, oFile, ptr, false, true, false, false, ioutil.Discard)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "input directory error") {
			t.Fatal(err)
		}
	}
	{
		d := "./noname"
		f := "foo.go"

		err := generate(d, f, oDir, oFile, ptr, false, true, false, false, ioutil.Discard)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "at same time") {
			t.Fatal(err)
		}
	}
	{
		d := "main.go"

		err := generate(d, iFile, oDir, oFile, ptr, false, true, false, false, ioutil.Discard)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "is not directory") {
			t.Fatal(err)
		}
	}
	{
		f := "foo.go"

		err := generate(iDir, f, oDir, oFile, ptr, false, true, false, false, ioutil.Discard)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "input file error") {
			t.Fatal(err)
		}
	}
	{
		f := "internal"

		err := generate(iDir, f, oDir, oFile, ptr, false, true, false, false, ioutil.Discard)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "is a directory") {
			t.Fatal(err)
		}
	}
	{
		f := "./testdata/test.sh"

		err := generate(iDir, f, oDir, oFile, ptr, false, true, false, false, ioutil.Discard)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "is not .go file") {
			t.Fatal(err)
		}
	}
	{
		d := "./noname"

		err := generate(iDir, iFile, d, oFile, ptr, false, true, false, false, ioutil.Discard)
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		d := "./main.go"

		err := generate(iDir, iFile, d, oFile, ptr, false, true, false, false, ioutil.Discard)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "path is not directory") {
			t.Fatal(err)
		}
	}
}

func TestGenerateCodeGoPathOutside(t *testing.T) {

	g := os.Getenv("GOPATH")
	path := os.Getenv("PATH")
	err := os.Setenv("GOPATH", path)
	if err != nil {
		t.Fatal(err)
	}

	err = generate(iDir, iFile, oDir, oFile, ptr, true, true, false, false, ioutil.Discard)
	if err == nil {
		t.Fatal("error has to return")
	}
	if !strings.Contains(err.Error(), "outside gopath") {
		t.Fatal(err)
	}

	err = os.Setenv("GOPATH", g)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateCodeDuplicateTag(t *testing.T) {

	f := "./testdata/def.go"

	err := generate(iDir, f, oDir, oFile, ptr, false, true, false, false, ioutil.Discard)
	if err == nil {
		t.Fatal("error has to return")
	}
	if !strings.Contains(err.Error(), "duplicate tags") {
		t.Fatal(err)
	}
}

func TestGenerateCodeDryRun(t *testing.T) {

	err := generate(iDir, iFile, "", oFile, -1, false, true, false, false, ioutil.Discard)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateCodeOK(t *testing.T) {
	var err error
	err = flag.CommandLine.Set("strict", "true")
	if err != nil {
		t.Fatal(err)
	}
	err = flag.CommandLine.Set("v", "true")
	if err != nil {
		t.Fatal(err)
	}
	err = flag.CommandLine.Set("pointer", "2")
	if err != nil {
		t.Fatal(err)
	}
	err = flag.CommandLine.Set("output-file", "resolver_test.go")
	if err != nil {
		t.Fatal(err)
	}

	// diff resolver_test.go main.go | wc -l
	main()

	// gopath
	err = generate(iDir, iFile, oDir, oFile, ptr, true, true, false, false, ioutil.Discard)
	if err != nil {
		t.Fatal(err)
	}
}
