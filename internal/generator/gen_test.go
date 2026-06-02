package generator

import (
	"fmt"
	"go/ast"
	"os"
	"os/exec"
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

func TestGenerateCodeProvidesPublicAPIs(t *testing.T) {
	oldAnalyzedStructs := analyzedStructs
	analyzedStructs = nil
	t.Cleanup(func() {
		analyzedStructs = oldAnalyzedStructs
	})

	g := generator{outputJenFilePath: "msgpackgen_test"}
	got := fmt.Sprintf("%#v", g.generateCode())

	for _, want := range []string{
		"func Marshal(v any) ([]byte, error)",
		"func MarshalAsMap(v any) ([]byte, error)",
		"func MarshalAsArray(v any) ([]byte, error)",
		"func Unmarshal(data []byte, v any) error",
		"func UnmarshalAsMap(data []byte, v any) error",
		"func UnmarshalAsArray(data []byte, v any) error",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("generated code does not contain %q:\n%s", want, got)
		}
	}
}

func TestGenerateCodeProvidesV1InternalEntrypoints(t *testing.T) {
	oldAnalyzedStructs := analyzedStructs
	analyzedStructs = nil
	t.Cleanup(func() {
		analyzedStructs = oldAnalyzedStructs
	})

	g := generator{outputJenFilePath: "msgpackgen_test"}
	got := fmt.Sprintf("%#v", g.generateCode())

	for _, want := range []string{
		"func ___marshalAsMapTo(v any, buf []byte) ([]byte, error)",
		"func ___marshalAsArrayTo(v any, buf []byte) ([]byte, error)",
		"func ___unmarshalAsMap(data []byte, v any) error",
		"func ___unmarshalAsArray(data []byte, v any) error",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("generated code does not contain %q:\n%s", want, got)
		}
	}
}

func TestGenerateCodeOmitsLegacyResolverAPIs(t *testing.T) {
	oldAnalyzedStructs := analyzedStructs
	analyzedStructs = nil
	t.Cleanup(func() {
		analyzedStructs = oldAnalyzedStructs
	})

	g := generator{outputJenFilePath: "msgpackgen_test"}
	got := fmt.Sprintf("%#v", g.generateCode())

	for _, unwanted := range []string{
		"func RegisterGeneratedResolver(",
		"func ___encode(",
		"func ___encodeAsArray(",
		"func ___encodeAsMap(",
		"func ___decode(",
		"func ___decodeAsArray(",
		"func ___decodeAsMap(",
	} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("generated code unexpectedly contains %q:\n%s", unwanted, got)
		}
	}
}

func TestGenerateCodeFallbackDependsOnStrict(t *testing.T) {
	oldAnalyzedStructs := analyzedStructs
	analyzedStructs = nil
	t.Cleanup(func() {
		analyzedStructs = oldAnalyzedStructs
	})

	nonStrictGenerator := generator{outputJenFilePath: "msgpackgen_test"}
	nonStrict := fmt.Sprintf("%#v", nonStrictGenerator.generateCode())
	for _, want := range []string{
		"errors.Is(err, msgpack.ErrUndefinedType)",
		"fallback.MarshalAsMap(v)",
		"fallback.MarshalAsArray(v)",
		"fallback.UnmarshalAsMap(data, v)",
		"fallback.UnmarshalAsArray(data, v)",
		"return nil, msgpack.ErrUndefinedType",
		"return msgpack.ErrUndefinedType",
	} {
		if !strings.Contains(nonStrict, want) {
			t.Fatalf("non-strict generated code does not contain %q:\n%s", want, nonStrict)
		}
	}

	strictGenerator := generator{outputJenFilePath: "msgpackgen_test", strict: true}
	strict := fmt.Sprintf("%#v", strictGenerator.generateCode())
	for _, want := range []string{
		"return nil, msgpack.ErrUndefinedType",
		"return msgpack.ErrUndefinedType",
	} {
		if !strings.Contains(strict, want) {
			t.Fatalf("strict generated code does not contain %q:\n%s", want, strict)
		}
	}
	for _, unwanted := range []string{
		"use strict option",
		"errors.Is(err, msgpack.ErrUndefinedType)",
		"fallback.MarshalAsMap(v)",
		"fallback.MarshalAsArray(v)",
		"fallback.UnmarshalAsMap(data, v)",
		"fallback.UnmarshalAsArray(data, v)",
		"fmt.Errorf(\"undefined type\")",
	} {
		if strings.Contains(strict, unwanted) {
			t.Fatalf("strict generated code unexpectedly contains %q:\n%s", unwanted, strict)
		}
	}
}

