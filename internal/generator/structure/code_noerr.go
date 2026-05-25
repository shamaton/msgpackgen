package structure

import (
	"strings"

	. "github.com/dave/jennifer/jen"
)

// CanCalcSizeNoErr reports whether a structure's size calculation cannot return
// a MessagePack length error. Slices, maps, and arrays stay on the error-return
// path because their container headers can exceed MessagePack's uint32 length.
func (st *Structure) CanCalcSizeNoErr() bool {
	return st.canCalcSizeNoErr(make(map[*Structure]bool))
}

func (st *Structure) canCalcSizeNoErr(seen map[*Structure]bool) bool {
	if seen[st] {
		return false
	}
	seen[st] = true
	defer delete(seen, st)

	for _, field := range st.Fields {
		if !st.canCalcNodeSizeNoErr(field.Node, seen) {
			return false
		}
	}
	return true
}

func (st *Structure) canCalcNodeSizeNoErr(node *Node, seen map[*Structure]bool) bool {
	switch {
	case node.IsIdentical():
		return true
	case node.IsStruct():
		if node.ImportPath == "time" {
			return true
		}
		ref := st.findStructure(node.ImportPath, node.StructName)
		return ref != nil && ref.canCalcSizeNoErr(seen)
	}
	return false
}

func (st *Structure) findStructure(importPath, name string) *Structure {
	for _, other := range st.Others {
		if other.ImportPath == importPath && other.Name == name {
			return other
		}
	}
	return nil
}

func (st *Structure) createFieldSizeNoErrCode(node *Node, encodeFieldName string, max bool) (cArray []Code, cMap []Code) {
	switch {
	case node.IsIdentical():
		funcSuffix := identCodeGen{}.toPascalCase(node.IdenticalName)
		funcName := "Calc" + funcSuffix
		if max {
			funcName += "Max"
		}
		codes := []Code{createAddSizeCode(funcName, Id(encodeFieldName))}
		return codes, codes

	case node.IsStruct():
		if node.ImportPath == "time" {
			funcName := "CalcTime"
			if max {
				funcName += "Max"
			}
			codes := []Code{createAddSizeCode(funcName, Id(encodeFieldName))}
			return codes, codes
		}
		return st.createNamedSizeNoErrCode(node, encodeFieldName, max)
	}

	return nil, nil
}

func (st *Structure) createNamedSizeNoErrCode(node *Node, fieldName string, max bool) (cArray []Code, cMap []Code) {
	sizeName := "size_" + fieldName
	if isRootField(fieldName) {
		sizeName = strings.ReplaceAll(sizeName, ".", "_")
	}

	arrayFuncName := "calcArraySizeNoErr"
	mapFuncName := "calcMapSizeNoErr"
	if max {
		arrayFuncName = "calcArraySizeMaxNoErr"
		mapFuncName = "calcMapSizeMaxNoErr"
	}

	cArray = createNamedSizeNoErrCode(node, fieldName, sizeName, arrayFuncName)
	cMap = createNamedSizeNoErrCode(node, fieldName, sizeName, mapFuncName)
	return
}

func createNamedSizeNoErrCode(node *Node, fieldName, sizeName, funcName string) []Code {
	return []Code{
		Id(sizeName).
			Op(":=").
			Id(createFuncName(funcName, node.StructName, node.ImportPath)).Call(Id(fieldName)),
		Id("size").Op("+=").Id(sizeName),
	}
}
