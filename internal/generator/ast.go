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
	//fieldTypeInterface
)

type analyzedASTFieldType struct {
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
	Key   *analyzedASTFieldType
	Value *analyzedASTFieldType

	Parent *analyzedASTFieldType
}

func (a analyzedASTFieldType) HasParent() bool { return a.Parent != nil }

func (a analyzedASTFieldType) IsIdentical() bool { return a.fieldType == fieldTypeIdent }
func (a analyzedASTFieldType) IsSlice() bool     { return a.fieldType == fieldTypeSlice }
func (a analyzedASTFieldType) IsArray() bool     { return a.fieldType == fieldTypeArray }
func (a analyzedASTFieldType) IsStruct() bool    { return a.fieldType == fieldTypeStruct }
func (a analyzedASTFieldType) IsMap() bool       { return a.fieldType == fieldTypeMap }

//func (a analyzedASTFieldType) IsInterface() bool { return a.fieldType == fieldTypeInterface }
func (a analyzedASTFieldType) IsPointer() bool { return a.fieldType == fieldTypePointer }

func (a analyzedASTFieldType) Elm() *analyzedASTFieldType { return a.Key }
func (a analyzedASTFieldType) KeyValue() (*analyzedASTFieldType, *analyzedASTFieldType) {
	return a.Key, a.Value
}

func (a analyzedASTFieldType) CanGenerate(sts []analyzedStruct) (bool, []string) {
	msgs := make([]string, 0)
	switch {
	case a.IsIdentical():
		return true, msgs

	case a.IsStruct():
		if a.ImportPath == "time" && a.StructName == "Time" {
			return true, msgs
		}
		// todo : performance
		for _, v := range sts {
			if v.PackageName == a.ImportPath && v.Name == a.StructName {
				return true, msgs
			}
		}
		return false, append(msgs, fmt.Sprintf("struct %s.%s is not generated.", a.ImportPath, a.StructName))

	case a.IsSlice():
		return a.Elm().CanGenerate(sts)

	case a.IsArray():
		return a.Elm().CanGenerate(sts)

	case a.IsMap():
		k, v := a.KeyValue()
		kb, kMsgs := k.CanGenerate(sts)
		vb, vMsgs := v.CanGenerate(sts)
		msgs = append(msgs, kMsgs...)
		msgs = append(msgs, vMsgs...)
		return kb && vb, msgs

	case a.IsPointer():
		return a.Elm().CanGenerate(sts)
	}
	return false, append(msgs, "unreachable code")
}

func (a analyzedASTFieldType) TypeJenChain(sts []analyzedStruct, s ...*Statement) *Statement {
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
				if v.PackageName == a.ImportPath && v.Name == a.StructName {
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

func (g *generator) checkFieldTypeRecursive(expr ast.Expr, parent *analyzedASTFieldType, importMap map[string]string, dotStructs map[string]analyzedStruct) (*analyzedASTFieldType, bool) {

	if i, ok := expr.(*ast.StructType); ok {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>", i.Fields)
		for _, f := range i.Fields.List {
			fmt.Println(*f)
		}
		// todo : nested struct not support comment
	}
	if i, ok := expr.(*ast.Ident); ok {

		// dot import
		if dot, found := dotStructs[i.Name]; found {
			return &analyzedASTFieldType{
				fieldType:   fieldTypeStruct,
				PackageName: dot.Name,
				StructName:  i.String(),
				ImportPath:  dot.PackageName,
				Parent:      parent,
			}, true
		}
		// time
		if i.Name == "Time" {
			return &analyzedASTFieldType{
				fieldType:   fieldTypeStruct,
				PackageName: "time",
				StructName:  i.String(),
				ImportPath:  "time",
				Parent:      parent,
			}, true
		}
		// same hierarchy struct
		if i.Obj != nil && i.Obj.Kind == ast.Typ {
			return &analyzedASTFieldType{
				fieldType:   fieldTypeStruct,
				PackageName: g.outputPackageName,
				StructName:  i.String(),
				ImportPath:  g.outputPackageFullName(),
				Parent:      parent,
			}, true
		}

		// todo : 型を念の為判定しておく必要がありそう
		if isPrimitive(i.Name) {
			return &analyzedASTFieldType{
				fieldType:     fieldTypeIdent,
				IdenticalName: i.Name,
				Parent:        parent,
			}, true
		}

		// can not generate
		// todo : error skip??
		// todo : 同じ階層で別ファイルのstructがここに来てしまう
		fmt.Println("shamoto0000000000000000000000000000000000", i.Obj, i.Name)
		return nil, false
	}
	if selector, ok := expr.(*ast.SelectorExpr); ok {
		pkgName := fmt.Sprint(selector.X) // todo : ok?
		return &analyzedASTFieldType{
			fieldType:   fieldTypeStruct,
			PackageName: pkgName,
			StructName:  selector.Sel.Name,
			ImportPath:  importMap[pkgName],
			Parent:      parent,
		}, true
	}
	if array, ok := expr.(*ast.ArrayType); ok {
		var node *analyzedASTFieldType
		if array.Len == nil {
			node = &analyzedASTFieldType{
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
			node = &analyzedASTFieldType{
				fieldType: fieldTypeArray,
				ArrayLen:  n.Uint64(),
				Parent:    parent,
			}
		}
		key, check := g.checkFieldTypeRecursive(array.Elt, node, importMap, dotStructs)
		node.Key = key
		return node, check
	}
	if mp, ok := expr.(*ast.MapType); ok {
		node := &analyzedASTFieldType{
			fieldType: fieldTypeMap,
			Parent:    parent,
		}
		key, c1 := g.checkFieldTypeRecursive(mp.Key, node, importMap, dotStructs)
		value, c2 := g.checkFieldTypeRecursive(mp.Value, node, importMap, dotStructs)
		node.Key = key
		node.Value = value
		return node, c1 && c2
	}
	if star, ok := expr.(*ast.StarExpr); ok {
		node := &analyzedASTFieldType{
			fieldType: fieldTypePointer,
			Parent:    parent,
		}
		key, check := g.checkFieldTypeRecursive(star.X, node, importMap, dotStructs)
		node.Key = key
		return node, check
	}
	if _, ok := expr.(*ast.InterfaceType); ok {
		return nil, false
	}

	// unreachable
	return nil, false
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
