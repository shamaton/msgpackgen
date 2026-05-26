package structure

import (
	. "github.com/dave/jennifer/jen"
	"github.com/shamaton/msgpackgen/internal/generator/ptn"
)

type pointerCodeGen struct {
}

func (st *Structure) createPointerCode(node *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	encodeChildName := encodeFieldName + "p"
	if isRootField(encodeFieldName) {
		encodeChildName = "vp"
	}

	var ca, cm, ea, em []Code
	if canPassPointerDirect(node.Elm()) {
		ca, cm, ea, em, _, _ = st.createFieldCode(node.Elm(), encodeFieldName, decodeFieldName)
	} else {
		ca, cm, ea, em, _, _ = st.createFieldCode(node.Elm(), encodeChildName, decodeFieldName)
	}
	_, _, _, _, da, dm := st.createFieldCode(node.Elm(), encodeChildName, decodeFieldName)

	g := pointerCodeGen{}
	cArray = g.createPointerCalcCode(encodeFieldName, encodeChildName, canPassPointerDirect(node.Elm()), ca)
	cMap = g.createPointerCalcCode(encodeFieldName, encodeChildName, canPassPointerDirect(node.Elm()), cm)

	eArray = g.createPointerEncCode(encodeFieldName, encodeChildName, canPassPointerDirect(node.Elm()), ea)
	eMap = g.createPointerEncCode(encodeFieldName, encodeChildName, canPassPointerDirect(node.Elm()), em)

	dArray = g.createPointerDecCode(node, da)
	dMap = g.createPointerDecCode(node, dm)

	return
}

func canPassPointerDirect(node *Node) bool {
	return node.IsStruct() && node.ImportPath != "time"
}

func (g pointerCodeGen) createPointerCalcCode(encodeFieldName, encodeChildName string, passPointerDirect bool, elmCodes []Code) []Code {
	codes := make([]Code, 0)
	nonNilCodes := elmCodes
	if !passPointerDirect {
		nonNilCodes = append([]Code{
			Id(encodeChildName).Op(":=").Op("*").Id(encodeFieldName),
		}, elmCodes...)
	}
	codes = append(codes, If(Id(encodeFieldName).Op("!=").Nil()).Block(
		nonNilCodes...,
	).Else().Block(
		createAddSizeCode("CalcNil"),
	))
	return codes
}

func (g pointerCodeGen) createPointerEncCode(encodeFieldName, encodeChildName string, passPointerDirect bool, elmCodes []Code) []Code {
	codes := make([]Code, 0)
	nonNilCodes := elmCodes
	if !passPointerDirect {
		nonNilCodes = append([]Code{
			Id(encodeChildName).Op(":=").Op("*").Id(encodeFieldName),
		}, elmCodes...)
	}
	codes = append(codes, If(Id(encodeFieldName).Op("!=").Nil()).Block(
		nonNilCodes...,
	).Else().Block(
		Id("offset").Op("=").Qual(ptn.PkEnc, "WriteNil").Call(Id("buf"), Id("offset")),
	))
	return codes
}

func (g pointerCodeGen) createPointerDecCode(node *Node, elmCodes []Code) []Code {
	var codes []Code
	if node.IsParentPointer() {
		codes = elmCodes
	} else {
		codes = []Code{createDecodeNilCheckedCode(elmCodes)}
	}
	return codes
}
