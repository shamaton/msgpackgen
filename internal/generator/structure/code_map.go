package structure

import (
	. "github.com/dave/jennifer/jen"
	"github.com/shamaton/msgpackgen/internal/generator/ptn"
)

type mapCodeGen struct {
}

func (st *Structure) createMapCode(node *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	key, value := node.KeyValue()

	encodeChildKey, encodeChildValue := encodeFieldName+"k", encodeFieldName+"v"
	if isRootField(encodeFieldName) {
		encodeChildKey = "kk"
		encodeChildValue = "vv"
	}

	decodeChildKey, decodeChildValue := decodeFieldName+"k", decodeFieldName+"v"
	if isRootField(decodeFieldName) {
		decodeChildKey = "kk"
		decodeChildValue = "vv"
	}

	caKey, cmKey, eaKey, emKey, daKey, dmKey := st.createFieldCode(key, encodeChildKey, decodeChildKey)
	caValue, cmValue, eaValue, emValue, daValue, dmValue := st.createFieldCode(value, encodeChildValue, decodeChildValue)

	g := mapCodeGen{}
	cArray = g.createCalcCode(encodeFieldName, encodeChildKey, encodeChildValue, caKey, caValue)
	cMap = g.createCalcCode(encodeFieldName, encodeChildKey, encodeChildValue, cmKey, cmValue)

	eArray = g.createEncCode(encodeFieldName, encodeChildKey, encodeChildValue, eaKey, eaValue)
	eMap = g.createEncCode(encodeFieldName, encodeChildKey, encodeChildValue, emKey, emValue)

	dArray = g.createDecCode(node, st.Others, decodeFieldName, decodeChildKey, decodeChildValue, daKey, daValue)
	dMap = g.createDecCode(node, st.Others, decodeFieldName, decodeChildKey, decodeChildValue, dmKey, dmValue)

	return
}

func (g mapCodeGen) createCalcCode(
	fieldName, childKeyName, childValueName string,
	elmKeyCodes, elmValueCodes []Code) []Code {
	calcCodes := CreateAddSizeErrCheckCode("CalcMapLength", Len(Id(fieldName)))
	calcCodes = append(calcCodes, For(List(Id(childKeyName), Id(childValueName)).Op(":=").Range().Id(fieldName)).Block(
		append(elmKeyCodes, elmValueCodes...)...,
	))

	var codes []Code
	codes = append(codes, If(Id(fieldName).Op("!=").Nil()).Block(
		calcCodes...,
	).Else().Block(
		CreateAddSizeCode("CalcNil"),
	))
	return codes
}

func (g mapCodeGen) createEncCode(
	fieldName, childKeyName, childValueName string,
	elmKeyCodes, elmValueCodes []Code) []Code {

	encCodes := make([]Code, 0)
	encCodes = append(encCodes, Id("offset").Op("=").Id(ptn.IdEncoder).Dot("WriteMapLength").Call(Len(Id(fieldName)), Id("offset")))
	encCodes = append(encCodes, For(List(Id(childKeyName), Id(childValueName)).Op(":=").Range().Id(fieldName)).Block(
		append(elmKeyCodes, elmValueCodes...)...,
	))

	var codes []Code
	codes = append(codes, If(Id(fieldName).Op("!=").Nil()).Block(
		encCodes...,
	).Else().Block(
		Id("offset").Op("=").Id(ptn.IdEncoder).Dot("WriteNil").Call(Id("offset")),
	))
	return codes
}

func (g mapCodeGen) createDecCode(
	ast *Node, structures []*Structure,
	fieldName, childKeyName, childValueName string,
	elmKeyCodes, elmValueCodes []Code) []Code {

	decCodes := make([]Code, 0)
	decCodes = append(decCodes, ast.TypeJenChain(structures, Var().Id(childValueName)))
	decCodes = append(decCodes, Var().Id(childValueName+"l").Int())
	decCodes = append(decCodes, List(Id(childValueName+"l"), Id("offset"), Err()).Op("=").Id(ptn.IdDecoder).Dot("MapLength").Call(Id("offset")))
	decCodes = append(decCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	decCodes = append(decCodes, Id(childValueName).Op("=").Make(ast.TypeJenChain(structures), Id(childValueName+"l")))

	da := []Code{ast.Key.TypeJenChain(structures, Var().Id(childKeyName+"v"))}
	da = append(da, elmKeyCodes...)
	da = append(da, ast.Value.TypeJenChain(structures, Var().Id(childValueName+"v")))
	da = append(da, elmValueCodes...)
	da = append(da, Id(childValueName).Index(Id(childKeyName+"v")).Op("=").Id(childValueName+"v"))

	decCodes = append(decCodes, For(Id(childValueName+"i").Op(":=").Lit(0).Op(";").Id(childValueName+"i").Op("<").Id(childValueName+"l").Op(";").Id(childValueName+"i").Op("++")).Block(
		da...,
	))

	name := childValueName
	andOp := ""
	prtCount, _ := ast.GetPointerInfo()

	if prtCount > 0 {
		andOp = "&"
	}

	for i := 0; i < prtCount-1; i++ {
		n := "_" + name
		decCodes = append(decCodes, Id(n).Op(":=").Op("&").Id(name))
		name = n
	}

	decCodes = append(decCodes, Id(fieldName).Op("=").Op(andOp).Id(name))

	var codes []Code
	if ast.HasParent() && ast.Parent.IsPointer() {
		codes = decCodes
	} else {

		codes = append(codes, If(Op("!").Id(ptn.IdDecoder).Dot("IsCodeNil").Call(Id("offset"))).Block(
			decCodes...,
		).Else().Block(
			Id("offset").Op("++"),
		))
	}
	return codes
}
