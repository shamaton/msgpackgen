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

	ca, cm, ea, em, da, dm := st.createFieldCode(node.Elm(), encodeChildName, decodeFieldName)

	g := pointerCodeGen{}
	cArray = g.createPointerCalcCode(encodeFieldName, encodeChildName, ca)
	cMap = g.createPointerCalcCode(encodeFieldName, encodeChildName, cm)

	eArray = g.createPointerEncCode(encodeFieldName, encodeChildName, ea)
	eMap = g.createPointerEncCode(encodeFieldName, encodeChildName, em)

	dArray = g.createPointerDecCode(node, da)
	dMap = g.createPointerDecCode(node, dm)

	return
}

func (g pointerCodeGen) createPointerCalcCode(encodeFieldName, encodeChildName string, elmCodes []Code) []Code {
	codes := make([]Code, 0)
	codes = append(codes, If(Id(encodeFieldName).Op("!=").Nil()).Block(
		append([]Code{
			Id(encodeChildName).Op(":=").Op("*").Id(encodeFieldName),
		}, elmCodes...)...,
	).Else().Block(
		Id("size").Op("+=").Id(ptn.IdEncoder).Dot("CalcNil").Call(),
	))
	return codes
}

func (g pointerCodeGen) createPointerEncCode(encodeFieldName, encodeChildName string, elmCodes []Code) []Code {
	codes := make([]Code, 0)
	codes = append(codes, If(Id(encodeFieldName).Op("!=").Nil()).Block(
		append([]Code{
			Id(encodeChildName).Op(":=").Op("*").Id(encodeFieldName),
		}, elmCodes...)...,
	).Else().Block(
		Id("offset").Op("=").Id(ptn.IdEncoder).Dot("WriteNil").Call(Id("offset")),
	))
	return codes
}

func (g pointerCodeGen) createPointerDecCode(node *Node, elmCodes []Code) []Code {
	var codes []Code
	if node.IsParentPointer() {
		codes = elmCodes
	} else {
		codes = make([]Code, 0)
		codes = append(codes, If(Op("!").Id(ptn.IdDecoder).Dot("IsCodeNil").Call(Id("offset"))).Block(
			elmCodes...,
		).Else().Block(
			Id("offset").Op("++"),
		))
	}
	return codes
}
