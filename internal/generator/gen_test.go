package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSearchGoMod(t *testing.T) {

	g := generator{}
	_, err := g.searchGoModFile("dummy", true)
	if err == nil {
		t.Fatal("error should occur")
	}

	sub := "no such file or directory"
	if runtime.GOOS == "windows" {
		sub = "The system cannot find"
	}
	if !strings.Contains(err.Error(), sub) {
		t.Fatal("something wrong", err)
	}

	sep := fmt.Sprintf("%s..", string(filepath.Separator))
	upper := filepath.Join(".", sep, sep, sep)
	_, err = g.searchGoModFile(upper, true)
	if err == nil {
		t.Fatal("error should occur")
	}

	if !strings.Contains(err.Error(), "not found go.mod") {
		t.Fatal("something wrong", err)
	}

	_, err = filepath.Abs("./fuga/hoge")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetModuleName(t *testing.T) {

	g := generator{}
	g.goModFilePath = "dummy"
	err := g.setModuleName()
	if err == nil {
		t.Fatal("error should occur")
	}

	sub := "no such file or directory"
	if runtime.GOOS == "windows" {
		sub = "The system cannot find"
	}
	if !strings.Contains(err.Error(), sub) {
		t.Fatal("something wrong", err)
	}

	g.goModFilePath = "gen.go"
	err = g.setModuleName()
	if err == nil {
		t.Fatal("error should occur")
	}
	if !strings.Contains(err.Error(), "not found module name in go.mod") {
		t.Fatal("something wrong", err)
	}
}

func TestOutput(t *testing.T) {
	g := generator{}

	err := g.output(nil, "")
	if err == nil {
		t.Fatal("error should occur")
	}

	sub := "no such file or directory"
	if runtime.GOOS == "windows" {
		sub = "The system cannot find"
	}
	if !strings.Contains(err.Error(), sub) {
		t.Fatal("something wrong", err)
	}

	g.outputDir = "tmp"
	err = g.output(nil, "")
	if err == nil {
		t.Error("error should occur")
	}
	if err != nil && !strings.Contains(err.Error(), "is a directory") {
		t.Error("something wrong", err)
	}

	err = os.Remove("tmp")
	if err != nil {
		t.Fatal(err)
	}
}
