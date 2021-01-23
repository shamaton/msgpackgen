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

func (a Node) CanGenerate(sts []analyzedStruct) (bool, []string) {
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

func (a Node) TypeJenChain(sts []analyzedStruct, s ...*Statement) *Statement {
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
			found := false
			asRef := analyzedStruct{}
			for _, v := range sts {
				if v.ImportPath == a.ImportPath && v.Name == a.StructName {
					found = true
					asRef = v
					break
				}
			}
			if !found {
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

func (g *generator) createNodeRecursive(expr ast.Expr, parent *Node, importMap map[string]string, dotStructs map[string]analyzedStruct, sameHierarchyStructs map[string]bool) (*Node, bool, []string) {

	reasons := make([]string, 0)
	if i, ok := expr.(*ast.Ident); ok {

		// dot import
		if dot, found := dotStructs[i.Name]; found {
			return &Node{
				fieldType:   fieldTypeStruct,
				PackageName: dot.Name,
				StructName:  i.Name,
				ImportPath:  dot.ImportPath,
				Parent:      parent,
			}, true, reasons
		}
		// time
		if i.Name == "Time" {
			return &Node{
				fieldType:   fieldTypeStruct,
				PackageName: "time",
				StructName:  i.Name,
				ImportPath:  "time",
				Parent:      parent,
			}, true, reasons
		}
		// same hierarchy struct in same file
		if i.Obj != nil && i.Obj.Kind == ast.Typ {
			return &Node{
				fieldType:   fieldTypeStruct,
				PackageName: g.outputPackageName,
				StructName:  i.Name,
				ImportPath:  g.outputPackageFullName(),
				Parent:      parent,
			}, true, reasons
		}

		// same hierarchy struct in other file
		if _, found := sameHierarchyStructs[i.Name]; found {
			return &Node{
				fieldType:   fieldTypeStruct,
				PackageName: g.outputPackageName,
				StructName:  i.Name,
				ImportPath:  g.outputPackageFullName(),
				Parent:      parent,
			}, true, reasons
		}

		if isPrimitive(i.Name) {
			return &Node{
				fieldType:     fieldTypeIdent,
				IdenticalName: i.Name,
				Parent:        parent,
			}, true, reasons
		}

		return nil, false, []string{fmt.Sprintf("identifier %s is not suppoted or unknown struct ", i.Name)}
	}

	if selector, ok := expr.(*ast.SelectorExpr); ok {
		pkgName := fmt.Sprint(selector.X)
		return &Node{
			fieldType:   fieldTypeStruct,
			PackageName: pkgName,
			StructName:  selector.Sel.Name,
			ImportPath:  importMap[pkgName],
			Parent:      parent,
		}, true, reasons
	}

	// slice or array
	if array, ok := expr.(*ast.ArrayType); ok {
		var node *Node
		if array.Len == nil {
			node = &Node{
				fieldType: fieldTypeSlice,
				Parent:    parent,
			}
		} else {
			lit := array.Len.(*ast.BasicLit)
			// todo : 処理されなかった場合はエラー
			// todo : box数値以外あればエラーでもいい
			// parse num
			n := new(big.Int)
			if litValie := strings.ToLower(lit.Value); strings.HasPrefix(litValie, "0b") {
				n.SetString(strings.ReplaceAll(litValie, "0b", ""), 2)
			} else if strings.HasPrefix(litValie, "0o") {
				n.SetString(strings.ReplaceAll(litValie, "0o", ""), 8)
			} else if strings.HasPrefix(litValie, "0x") {
				n.SetString(strings.ReplaceAll(litValie, "0x", ""), 16)
			} else {
				n.SetString(litValie, 10)
			}
			node = &Node{
				fieldType: fieldTypeArray,
				ArrayLen:  n.Uint64(),
				Parent:    parent,
			}
		}
		key, check, rs := g.createNodeRecursive(array.Elt, node, importMap, dotStructs, sameHierarchyStructs)
		node.Key = key
		reasons = append(reasons, rs...)
		return node, check, reasons
	}

	// map
	if mp, ok := expr.(*ast.MapType); ok {
		node := &Node{
			fieldType: fieldTypeMap,
			Parent:    parent,
		}
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
		node := &Node{
			fieldType: fieldTypePointer,
			Parent:    parent,
		}
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
