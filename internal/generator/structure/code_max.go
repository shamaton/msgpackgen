package structure

import (
	"strings"

	. "github.com/dave/jennifer/jen"
)

func (st *Structure) createFieldMaxCode(node *Node, encodeFieldName string) (cArray []Code, cMap []Code) {
	switch {
	case node.IsIdentical():
		funcSuffix := identCodeGen{}.toPascalCase(node.IdenticalName)
		codes := []Code{createAddSizeMaxCode("Calc"+funcSuffix, Id(encodeFieldName))}
		return codes, codes

	case node.IsSlice():
		return st.createSliceMaxCode(node, encodeFieldName)

	case node.IsArray():
		return st.createArrayMaxCode(node, encodeFieldName)

	case node.IsMap():
		return st.createMapMaxCode(node, encodeFieldName)

	case node.IsPointer():
		return st.createPointerMaxCode(node, encodeFieldName)

	case node.IsStruct():
		if node.ImportPath == "time" {
			codes := []Code{createAddSizeMaxCode("CalcTime", Id(encodeFieldName))}
			return codes, codes
		}
		return st.createNamedMaxCode(node, encodeFieldName)
	}

	return nil, nil
}

func (st *Structure) createSliceMaxCode(node *Node, fieldName string) (cArray []Code, cMap []Code) {
	childName := fieldName + "v"
	if isRootField(fieldName) {
		childName = "vv"
	}

	childArray, childMap := st.createFieldMaxCode(node.Elm(), childName)
	isChildByte := node.Elm().IsIdentical() && node.Elm().IdenticalName == "byte"

	cArray = createNullableSequenceMaxCode(fieldName, childName, "CalcSliceLength", isChildByte, childArray)
	cMap = createNullableSequenceMaxCode(fieldName, childName, "CalcSliceLength", isChildByte, childMap)
	return
}

func (st *Structure) createArrayMaxCode(node *Node, fieldName string) (cArray []Code, cMap []Code) {
	childName := fieldName + "v"
	if isRootField(fieldName) {
		childName = "vv"
	}

	childArray, childMap := st.createFieldMaxCode(node.Elm(), childName)
	isChildByte := node.Elm().IsIdentical() && node.Elm().IdenticalName == "byte"

	cArray = createSequenceMaxCode(fieldName, childName, "CalcSliceLength", isChildByte, childArray)
	cMap = createSequenceMaxCode(fieldName, childName, "CalcSliceLength", isChildByte, childMap)
	return
}

func (st *Structure) createMapMaxCode(node *Node, fieldName string) (cArray []Code, cMap []Code) {
	key, value := node.KeyValue()

	childKeyName, childValueName := fieldName+"k", fieldName+"v"
	if isRootField(fieldName) {
		childKeyName = "kk"
		childValueName = "vv"
	}

	keyArray, keyMap := st.createFieldMaxCode(key, childKeyName)
	valueArray, valueMap := st.createFieldMaxCode(value, childValueName)

	cArray = createNullableMapMaxCode(fieldName, childKeyName, childValueName, keyArray, valueArray)
	cMap = createNullableMapMaxCode(fieldName, childKeyName, childValueName, keyMap, valueMap)
	return
}

func (st *Structure) createPointerMaxCode(node *Node, fieldName string) (cArray []Code, cMap []Code) {
	childName := fieldName + "p"
	if isRootField(fieldName) {
		childName = "vp"
	}

	childArray, childMap := st.createFieldMaxCode(node.Elm(), childName)

	cArray = createNullablePointerMaxCode(fieldName, childName, childArray)
	cMap = createNullablePointerMaxCode(fieldName, childName, childMap)
	return
}

func (st *Structure) createNamedMaxCode(node *Node, fieldName string) (cArray []Code, cMap []Code) {
	sizeName := "size_" + fieldName
	if isRootField(fieldName) {
		sizeName = strings.ReplaceAll(sizeName, ".", "_")
	}

	cArray = createNamedMaxSizeCode(node, fieldName, sizeName, "calcArraySizeMax")
	cMap = createNamedMaxSizeCode(node, fieldName, sizeName, "calcMapSizeMax")
	return
}

func createNullableSequenceMaxCode(fieldName, childName, funcName string, isChildByte bool, childCodes []Code) []Code {
	blockCodes := createSequenceMaxCode(fieldName, childName, funcName, isChildByte, childCodes)
	return []Code{
		If(Id(fieldName).Op("!=").Nil()).Block(
			blockCodes...,
		).Else().Block(
			createAddSizeCode("CalcNil"),
		),
	}
}

func createSequenceMaxCode(fieldName, childName, funcName string, isChildByte bool, childCodes []Code) []Code {
	blockCodes := createAddSizeMaxErrCheckCode(funcName, Len(Id(fieldName)), Lit(isChildByte))
	blockCodes = append(blockCodes, For(List(Id("_"), Id(childName)).Op(":=").Range().Id(fieldName)).Block(
		childCodes...,
	))
	return []Code{Block(blockCodes...)}
}

func createNullableMapMaxCode(fieldName, childKeyName, childValueName string, keyCodes, valueCodes []Code) []Code {
	calcCodes := createAddSizeMaxErrCheckCode("CalcMapLength", Len(Id(fieldName)))
	calcCodes = append(calcCodes, For(List(Id(childKeyName), Id(childValueName)).Op(":=").Range().Id(fieldName)).Block(
		append(keyCodes, valueCodes...)...,
	))

	return []Code{
		If(Id(fieldName).Op("!=").Nil()).Block(
			calcCodes...,
		).Else().Block(
			createAddSizeCode("CalcNil"),
		),
	}
}

func createNullablePointerMaxCode(fieldName, childName string, childCodes []Code) []Code {
	return []Code{
		If(Id(fieldName).Op("!=").Nil()).Block(
			append([]Code{
				Id(childName).Op(":=").Op("*").Id(fieldName),
			}, childCodes...)...,
		).Else().Block(
			createAddSizeCode("CalcNil"),
		),
	}
}

func createNamedMaxSizeCode(node *Node, fieldName, sizeName, funcName string) []Code {
	return []Code{
		List(Id(sizeName), Err()).
			Op(":=").
			Id(createFuncName(funcName, node.StructName, node.ImportPath)).Call(Id(fieldName)),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id("size").Op("+=").Id(sizeName),
	}
}
