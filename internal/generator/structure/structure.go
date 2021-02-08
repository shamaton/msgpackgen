package structure

import (
	"fmt"
	"go/ast"
	"math"

	. "github.com/dave/jennifer/jen"
	"github.com/shamaton/msgpackgen/internal/generator/ptn"
)

type Structure struct {
	ImportPath string
	Package    string
	Name       string
	Fields     []Field
	NoUseQual  bool

	Others []*Structure
	File   *ast.File

	CanGen  bool
	Reasons []string
}

type Field struct {
	Name string
	Tag  string
	Node *Node
}

func (st *Structure) CalcArraySizeFuncName() string {
	return st.createFuncName("calcArraySize")
}

func (st *Structure) CalcMapSizeFuncName() string {
	return st.createFuncName("calcMapSize")
}

func (st *Structure) EncodeArrayFuncName() string {
	return st.createFuncName("encodeArray")
}

func (st *Structure) EncodeMapFuncName() string {
	return st.createFuncName("encodeMap")
}

func (st *Structure) DecodeArrayFuncName() string {
	return st.createFuncName("decodeArray")
}

func (st *Structure) DecodeMapFuncName() string {
	return st.createFuncName("decodeMap")
}

func (st *Structure) createFuncName(prefix string) string {
	return createFuncName(prefix, st.Name, st.ImportPath)
}

