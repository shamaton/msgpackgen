package generator

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/shamaton/msgpackgen/internal/generator/ptn"
	"github.com/shamaton/msgpackgen/internal/generator/structure"
)

var (
	analyzedStructs []*structure.Structure
	structsInBrace  []string
)

type generator struct {
	fileSet               *token.FileSet
	targetPackages        map[string]bool
	parseFiles            []*ast.File
	importPath2ParseFiles map[string][]*ast.File
	parseFile2ImportPath  map[*ast.File]string
	importPath2package    map[string]string
	noUserQualMap         map[string]bool

	parseFile2ImportMap    map[*ast.File]map[string]string
	parseFile2DotImportMap map[*ast.File]map[string]*structure.Structure

	outputDir         string
	outputImportPath  string
	outputJenFilePath string

	goModFilePath   string
	goModModuleName string

	pointer   int
	useGopath bool
	verbose   bool
	strict    bool
}

// Run runs code analyzing and generation.
func Run(inputDir, inputFile, outDir, fileName string, pointer int, useGopath, dryRun, strict, verbose bool, w io.Writer) error {

	// can not input at same time
	if len(inputFile) > 0 && inputDir != "." {
		return fmt.Errorf("can not input directory and file at same time")
	}

	input := inputDir
	isInputDir := true
	if len(inputFile) < 1 {
		fi, err := os.Stat(inputDir)
		if err != nil {
			return fmt.Errorf("input directory error. os.Stat says %v", err)
		}
		if !fi.IsDir() {
			return fmt.Errorf("this(%s) path is not directory", inputDir)
		}
	} else {
		fi, err := os.Stat(inputFile)
		if err != nil {
			return fmt.Errorf("input file error. os.Stat says %v", err)
		}
		if fi.IsDir() {
			return fmt.Errorf("this(%s) is a directory", inputFile)
		}
		if !strings.HasSuffix(inputFile, ".go") {
			return fmt.Errorf("this(%s) is not .go file", inputFile)
		}
		input = inputFile
		isInputDir = false
	}

	if outDir == "" {
		outDir = inputDir
	}

	if pointer < 0 {
		pointer = 1
	}

	analyzedStructs = make([]*structure.Structure, 0)
	structsInBrace = make([]string, 0)
	g := generator{
		useGopath:             useGopath,
		pointer:               pointer,
		strict:                strict,
		verbose:               verbose,
		targetPackages:        map[string]bool{},
		parseFiles:            []*ast.File{},
		importPath2package:    map[string]string{},
		importPath2ParseFiles: map[string][]*ast.File{},
		parseFile2ImportPath:  map[*ast.File]string{},
		noUserQualMap:         map[string]bool{},

		parseFile2ImportMap:    map[*ast.File]map[string]string{},
		parseFile2DotImportMap: map[*ast.File]map[string]*structure.Structure{},
	}
	return g.run(input, outDir, fileName, isInputDir, dryRun, w)
}

func (g *generator) run(input, out, fileName string, isInputDir, dryRun bool, w io.Writer) error {
	g.fileSet = token.NewFileSet()

	if !g.useGopath {
		modFilePath, err := g.searchGoModFile(input, isInputDir)
		if err != nil {
			return err
		}
		g.goModFilePath = modFilePath

		err = g.setModuleName()
		if err != nil {
			return err
		}
	}

	err := g.setOutputInfo(out)
	if err != nil {
		return err
	}

	var filePaths []string
	if isInputDir {
		filePaths, err = g.getTargetFiles(input, true)
		if err != nil {
			return err
		}
		if len(filePaths) < 1 {
			return fmt.Errorf("not found go File")
		}
	} else {
		filePaths, err = g.getAbsolutePaths([]string{input})
		if err != nil {
			return err
		}
	}

	err = g.getPackages(filePaths)
	if err != nil {
		return err
	}

	err = g.analyze()
	if err != nil {
		return err
	}

	var reasons []string
	analyzedStructs, reasons = g.filter(analyzedStructs, reasons)

	g.printAnalyzedResult(reasons)

	g.setOthers()
	f := g.generateCode()

	if dryRun {
		fmt.Fprintf(w, "%#v", f)
		return nil
	}
	err = g.output(f, fileName)
	if err != nil {
		return err
	}
	return nil
}