func TestGeneratedNonStrictFallbackCompilesAndRuns(t *testing.T) {
	oldAnalyzedStructs := analyzedStructs
	analyzedStructs = []*structure.Structure{
		{
			ImportPath: "example.com/nonstrict",
			Name:       "Known",
			NoUseQual:  true,
			Fields: []structure.Field{
				{
					Name: "I",
					Tag:  "I",
					Node: structure.CreateIdentNode(ast.NewIdent("int"), nil),
				},
			},
		},
	}
	for _, st := range analyzedStructs {
		st.Others = analyzedStructs
	}
	t.Cleanup(func() {
		analyzedStructs = oldAnalyzedStructs
	})

	g := generator{outputJenFilePath: "example.com/nonstrict"}
	generated := fmt.Sprintf("%#v", g.generateCode())

	dir := t.TempDir()
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}

	files := map[string]string{
		"go.mod": fmt.Sprintf(`module example.com/nonstrict

go 1.23

require (
	github.com/shamaton/msgpack/v3 v3.1.2
	github.com/shamaton/msgpackgen v0.0.0
)

replace github.com/shamaton/msgpackgen => %s
`, filepath.ToSlash(repoRoot)),
		"model.go": `package nonstrict

type Known struct {
	I int
}

type Unknown struct {
	S string
}
`,
		"msgpack.msgpackgen.go": generated,
		"nonstrict_test.go": `package nonstrict

import (
	"errors"
	"testing"

	"github.com/shamaton/msgpackgen/msgpack"
)

func TestNonStrictFallback(t *testing.T) {
	if _, err := MarshalAsMap(Known{I: 1}); err != nil {
		t.Fatal(err)
	}
	b, err := MarshalAsMap(Unknown{S: "fallback"})
	if err != nil {
		t.Fatal(err)
	}
	var got Unknown
	if err := UnmarshalAsMap(b, &got); err != nil {
		t.Fatal(err)
	}
	if got.S != "fallback" {
		t.Fatalf("got %q", got.S)
	}
	if _, err := ___marshalAsMapTo(Unknown{}, nil); !errors.Is(err, msgpack.ErrUndefinedType) {
		t.Fatalf("private encode error = %v, want ErrUndefinedType", err)
	}
	if err := ___unmarshalAsMap(b, &got); !errors.Is(err, msgpack.ErrUndefinedType) {
		t.Fatalf("private decode error = %v, want ErrUndefinedType", err)
	}
}
`,
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0666); err != nil {
			t.Fatal(err)
		}
	}

	cmd := exec.Command("go", "test", ".")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOCACHE="+filepath.Join(dir, "gocache"), "GOFLAGS=-mod=mod")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go test failed: %v\n%s", err, out)
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

	g := generator{outputJenFilePath: "msgpackgen_test"}
	got := fmt.Sprintf("%#v", g.generateCode())

	for _, want := range []string{
		"func ___calcArraySizegeneratedFixture_",
		"(v *generatedFixture) (int, error)",
		"func ___encodeArraygeneratedFixture_",
		"v *generatedFixture, buf []byte, offset int",
		"enc.CalcIntMax(v.Value)",
		"enc.WriteInt(buf, v.Value, offset)",
		"enc.CalcTimeMax(v.At)",
		"enc.WriteTime(buf, v.At, offset)",
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
		"enc.WriteStringFix(buf, \"Value\"",
		"enc.WriteStringFix(buf, \"At\"",
	} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("generated code unexpectedly contains %q:\n%s", unwanted, got)
		}
	}
}

