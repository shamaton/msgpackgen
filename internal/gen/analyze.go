package gen

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"

	. "github.com/dave/jennifer/jen"
)

func (g *Generator) GetPackages(files []string) error {
	g.fileSet = token.NewFileSet()
	for _, file := range files {

		dir := filepath.Dir(file)
		paths := strings.SplitN(dir, "src/", 2)
		if len(paths) != 2 {
			return fmt.Errorf("%s get import path failed", file)
		}
		prefix := paths[1]

		parseFile, err := parser.ParseFile(g.fileSet, file, nil, 0)
		if err != nil {
			return err
		}

		var packageName string
		ast.Inspect(parseFile, func(n ast.Node) bool {

			switch x := n.(type) {
			case *ast.File:
				packageName = x.Name.String()
				//fmt.Println(x.Name)
			}

			return true
		})

		if dir == g.outputDir {
			g.outputPackageName = packageName
			g.noUserQualMap[file] = true
		} else if packageName == "main" {
			// todo : verbose
			continue
		}

		g.parseFiles = append(g.parseFiles, parseFile)
		g.fileNames = append(g.fileNames, file)
		g.file2PackageName[file] = packageName
		g.file2FullPackageName[file] = prefix
		g.targetPackages[packageName] = true
	}
	return nil
}

func (g *Generator) CreateAnalyzedStructs() error {

	for i, parseFile := range g.parseFiles {
		fileName := g.fileNames[i]
		importMap := map[string]string{}

		for _, imp := range parseFile.Imports {

			value := strings.ReplaceAll(imp.Path.Value, "\"", "")

			if imp.Name == nil || imp.Name.Name == "" {
				key := strings.Split(value, "/")
				importMap[key[len(key)-1]] = value
			} else {
				key := strings.ReplaceAll(imp.Name.Name, "\"", "")
				importMap[key] = value
			}
		}

		structNames := make([]string, 0)
		analyzedFieldMap := map[string]*analyzedASTFieldType{}
		ast.Inspect(parseFile, func(n ast.Node) bool {

			x, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			if st, ok := x.Type.(*ast.StructType); ok {

				// todo : 出力パッケージの場所と同じならLowerでもOK

				if g.file2FullPackageName[fileName] != g.OutputPackageFullName() && !unicode.IsUpper(rune(x.Name.String()[0])) {
					return true
				}

				canGen := true
				for _, field := range st.Fields.List {

					key := ""
					for _, name := range field.Names {
						key = name.Name
					}

					value, ok := g.checkFieldTypeRecursive(field.Type, nil, importMap)
					canGen = canGen && ok
					if ok {
						analyzedFieldMap[key+"@"+x.Name.String()] = value
					}
				}
				if canGen {
					structNames = append(structNames, x.Name.String())
				}
			}
			return true
		})

		for _, structName := range structNames {
			fmt.Println()
			fmt.Println()
			fmt.Println(structName, ".........................................", g.noUserQualMap[fileName])
			fields := g.createAnalyzedFields(g.file2PackageName[fileName], structName, analyzedFieldMap, g.fileSet, parseFile)
			analyzedStructs = append(analyzedStructs, analyzedStruct{
				PackageName: g.file2FullPackageName[fileName],
				Name:        structName,
				Fields:      fields,
				NoUseQual:   g.noUserQualMap[fileName],
			})

		}
	}
	return nil
}

func (g *Generator) createAnalyzedFields(packageName, structName string, analyzedFieldMap map[string]*analyzedASTFieldType, fset *token.FileSet, file *ast.File) []analyzedField {

	// todo : ここなにか解決策あれば
	imp := importer.Default()
	_, err := imp.Import("github.com/shamaton/tetest/example/item")
	if err != nil {
		fmt.Println("import error", err)
	}
	conf := types.Config{
		Importer: imp,
		Error: func(err error) {
			//fmt.Printf("!!! %#v\n", err)
		},
	}

	pkg, err := conf.Check(packageName, fset, []*ast.File{file}, nil)
	if err != nil {
		fmt.Println(err)
	}

	// todo : FullNameとかQual使って重複を回避する必要がある

	S := pkg.Scope().Lookup(structName)
	internal := S.Type().Underlying().(*types.Struct)

	analyzedFields := make([]analyzedField, 0)
	for i := 0; i < internal.NumFields(); i++ {
		field := internal.Field(i)

		// fmt.Println(field.Id(), field.Type(), field.IsField())

		if field.IsField() && field.Exported() {
			tagName, _ := reflect.StructTag(internal.Tag(i)).Lookup("msgpack")
			if tagName == "ignore" {
				continue
			}
			name := field.Id()
			tag := name
			if len(tagName) > 0 {
				tag = tagName
			}

			//fmt.Println("hogehoge", reflect.TypeOf(field.Type()))

			// todo : type.Namedの場合、解析対象に含まれてないものがあったら、スキップする？
			// todo : タグが重複してたら、エラー

			analyzedFields = append(analyzedFields, analyzedField{
				Name: name,
				Tag:  tag,
				Type: field.Type(),
				Ast:  analyzedFieldMap[name+"@"+structName],
			})
		}
	}

	// todo : msgpackresolverとして出力
	return analyzedFields
}

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
	fmt.Println(a.ImportPath, a.StructName, a.fieldType, a.IdenticalName)
	switch {
	case a.IsIdentical():
		return true, msgs

	case a.IsStruct():
		fmt.Println("GGGGGGGGGGGGGGGGGGGGGGG", a.ImportPath, a.StructName)
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

func (g *Generator) checkFieldTypeRecursive(expr ast.Expr, parent *analyzedASTFieldType, importMap map[string]string) (*analyzedASTFieldType, bool) {
	if i, ok := expr.(*ast.Ident); ok {
		fmt.Println("HHHHHHHHHHHHHHH", i.String(), i.Obj)
		// todo : 整理
		// same hierarchy struct
		if i.Obj != nil && i.Obj.Kind == ast.Typ {
			pkgName := g.OutputPackageFullName()
			return &analyzedASTFieldType{
				fieldType:   fieldTypeStruct,
				PackageName: g.OutputPackageFullName(),
				StructName:  i.String(),
				ImportPath:  importMap[pkgName],
				Parent:      parent,
			}, true
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
		key, check := g.checkFieldTypeRecursive(array.Elt, node, importMap)
		node.Key = key
		return node, check
	}
	if mp, ok := expr.(*ast.MapType); ok {
		node := &analyzedASTFieldType{
			fieldType: fieldTypeMap,
			Parent:    parent,
		}
		key, c1 := g.checkFieldTypeRecursive(mp.Key, node, importMap)
		value, c2 := g.checkFieldTypeRecursive(mp.Value, node, importMap)
		node.Key = key
		node.Value = value
		return node, c1 && c2
	}
	if star, ok := expr.(*ast.StarExpr); ok {
		node := &analyzedASTFieldType{
			fieldType: fieldTypePointer,
			Parent:    parent,
		}
		key, check := g.checkFieldTypeRecursive(star.X, node, importMap)
		node.Key = key
		return node, check
	}
	if _, ok := expr.(*ast.InterfaceType); ok {
		return nil, false
	}
	return nil, false
}
