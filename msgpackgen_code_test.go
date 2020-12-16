package msgpackgen_test

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"reflect"
	"testing"

	. "github.com/dave/jennifer/jen"
)

func TestGenerate(t *testing.T) {
	fmt.Println(findStructs("msgpackgen_test.go"))
	jen()

	return

}

func jen() {
	packageName := "msgpackgen"

	f := NewFilePath("a.b/c")
	f.Func().Id("init").Params().Block(
		Qual("a.b/c", "Foo").Call().Comment("Local package - name is omitted."),
		Qual("d.e/f", "Bar").Call().Comment("Import is automatically added."),
		Qual("g.h/f", "Baz").Call().Comment("Colliding package name is renamed."),
	)

	decodeTopTemplate("decode", f).Block(
		// todo : Qual
		If(Id(packageName + ".StructAsArray").Call()).Block(
			Return(Id("decodeAsArray").Call(Id("data"), Id("i"))),
		).Else().Block(
			Return(Id("decodeAsMap").Call(Id("data"), Id("i"))),
		),
	)

	decodeTopTemplate("decodeAsArray", f).Block(
		Switch(Id("v").Op(":=").Id("i").Assert(Type())).Block(
			cases()...,
		),
		Return(Nil(), Nil()),
	)

	decodeTopTemplate("decodeAsMap", f).Block(
		Return(Nil(), Nil()),
	)

	f.Func().Id("encode").Params(Id("i").Interface()).Params(Id("data").Index().Byte(), Error()).Block(
		If(Id(packageName + ".StructAsArray").Call()).Block(
			Return(Id("encodeAsArray").Call(Id("i"))),
		).Else().Block(
			Return(Id("encodeAsMap").Call(Id("i"))),
		),
	)

	fmt.Printf("%#v", f)
}

func decodeTopTemplate(name string, f *File) *Statement {
	return f.Func().Id(name).Params(Id("data").Index().Byte(), Id("i").Interface()).Params(Bool(), Error())
}

func cases() []Code {
	var states []Code
	for _, v := range analyzedStructs {
		states = append(states, Case(Id(v.Name)).Block(
			Return(Id("_"), Err())))
	}
	return states
}

func findStructs(fileName string) (string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		return "", err
	}

	// 最初に見つかったstruct名を入れる
	var name string
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
			fmt.Println(x.Name.String())
			if name == "" {
				name = x.Name.String()
			}

			//for _, field := range st.Fields.List {
			//
			//	if i, ok := field.Type.(*ast.Ident); ok {
			//		fieldType := i.Name
			//
			//		for _, name := range field.Names {
			//			fmt.Printf("\tField: name=%s type=%s\n", name.Name, fieldType)
			//		}
			//	}
			//	if a, ok := field.Type.(*ast.ArrayType); ok {
			//		if i, ok := a.Elt.(*ast.Ident); ok {
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
	a := aaa(packageName, name, fset, f)
	analyzedStructs = append(analyzedStructs, a)
	fmt.Println(analyzedStructs)
	return name, nil
}

var analyzedStructs []analyzedStruct

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
