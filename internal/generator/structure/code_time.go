package structure

import (
	. "github.com/dave/jennifer/jen"
	"github.com/shamaton/msgpackgen/internal/generator/ptn"
)

type timeCodeGen struct {
}

func (st *Structure) createTimeCode(encodeFieldName, decodeFieldName string, node *Node) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {
	g := timeCodeGen{}
	cArray = g.createCalcCode("CalcTime", Id(encodeFieldName))
	cMap = g.createCalcCode("CalcTime", Id(encodeFieldName))

	eArray = g.createEncCode("WriteTime", Id(encodeFieldName), Id("offset"))
	eMap = g.createEncCode("WriteTime", Id(encodeFieldName), Id("offset"))

	dArray = g.createDecCode(node, st.Others, decodeFieldName, "AsDateTime")
	dMap = g.createDecCode(node, st.Others, decodeFieldName, "AsDateTime")
	return
}

func (g timeCodeGen) createCalcCode(funcName string, params ...Code) []Code {
	return []Code{
		CreateAddSizeCode(funcName, params...),
	}
}

func (g timeCodeGen) createEncCode(funcName string, params ...Code) []Code {

	return []Code{
		Id("offset").Op("=").Id(ptn.IdEncoder).Dot(funcName).Call(params...),
	}
}

func (g timeCodeGen) createDecCode(node *Node, structures []*Structure, fieldName, funcName string) []Code {
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
