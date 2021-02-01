package structure

import (
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/shamaton/msgpackgen/internal/generator/ptn"
)

type namedCodeGen struct {
}

func (st *Structure) createNamedCode(encodeFieldName, decodeFieldName string, ast *Node) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	sizeName := "size_" + encodeFieldName
	if isRootField(encodeFieldName) {
		sizeName = strings.ReplaceAll(sizeName, ".", "_")
	}

	g := namedCodeGen{}
	cArray = g.createCalcCode(ast, encodeFieldName, sizeName, "calcArraySize")
	cMap = g.createCalcCode(ast, encodeFieldName, sizeName, "calcMapSize")

	eArray = g.createEncCode(ast, encodeFieldName, "encodeArray")
	eMap = g.createEncCode(ast, encodeFieldName, "encodeMap")

	dArray = g.createDecCode(ast, st.Others, decodeFieldName, "decodeArray")
	dMap = g.createDecCode(ast, st.Others, decodeFieldName, "decodeMap")

	return
}

func (g namedCodeGen) createCalcCode(node *Node, fieldName, sizeName, funcName string) []Code {

	return []Code{
		List(Id(sizeName), Err()).
			Op(":=").
			Id(createFuncName(funcName, node.StructName, node.ImportPath)).Call(Id(fieldName), Id(ptn.IdEncoder)),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id("size").Op("+=").Id(sizeName),
	}
}

func (g namedCodeGen) createEncCode(node *Node, fieldName, funcName string) []Code {

	return []Code{
		List(Id("_"), Id("offset"), Err()).
			Op("=").
			Id(createFuncName(funcName, node.StructName, node.ImportPath)).Call(Id(fieldName), Id(ptn.IdEncoder), Id("offset")),
		If(Err().Op("!=").Nil()).Block(
			Return(Nil(), Lit(0), Err()),
		),
	}
}

func (g namedCodeGen) createDecCode(node *Node, structures []*Structure, fieldName, funcName string) []Code {

	varName := fieldName + "v"
	if isRootField(fieldName) {
		varName = "vv"
	}

	_, isParentTypeArrayOrMap := node.GetPointerInfo()

	codes, receiverName := createDecodeDefineVarCode(node, structures, varName)

	codes = append(codes,
		List(Id("offset"), Err()).Op("=").Id(
			createFuncName(funcName, node.StructName, node.ImportPath)).Call(Op("&").Id(receiverName),
			Id(ptn.IdDecoder), Id("offset")),
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
