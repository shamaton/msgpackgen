package structure

import (
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/shamaton/msgpackgen/internal/generator/ptn"
)

type identCodeGen struct {
}

func (st *Structure) createIdentCode(node *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	funcSuffix := strings.Title(node.IdenticalName)

	g := identCodeGen{}

	cArray = g.createCalcCode("Calc"+funcSuffix, Id(encodeFieldName))
	cMap = g.createCalcCode("Calc"+funcSuffix, Id(encodeFieldName))

	eArray = g.createEncCode("Write"+funcSuffix, Id(encodeFieldName), Id("offset"))
	eMap = g.createEncCode("Write"+funcSuffix, Id(encodeFieldName), Id("offset"))

	dArray = g.createDecCode(node, st.Others, decodeFieldName, "As"+funcSuffix)
	dMap = g.createDecCode(node, st.Others, decodeFieldName, "As"+funcSuffix)

	return
}

func (g identCodeGen) createCalcCode(funcName string, params ...Code) []Code {
	return []Code{
		createAddSizeCode(funcName, params...),
	}
}

func (g identCodeGen) createEncCode(funcName string, params ...Code) []Code {
	return []Code{
		Id("offset").Op("=").Id(ptn.IdEncoder).Dot(funcName).Call(params...),
	}
}

func (g identCodeGen) createDecCode(node *Node, structures []*Structure, fieldName, funcName string) []Code {

	varName := fieldName + "v"
	if isRootField(fieldName) {
		varName = "vv"
	}

	_, isParentTypeArrayOrMap := node.GetPointerInfo()

	codes, receiverName := createDecodeDefineVarCode(node, structures, varName)

	codes = append(codes,
		List(Id(receiverName), Id("offset"), Err()).Op("=").Id(ptn.IdDecoder).Dot(funcName).Call(Id("offset")),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
	)

	codes = append(codes, createDecodeSetValueCode(node, varName, fieldName)...)

	// array or map
	if isParentTypeArrayOrMap {
		return codes
	}

	return []Code{Block(codes...)}
}
