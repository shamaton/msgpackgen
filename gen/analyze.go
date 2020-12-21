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
	for _, file := range files {

		dir := filepath.Dir(file)
		paths := strings.SplitN(dir, "src/", 2)
		if len(paths) != 2 {
			return fmt.Errorf("%s get import path failed", file)
		}
		prefix := paths[1]

		fileSet := token.NewFileSet()
		parseFile, err := parser.ParseFile(fileSet, file, nil, 0)
		if err != nil {
			return err
		}

		var packageName string
		ast.Inspect(parseFile, func(n ast.Node) bool {

			switch x := n.(type) {
			case *ast.File:
				packageName = prefix + "/" + x.Name.String()
				//fmt.Println(x.Name)
			}

			return true
		})

		g.file2Parse[file] = parseFile
		g.file2PackageName[file] = packageName
		g.targetPackages[packageName] = true
	}
	return nil
}

func (g *generator) createAnalyzedStructs() error {
	for _, parseFile := range g.file2Parse {
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
		ast.Inspect(parseFile, func(n ast.Node) bool {

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
						fmt.Println("1111111111", stt.X, stt.Sel.Name)
					}
					if analyzedSt, ok := field.Type.(*ast.ArrayType); ok {
						if i, ok := analyzedSt.Elt.(*ast.Ident); ok {

							fieldType := i.Name

							for _, name := range field.Names {
								fmt.Printf("\tField: name=%s type=[]%s\n", name.Name, fieldType)
							}
						}
					}

				}
			}
			return true
		})
	}
	return nil
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
					fmt.Println("1111111111", stt.X, stt.Sel.Name, field.Names)

				}
				if analyzedSt, ok := field.Type.(*ast.ArrayType); ok {
					if i, ok := analyzedSt.Elt.(*ast.Ident); ok {

						fieldType := i.Name

						for _, name := range field.Names {
							fmt.Printf("\tField: name=%s type=[]%s\n", name.Name, fieldType)
						}
					}
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
