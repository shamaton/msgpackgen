package main

import (
	"crypto/sha256"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/dave/jennifer/jen"
)

var analyzedStructs []analyzedStruct

const (
	pkTop = "github.com/shamaton/msgpackgen"
	pkEnc = "github.com/shamaton/msgpackgen/enc"
	pkDec = "github.com/shamaton/msgpackgen/dec"

	idEncoder = "encoder"
	idDecoder = "decoder"
)

// todo : tagをmapのcaseに使いつつ、変数に代入するようにしないといけない

var funcIdMap = map[string]string{}

type generator struct {
	fileSet              *token.FileSet
	targetPackages       map[string]bool
	file2Parse           map[string]*ast.File
	file2FullPackageName map[string]string
	file2PackageName     map[string]string
}

type analyzedStruct struct {
	PackageName string
	Name        string
	Fields      []analyzedField
}

type analyzedField struct {
	Name string
	Type types.Type
	Ast  *analyzedASTFieldType
}

func main() {
	dir := "../tetest/example"
	files := dirwalk(dir)

	// 最初にgenerate対象のパッケージをすべて取得

	// 構造体の解析時にgenerate対象でないパッケージを含んだ構造体がある場合
	// 出力対象にしない

	// 出力対象にしない構造体が見つからなくなるまで実行する

	g := generator{
		targetPackages:       map[string]bool{},
		file2Parse:           map[string]*ast.File{},
		file2FullPackageName: map[string]string{},
		file2PackageName:     map[string]string{},
	}
	g.getPackages(files)
	g.createAnalyzedStructs()
	g.generate(dir)
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	var abss []string
	for _, path := range paths {
		abs, err := filepath.Abs(path)
		if err != nil {
			panic(err)
		}
		abss = append(abss, abs)
	}
	return abss
}

func (g *generator) generate(dir string) {

	for _, st := range analyzedStructs {
		funcIdMap[st.PackageName] = fmt.Sprintf("%x", sha256.Sum256([]byte(st.PackageName)))
	}

	fmt.Println(funcIdMap)

	path := "msgpackgen/resolver"
	f := NewFilePath(path)

	f.Func().Id("init").Params().Block(
		Qual(pkTop, "SetResolver").Call(Id("encode"), Id("decode")),
	)

	g.decodeTopTemplate("decode", f).Block(
		If(Qual(pkTop, "StructAsArray").Call()).Block(
			Return(Id("decodeAsArray").Call(Id("data"), Id("i"))),
		).Else().Block(
			Return(Id("decodeAsMap").Call(Id("data"), Id("i"))),
		),
	)

	g.decodeTopTemplate("decodeAsArray", f).Block(
		Switch(Id("v").Op(":=").Id("i").Assert(Type())).Block(
			g.decodeAsArrayCases()...,
		),
		Return(False(), Nil()),
	)

	g.decodeTopTemplate("decodeAsMap", f).Block(
		Switch(Id("v").Op(":=").Id("i").Assert(Type())).Block(
			g.decodeAsMapCases()...,
		),
		Return(False(), Nil()),
	)

	g.encodeTopTemplate("encode", f).Block(
		If(Qual(pkTop, "StructAsArray").Call()).Block(
			Return(Id("encodeAsArray").Call(Id("i"))),
		).Else().Block(
			Return(Id("encodeAsMap").Call(Id("i"))),
		),
	)

	g.encodeTopTemplate("encodeAsArray", f).Block(
		Switch(Id("v").Op(":=").Id("i").Assert(Type())).Block(
			g.encodeAsArrayCases()...,
		),
		Return(Nil(), Nil()),
	)

	g.encodeTopTemplate("encodeAsMap", f).Block(
		Switch(Id("v").Op(":=").Id("i").Assert(Type())).Block(
			g.encodeAsMapCases()...,
		),
		Return(Nil(), Nil()),
	)

	for _, st := range analyzedStructs {
		st.calcFunction(f)
	}

	fmt.Printf("%#v", f)

	d := dir + "/" + path
	if err := os.MkdirAll(d, 0777); err != nil {
		fmt.Println(err)
	}

	file, err := os.Create(d + "/resolver.go")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%#v", f)
	fmt.Println(err)
}

func (g *generator) decodeTopTemplate(name string, f *File) *Statement {
	return f.Func().Id(name).Params(Id("data").Index().Byte(), Id("i").Interface()).Params(Bool(), Error())
}

func (g *generator) encodeTopTemplate(name string, f *File) *Statement {
	return f.Func().Id(name).Params(Id("i").Interface()).Params(Index().Byte(), Error())
}

