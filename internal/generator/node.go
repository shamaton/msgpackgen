package generator

import (
	"fmt"
	"go/ast"
	"math/big"
	"strings"

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

func (g *generator) createNodeRecursive(expr ast.Expr, parent *Node, importMap map[string]string, dotStructs map[string]*analyzedStruct, sameHierarchyStructs map[string]bool) (*Node, bool, []string) {

	reasons := make([]string, 0)
	if ident, ok := expr.(*ast.Ident); ok {
		// dot import
		if dot, found := dotStructs[ident.Name]; found {
			return CreateStructNode(dot.ImportPath, dot.Name, ident.Name, parent), true, reasons
		}
		// time
		if ident.Name == "Time" {
			return CreateStructNode("time", "time", ident.Name, parent), true, reasons
		}
		// same hierarchy struct in same File
		if ident.Obj != nil && ident.Obj.Kind == ast.Typ {
			return CreateStructNode(g.outputPackageFullName(), g.outputPackageName, ident.Name, parent), true, reasons
		}

		// same hierarchy struct in other File
		if _, found := sameHierarchyStructs[ident.Name]; found {
			return CreateStructNode(g.outputPackageFullName(), g.outputPackageName, ident.Name, parent), true, reasons
		}

		if isPrimitive(ident.Name) {
			return CreateIdentNode(ident, parent), true, reasons
		}
		return nil, false, []string{fmt.Sprintf("identifier %s is not suppoted or unknown struct ", ident.Name)}
	}

	if selector, ok := expr.(*ast.SelectorExpr); ok {
		pkgName := fmt.Sprint(selector.X)
		return CreateStructNode(importMap[pkgName], pkgName, selector.Sel.Name, parent), true, reasons
	}

	// slice or array
	if array, ok := expr.(*ast.ArrayType); ok {
		var node *Node
		if array.Len == nil {
			node = CreateSliceNode(parent)
		} else {
			lit := array.Len.(*ast.BasicLit)
			// todo : 処理されなかった場合はエラー
			// todo : box数値以外あればエラーでもいい
			// parse num
			n := new(big.Int)
			if litValue := strings.ToLower(lit.Value); strings.HasPrefix(litValue, "0b") {
				n.SetString(strings.ReplaceAll(litValue, "0b", ""), 2)
			} else if strings.HasPrefix(litValue, "0o") {
				n.SetString(strings.ReplaceAll(litValue, "0o", ""), 8)
			} else if strings.HasPrefix(litValue, "0x") {
				n.SetString(strings.ReplaceAll(litValue, "0x", ""), 16)
			} else {
				n.SetString(litValue, 10)
			}
			node = CreateArrayNode(n.Uint64(), parent)
		}
		key, check, rs := g.createNodeRecursive(array.Elt, node, importMap, dotStructs, sameHierarchyStructs)
		node.Key = key
		reasons = append(reasons, rs...)
		return node, check, reasons
	}

	// map
	if mp, ok := expr.(*ast.MapType); ok {
		node := CreateMapNode(parent)
		key, c1, krs := g.createNodeRecursive(mp.Key, node, importMap, dotStructs, sameHierarchyStructs)
		value, c2, vrs := g.createNodeRecursive(mp.Value, node, importMap, dotStructs, sameHierarchyStructs)
		node.Key = key
		node.Value = value
		reasons = append(reasons, krs...)
		reasons = append(reasons, vrs...)
		return node, c1 && c2, reasons
	}

	// *
	if star, ok := expr.(*ast.StarExpr); ok {
		node := CreatePointerNode(parent)
		key, check, rs := g.createNodeRecursive(star.X, node, importMap, dotStructs, sameHierarchyStructs)
		node.Key = key
		reasons = append(reasons, rs...)
		return node, check, reasons
	}

	// not supported
	if _, ok := expr.(*ast.InterfaceType); ok {
		return nil, false, []string{fmt.Sprintf("interface type is not supported")}
	}
	if _, ok := expr.(*ast.StructType); ok {
		return nil, false, []string{fmt.Sprintf("inner struct is not supported")}
	}
	if _, ok := expr.(*ast.ChanType); ok {
		return nil, false, []string{fmt.Sprintf("chan type is not supported")}
	}
	if _, ok := expr.(*ast.FuncType); ok {
		return nil, false, []string{fmt.Sprintf("func type is not supported")}
	}

	// unreachable
	return nil, false, []string{fmt.Sprintf("this field is unknown field")}
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