func (g *generator) searchGoModFile(input string, isInputDir bool) (string, error) {
	goModFilePath := ""

	dir := input
	if !isInputDir {
		dir = filepath.Dir(input)
	}

	path, err := filepath.Abs(dir)
	if err != nil {
		return goModFilePath, err
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return goModFilePath, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if file.Name() == "go.mod" {
			goModFilePath = filepath.Join(path, file.Name())
		}
	}

	// recursive upper
	if goModFilePath == "" {
		// upper path
		sep := string(filepath.Separator)
		upper := filepath.Join(path, fmt.Sprintf("%s..%s", sep, sep))

		// reached root
		if path == upper {
			return goModFilePath, fmt.Errorf("not found go.mod")
		}
		return g.searchGoModFile(upper, true)
	}
	return goModFilePath, nil
}

func (g *generator) setModuleName() error {

	file, err := os.Open(g.goModFilePath)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println("File close error", err)
		}
	}()

	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		text := scanner.Text()
		if !strings.HasPrefix(text, "module") {
			return fmt.Errorf("not found module name in go.mod")
		}

		results := strings.Split(text, " ")
		if len(results) != 2 {
			return fmt.Errorf("something wrong in go.mod \n %s", text)
		}
		g.goModModuleName = results[1]
	}

	if err = scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (g *generator) setOutputInfo(out string) error {

	outAbs, err := filepath.Abs(out)
	if err != nil {
		return err
	}
	g.outputDir = outAbs

	importPath, err := g.getImportPath(g.outputDir)
	if err != nil {
		return err
	}
	g.outputImportPath = importPath
	g.outputJenFilePath = fmt.Sprintf("%s/%s", filepath.Dir(importPath), filepath.Base(importPath))

	// if exist go file
	fi, err := os.Stat(outAbs)
	if err != nil {
		// end proc.
		return nil
	}
	if !fi.IsDir() {
		return fmt.Errorf("this(%s) path is not directory", out)
	}

	files, err := g.getTargetFiles(outAbs, false)
	if err != nil {
		return err
	}
	if len(files) < 1 {
		return nil
	}
	_, packageName, _, err := g.getImportPathAndParseFile(files[0])
	if err != nil {
		return err
	}
	g.outputJenFilePath = fmt.Sprintf("%s/%s", filepath.Dir(importPath), packageName)
	return nil
}

func (g *generator) getImportPath(path string) (string, error) {
	if !g.useGopath {
		rep := strings.Replace(path, filepath.Dir(g.goModFilePath), g.goModModuleName, 1)
		return filepath.ToSlash(rep), nil
	}

	// use GOPATH option
	goPathAll := os.Getenv("GOPATH")
	sep := ":"
	if runtime.GOOS == "windows" {
		sep = ";"
	}
	goPaths := strings.Split(goPathAll, sep)

	p := filepath.ToSlash(path)
	for _, goPath := range goPaths {
		gp := filepath.ToSlash(goPath) + "/src/"
		if !strings.HasPrefix(p, gp) {
			continue
		}
		paths := strings.SplitN(p, gp, 2)
		return paths[1], nil
	}
	return "", fmt.Errorf("path %s is outside gopath", path)
}

func (g *generator) getTargetFiles(dir string, recursive bool) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() && recursive {
			if n := file.Name(); strings.HasPrefix(n, ".") ||
				strings.HasPrefix(n, "_") ||
				n == "testdata" ||
				n == "vendor" {
				if g.verbose {
					fmt.Printf("%s is not covered directory. skipping. \n", n)
				}
				continue
			}

			path, err := g.getTargetFiles(filepath.Join(dir, file.Name()), recursive)
			if err != nil {
				return nil, err
			}
			paths = append(paths, path...)
			continue
		}
		if filepath.Ext(file.Name()) == ".go" && !strings.HasSuffix(file.Name(), "_test.go") {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}

	absPaths, err := g.getAbsolutePaths(paths)
	if err != nil {
		return nil, err
	}
	return absPaths, nil
}

func (g *generator) getAbsolutePaths(paths []string) ([]string, error) {
	absPaths := make([]string, len(paths))
	for i, path := range paths {
		abs, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		absPaths[i] = abs
	}
	return absPaths, nil
}