func (g *generator) encodeAsArrayCases() []Code {
	var states []Code
	for _, v := range analyzedStructs {
		for _, vv := range []string{"", "*"} {
			states = append(states, Case(Op(vv).Qual(v.PackageName, v.Name)).Block(
				Id(idEncoder).Op(":=").Qual(pkEnc, "NewEncoder").Call(),
				List(Id("size"), Err()).Op(":=").Id(v.calcArraySizeFuncName()).Call(Id(vv+"v"), Id(idEncoder)),
				If(Err().Op("!=").Nil()).Block(
					Return(Nil(), Err()),
				),
				Id(idEncoder).Dot("MakeBytes").Call(Id("size")),
				List(Id("b"), Id("offset"), Err()).Op(":=").Id(v.encodeArrayFuncName()).Call(Id(vv+"v"), Id(idEncoder), Lit(0)),
				If(Err().Op("!=").Nil()).Block(
					Return(Nil(), Err()),
				),
				If(Id("size").Op("!=").Id("offset")).Block(
					Return(Nil(), Qual("fmt", "Errorf").Call(Lit(v.Name+" size / offset different %d : %d"), Id("size"), Id("offset"))),
				),
			))
		}
	}
	return states
}

func (g *generator) encodeAsMapCases() []Code {
	var states []Code
	for _, v := range analyzedStructs {
		for _, vv := range []string{"", "*"} {
			states = append(states, Case(Op(vv).Qual(v.PackageName, v.Name)).Block(
				Id(idEncoder).Op(":=").Qual(pkEnc, "NewEncoder").Call(),
				List(Id("size"), Err()).Op(":=").Id(v.calcMapSizeFuncName()).Call(Id(vv+"v"), Id(idEncoder)),
				If(Err().Op("!=").Nil()).Block(
					Return(Nil(), Err()),
				),
				Id(idEncoder).Dot("MakeBytes").Call(Id("size")),
				List(Id("b"), Id("offset"), Err()).Op(":=").Id(v.encodeMapFuncName()).Call(Id(vv+"v"), Id(idEncoder), Lit(0)),
				If(Err().Op("!=").Nil()).Block(
					Return(Nil(), Err()),
				),
				If(Id("size").Op("!=").Id("offset")).Block(
					Return(Nil(), Qual("fmt", "Errorf").Call(Lit(v.Name+" size / offset different %d : %d"), Id("size"), Id("offset"))),
				),
			))
		}
	}
	return states
}

func (g *generator) decodeAsArrayCases() []Code {
	var states []Code
	for _, v := range analyzedStructs {
		for _, vv := range []string{"", "*"} {
			states = append(states, Case(Op("*"+vv).Qual(v.PackageName, v.Name)).Block(
				List(Id("_"), Err()).Op(":=").Id(v.decodeArrayFuncName()).Call(Id(vv+"v"), Qual(pkDec, "NewDecoder").Call(Id("data")), Id("0")),
				Return(True(), Err())))
		}
	}
	return states
}

func (g *generator) decodeAsMapCases() []Code {
	var states []Code
	for _, v := range analyzedStructs {
		for _, vv := range []string{"", "*"} {
			states = append(states, Case(Op("*"+vv).Qual(v.PackageName, v.Name)).Block(
				List(Id("_"), Err()).Op(":=").Id(v.decodeMapFuncName()).Call(Id(vv+"v"), Qual(pkDec, "NewDecoder").Call(Id("data")), Id("0")),
				Return(True(), Err())))
		}
	}
	return states
}

func (as *analyzedStruct) calcArraySizeFuncName() string {
	return fmt.Sprintf("decodeArray%s_%s", as.Name, funcIdMap[as.PackageName])
}

func (as *analyzedStruct) calcMapSizeFuncName() string {
	return fmt.Sprintf("decodeMap%s_%s", as.Name, funcIdMap[as.PackageName])
}

func (as *analyzedStruct) encodeArrayFuncName() string {
	return fmt.Sprintf("encodeArray%s_%s", as.Name, funcIdMap[as.PackageName])
}

func (as *analyzedStruct) encodeMapFuncName() string {
	return fmt.Sprintf("encodeMap%s_%s", as.Name, funcIdMap[as.PackageName])
}

func (as *analyzedStruct) decodeArrayFuncName() string {
	return fmt.Sprintf("decodeArray%s_%s", as.Name, funcIdMap[as.PackageName])
}

func (as *analyzedStruct) decodeMapFuncName() string {
	return fmt.Sprintf("decodeMap%s_%s", as.Name, funcIdMap[as.PackageName])
}
