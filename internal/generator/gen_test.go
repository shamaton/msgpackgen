package generator

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/shamaton/msgpackgen/internal/generator/structure"
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

	err = os.MkdirAll("tmp/resolver.go", 0777)
	if err != nil {
		t.Fatal("unexpected error", err)
	}

	g.outputDir = "tmp"
	err = g.output(nil, "resolver.go")
	if err == nil {
		t.Error("error should occur")
	}
	if err != nil && !strings.Contains(err.Error(), "is a directory") {
		t.Error("something wrong", err)
	}

	err = os.RemoveAll("tmp/resolver.go")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateCodeRegistersToResolver(t *testing.T) {
	oldAnalyzedStructs := analyzedStructs
	analyzedStructs = nil
	t.Cleanup(func() {
		analyzedStructs = oldAnalyzedStructs
	})

	g := generator{outputJenFilePath: "resolver_test"}
	got := fmt.Sprintf("%#v", g.generateCode())

	for _, want := range []string{
		"msgpack.SetResolver(___encodeAsMap, ___encodeAsArray, ___decodeAsMap, ___decodeAsArray)",
		"msgpack.SetToResolver(___encodeAsMapTo, ___encodeAsArrayTo)",
		"func ___encodeAsMapTo(i any, buf []byte) ([]byte, bool, error)",
		"func ___encodeAsArrayTo(i any, buf []byte) ([]byte, bool, error)",
		"b, handled, err := ___encodeAsMapTo(i, nil)",
		"return buf, false, nil",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("generated code does not contain %q:\n%s", want, got)
		}
	}
}

func TestGenerateCodeUsesStatelessStructEncoder(t *testing.T) {
	oldAnalyzedStructs := analyzedStructs
	analyzedStructs = []*structure.Structure{
		{
			ImportPath: "github.com/shamaton/msgpackgen",
			Name:       "generatedFixture",
			NoUseQual:  true,
			Fields: []structure.Field{
				{
					Name: "Value",
					Tag:  "Value",
					Node: structure.CreateIdentNode(ast.NewIdent("int"), nil),
				},
				{
					Name: "At",
					Tag:  "At",
					Node: structure.CreateStructNode("time", "time", "Time", nil),
				},
			},
		},
	}
	t.Cleanup(func() {
		analyzedStructs = oldAnalyzedStructs
	})

	g := generator{outputJenFilePath: "resolver_test"}
	got := fmt.Sprintf("%#v", g.generateCode())

	for _, want := range []string{
		"func ___calcArraySizegeneratedFixture_",
		"(v generatedFixture) (int, error)",
		"func ___encodeArraygeneratedFixture_",
		"v generatedFixture, buf []byte, offset int",
		"enc.CalcIntMax(v.Value)",
		"enc.WriteIntTo(buf, v.Value, offset)",
		"enc.CalcTimeMax(v.At)",
		"enc.WriteTimeTo(buf, v.At, offset)",
		"offset += copy(buf[offset:offset+6], \"\\xa5Value\")",
		"offset += copy(buf[offset:offset+3], \"\\xa2At\")",
		"enc.RequireAt(buf, start, size)",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("generated code does not contain %q:\n%s", want, got)
		}
	}

	for _, unwanted := range []string{
		"enc.NewEncoder()",
		"MakeBytes",
		"EncodedBytes",
		"enc.WriteStringFixTo(buf, \"Value\"",
		"enc.WriteStringFixTo(buf, \"At\"",
	} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("generated code unexpectedly contains %q:\n%s", unwanted, got)
		}
	}
}
