package structure

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
	"github.com/shamaton/msgpackgen/internal/generator/ptn"
)

type arrayCodeGen struct {
}

func (st *Structure) createArrayCode(node *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	encodeChildName := encodeFieldName + "v"
	if isRootField(encodeFieldName) {
		encodeChildName = "vv"
	}

	decodeChildName := decodeFieldName + "v"
	if isRootField(decodeFieldName) {
		decodeChildName = "vv"
	}

	ca, cm, ea, em, da, dm := st.createFieldCode(node.Elm(), encodeChildName, decodeChildName)
	// double or more ...
	if isRecursiveChildArraySliceMap(node) {
		_, _, _, _, da, dm = st.createFieldCode(node.Elm(), encodeChildName, decodeChildName+"v")
	}
	isChildByte := node.Elm().IsIdentical() && node.Elm().IdenticalName == "byte" && node.Elm().AliasName == ""
	passChildPointer := isPointerLoopElement(node.Elm())

	g := arrayCodeGen{}
	cArray = g.createCalcCode(encodeFieldName, encodeChildName, isChildByte, passChildPointer, ca)
	cMap = g.createCalcCode(encodeFieldName, encodeChildName, isChildByte, passChildPointer, cm)

	eArray = g.createEncCode(encodeFieldName, encodeChildName, isChildByte, passChildPointer, ea)
	eMap = g.createEncCode(encodeFieldName, encodeChildName, isChildByte, passChildPointer, em)

	dArray = g.createDecCode(node, st.Others, decodeFieldName, decodeChildName, "decodeArray", da)
	dMap = g.createDecCode(node, st.Others, decodeFieldName, decodeChildName, "decodeMap", dm)

	return
}

func (g arrayCodeGen) createCalcCode(fieldName, childName string, isChildByte, passChildPointer bool, elmCodes []Code) []Code {
	blockCodes := createAddSizeErrCheckCode("CalcSliceLength", Len(Id(fieldName)), Lit(isChildByte))
	if isChildByte {
		blockCodes = append(blockCodes, Id("size").Op("+=").Len(Id(fieldName)))
	} else {
		blockCodes = append(blockCodes, createSequenceRangeCode(fieldName, childName, passChildPointer, elmCodes))
	}

	codes := make([]Code, 0)
	codes = append(codes, Block(
		blockCodes...,
	))
	return codes
}

func (g arrayCodeGen) createEncCode(fieldName, childName string, isChildByte, passChildPointer bool, elmCodes []Code) []Code {

	blockCodes := make([]Code, 0)
	blockCodes = append(blockCodes, Id("offset").Op("=").Qual(ptn.PkEnc, "WriteSliceLength").Call(Id("buf"), Len(Id(fieldName)), Id("offset"), Lit(isChildByte)))
	if isChildByte {
		blockCodes = append(blockCodes,
			Id("offset").Op("+=").Id("copy").Call(
				Id("buf").Index(Id("offset").Op(":").Id("offset").Op("+").Len(Id(fieldName))),
				Id(fieldName).Index(Op(":")),
			),
		)
	} else {
		blockCodes = append(blockCodes, createSequenceRangeCode(fieldName, childName, passChildPointer, elmCodes))
	}

	codes := make([]Code, 0)
	codes = append(codes, Block(
		blockCodes...,
	))
	return codes
}

func (g arrayCodeGen) createDecCode(node *Node, structures []*Structure, fieldName, childName, funcName string, elmCodes []Code) []Code {

	blockCodes := make([]Code, 0)
	blockCodes = append(blockCodes, node.TypeJenChain(structures, Var().Id(childName)))
	blockCodes = append(blockCodes, Var().Id(childName+"l").Int())
	blockCodes = append(blockCodes, List(Id(childName+"l"), Id("offset"), Err()).Op("=").Id(ptn.IdDecoder).Dot("SliceLength").Call(Id("offset")))
	blockCodes = append(blockCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	blockCodes = append(blockCodes, If(Id(childName+"l").Op(">").Id(fmt.Sprint(node.ArrayLen))).Block(
		Return(Lit(0), Qual("fmt", "Errorf").Call(Lit("length size(%d) is over array size(%d)"), Id(childName+"l"), Id(fmt.Sprint(node.ArrayLen)))),
	))

	if directCodes, ok := createDirectSequenceDecodeCode(node, childName, childName+"i", funcName); ok {
		elmCodes = directCodes
	} else {
		elmCodes = append([]Code{node.Elm().TypeJenChain(structures, Var().Id(childName+"v"))}, elmCodes...)
		elmCodes = append(elmCodes, Id(childName).Index(Id(childName+"i")).Op("=").Id(childName+"v"))
	}

	blockCodes = append(blockCodes, For(Id(childName+"i").Op(":=").Range().Id(childName).Index(Id(":"+childName+"l"))).Block(
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
		codes = []Code{createDecodeNilCheckedCode(blockCodes)}
	}
	return codes
}
