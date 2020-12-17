package main

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
)

var analyzedStructs []analyzedStruct

func main() {
	fmt.Println(findStructs("msgpackgen_test.go"))
	GenerateCode()
}

func GenerateCode() {
	packageName := "msgpackgen"

	f := NewFilePath("a.b/c")
	f.Func().Id("init").Params().Block(
		Qual("a.b/c", "Foo").Call().Comment("Local package - name is omitted."),
		Qual("d.e/f", "Bar").Call().Comment("Import is automatically added."),
		Qual("g.h/f", "Baz").Call().Comment("Colliding package name is renamed."),
	)

	decodeTopTemplate("decode", f).Block(
		// todo : Qual
		If(Id(packageName + ".StructAsArray").Call()).Block(
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
		If(Id(packageName + ".StructAsArray").Call()).Block(
			Return(Id("encodeAsArray").Call(Id("i"))),
		).Else().Block(
			Return(Id("encodeAsMap").Call(Id("i"))),
		),
	)

	encodeTopTemplate("encodeAsArray", f).Block(
		Return(Nil(), Nil()),
	)

	decodeTopTemplate("encodeAsMap", f).Block(
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
	return f.Func().Id(name).Params(Id("i").Interface()).Params(Id("data").Index().Byte(), Error())
}

func cases() []Code {
	var states []Code
	for _, v := range analyzedStructs {
		states = append(states, Case(Id("*"+v.Name)).Block(
			// 			_, err := decodestructTest(v, dec.NewDecoder(data), 0)
			List(Id("_"), Err()).Op(":=").Id("decode"+v.Name).Call(Id("v"), Qual("msgpackgen", "NewDecoder").Call(Id("data")), Id("0")),
			Return(True(), Err())))
	}
	return states
}

func calcFunction(st analyzedStruct, f *File) {

	f.Func().Id("calcArraySize"+st.Name).Params(Id("v").Id(st.Name), Id("encoder").Op("*").Qual("github.com/shamaton/msgpackgen/enc", "Encoder")).Params(Int(), Error()).Block(
		Id("size").Op(":=").Lit(0),
	)
}
