package generator

import (
	"fmt"
	"go/ast"

	. "github.com/dave/jennifer/jen"
)

const (
	fieldTypeIdent = iota + 1
	fieldTypeSlice
	fieldTypeArray
	fieldTypeStruct
	fieldTypeMap
	fieldTypePointer
)

type Node struct {
	fieldType int

	// for identical
	IdenticalName string

	// for array
	ArrayLen uint64

	// for struct
	ImportPath  string
	PackageName string
	StructName  string

	// for array / map / pointer
	Key   *Node
	Value *Node

	Parent *Node
}

func (a Node) HasParent() bool { return a.Parent != nil }

func (a Node) IsIdentical() bool { return a.fieldType == fieldTypeIdent }
func (a Node) IsSlice() bool     { return a.fieldType == fieldTypeSlice }
func (a Node) IsArray() bool     { return a.fieldType == fieldTypeArray }
func (a Node) IsStruct() bool    { return a.fieldType == fieldTypeStruct }
func (a Node) IsMap() bool       { return a.fieldType == fieldTypeMap }

func (a Node) IsPointer() bool { return a.fieldType == fieldTypePointer }

func (a Node) Elm() *Node { return a.Key }
func (a Node) KeyValue() (*Node, *Node) {
	return a.Key, a.Value
}

func (a Node) CanGenerate(sts []*analyzedStruct) (bool, []string) {
	messages := make([]string, 0)
	switch {
	case a.IsIdentical():
		return true, messages

	case a.IsStruct():
		if a.ImportPath == "time" && a.StructName == "Time" {
			return true, messages
		}
		for _, v := range sts {
			if v.ImportPath == a.ImportPath && v.Name == a.StructName {
				return true, messages
			}
		}
		return false, append(messages, fmt.Sprintf("struct %s.%s is not generated.", a.ImportPath, a.StructName))

	case a.IsSlice():
		return a.Elm().CanGenerate(sts)

	case a.IsArray():
		return a.Elm().CanGenerate(sts)

	case a.IsMap():
		k, v := a.KeyValue()
		kb, kMessages := k.CanGenerate(sts)
		vb, vMessages := v.CanGenerate(sts)
		messages = append(messages, kMessages...)
		messages = append(messages, vMessages...)
		return kb && vb, messages

	case a.IsPointer():
		return a.Elm().CanGenerate(sts)
	}
	return false, append(messages, "unreachable code")
}

func (a Node) TypeJenChain(sts []*analyzedStruct, s ...*Statement) *Statement {
	var str *Statement
	if len(s) > 0 {
		str = s[0]
	} else {
		str = Id("")
	}

	switch {
	case a.IsIdentical():
		str = str.Id(a.IdenticalName)

	case a.IsStruct():
		if a.ImportPath == "time" && a.StructName == "Time" {
			str = str.Qual(a.ImportPath, a.StructName)
		} else {
			// todo : performance
			var asRef *analyzedStruct
			for _, v := range sts {
				if v.ImportPath == a.ImportPath && v.Name == a.StructName {
					asRef = v
					break
				}
			}
			if asRef == nil {
				// unreachable
				panic(fmt.Sprintf("not found struct %s.%s", a.ImportPath, a.StructName))
			}

			if asRef.NoUseQual {
				str = str.Id(a.StructName)
			} else {
				str = str.Qual(a.ImportPath, a.StructName)
			}
		}

	case a.IsSlice():
		str = str.Id("[]")
		str = a.Elm().TypeJenChain(sts, str)

	case a.IsArray():
		str = str.Id(fmt.Sprintf("[%d]", a.ArrayLen))
		str = a.Elm().TypeJenChain(sts, str)

	case a.IsMap():
		str = str.Id("map[")
		k, v := a.KeyValue()
		str = k.TypeJenChain(sts, str)
		str = str.Id("]")
		str = v.TypeJenChain(sts, str)

	case a.IsPointer():
		str = str.Id("*")
		str = a.Elm().TypeJenChain(sts, str)
	}
	return str
}

func (n *Node) SetKeyNode(key *Node) {
	n.Key = key
}

func (n *Node) SetValueNode(value *Node) {
	n.Value = value
}

func CreateIdentNode(ident *ast.Ident, parent *Node) *Node {
	return &Node{
		fieldType:     fieldTypeIdent,
		IdenticalName: ident.Name,
		Parent:        parent,
	}
}

func CreateStructNode(importPath, packageName, structName string, parent *Node) *Node {
	return &Node{
		fieldType:   fieldTypeStruct,
		ImportPath:  importPath,
		PackageName: packageName,
		StructName:  structName,
		Parent:      parent,
	}
}

func CreateSliceNode(parent *Node) *Node {
	return &Node{
		fieldType: fieldTypeSlice,
		Parent:    parent,
	}
}

func CreateArrayNode(len uint64, parent *Node) *Node {
	return &Node{
		fieldType: fieldTypeArray,
		ArrayLen:  len,
		Parent:    parent,
	}
}

func CreateMapNode(parent *Node) *Node {
	return &Node{
		fieldType: fieldTypeMap,
		Parent:    parent,
	}
}

func CreatePointerNode(parent *Node) *Node {
	return &Node{
		fieldType: fieldTypePointer,
		Parent:    parent,
	}
}

func isPrimitive(name string) bool {
	switch name {
	case "int", "int8", "int16", "int32", "int64":
		return true

	case "uint", "uint8", "uint16", "uint32", "uint64":
		return true

	case "float32", "float64":
		return true

	case "string", "rune":
		return true

	case "bool", "byte", "complex64", "complex128":
		return true
	}
	return false
}