func (g *generator) filter(structures []*structure.Structure, reasons []string) ([]*structure.Structure, []string) {
	newStructs := make([]*structure.Structure, 0)
	allOk := true
	for _, v := range structures {
		ok := true

		var rs []string
		if v.CanGen {
			for _, field := range v.Fields {
				if canGen, msgs := field.Node.CanGenerate(structures); !canGen {
					ok = false
					rs = msgs
				}
			}
			if ok {
				newStructs = append(newStructs, v)
			}
		} else {
			ok = false
			rs = v.Reasons
		}

		if !ok {
			reasons = append(reasons, fmt.Sprintf("notgen:%s.%s", v.ImportPath, v.Name))
			reasons = append(reasons, "reason")
			for i, s := range rs {
				mark := " |-  "
				if i == len(rs)-1 {
					mark = " `-  "
				}
				reasons = append(reasons, mark+s)
			}
		}

		allOk = allOk && ok
	}
	if !allOk {
		return g.filter(newStructs, reasons)
	}
	return newStructs, reasons
}

func (g *generator) setOthers() {
	for i := range analyzedStructs {
		analyzedStructs[i].Others = analyzedStructs
	}
}

func (g *generator) generateCode() *File {

	f := NewFilePath(g.outputJenFilePath)

	registerName := "RegisterGeneratedResolver"
	f.HeaderComment("// Code generated by msgpackgen. DO NOT EDIT.")
	f.Comment(fmt.Sprintf("// %s registers generated resolver.\n", registerName)).
		Func().Id(registerName).Params().Block(
		Qual(ptn.PkTop, "SetResolver").Call(
			Id(ptn.PrivateFuncName("encodeAsMap")),
			Id(ptn.PrivateFuncName("encodeAsArray")),
			Id(ptn.PrivateFuncName("decodeAsMap")),
			Id(ptn.PrivateFuncName("decodeAsArray")),
		),
	)

	encReturn := Return(Nil(), Nil())
	decReturn := Return(False(), Nil())
	if g.strict {
		encReturn = Return(Nil(), Qual("fmt", "Errorf").Call(Lit("use strict option : undefined type")))
		decReturn = Return(False(), Qual("fmt", "Errorf").Call(Lit("use strict option : undefined type")))
	}

	encodeAsArrayCode := []Code{encReturn}
	encodeAsMapCode := []Code{encReturn}
	decodeAsArrayCode := []Code{decReturn}
	decodeAsMapCode := []Code{decReturn}
	if len(analyzedStructs) > 0 {
		encodeAsArrayCode = append([]Code{
			Switch(Id("v").Op(":=").Id("i").Assert(Type())).Block(
				g.encodeAsArrayCases()...,
			)},
			encodeAsArrayCode...,
		)
		encodeAsMapCode = append([]Code{
			Switch(Id("v").Op(":=").Id("i").Assert(Type())).Block(
				g.encodeAsMapCases()...,
			)},
			encodeAsMapCode...,
		)
		decodeAsArrayCode = append([]Code{
			Switch(Id("v").Op(":=").Id("i").Assert(Type())).Block(
				g.decodeAsArrayCases()...,
			)},
			decodeAsArrayCode...,
		)
		decodeAsMapCode = append([]Code{
			Switch(Id("v").Op(":=").Id("i").Assert(Type())).Block(
				g.decodeAsMapCases()...,
			)},
			decodeAsMapCode...,
		)
	}

	g.encodeTopTemplate("encode", f).Block(
		If(Qual(ptn.PkTop, "StructAsArray").Call()).Block(
			Return(Id(ptn.PrivateFuncName("encodeAsArray")).Call(Id("i"))),
		).Else().Block(
			Return(Id(ptn.PrivateFuncName("encodeAsMap")).Call(Id("i"))),
		),
	)

	g.encodeTopTemplate("encodeAsArray", f).Block(encodeAsArrayCode...)
	g.encodeTopTemplate("encodeAsMap", f).Block(encodeAsMapCode...)

	g.decodeTopTemplate("decode", f).Block(
		If(Qual(ptn.PkTop, "StructAsArray").Call()).Block(
			Return(Id(ptn.PrivateFuncName("decodeAsArray")).Call(Id("data"), Id("i"))),
		).Else().Block(
			Return(Id(ptn.PrivateFuncName("decodeAsMap")).Call(Id("data"), Id("i"))),
		),
	)

	g.decodeTopTemplate("decodeAsArray", f).Block(decodeAsArrayCode...)
	g.decodeTopTemplate("decodeAsMap", f).Block(decodeAsMapCode...)

	for _, st := range analyzedStructs {
		st.CreateCode(f)
	}

	return f
}

