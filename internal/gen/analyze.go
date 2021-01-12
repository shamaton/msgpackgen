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

func (g *Generator) getPackages(files []string) error {
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
			g.noUserQualMap[prefix] = true
		} else if packageName == "main" {
			// todo : verbose
			// todo : 下の処理なくてもいいか
			continue
		}

		g.parseFiles = append(g.parseFiles, parseFile)
		g.fileNames = append(g.fileNames, file)
		g.parseFile2fullPackage[parseFile] = prefix
		g.fullPackage2package[prefix] = packageName
		g.targetPackages[packageName] = true
		if _, ok := g.fullPackage2ParseFiles[prefix]; !ok {
			g.fullPackage2ParseFiles[prefix] = make([]*ast.File, 0)
		}
		g.fullPackage2ParseFiles[prefix] = append(g.fullPackage2ParseFiles[prefix], parseFile)
	}
	return nil
}

func (g *Generator) createAnalyzedStructs() error {
	analyzedMap := map[*ast.File]bool{}
	for _, parseFile := range g.parseFiles {
		// done analysis
		if _, ok := analyzedMap[parseFile]; ok {
			continue
		}

		fullPackageName, ok := g.parseFile2fullPackage[parseFile]
		if !ok {
			return fmt.Errorf("not found fullPackageName")
		}
		packageName, ok := g.fullPackage2package[fullPackageName]
		if !ok {
			return fmt.Errorf("not found package name")
		}

		err := g.hogehoge(parseFile, packageName, fullPackageName, analyzedMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) hogehoge(parseFile *ast.File, packageName, fullPackageName string, analyzedMap map[*ast.File]bool) error {

	importMap, dotImports := g.createImportMap(parseFile)
	// todo : ドットインポートが見つかった場合先にそのファイルを解析してしまうようにする
	//
	for _, dotImport := range dotImports {
		pfs, ok := g.fullPackage2ParseFiles[dotImport]
		if !ok {
			return fmt.Errorf("%s not found parse files", dotImport)
		}
		name, ok := g.fullPackage2package[dotImport]
		if !ok {
			return fmt.Errorf("not found package name")
		}

		for _, pf := range pfs {
			// todo : 前後のところ含めた関数化が必要
			err := g.hogehoge(pf, name, dotImport, analyzedMap)
			if err != nil {
				return err
			}
			analyzedMap[pf] = true
		}
	}

	structs := g.createAnalyzedStructsPerFile(parseFile, packageName, fullPackageName, importMap)
	analyzedStructs = append(analyzedStructs, structs...)
	analyzedMap[parseFile] = true
	return nil
}

func (g *Generator) createAnalyzedStructsPerFile(parseFile *ast.File, packageName, fullPackageName string, importMap map[string]string) []analyzedStruct {

	structNames := make([]string, 0)
	analyzedFieldMap := map[string]*analyzedASTFieldType{}
	ast.Inspect(parseFile, func(n ast.Node) bool {

		x, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if st, ok := x.Type.(*ast.StructType); ok {

			// todo : 出力パッケージの場所と同じならLowerでもOK

			if fullPackageName != g.outputPackageFullName() && !unicode.IsUpper(rune(x.Name.String()[0])) {
				return true
			}

			canGen := true
			for _, field := range st.Fields.List {

				key := ""
				for _, name := range field.Names {
					key = name.Name
				}

				// todo : dotImportMapが必要
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

	structs := make([]analyzedStruct, len(structNames))
	for i, structName := range structNames {
		fmt.Println()
		fmt.Println()
		fmt.Println(structName, ".........................................", g.noUserQualMap[fullPackageName])
		fields := g.createAnalyzedFields(packageName, structName, analyzedFieldMap, g.fileSet, parseFile)
		structs[i] = analyzedStruct{
			PackageName: fullPackageName,
			Name:        structName,
			Fields:      fields,
			NoUseQual:   g.noUserQualMap[fullPackageName],
		}

	}
	return structs
}

func (g *Generator) createImportMap(parseFile *ast.File) (map[string]string, []string) {

	importMap := map[string]string{}
	dotImports := make([]string, 0)

	for _, imp := range parseFile.Imports {

		fmt.Println(imp.Name, imp.Path.Value)

		value := strings.ReplaceAll(imp.Path.Value, "\"", "")

		if imp.Name == nil || imp.Name.Name == "" {
			key := strings.Split(value, "/")
			importMap[key[len(key)-1]] = value
		} else if imp.Name.Name == "." {
			dotImports = append(dotImports, value)
		} else {
			key := strings.ReplaceAll(imp.Name.Name, "\"", "")
			importMap[key] = value
		}
	}
	return importMap, dotImports
}

func (g *Generator) createAnalyzedFields(packageName, structName string, analyzedFieldMap map[string]*analyzedASTFieldType, fset *token.FileSet, file *ast.File) []analyzedField {

	// todo : ここなにか解決策あれば
	imp := importer.Default()
	//_, err := imp.Import("github.com/shamaton/msgpackgen/internal/tst/tst")
	//if err != nil {
	//	fmt.Println("import error", err)
	//}
	conf := types.Config{
		Importer: imp,
		Error: func(err error) {
			// fmt.Printf("!!! %#v\n", err)
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
