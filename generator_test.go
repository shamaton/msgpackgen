package main

import (
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	iDir  = "."
	iFile = ""
	oDir  = "."
	oFile = defaultFileName
	ptr   = defaultPointerLevel
)

func TestGenerateCodeErrorInput(t *testing.T) {
	{
		d := "./noname"

		err := generate(d, iFile, oDir, oFile, ptr, false, true, false, false, io.Discard)
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

		err := generate(d, f, oDir, oFile, ptr, false, true, false, false, io.Discard)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "at same time") {
			t.Fatal(err)
		}
	}
	{
		d := "main.go"

		err := generate(d, iFile, oDir, oFile, ptr, false, true, false, false, io.Discard)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "is not directory") {
			t.Fatal(err)
		}
	}
	{
		f := "foo.go"

		err := generate(iDir, f, oDir, oFile, ptr, false, true, false, false, io.Discard)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "input file error") {
			t.Fatal(err)
		}
	}
	{
		f := "internal"

		err := generate(iDir, f, oDir, oFile, ptr, false, true, false, false, io.Discard)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "is a directory") {
			t.Fatal(err)
		}
	}
	{
		f := "./testdata/test.sh"

		err := generate(iDir, f, oDir, oFile, ptr, false, true, false, false, io.Discard)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "is not .go file") {
			t.Fatal(err)
		}
	}
	{
		d := "./noname"

		err := generate(iDir, iFile, d, oFile, ptr, false, true, false, false, io.Discard)
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		d := "./main.go"

		err := generate(iDir, iFile, d, oFile, ptr, false, true, false, false, io.Discard)
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

	err = generate(iDir, iFile, oDir, oFile, ptr, true, true, false, false, io.Discard)
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

	err := generate(iDir, f, oDir, oFile, ptr, false, true, false, false, io.Discard)
	if err == nil {
		t.Fatal("error has to return")
	}
	if !strings.Contains(err.Error(), "duplicate tags") {
		t.Fatal(err)
	}
}

func TestGenerateCodeDryRun(t *testing.T) {

	err := generate(iDir, iFile, "", oFile, -1, false, true, false, false, io.Discard)
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
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	slashWD := filepath.ToSlash(wd)
	if !strings.Contains(slashWD, "/src/") {
		return
	}
	goPath := strings.SplitN(slashWD, "/src", 2)[0]
	err = os.Setenv("GOPATH", filepath.FromSlash(goPath))
	if err != nil {
		t.Fatal(err)
	}
	err = generate(iDir, iFile, oDir, oFile, ptr, true, true, false, false, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	generated, err := os.ReadFile("resolver_test.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(generated), ".IsCodeNilChecked(") {
		t.Fatal("generated code must use IsCodeNilChecked")
	}
	if strings.Contains(string(generated), ".IsCodeNil(") {
		t.Fatal("generated code must not use legacy IsCodeNil")
	}
}