func TestGenerateCodeUsesNoErrSizeForEligibleNamedStructs(t *testing.T) {
	oldAnalyzedStructs := analyzedStructs

	child := &structure.Structure{
		ImportPath: "github.com/shamaton/msgpackgen",
		Name:       "ChildFixture",
		NoUseQual:  true,
		Fields: []structure.Field{
			{
				Name: "ID",
				Tag:  "ID",
				Node: structure.CreateIdentNode(ast.NewIdent("int"), nil),
			},
			{
				Name: "Name",
				Tag:  "Name",
				Node: structure.CreateIdentNode(ast.NewIdent("string"), nil),
			},
		},
	}

	childNode := structure.CreateStructNode("github.com/shamaton/msgpackgen", "generator", "ChildFixture", nil)
	itemsNode := structure.CreateSliceNode(nil)
	itemsNode.SetKeyNode(structure.CreateStructNode("github.com/shamaton/msgpackgen", "generator", "ChildFixture", itemsNode))
	childPtrNode := structure.CreatePointerNode(nil)
	childPtrNode.SetKeyNode(structure.CreateStructNode("github.com/shamaton/msgpackgen", "generator", "ChildFixture", childPtrNode))
	parent := &structure.Structure{
		ImportPath: "github.com/shamaton/msgpackgen",
		Name:       "ParentFixture",
		NoUseQual:  true,
		Fields: []structure.Field{
			{
				Name: "Child",
				Tag:  "Child",
				Node: childNode,
			},
			{
				Name: "Items",
				Tag:  "Items",
				Node: itemsNode,
			},
			{
				Name: "ChildPtr",
				Tag:  "ChildPtr",
				Node: childPtrNode,
			},
		},
	}

	analyzedStructs = []*structure.Structure{child, parent}
	for _, st := range analyzedStructs {
		st.Others = analyzedStructs
	}
	t.Cleanup(func() {
		analyzedStructs = oldAnalyzedStructs
	})

	g := generator{outputJenFilePath: "msgpackgen_test"}
	got := fmt.Sprintf("%#v", g.generateCode())

	for _, want := range []string{
		"func ___calcArraySizeNoErrChildFixture_",
		"func ___calcArraySizeMaxNoErrChildFixture_",
		"func ___calcMapSizeNoErrChildFixture_",
		"func ___calcMapSizeMaxNoErrChildFixture_",
		"___calcArraySizeNoErrChildFixture_",
		"___calcArraySizeMaxNoErrChildFixture_",
		"___calcMapSizeNoErrChildFixture_",
		"___calcMapSizeMaxNoErrChildFixture_",
		"___calcArraySizeNoErrChildFixture_",
		"(&v.Child)",
		"vv := &v.Items[vvi]",
		"___encodeArrayChildFixture_",
		"(vv, buf, offset)",
		"size_v_ChildPtr := ___calcArraySizeNoErrChildFixture_",
		"size_v_ChildPtr := ___calcArraySizeMaxNoErrChildFixture_",
		"(v.ChildPtr)",
		"___encodeArrayChildFixture_",
		"(v.ChildPtr, buf, offset)",
		"(&vv[vvi], decoder, offset)",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("generated code does not contain %q:\n%s", want, got)
		}
	}

	if strings.Contains(got, "ChildPtrp := *v.ChildPtr") || strings.Contains(got, "vp := *v.ChildPtr") {
		t.Fatalf("generated code unexpectedly copies pointer child:\n%s", got)
	}
	if strings.Contains(got, "var vvv ChildFixture") || strings.Contains(got, "vv[vvi] = vvv") {
		t.Fatalf("generated code unexpectedly decodes slice child through a temporary:\n%s", got)
	}
}

func TestGenerateCodeOptimizesMapDecodeDispatch(t *testing.T) {
	oldAnalyzedStructs := analyzedStructs
	scoresNode := structure.CreateSliceNode(nil)
	scoresNode.SetKeyNode(structure.CreateIdentNode(ast.NewIdent("int"), scoresNode))
	valuesNode := structure.CreateMapNode(nil)
	valuesNode.SetKeyNode(structure.CreateIdentNode(ast.NewIdent("string"), valuesNode))
	valuesNode.SetValueNode(structure.CreateIdentNode(ast.NewIdent("uint"), valuesNode))
	analyzedStructs = []*structure.Structure{
		{
			ImportPath: "github.com/shamaton/msgpackgen",
			Name:       "MapDecodeDispatchFixture",
			NoUseQual:  true,
			Fields: []structure.Field{
				{
					Name: "ID",
					Tag:  "ID",
					Node: structure.CreateIdentNode(ast.NewIdent("int"), nil),
				},
				{
					Name: "Name",
					Tag:  "Name",
					Node: structure.CreateIdentNode(ast.NewIdent("string"), nil),
				},
				{
					Name: "Scores",
					Tag:  "Scores",
					Node: scoresNode,
				},
				{
					Name: "Values",
					Tag:  "Values",
					Node: valuesNode,
				},
			},
		},
	}
	t.Cleanup(func() {
		analyzedStructs = oldAnalyzedStructs
	})

	g := generator{outputJenFilePath: "msgpackgen_test"}
	got := fmt.Sprintf("%#v", g.generateCode())

	for _, want := range []string{
		"switch string(dataKey)",
		"case \"ID\":",
		"case \"Name\":",
		"case \"Scores\":",
		"vv[vvi], offset, err = decoder.AsInt(offset)",
		"case \"Values\":",
		"vv[kkv], offset, err = decoder.AsUint(offset)",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("generated code does not contain %q:\n%s", want, got)
		}
	}

	for _, unwanted := range []string{
		"keys := [][]byte",
		"fieldIndex",
		"for i, key := range keys",
		"var vvv int",
		"vv[vvi] = vvv",
		"var vvv uint",
		"vv[kkv] = vvv",
	} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("generated code unexpectedly contains %q:\n%s", unwanted, got)
		}
	}
}