func (g *generator) output(f *File, genFileName string) error {

	if err := os.MkdirAll(g.outputDir, 0777); err != nil {
		return err
	}

	path := g.outputDir + string(filepath.Separator) + genFileName
	file, err := os.Create(filepath.FromSlash(path))
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println("File close error", err)
		}
	}()

	_, err = fmt.Fprintf(file, "%#v", f)
	if err != nil {
		return err
	}

	p := genFileName
	if g.verbose {
		p = path
	}
	fmt.Println(p, "generated.")
	return err
}

func (g *generator) printAnalyzedResult(reasons []string) {
	if !g.verbose {
		return
	}

	fmt.Printf("=========== %d generated ==========\n", len(analyzedStructs))
	for _, v := range analyzedStructs {
		fmt.Println(v.ImportPath, v.Name)
	}

	notGen := 0
	for _, s := range reasons {
		if strings.Contains(s, "notgen:") {
			notGen++
		}
	}
	fmt.Printf("=========== %d not generated ==========\n", notGen+len(structsInBrace))
	for _, s := range structsInBrace {
		fmt.Println(s)
		fmt.Println("reason")
		fmt.Println(" └  defined in function")
	}
	for _, s := range reasons {
		if strings.Contains(s, "notgen:") {
			fmt.Println(strings.ReplaceAll(s, "notgen:", ""))
		} else {
			fmt.Println(s)
		}
	}
	fmt.Println("=========================================")
}

func (g *generator) decodeTopTemplate(name string, f *File) *Statement {
	return f.Comment(fmt.Sprintf("// %s\n", name)).
		Func().Id(ptn.PrivateFuncName(name)).Params(Id("data").Index().Byte(), Id("i").Interface()).Params(Bool(), Error())
}

func (g *generator) encodeTopTemplate(name string, f *File) *Statement {
	return f.Comment(fmt.Sprintf("// %s\n", name)).
		Func().Id(ptn.PrivateFuncName(name)).Params(Id("i").Interface()).Params(Index().Byte(), Error())
}

func (g *generator) encodeAsArrayCases() []Code {
	var states, pointers []Code
	for _, v := range analyzedStructs {
		s, p := g.encodeCaseCode(v, true)
		states = append(states, s...)
		pointers = append(pointers, p...)
	}
	return append(states, pointers...)
}

func (g *generator) encodeAsMapCases() []Code {
	var states, pointers []Code
	for _, v := range analyzedStructs {
		s, p := g.encodeCaseCode(v, false)
		states = append(states, s...)
		pointers = append(pointers, p...)
	}
	return append(states, pointers...)
}

func (g *generator) encodeCaseCode(v *structure.Structure, asArray bool) (states []Code, pointers []Code) {

	var caseStatement func(string) *Statement
	var errID *Statement
	if v.NoUseQual {
		caseStatement = func(op string) *Statement { return Op(op).Id(v.Name) }
		errID = Lit(v.Name)
	} else {
		caseStatement = func(op string) *Statement { return Op(op).Qual(v.ImportPath, v.Name) }
		errID = Lit(v.ImportPath + "." + v.Name)
	}

	var calcFuncName, encodeFuncName, pointerFuncName string
	if asArray {
		calcFuncName = v.CalcArraySizeFuncName()
		encodeFuncName = v.EncodeArrayFuncName()
		pointerFuncName = "encodeAsArray"
	} else {
		calcFuncName = v.CalcMapSizeFuncName()
		encodeFuncName = v.EncodeMapFuncName()
		pointerFuncName = "encodeAsMap"
	}

	f := func(ptr string) *Statement {
		return Case(caseStatement(ptr)).Block(
			Id(ptn.IdEncoder).Op(":=").Qual(ptn.PkEnc, "NewEncoder").Call(),
			List(Id("size"), Err()).Op(":=").Id(calcFuncName).Call(Id(ptr+"v"), Id(ptn.IdEncoder)),
			If(Err().Op("!=").Nil()).Block(
				Return(Nil(), Err()),
			),
			Id(ptn.IdEncoder).Dot("MakeBytes").Call(Id("size")),
			List(Id("b"), Id("offset"), Err()).Op(":=").Id(encodeFuncName).Call(Id(ptr+"v"), Id(ptn.IdEncoder), Lit(0)),
			If(Err().Op("!=").Nil()).Block(
				Return(Nil(), Err()),
			),
			If(Id("size").Op("!=").Id("offset")).Block(
				Return(Nil(), Qual("fmt", "Errorf").Call(Lit("%s size / offset different %d : %d"), errID, Id("size"), Id("offset"))),
			),
			Return(Id("b"), Err()),
		)
	}

	states = append(states, f(""))

	if g.pointer > 0 {
		states = append(states, f("*"))
	}

	for i := 0; i < g.pointer-1; i++ {
		ptr := strings.Repeat("*", i+2)
		pointers = append(pointers, Case(caseStatement(ptr)).Block(
			Return(Id(ptn.PrivateFuncName(pointerFuncName)).Call(Id("*v"))),
		))
	}
	return
}

