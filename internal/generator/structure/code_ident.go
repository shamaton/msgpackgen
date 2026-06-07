package structure

import (
	"strings"
	"unicode"

	. "github.com/dave/jennifer/jen"
	"github.com/shamaton/msgpackgen/internal/generator/ptn"
)

type identCodeGen struct {
}

func (st *Structure) createIdentCode(node *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	g := identCodeGen{}

	funcSuffix := g.toPascalCase(node.IdenticalName)

	encodeValue := createIdentEncodeValue(node, encodeFieldName)
	cArray = g.createCalcCode("Calc"+funcSuffix, encodeValue)
	cMap = g.createCalcCode("Calc"+funcSuffix, encodeValue)

	eArray = g.createEncCode("Write"+funcSuffix, encodeValue, Id("offset"))
	eMap = g.createEncCode("Write"+funcSuffix, encodeValue, Id("offset"))

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
		createWriteCode(funcName, params[:len(params)-1]...),
	}
}

func (g identCodeGen) createDecCode(node *Node, structures []*Structure, fieldName, funcName string) []Code {

	varName := fieldName + "v"
	if isRootField(fieldName) {
		varName = "vv"
	}

	_, isParentTypeArrayOrMap := node.GetPointerInfo()

	codes, receiverName := createDecodeDefineVarCode(node, structures, varName)

	if node.AliasName != "" {
		decodedName := receiverName + "d"
		codes = append(codes,
			Var().Id(decodedName).Id(node.IdenticalName),
			List(Id(decodedName), Id("offset"), Err()).Op("=").Id(ptn.IdDecoder).Dot(funcName).Call(Id("offset")),
			If(Err().Op("!=").Nil()).Block(
				Return(Lit(0), Err()),
			),
			Id(receiverName).Op("=").Add(node.AliasTypeJenChain().Call(Id(decodedName))),
		)
	} else {
		codes = append(codes,
			List(Id(receiverName), Id("offset"), Err()).Op("=").Id(ptn.IdDecoder).Dot(funcName).Call(Id("offset")),
			If(Err().Op("!=").Nil()).Block(
				Return(Lit(0), Err()),
			),
		)
	}

	codes = append(codes, createDecodeSetValueCode(node, varName, fieldName)...)

	// array or map
	if isParentTypeArrayOrMap {
		return codes
	}

	return []Code{Block(codes...)}
}

func (g identCodeGen) toPascalCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		// character that is not a letter or digit
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	if len(parts) == 0 {
		return ""
	}

	result := ""
	for _, p := range parts {
		if p == "" {
			continue
		}
		// Convert first character to uppercase and the rest to lowercase
		result += strings.ToUpper(string(p[0])) + strings.ToLower(p[1:])
	}

	return result
}
