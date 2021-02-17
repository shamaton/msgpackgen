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

// Node is a collection of field information.
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

// Elm gets the child node, if field-type is slice, array or pointer.
func (n Node) Elm() *Node { return n.Key }

// KeyValue gets child nodes, if field-type is map.
func (n Node) KeyValue() (*Node, *Node) { return n.Key, n.Value }

// IsIdentical returns true, if field-type is ident.
func (n Node) IsIdentical() bool { return n.fieldType == fieldTypeIdent }

// IsSlice returns true, if field-type is slice.
func (n Node) IsSlice() bool { return n.fieldType == fieldTypeSlice }

// IsArray returns true, if field-type is array.
func (n Node) IsArray() bool { return n.fieldType == fieldTypeArray }

// IsStruct returns true, if field-type is struct.
func (n Node) IsStruct() bool { return n.fieldType == fieldTypeStruct }

// IsMap returns true, if field-type is map.
func (n Node) IsMap() bool { return n.fieldType == fieldTypeMap }

// IsPointer returns true, if field-type is pointer.
func (n Node) IsPointer() bool { return n.fieldType == fieldTypePointer }

// HasParent returns true, if parent node exists.
func (n Node) HasParent() bool { return n.Parent != nil }

// IsParentPointer returns true, if HasParent is true and parent node is pointer.
func (n Node) IsParentPointer() bool { return n.HasParent() && n.Parent.IsPointer() }

// SetKeyNode sets node to key field.
func (n *Node) SetKeyNode(key *Node) { n.Key = key }

// SetValueNode sets node to value field.
func (n *Node) SetValueNode(value *Node) { n.Value = value }

// GetPointerInfo gets some pointer information to create codes.
func (n *Node) GetPointerInfo() (ptrCount int, isParentTypeArrayOrMap bool) {
	node := n
	for node.HasParent() {
		node = node.Parent
		if node.IsPointer() {
			// pointer
			ptrCount++
		} else {
			// slice / array / map
			isParentTypeArrayOrMap = true
			break
		}
	}
	return
}

// CanGenerate return true, if it satisfied conditions by this node.
func (n Node) CanGenerate(structures []*Structure) (bool, []string) {
	messages := make([]string, 0)
	switch {
	case n.IsIdentical():
		return true, messages

	case n.IsStruct():
		if n.ImportPath == "time" && n.StructName == "Time" {
			return true, messages
		}
		for _, v := range structures {
			if v.ImportPath == n.ImportPath && v.Name == n.StructName {
				return true, messages
			}
		}
		return false, append(messages, fmt.Sprintf("struct %s.%s is not generated.", n.ImportPath, n.StructName))

	case n.IsSlice():
		return n.Elm().CanGenerate(structures)

	case n.IsArray():
		return n.Elm().CanGenerate(structures)

	case n.IsMap():
		k, v := n.KeyValue()
		kb, kMessages := k.CanGenerate(structures)
		vb, vMessages := v.CanGenerate(structures)
		messages = append(messages, kMessages...)
		messages = append(messages, vMessages...)
		return kb && vb, messages

	case n.IsPointer():
		return n.Elm().CanGenerate(structures)
	}
	return false, append(messages, "unreachable code")
}

// TypeJenChain is a helper method to create code.
func (n Node) TypeJenChain(structures []*Structure, statements ...*Statement) *Statement {
	var str *Statement
	if len(statements) > 0 {
		str = statements[0]
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
			var asRef *Structure
			for _, v := range structures {
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
		str = n.Elm().TypeJenChain(structures, str)

	case n.IsArray():
		str = str.Id(fmt.Sprintf("[%d]", n.ArrayLen))
		str = n.Elm().TypeJenChain(structures, str)

	case n.IsMap():
		str = str.Id("map[")
		k, v := n.KeyValue()
		str = k.TypeJenChain(structures, str)
		str = str.Id("]")
		str = v.TypeJenChain(structures, str)

	case n.IsPointer():
		str = str.Id("*")
		str = n.Elm().TypeJenChain(structures, str)
	}
	return str
}

// CreateIdentNode creates a node of ident type.
func CreateIdentNode(ident *ast.Ident, parent *Node) *Node {
	return &Node{
		fieldType:     fieldTypeIdent,
		IdenticalName: ident.Name,
		Parent:        parent,
	}
}

// CreateStructNode creates a node of struct type.
func CreateStructNode(importPath, packageName, structName string, parent *Node) *Node {
	return &Node{
		fieldType:   fieldTypeStruct,
		ImportPath:  importPath,
		PackageName: packageName,
		StructName:  structName,
		Parent:      parent,
	}
}

// CreateSliceNode creates a node of slice type.
func CreateSliceNode(parent *Node) *Node {
	return &Node{
		fieldType: fieldTypeSlice,
		Parent:    parent,
	}
}

// CreateArrayNode creates a node of array type.
func CreateArrayNode(len uint64, parent *Node) *Node {
	return &Node{
		fieldType: fieldTypeArray,
		ArrayLen:  len,
		Parent:    parent,
	}
}

// CreateMapNode creates a node of map type.
func CreateMapNode(parent *Node) *Node {
	return &Node{
		fieldType: fieldTypeMap,
		Parent:    parent,
	}
}

// CreatePointerNode creates a node of pointer type.
func CreatePointerNode(parent *Node) *Node {
	return &Node{
		fieldType: fieldTypePointer,
		Parent:    parent,
	}
}

// IsPrimitive returns true, if name matches any case.
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