func (g *generator) decodeAsArrayCases() []Code {
	var states, pointers []Code
	for _, v := range analyzedStructs {
		s, p := g.decodeCaseCode(v, true)
		states = append(states, s...)
		pointers = append(pointers, p...)
	}
	return append(states, pointers...)
}

func (g *generator) decodeAsMapCases() []Code {
	var states, pointers []Code
	for _, v := range analyzedStructs {
		s, p := g.decodeCaseCode(v, false)
		states = append(states, s...)
		pointers = append(pointers, p...)
	}
	return append(states, pointers...)
}

func (g *generator) decodeCaseCode(v *structure.Structure, asArray bool) (states []Code, pointers []Code) {

	var caseStatement func(string) *Statement
	if v.NoUseQual {
		caseStatement = func(op string) *Statement { return Op(op).Id(v.Name) }
	} else {
		caseStatement = func(op string) *Statement { return Op(op).Qual(v.ImportPath, v.Name) }
	}

	var decodeFuncName, pointerFuncName string
	if asArray {
		decodeFuncName = v.DecodeArrayFuncName()
		pointerFuncName = "decodeAsArray"
	} else {

		decodeFuncName = v.DecodeMapFuncName()
		pointerFuncName = "decodeAsMap"
	}

	states = append(states, Case(caseStatement("*")).Block(
		Id(ptn.IdDecoder).Op(":=").Qual(ptn.PkDec, "NewDecoder").Call(Id("data")),
		List(Id("offset"), Err()).Op(":=").Id(decodeFuncName).Call(Id("v"), Id(ptn.IdDecoder), Id("0")),
		If(Err().Op("==").Nil().Op("&&").Id("offset").Op("!=").Id(ptn.IdDecoder).Dot("Len").Call()).Block(
			Return(True(), Qual("fmt", "Errorf").Call(Lit("read length is different [%d] [%d] "), Id("offset"), Id(ptn.IdDecoder).Dot("Len").Call())),
		),
		Return(True(), Err())))

	if g.pointer > 0 {
		states = append(states, Case(caseStatement("**")).Block(
			Id(ptn.IdDecoder).Op(":=").Qual(ptn.PkDec, "NewDecoder").Call(Id("data")),
			List(Id("offset"), Err()).Op(":=").Id(decodeFuncName).Call(Id("*v"), Id(ptn.IdDecoder), Id("0")),
			If(Err().Op("==").Nil().Op("&&").Id("offset").Op("!=").Id(ptn.IdDecoder).Dot("Len").Call()).Block(
				Return(True(), Qual("fmt", "Errorf").Call(Lit("read length is different [%d] [%d] "), Id("offset"), Id(ptn.IdDecoder).Dot("Len").Call())),
			),
			Return(True(), Err())))
	}

	for i := 0; i < g.pointer-1; i++ {
		ptr := strings.Repeat("*", i+3)
		pointers = append(pointers, Case(caseStatement(ptr)).Block(
			Return(Id(ptn.PrivateFuncName(pointerFuncName)).Call(Id("data"), Id("*v"))),
		))
	}
	return
}
