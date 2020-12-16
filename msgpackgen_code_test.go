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

	f := NewFilePath("a.b/c")
	f.Func().Id("init").Params().Block(
		Qual("a.b/c", "Foo").Call().Comment("Local package - name is omitted."),
		Qual("d.e/f", "Bar").Call().Comment("Import is automatically added."),
		Qual("g.h/f", "Baz").Call().Comment("Colliding package name is renamed."),
	)

	f.Func().Id("decode").Params(Id("data").Index().Byte(), Id("i").Interface()).Params(Bool(), Error()).Block()
	fmt.Printf("%#v", f)
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
	aaa(packageName, name, fset, f)
	return name, nil
}

func aaa(packageName, structName string, fset *token.FileSet, file *ast.File) {

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

	S := pkg.Scope().Lookup(structName)
	internal := S.Type().Underlying().(*types.Struct)

	for i := 0; i < internal.NumFields(); i++ {
		tag, found := reflect.StructTag(internal.Tag(i)).Lookup("msgpack")
		field := internal.Field(i)
		fmt.Printf("%v (exported=%t, tag=%s, found=%t)\n", field, field.Exported(), tag, found)
	}

	// todo : msgpackresolverとして出力
}
