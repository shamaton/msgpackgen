package structure

import (
	. "github.com/dave/jennifer/jen"
	"github.com/shamaton/msgpackgen/internal/generator/ptn"
)

type sliceCodeGen struct {
}

func (st *Structure) createSliceCode(node *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	encodeChildName, decodeChildName := encodeFieldName+"v", decodeFieldName+""
	if isRootField(encodeFieldName) {
		encodeChildName = "vv"
	}
	if isRootField(decodeFieldName) {
		decodeChildName = "vv"
	}

	ca, cm, ea, em, da, dm := st.createFieldCode(node.Elm(), encodeChildName, decodeChildName)
	isChildByte := node.Elm().IsIdentical() && node.Elm().IdenticalName == "byte"

	g := sliceCodeGen{}

	cArray = g.createCalcCode(encodeFieldName, encodeChildName, isChildByte, ca)
	cMap = g.createCalcCode(encodeFieldName, encodeChildName, isChildByte, cm)

	eArray = g.createEncCode(encodeFieldName, encodeChildName, isChildByte, ea)
	eMap = g.createEncCode(encodeFieldName, encodeChildName, isChildByte, em)

	dArray = g.createDecCode(node, st.Others, decodeFieldName, decodeChildName, da)
	dMap = g.createDecCode(node, st.Others, decodeFieldName, decodeChildName, dm)
	return
}

func (g sliceCodeGen) createCalcCode(fieldName, childName string, isChildTypeByte bool, elmCodes []Code) []Code {

	blockCodes := createAddSizeErrCheckCode("CalcSliceLength", Len(Id(fieldName)), Lit(isChildTypeByte))
	blockCodes = append(blockCodes, For(List(Id("_"), Id(childName)).Op(":=").Range(). /*Op(ptrOp).*/ Id(fieldName)).Block(
		elmCodes...,
	))

	codes := make([]Code, 0)
	codes = append(codes, If(Id(fieldName).Op("!=").Nil()).Block(
		blockCodes...,
	).Else().Block(
		createAddSizeCode("CalcNil"),
	))
	return codes
}

func (g sliceCodeGen) createEncCode(fieldName, childName string, isChildTypeByte bool, elmCodes []Code) []Code {

	blockCodes := make([]Code, 0)
	blockCodes = append(blockCodes, Id("offset").Op("=").Id(ptn.IdEncoder).Dot("WriteSliceLength").Call(Len(Id(fieldName)), Id("offset"), Lit(isChildTypeByte)))
	blockCodes = append(blockCodes, For(List(Id("_"), Id(childName)).Op(":=").Range().Id(fieldName)).Block(
		elmCodes...,
	))

	codes := make([]Code, 0)
	codes = append(codes, If(Id(fieldName).Op("!=").Nil()).Block(
		blockCodes...,
	).Else().Block(
		Id("offset").Op("=").Id(ptn.IdEncoder).Dot("WriteNil").Call(Id("offset")),
	))
	return codes
}

func (g sliceCodeGen) createDecCode(node *Node, structures []*Structure, fieldName, childName string, elmCodes []Code) []Code {

	childLengthName := childName + "l"
	childIndexName := childName + "i"
	childChildName := childName + "v"

	blockCodes := make([]Code, 0)
	blockCodes = append(blockCodes, node.TypeJenChain(structures, Var().Id(childName)))
	blockCodes = append(blockCodes, Var().Id(childLengthName).Int())
	blockCodes = append(blockCodes, List(Id(childLengthName), Id("offset"), Err()).Op("=").Id(ptn.IdDecoder).Dot("SliceLength").Call(Id("offset")))
	blockCodes = append(blockCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	blockCodes = append(blockCodes, Id(childName).Op("=").Make(node.TypeJenChain(structures), Id(childLengthName)))

	elmCodes = append([]Code{node.Elm().TypeJenChain(structures, Var().Id(childChildName))}, elmCodes...)
	elmCodes = append(elmCodes, Id(childName).Index(Id(childIndexName)).Op("=").Id(childChildName))

	blockCodes = append(blockCodes, For(Id(childIndexName).Op(":=").Range().Id(childName)).Block(
		elmCodes...,
	))

	name := childName
	andOp := ""
	prtCount, _ := node.GetPointerInfo()
	if prtCount > 0 {
		andOp = "&"
	}
	for i := 0; i < prtCount-1; i++ {
		n := "_" + name
		blockCodes = append(blockCodes, Id(n).Op(":=").Op("&").Id(name))
		name = n
	}

	blockCodes = append(blockCodes, Id(fieldName).Op("=").Op(andOp).Id(name))

	var codes []Code
	if node.HasParent() && node.Parent.IsPointer() {
		codes = blockCodes
	} else {
		codes = append(codes, If(Op("!").Id(ptn.IdDecoder).Dot("IsCodeNil").Call(Id("offset"))).Block(
			blockCodes...,
		).Else().Block(
			Id("offset").Op("++"),
		))
	}
	return codes
}
