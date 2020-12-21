package main

import (
	"crypto/sha256"
	"fmt"
	"go/types"
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
	FullPath string
}

type analyzedStruct struct {
	PackageName string
	Name        string
	Fields      []analyzedField
}

type analyzedField struct {
	Name string
	Type types.Type
}

func main() {
	g := new(generator)
	// todo : ここで対象のフォルダを再帰的に見て、収集
	fileName := "msgpackgen_struct.go"
	path, err := filepath.Abs(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	g.FullPath = path
	g.findStructs(path)
	g.generate()
}

func (g *generator) generate() {

	for _, st := range analyzedStructs {
		funcIdMap[st.PackageName] = fmt.Sprintf("%x", sha256.Sum256([]byte(st.PackageName)))
	}

	fmt.Println(funcIdMap)

	f := NewFilePath("msgpackgen/resolver")

	f.Func().Id("init").Params().Block(
		Qual(pkTop, "SetResolver").Call(Id("encode"), Id("decode")),
	)

	decodeTopTemplate("decode", f).Block(
		If(Qual(pkTop, "StructAsArray").Call()).Block(
			Return(Id("decodeAsArray").Call(Id("data"), Id("i"))),
		).Else().Block(
			Return(Id("decodeAsMap").Call(Id("data"), Id("i"))),
		),
	)

	decodeTopTemplate("decodeAsArray", f).Block(
		Switch(Id("v").Op(":=").Id("i").Assert(Type())).Block(
			cases()...,
		),
		Return(False(), Nil()),
	)

	decodeTopTemplate("decodeAsMap", f).Block(
		Return(False(), Nil()),
	)

	encodeTopTemplate("encode", f).Block(
		If(Qual(pkTop, "StructAsArray").Call()).Block(
			Return(Id("encodeAsArray").Call(Id("i"))),
		).Else().Block(
			Return(Id("encodeAsMap").Call(Id("i"))),
		),
	)

	encodeTopTemplate("encodeAsArray", f).Block(
		Return(Nil(), Nil()),
	)

	encodeTopTemplate("encodeAsMap", f).Block(
		Return(Nil(), Nil()),
	)

	for _, st := range analyzedStructs {
		st.calcFunction(f)
	}

	fmt.Printf("%#v", f)
}

func decodeTopTemplate(name string, f *File) *Statement {
	return f.Func().Id(name).Params(Id("data").Index().Byte(), Id("i").Interface()).Params(Bool(), Error())
}

func encodeTopTemplate(name string, f *File) *Statement {
	return f.Func().Id(name).Params(Id("i").Interface()).Params(Index().Byte(), Error())
}

func cases() []Code {
	var states []Code
	for _, v := range analyzedStructs {
		states = append(states, Case(Op("*").Qual(v.PackageName, v.Name)).Block(
			List(Id("_"), Err()).Op(":=").Id(v.decodeArrayFuncName()).Call(Id("v"), Qual(pkDec, "NewDecoder").Call(Id("data")), Id("0")),
			Return(True(), Err())))
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
