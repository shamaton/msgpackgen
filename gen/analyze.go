package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"reflect"
)

func findStructs(fileName string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		return err
	}

	structNames := make([]string, 0)
	var packageName string
	ast.Inspect(f, func(n ast.Node) bool {

		switch x := n.(type) {
		case *ast.File:
			packageName = x.Name.String()
			fmt.Println(x.Name)
		}

		x, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if _, ok := x.Type.(*ast.StructType); ok {
			structNames = append(structNames, x.Name.String())

			//for _, field := range st.Fields.List {
			//
			//	if i, ok := field.Type.(*ast.Ident); ok {
			//		fieldType := i.Name
			//
			//		for _, name := range field.Names {
			//			fmt.Printf("\tField: name=%s type=%s\n", name.Name, fieldType)
			//		}
			//	}
			//	if analyzedSt, ok := field.Type.(*ast.ArrayType); ok {
			//		if i, ok := analyzedSt.Elt.(*ast.Ident); ok {
			//
			//			fieldType := i.Name
			//
			//			for _, name := range field.Names {
			//				fmt.Printf("\tField: name=%s type=[]%s\n", name.Name, fieldType)
			//			}
			//		}
			//	}
			//
			//}
		}
		return true
	})

	for _, name := range structNames {
		analyzedSt := aaa(packageName, name, fset, f)
		analyzedStructs = append(analyzedStructs, analyzedSt)
	}
	return nil
}

type analyzedStruct struct {
	Name   string
	Fields []analyzedField
}

type analyzedField struct {
	Name string
	Type types.Type
}

func aaa(packageName, structName string, fset *token.FileSet, file *ast.File) analyzedStruct {

	conf := types.Config{
		Importer: importer.Default(),
		Error: func(err error) {
			fmt.Printf("!!! %#v\n", err)
		},
	}

	pkg, err := conf.Check(packageName, fset, []*ast.File{file}, nil)
	if err != nil {
		fmt.Println(err)
	}

	// todo : FullNameとかQual使って重複を回避する必要がある

	S := pkg.Scope().Lookup(structName)
	internal := S.Type().Underlying().(*types.Struct)

	analyzed := analyzedStruct{Name: structName}

	for i := 0; i < internal.NumFields(); i++ {
		field := internal.Field(i)

		//fmt.Println(field.Id(), field.Type().Underlying(), field.IsField())

		if field.IsField() && field.Exported() {
			name, _ := reflect.StructTag(internal.Tag(i)).Lookup("msgpack")
			if len(name) < 1 {
				name = field.Id()
			}

			analyzed.Fields = append(analyzed.Fields, analyzedField{
				Name: name,
				Type: field.Type(),
			})
		}
	}

	// todo : msgpackresolverとして出力
	return analyzed
}
