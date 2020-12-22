package main

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
)

func (g *generator) getPackages(files []string) error {
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

		g.file2Parse[file] = parseFile
		g.file2PackageName[file] = packageName
		g.file2FullPackageName[file] = prefix + "/" + packageName
		g.targetPackages[packageName] = true
	}
	return nil
}

func (g *generator) createAnalyzedStructs() error {

	for fileName, parseFile := range g.file2Parse {
		importMap := map[string]string{}

		for _, imp := range parseFile.Imports {

			if imp.Name == nil || imp.Name.Name == "" {
				vs := strings.Split(imp.Path.Value, "/")
				importMap[vs[len(vs)-1]] = imp.Path.Value
			} else {
				importMap[imp.Name.Name] = imp.Path.Value
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

				canGen := true
				for _, field := range st.Fields.List {

					key := ""
					for _, name := range field.Names {
						key = name.Name
					}

					value, ok := g.checkFieldTypeRecursive(field.Type)
					canGen = canGen && ok
					if ok {
						if value.IsStruct() {
							value.ImportPath = importMap[value.PackageName]
						}
						analyzedFieldMap[key] = value
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
			fmt.Println(structName, ".........................................")
			fields := g.createAnalyzedFields(g.file2PackageName[fileName], structName, analyzedFieldMap, g.fileSet, parseFile)
			if len(fields) > 0 {
				analyzedStructs = append(analyzedStructs, analyzedStruct{
					PackageName: g.file2FullPackageName[fileName],
					Name:        structName,
					Fields:      fields,
				})
			}
		}
	}
	return nil
}

func (g *generator) createAnalyzedFields(packageName, structName string, analyzedFieldMap map[string]*analyzedASTFieldType, fset *token.FileSet, file *ast.File) []analyzedField {

	conf := types.Config{
		Importer: importer.Default(),
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
			if len(tagName) > 0 {
				name = tagName
			}

			//fmt.Println("hogehoge", reflect.TypeOf(field.Type()))

			// todo : type.Namedの場合、解析対象に含まれてないものがあったら、スキップする？

			analyzedFields = append(analyzedFields, analyzedField{
				Name: name,
				Type: field.Type(),
				Ast:  analyzedFieldMap[name],
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

	ImportPath  string
	PackageName string
	StructName  string

	Key   *analyzedASTFieldType
	Value *analyzedASTFieldType
}

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

func (g *generator) checkFieldTypeRecursive(expr ast.Expr) (*analyzedASTFieldType, bool) {
	if _, ok := expr.(*ast.Ident); ok {
		return &analyzedASTFieldType{fieldType: fieldTypeIdent}, true
	}
	if selector, ok := expr.(*ast.SelectorExpr); ok {
		return &analyzedASTFieldType{
			fieldType:   fieldTypeStruct,
			PackageName: fmt.Sprint(selector.X), // todo : ok?
			StructName:  selector.Sel.Name,
		}, true
	}
	if array, ok := expr.(*ast.ArrayType); ok {
		key, check := g.checkFieldTypeRecursive(array.Elt)
		return &analyzedASTFieldType{
			fieldType: fieldTypeArray,
			Key:       key,
		}, check
	}
	if mp, ok := expr.(*ast.MapType); ok {
		key, c1 := g.checkFieldTypeRecursive(mp.Key)
		value, c2 := g.checkFieldTypeRecursive(mp.Value)
		return &analyzedASTFieldType{
			fieldType: fieldTypeMap,
			Key:       key,
			Value:     value,
		}, c1 && c2
	}
	if star, ok := expr.(*ast.StarExpr); ok {
		key, check := g.checkFieldTypeRecursive(star.X)
		return &analyzedASTFieldType{
			fieldType: fieldTypePointer,
			Key:       key,
		}, check
	}
	if _, ok := expr.(*ast.InterfaceType); ok {
		return nil, false
	}
	return nil, false
}

func (g *generator) findStructs(fileName string) error {
	dir := filepath.Dir(fileName)
	paths := strings.SplitN(dir, "src/", 2)
	if len(paths) != 2 {
		return fmt.Errorf("error...")
	}

	prefix := paths[1]

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		return err
	}

	for _, v := range f.Imports {
		fmt.Println(v.Name, v.Path.Value)
	}

	structNames := make([]string, 0)
	var packageName string
	ast.Inspect(f, func(n ast.Node) bool {

		switch x := n.(type) {
		case *ast.File:
			packageName = x.Name.String()
			//fmt.Println(x.Name)
		}

		x, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if st, ok := x.Type.(*ast.StructType); ok {
			structNames = append(structNames, x.Name.String())

			for _, field := range st.Fields.List {

				fmt.Println(reflect.TypeOf(field.Type))

				if i, ok := field.Type.(*ast.Ident); ok {
					fieldType := i.Name

					for _, name := range field.Names {
						fmt.Printf("\tField:  \n")
						fmt.Printf("\tField: name=%s type=%s\n", name.Name, fieldType)
					}
				}
				if stt, ok := field.Type.(*ast.SelectorExpr); ok {
					fmt.Println("1111111111", stt.X, stt.Sel.Name, field.Names, reflect.TypeOf(stt.X), fmt.Sprint(stt.X))

				}
				if analyzedSt, ok := field.Type.(*ast.ArrayType); ok {
					if i, ok := analyzedSt.Elt.(*ast.Ident); ok {

						fieldType := i.Name

						for _, name := range field.Names {
							fmt.Printf("\tField: name=%s type=[]%s\n", name.Name, fieldType)
						}
					}
				}

				if star, ok := field.Type.(*ast.StarExpr); ok {
					fmt.Println("starrrrrrrrr", star.X, reflect.TypeOf(star.X))
				}

			}
		}
		return true
	})

	for _, name := range structNames {
		analyzedSt := aaa(packageName, name, fset, f)
		analyzedSt.PackageName = prefix + "/" + packageName
		analyzedStructs = append(analyzedStructs, analyzedSt)
	}
	return nil
}

func aaa(packageName, structName string, fset *token.FileSet, file *ast.File) analyzedStruct {

	conf := types.Config{
		Importer: importer.Default(),
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

	analyzed := analyzedStruct{
		PackageName: packageName,
		Name:        structName,
	}

	for i := 0; i < internal.NumFields(); i++ {
		field := internal.Field(i)

		fmt.Printf("gugug %v\n", field)
		fmt.Println(field.Id(), field.Type(), field.IsField())

		if field.IsField() && field.Exported() {
			tagName, _ := reflect.StructTag(internal.Tag(i)).Lookup("msgpack")
			if tagName == "ignore" {
				continue
			}
			name := field.Id()
			if len(tagName) > 0 {
				name = tagName
			}

			//fmt.Println("hogehoge", reflect.TypeOf(field.Type()))

			// todo : type.Namedの場合、解析対象に含まれてないものがあったら、スキップする？

			analyzed.Fields = append(analyzed.Fields, analyzedField{
				Name: name,
				Type: field.Type(),
			})
		}
	}

	// todo : msgpackresolverとして出力
	return analyzed
}
