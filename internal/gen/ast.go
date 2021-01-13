package gen

import (
	"fmt"
	"go/ast"

	. "github.com/dave/jennifer/jen"
)

const (
	fieldTypeIdent = iota + 1
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

func (a analyzedASTFieldType) TypeJenChain(s ...*Statement) *Statement {
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
		// todo : performance
		var asRef analyzedStruct
		for _, v := range analyzedStructs {
			if v.PackageName == a.ImportPath && v.Name == a.StructName {
				asRef = v
			}
		}
		if asRef.NoUseQual {
			str = str.Id(a.StructName)
		} else {
			str = str.Qual(a.ImportPath, a.StructName)
		}

	case a.IsArray():
		str = str.Id("[]")
		str = a.Elm().TypeJenChain(str)

	case a.IsMap():
		str = str.Id("map[")
		k, v := a.KeyValue()
		str = k.TypeJenChain(str)
		str = str.Id("]")
		str = v.TypeJenChain(str)

	case a.IsPointer():
		str = str.Id("*")
		str = a.Elm().TypeJenChain(str)
	}
	return str
}

func (a analyzedASTFieldType) TypeString(s ...string) string {
	str := ""
	if len(s) > 0 {
		str += s[0]
	}

	switch {
	case a.IsIdentical():
		str += a.IdenticalName

	case a.IsStruct():

	case a.IsArray():
		str += "[]"
		str = a.Elm().TypeString(str)

	case a.IsMap():
		str += "map["
		k, v := a.KeyValue()
		str = k.TypeString(str)
		str += "]"
		str = v.TypeString(str)

	case a.IsPointer():
		str += "*"
		str = a.Elm().TypeString(str)
	}
	return str
}

func (g *Generator) checkFieldTypeRecursive(expr ast.Expr, parent *analyzedASTFieldType, importMap map[string]string, dotStructs map[string]analyzedStruct) (*analyzedASTFieldType, bool) {
	if i, ok := expr.(*ast.Ident); ok {
		// todo : 整理
		fmt.Println(i.Name, i.Obj, expr)

		if dot, found := dotStructs[i.Name]; found {
			return &analyzedASTFieldType{
				fieldType:   fieldTypeStruct,
				PackageName: dot.Name,
				StructName:  i.String(),
				ImportPath:  dot.PackageName,
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
		// can not generate
		if i.Name == "uintptr" || i.Name == "error" {
			return nil, false
		}

		// todo : complex
		if i.Name == "complex64" || i.Name == "complex128" {
			return nil, false
		}

		return &analyzedASTFieldType{
			fieldType:     fieldTypeIdent,
			IdenticalName: i.Name,
			Parent:        parent,
		}, true
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
		node := &analyzedASTFieldType{
			fieldType: fieldTypeArray,
			Parent:    parent,
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
	return nil, false
}
