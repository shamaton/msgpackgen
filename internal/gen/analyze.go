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
			g.outputPackagePrefix = filepath.Dir(prefix)
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
