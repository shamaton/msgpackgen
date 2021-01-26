package structure

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

func (n Node) Elm() *Node               { return n.Key }
func (n Node) KeyValue() (*Node, *Node) { return n.Key, n.Value }

func (n Node) IsIdentical() bool { return n.fieldType == fieldTypeIdent }
func (n Node) IsSlice() bool     { return n.fieldType == fieldTypeSlice }
func (n Node) IsArray() bool     { return n.fieldType == fieldTypeArray }
func (n Node) IsStruct() bool    { return n.fieldType == fieldTypeStruct }
func (n Node) IsMap() bool       { return n.fieldType == fieldTypeMap }
func (n Node) IsPointer() bool   { return n.fieldType == fieldTypePointer }

func (n Node) HasParent() bool       { return n.Parent != nil }
func (n Node) IsParentPointer() bool { return n.HasParent() && n.Parent.IsPointer() }

func (n *Node) SetKeyNode(key *Node)     { n.Key = key }
func (n *Node) SetValueNode(value *Node) { n.Value = value }

func (n Node) CanGenerate(sts []*Structure) (bool, []string) {
	messages := make([]string, 0)
	switch {
	case n.IsIdentical():
		return true, messages

	case n.IsStruct():
		if n.ImportPath == "time" && n.StructName == "Time" {
			return true, messages
		}
		for _, v := range sts {
			if v.ImportPath == n.ImportPath && v.Name == n.StructName {
				return true, messages
			}
		}
		return false, append(messages, fmt.Sprintf("struct %s.%s is not generated.", n.ImportPath, n.StructName))

	case n.IsSlice():
		return n.Elm().CanGenerate(sts)

	case n.IsArray():
		return n.Elm().CanGenerate(sts)

	case n.IsMap():
		k, v := n.KeyValue()
		kb, kMessages := k.CanGenerate(sts)
		vb, vMessages := v.CanGenerate(sts)
		messages = append(messages, kMessages...)
		messages = append(messages, vMessages...)
		return kb && vb, messages

	case n.IsPointer():
		return n.Elm().CanGenerate(sts)
	}
	return false, append(messages, "unreachable code")
}

func (n Node) TypeJenChain(sts []*Structure, s ...*Statement) *Statement {
	var str *Statement
	if len(s) > 0 {
		str = s[0]
	} else {
		str = Id("")
	}

	switch {
	case n.IsIdentical():
		str = str.Id(n.IdenticalName)

	case n.IsStruct():
		if n.ImportPath == "time" && n.StructName == "Time" {
			str = str.Qual(n.ImportPath, n.StructName)
		} else {
			// todo : performance
			var asRef *Structure
			for _, v := range sts {
				if v.ImportPath == n.ImportPath && v.Name == n.StructName {
					asRef = v
					break
				}
			}
			if asRef == nil {
				// unreachable
				panic(fmt.Sprintf("not found struct %s.%s", n.ImportPath, n.StructName))
			}

			if asRef.NoUseQual {
				str = str.Id(n.StructName)
			} else {
				str = str.Qual(n.ImportPath, n.StructName)
			}
		}

	case n.IsSlice():
		str = str.Id("[]")
		str = n.Elm().TypeJenChain(sts, str)

	case n.IsArray():
		str = str.Id(fmt.Sprintf("[%d]", n.ArrayLen))
		str = n.Elm().TypeJenChain(sts, str)

	case n.IsMap():
		str = str.Id("map[")
		k, v := n.KeyValue()
		str = k.TypeJenChain(sts, str)
		str = str.Id("]")
		str = v.TypeJenChain(sts, str)

	case n.IsPointer():
		str = str.Id("*")
		str = n.Elm().TypeJenChain(sts, str)
	}
	return str
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

func IsPrimitive(name string) bool {
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
