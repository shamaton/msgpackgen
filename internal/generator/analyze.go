package generator

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
)

func (g *generator) getPackages(files []string) error {
	g.fileSet = token.NewFileSet()
	for _, file := range files {

		dir := filepath.Dir(file)
		paths := strings.SplitN(filepath.ToSlash(dir), "src/", 2)
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

func (g *generator) analyze() error {
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

		err := g.createAnalyzedStructs(parseFile, packageName, fullPackageName, analyzedMap)
		if err != nil {
			return err
		}
	}
	g.setFieldToStruct()
	return nil
}

func (g *generator) createAnalyzedStructs(parseFile *ast.File, packageName, importPath string, analyzedMap map[*ast.File]bool) error {

	importMap, dotImports := g.createImportMap(parseFile)
	// dot imports
	dotStructs := map[string]analyzedStruct{}
	for _, dotImport := range dotImports {
		pfs, ok := g.fullPackage2ParseFiles[dotImport]
		if !ok {
			continue
		}
		name, ok := g.fullPackage2package[dotImport]
		if !ok {
			continue
		}

		for _, pf := range pfs {
			err := g.createAnalyzedStructs(pf, name, dotImport, analyzedMap)
			if err != nil {
				return err
			}
			analyzedMap[pf] = true
		}

		for _, st := range analyzedStructs {
			if st.ImportPath == dotImport {
				dotStructs[st.Name] = st
			}
		}
	}

	structNames := make([]string, 0)
	ast.Inspect(parseFile, func(n ast.Node) bool {

		x, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if _, ok := x.Type.(*ast.StructType); ok {

			structName := x.Name.String()
			if importPath != g.outputPackageFullName() && !unicode.IsUpper(rune(structName[0])) {
				return true
			}
			structNames = append(structNames, structName)
		}
		return true
	})

	structs := make([]analyzedStruct, len(structNames))
	for i, structName := range structNames {
		structs[i] = analyzedStruct{
			ImportPath: importPath,
			Package:    packageName,
			Name:       structName,
			NoUseQual:  g.noUserQualMap[importPath],
			file:       parseFile,
		}
	}
	analyzedStructs = append(analyzedStructs, structs...)
	analyzedMap[parseFile] = true

	g.parseFile2ImportMap[parseFile] = importMap
	g.parseFile2DotImportMap[parseFile] = dotStructs
	return nil
}

func (g *generator) setFieldToStruct() {
	for i, analyzedStruct := range analyzedStructs {

		importMap := g.parseFile2ImportMap[analyzedStruct.file]
		dotStructs := g.parseFile2DotImportMap[analyzedStruct.file]

		sameHierarchyStructs := map[string]bool{}
		for _, aast := range analyzedStructs {
			if analyzedStruct.ImportPath == aast.ImportPath {
				sameHierarchyStructs[aast.Name] = true
			}
		}

		analyzedFieldMap := map[string]*Node{}
		ast.Inspect(analyzedStruct.file, func(n ast.Node) bool {

			x, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			if st, ok := x.Type.(*ast.StructType); ok {
				if x.Name.String() != analyzedStruct.Name {
					return true
				}

				canGen := true
				reasons := make([]string, 0)
				for i, field := range st.Fields.List {

					key := fmt.Sprint(i)

					value, ok, rs := g.createNodeRecursive(field.Type, nil, importMap, dotStructs, sameHierarchyStructs)
					canGen = canGen && ok
					if ok {
						analyzedFieldMap[key+"@"+x.Name.String()] = value
					}
					reasons = append(reasons, rs...)
				}

				if canGen {
					analyzedStructs[i].CanGen = true
					analyzedStructs[i].Fields = g.createAnalyzedFields(analyzedStruct.Package, analyzedStruct.Name, analyzedFieldMap, g.fileSet, analyzedStruct.file)
				} else {
					analyzedStructs[i].CanGen = false
					analyzedStructs[i].Reasons = reasons
				}
			}
			return true
		})

	}
}

func (g *generator) createImportMap(parseFile *ast.File) (map[string]string, []string) {

	importMap := map[string]string{}
	dotImports := make([]string, 0)

	for _, imp := range parseFile.Imports {

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

func (g *generator) createAnalyzedFields(packageName, structName string, analyzedFieldMap map[string]*Node, fset *token.FileSet, file *ast.File) []Field {

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

	analyzedFields := make([]Field, 0)
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

			analyzedFields = append(analyzedFields, Field{
				Name: name,
				Tag:  tag,
				Node: analyzedFieldMap[fmt.Sprint(i)+"@"+structName],
			})
		}
	}

	// todo : msgpackresolverとして出力
	return analyzedFields
}
