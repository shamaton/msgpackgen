package main

import (
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
)

func main() {
	fileName := "msgpackgen_struct.go"
	path, err := filepath.Abs(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(findStructs(path))
	GenerateCode()
}

func GenerateCode() {

	f := NewFilePath("msgpackgen/resolver")

	f.Func().Id("init").Params().Block(
		Qual(pkTop, "SetResolver").Call(Id("encode"), Id("decode")),
	)

	decodeTopTemplate("decode", f).Block(
		// todo : Qual
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
		calcFunction(st, f)
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
		states = append(states, Case(Id("*"+v.Name)).Block(
			// 			_, err := decodestructTest(v, dec.NewDecoder(data), 0)
			List(Id("_"), Err()).Op(":=").Id("decode"+v.Name).Call(Id("v"), Qual(pkDec, "NewDecoder").Call(Id("data")), Id("0")),
			Return(True(), Err())))
	}
	return states
}

func calcFunction(st analyzedStruct, f *File) {
	v := "v"
	enc := "encoder"

	errTemplate := If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	)

	addSizePattern1 := func(funcName string, params ...Code) Code {
		return Id("size").Op("+=").Id(enc).Dot(funcName).Call(params...)
	}

	addSizePattern2 := func(funcName string, params ...Code) []Code {
		return []Code{
			List(Id("s"), Err()).Op(":=").Id(enc).Dot(funcName).Call(params...),
			errTemplate,
			Id("size").Op("+=").Id("s"),
		}

	}

	codes := make([]Code, 0)
	codes = append(codes, Id("size").Op(":=").Lit(0))
	codes = append(codes, Block(addSizePattern2("CalcStructHeader", Lit(len(st.Fields)))...))

	for _, field := range st.Fields {
		switch {
		case types.Identical(types.Typ[types.Int], field.Type):
			codes = append(codes, addSizePattern1("CalcInt", Id("int64").Call(Id("v").Dot(field.Name))))

			//case types.Identical(types.Typ[types.], field.Type):

		}
	}

	f.Func().Id("calcArraySize"+st.Name).Params(Id(v).Qual(st.PackageName, st.Name), Id(enc).Op("*").Qual(pkEnc, "Encoder")).Params(Int(), Error()).Block(
		codes...,
	)
}

type int struct {
}

var a int
