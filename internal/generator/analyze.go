package generator

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"math/big"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"

	"github.com/shamaton/msgpackgen/internal/generator/structure"
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
	g.setFieldToStructs()
	return nil
}

func (g *generator) createAnalyzedStructs(parseFile *ast.File, packageName, importPath string, analyzedMap map[*ast.File]bool) error {

	importMap, dotImports := g.createImportMap(parseFile)
	// dot imports
	dotStructs := map[string]*structure.Structure{}
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

	structs := make([]*structure.Structure, len(structNames))
	for i, structName := range structNames {
		structs[i] = &structure.Structure{
			ImportPath: importPath,
			Package:    packageName,
			Name:       structName,
			NoUseQual:  g.noUserQualMap[importPath],
			File:       parseFile,
		}
	}
	analyzedStructs = append(analyzedStructs, structs...)
	analyzedMap[parseFile] = true

	g.parseFile2ImportMap[parseFile] = importMap
	g.parseFile2DotImportMap[parseFile] = dotStructs
	return nil
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

func (g *generator) setFieldToStructs() {
	for _, analyzedStruct := range analyzedStructs {

		importMap := g.parseFile2ImportMap[analyzedStruct.File]
		dotStructs := g.parseFile2DotImportMap[analyzedStruct.File]

		sameHierarchyStructs := map[string]bool{}
		for _, aast := range analyzedStructs {
			if analyzedStruct.ImportPath == aast.ImportPath {
				sameHierarchyStructs[aast.Name] = true
			}
		}

		g.setFieldToStruct(analyzedStruct, importMap, dotStructs, sameHierarchyStructs)
	}
}

func (g *generator) setFieldToStruct(target *structure.Structure,
	importMap map[string]string, dotStructs map[string]*structure.Structure, sameHierarchyStructs map[string]bool,
) {

	analyzedFieldMap := map[string]*structure.Node{}
	ast.Inspect(target.File, func(n ast.Node) bool {

		x, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if st, ok := x.Type.(*ast.StructType); ok {
			if x.Name.String() != target.Name {
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
				target.CanGen = true
				target.Fields = g.createAnalyzedFields(target.Package, target.Name, analyzedFieldMap, g.fileSet, target.File)
			} else {
				target.CanGen = false
				target.Reasons = reasons
			}
		}
		return true
	})

}
func (g *generator) createNodeRecursive(expr ast.Expr, parent *structure.Node, importMap map[string]string, dotStructs map[string]*structure.Structure, sameHierarchyStructs map[string]bool) (*structure.Node, bool, []string) {

	reasons := make([]string, 0)
	if ident, ok := expr.(*ast.Ident); ok {
		// dot import
		if dot, found := dotStructs[ident.Name]; found {
			return structure.CreateStructNode(dot.ImportPath, dot.Name, ident.Name, parent), true, reasons
		}
		// time
		if ident.Name == "Time" {
			return structure.CreateStructNode("time", "time", ident.Name, parent), true, reasons
		}
		// same hierarchy struct in same File
		if ident.Obj != nil && ident.Obj.Kind == ast.Typ {
			return structure.CreateStructNode(g.outputPackageFullName(), g.outputPackageName, ident.Name, parent), true, reasons
		}

		// same hierarchy struct in other File
		if _, found := sameHierarchyStructs[ident.Name]; found {
			return structure.CreateStructNode(g.outputPackageFullName(), g.outputPackageName, ident.Name, parent), true, reasons
		}

		if structure.IsPrimitive(ident.Name) {
			return structure.CreateIdentNode(ident, parent), true, reasons
		}
		return nil, false, []string{fmt.Sprintf("identifier %s is not suppoted or unknown struct ", ident.Name)}
	}

	if selector, ok := expr.(*ast.SelectorExpr); ok {
		pkgName := fmt.Sprint(selector.X)
		return structure.CreateStructNode(importMap[pkgName], pkgName, selector.Sel.Name, parent), true, reasons
	}

	// slice or array
	if array, ok := expr.(*ast.ArrayType); ok {
		var node *structure.Node
		if array.Len == nil {
			node = structure.CreateSliceNode(parent)
		} else {
			lit := array.Len.(*ast.BasicLit)
			// parse num
			n := new(big.Int)
			if litValue := strings.ToLower(lit.Value); strings.HasPrefix(litValue, "0b") {
				n.SetString(strings.ReplaceAll(litValue, "0b", ""), 2)
			} else if strings.HasPrefix(litValue, "0o") {
				n.SetString(strings.ReplaceAll(litValue, "0o", ""), 8)
			} else if strings.HasPrefix(litValue, "0x") {
				n.SetString(strings.ReplaceAll(litValue, "0x", ""), 16)
			} else {
				n.SetString(litValue, 10)
			}
			node = structure.CreateArrayNode(n.Uint64(), parent)
		}
		key, check, rs := g.createNodeRecursive(array.Elt, node, importMap, dotStructs, sameHierarchyStructs)
		node.SetKeyNode(key)
		reasons = append(reasons, rs...)
		return node, check, reasons
	}

	// map
	if mp, ok := expr.(*ast.MapType); ok {
		node := structure.CreateMapNode(parent)
		key, c1, krs := g.createNodeRecursive(mp.Key, node, importMap, dotStructs, sameHierarchyStructs)
		value, c2, vrs := g.createNodeRecursive(mp.Value, node, importMap, dotStructs, sameHierarchyStructs)
		node.SetKeyNode(key)
		node.SetValueNode(value)
		reasons = append(reasons, krs...)
		reasons = append(reasons, vrs...)
		return node, c1 && c2, reasons
	}

	// *
	if star, ok := expr.(*ast.StarExpr); ok {
		node := structure.CreatePointerNode(parent)
		key, check, rs := g.createNodeRecursive(star.X, node, importMap, dotStructs, sameHierarchyStructs)
		node.SetKeyNode(key)
		reasons = append(reasons, rs...)
		return node, check, reasons
	}

	// not supported
	if _, ok := expr.(*ast.InterfaceType); ok {
		return nil, false, []string{fmt.Sprintf("interface type is not supported")}
	}
	if _, ok := expr.(*ast.StructType); ok {
		return nil, false, []string{fmt.Sprintf("inner struct is not supported")}
	}
	if _, ok := expr.(*ast.ChanType); ok {
		return nil, false, []string{fmt.Sprintf("chan type is not supported")}
	}
	if _, ok := expr.(*ast.FuncType); ok {
		return nil, false, []string{fmt.Sprintf("func type is not supported")}
	}

	// unreachable
	return nil, false, []string{fmt.Sprintf("this field is unknown field")}
}

func (g *generator) createAnalyzedFields(packageName, structName string, analyzedFieldMap map[string]*structure.Node, fset *token.FileSet, file *ast.File) []structure.Field {

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

	analyzedFields := make([]structure.Field, 0)
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

			analyzedFields = append(analyzedFields, structure.Field{
				Name: name,
				Tag:  tag,
				Node: analyzedFieldMap[fmt.Sprint(i)+"@"+structName],
			})
		}
	}

	// todo : msgpackresolverとして出力
	return analyzedFields
}