func (st *Structure) CreateCode(f *File) {
	v := "v"

	calcStruct, encStructArray, encStructMap := st.createStructCode(len(st.Fields))

	calcArraySizeCodes := make([]Code, 0)
	calcArraySizeCodes = append(calcArraySizeCodes, Id("size").Op(":=").Lit(0))
	calcArraySizeCodes = append(calcArraySizeCodes, calcStruct)

	calcMapSizeCodes := make([]Code, 0)
	calcMapSizeCodes = append(calcMapSizeCodes, Id("size").Op(":=").Lit(0))
	calcMapSizeCodes = append(calcMapSizeCodes, calcStruct)

	encArrayCodes := make([]Code, 0)
	encArrayCodes = append(encArrayCodes, Var().Err().Error())
	encArrayCodes = append(encArrayCodes, encStructArray)

	encMapCodes := make([]Code, 0)
	encMapCodes = append(encMapCodes, Var().Err().Error())
	encMapCodes = append(encMapCodes, encStructMap)

	decArrayCodes := make([]Code, 0)
	decArrayCodes = append(decArrayCodes, List(Id("offset"), Err()).Op(":=").Id(ptn.IdDecoder).Dot("CheckStructHeader").Call(Lit(len(st.Fields)), Id("offset")))
	decArrayCodes = append(decArrayCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))

	decMapCodeSwitchCases := make([]Code, 0)

	for _, field := range st.Fields {
		fieldName := "v." + field.Name

		calcKeyStringCode, writeKeyStringCode := st.createKeyStringCode(field.Tag)
		calcMapSizeCodes = append(calcMapSizeCodes, calcKeyStringCode)
		encMapCodes = append(encMapCodes, writeKeyStringCode)

		cArray, cMap, eArray, eMap, dArray, dMap := st.createFieldCode(field.Node, fieldName, fieldName)
		calcArraySizeCodes = append(calcArraySizeCodes, cArray...)

		calcMapSizeCodes = append(calcMapSizeCodes, cMap...)

		encArrayCodes = append(encArrayCodes, eArray...)
		encMapCodes = append(encMapCodes, eMap...)

		decArrayCodes = append(decArrayCodes, dArray...)

		decMapCodeSwitchCases = append(decMapCodeSwitchCases, Case(Lit(field.Tag)).Block(
			append(dMap, Id("count").Op("++"))...,
		// dMap...,
		),
		)
	}

	// not use jump offset
	//decMapCodeSwitchCases = append(decMapCodeSwitchCases, Default().Block(
	//	Id("offset").Op("=").Id(ptn.IdDecoder).Dot("JumpOffset").Call(Id("offset")),
	//	),
	//)
	decMapCodeSwitchCases = append(decMapCodeSwitchCases, Default().Block(
		Return(Lit(0), Qual("fmt", "Errorf").Call(Lit("unknown key[%s] found"), Id("s"))),
	),
	)

	decMapCodes := make([]Code, 0)
	decMapCodes = append(decMapCodes, List(Id("offset"), Err()).Op(":=").Id(ptn.IdDecoder).Dot("CheckStructHeader").Call(Lit(len(st.Fields)), Id("offset")))
	decMapCodes = append(decMapCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	//decMapCodes = append(decMapCodes, Id("dataLen").Op(":=").Id(ptn.IdDecoder).Dot("Len").Call())
	//decMapCodes = append(decMapCodes, For(Id("count").Op("<").Id("dataLen").Block(
	decMapCodes = append(decMapCodes, Id("count").Op(":=").Lit(0))
	decMapCodes = append(decMapCodes, For(Id("count").Op("<").Lit(len(st.Fields)).Block(
		Var().Id("s").String(),
		List(Id("s"), Id("offset"), Err()).Op("=").Id(ptn.IdDecoder).Dot("AsString").Call(Id("offset")),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Switch(Id("s")).Block(
			decMapCodeSwitchCases...,
		),
	)))

	var firstEncParam, firstDecParam *Statement
	if st.NoUseQual {
		firstEncParam = Id(v).Id(st.Name)
		firstDecParam = Id(v).Op("*").Id(st.Name)
	} else {
		firstEncParam = Id(v).Qual(st.ImportPath, st.Name)
		firstDecParam = Id(v).Op("*").Qual(st.ImportPath, st.Name)
	}

	f.Comment(fmt.Sprintf("// calculate size from %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.CalcArraySizeFuncName()).Params(firstEncParam, Id(ptn.IdEncoder).Op("*").Qual(ptn.PkEnc, "Encoder")).Params(Int(), Error()).Block(
		append(calcArraySizeCodes, Return(Id("size"), Nil()))...,
	)

	f.Comment(fmt.Sprintf("// calculate size from %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.CalcMapSizeFuncName()).Params(firstEncParam, Id(ptn.IdEncoder).Op("*").Qual(ptn.PkEnc, "Encoder")).Params(Int(), Error()).Block(
		append(calcMapSizeCodes, Return(Id("size"), Nil()))...,
	)

	f.Comment(fmt.Sprintf("// encode from %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.EncodeArrayFuncName()).Params(firstEncParam, Id(ptn.IdEncoder).Op("*").Qual(ptn.PkEnc, "Encoder"), Id("offset").Int()).Params(Index().Byte(), Int(), Error()).Block(
		append(encArrayCodes, Return(Id(ptn.IdEncoder).Dot("EncodedBytes").Call(), Id("offset"), Err()))...,
	)

	f.Comment(fmt.Sprintf("// encode from %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.EncodeMapFuncName()).Params(firstEncParam, Id(ptn.IdEncoder).Op("*").Qual(ptn.PkEnc, "Encoder"), Id("offset").Int()).Params(Index().Byte(), Int(), Error()).Block(
		append(encMapCodes, Return(Id(ptn.IdEncoder).Dot("EncodedBytes").Call(), Id("offset"), Err()))...,
	)

	f.Comment(fmt.Sprintf("// decode to %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.DecodeArrayFuncName()).Params(firstDecParam, Id(ptn.IdDecoder).Op("*").Qual(ptn.PkDec, "Decoder"), Id("offset").Int()).Params(Int(), Error()).Block(
		append(decArrayCodes, Return(Id("offset"), Err()))...,
	)

	f.Comment(fmt.Sprintf("// decode to %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.DecodeMapFuncName()).Params(firstDecParam, Id(ptn.IdDecoder).Op("*").Qual(ptn.PkDec, "Decoder"), Id("offset").Int()).Params(Int(), Error()).Block(

		append(decMapCodes, Return(Id("offset"), Err()))...,
	)
}

func (st *Structure) createStructCode(fieldNum int) (Code, Code, Code) {

	suffix := ""
	if fieldNum <= 0x0f {
		suffix = "Fix"
	} else if fieldNum <= math.MaxUint16 {
		suffix = "16"
	} else if uint(fieldNum) <= math.MaxUint32 {
		suffix = "32"
	}

	return Id("size").Op("+=").Id(ptn.IdEncoder).Dot("CalcStructHeader" + suffix).Call(Lit(fieldNum)),
		Id("offset").Op("=").Id(ptn.IdEncoder).Dot(" WriteStructHeader"+suffix+"AsArray").Call(Lit(fieldNum), Id("offset")),
		Id("offset").Op("=").Id(ptn.IdEncoder).Dot(" WriteStructHeader"+suffix+"AsMap").Call(Lit(fieldNum), Id("offset"))
}

func (st *Structure) createKeyStringCode(v string) (Code, Code) {
	l := len(v)
	suffix := ""
	if l < 32 {
		suffix = "Fix"
	} else if l <= math.MaxUint8 {
		suffix = "8"
	} else if l <= math.MaxUint16 {
		suffix = "16"
	} else {
		suffix = "32"
	}

	return Id("size").Op("+=").Id(ptn.IdEncoder).Dot("CalcString" + suffix).Call(Lit(l)),
		Id("offset").Op("=").Id(ptn.IdEncoder).Dot("WriteString"+suffix).Call(Lit(v), Lit(l), Id("offset"))
}

func (st *Structure) createFieldCode(node *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	switch {
	case node.IsIdentical():
		cArray, cMap, eArray, eMap, dArray, dMap = st.createIdentCode(node, encodeFieldName, decodeFieldName)

	case node.IsSlice():
		cArray, cMap, eArray, eMap, dArray, dMap = st.createSliceCode(node, encodeFieldName, decodeFieldName)

	case node.IsArray():
		cArray, cMap, eArray, eMap, dArray, dMap = st.createArrayCode(node, encodeFieldName, decodeFieldName)

	case node.IsMap():
		cArray, cMap, eArray, eMap, dArray, dMap = st.createMapCode(node, encodeFieldName, decodeFieldName)

	case node.IsPointer():
		cArray, cMap, eArray, eMap, dArray, dMap = st.createPointerCode(node, encodeFieldName, decodeFieldName)

	case node.IsStruct():

		if node.ImportPath == "time" {
			cArray, cMap, eArray, eMap, dArray, dMap = st.createTimeCode(encodeFieldName, decodeFieldName, node)
		} else {
			cArray, cMap, eArray, eMap, dArray, dMap = st.createNamedCode(encodeFieldName, decodeFieldName, node)
		}
	}

	return
}
