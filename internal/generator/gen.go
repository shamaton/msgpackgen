package generator

import (
	"fmt"
	"go/ast"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/shamaton/msgpackgen/internal/generator/ptn"
	"github.com/shamaton/msgpackgen/internal/generator/structure"
)

var analyzedStructs []*structure.Structure

// todo : ドットインポートのコンフリクトチェック

// todo : complexのext値を変更できるようにする

// todo : full package name -> import path

type generator struct {
	fileSet                *token.FileSet
	targetPackages         map[string]bool
	parseFiles             []*ast.File
	fullPackage2ParseFiles map[string][]*ast.File
	parseFile2fullPackage  map[*ast.File]string
	fullPackage2package    map[string]string
	noUserQualMap          map[string]bool

	parseFile2ImportMap    map[*ast.File]map[string]string
	parseFile2DotImportMap map[*ast.File]map[string]*structure.Structure

	outputDir           string
	outputPackageName   string
	outputPackagePrefix string

	pointer int
	verbose bool
	strict  bool
}

func (g *generator) outputPackageFullName() string {
	return fmt.Sprintf("%s/%s", g.outputPackagePrefix, g.outputPackageName)
}

func Run(input, out, fileName string, pointer int, strict, verbose bool) error {

	_, err := os.Stat(input)
	if err != nil {
		return err
	}

	if out == "" {
		out = input
	}

	if pointer < 0 {
		pointer = 1
	}

	g := generator{
		pointer:                pointer,
		strict:                 strict,
		verbose:                verbose,
		targetPackages:         map[string]bool{},
		parseFiles:             []*ast.File{},
		fullPackage2package:    map[string]string{},
		fullPackage2ParseFiles: map[string][]*ast.File{},
		parseFile2fullPackage:  map[*ast.File]string{},
		noUserQualMap:          map[string]bool{},

		parseFile2ImportMap:    map[*ast.File]map[string]string{},
		parseFile2DotImportMap: map[*ast.File]map[string]*structure.Structure{},
	}
	return g.run(input, out, fileName)
}

func (g *generator) run(input, out, fileName string) error {

	outAbs, err := filepath.Abs(out)
	if err != nil {
		return err
	}

	g.outputDir = outAbs
	paths := strings.SplitN(filepath.ToSlash(g.outputDir), "src/", 2)
	if len(paths) != 2 {
		return fmt.Errorf("%s get import path failed", outAbs)
	}
	g.outputPackageName = paths[1]

	targets, err := g.getTargetFiles(input)
	if err != nil {
		return err
	}
	if len(targets) < 1 {
		return fmt.Errorf("not found go File")
	}

	err = g.getPackages(targets)
	if err != nil {
		return err
	}

	err = g.analyze()
	if err != nil {
		return err
	}

	var reasons []string
	analyzedStructs, reasons = g.filter(analyzedStructs, reasons)

	if g.verbose {
		fmt.Println("=========== generated ==========")
		for _, v := range analyzedStructs {
			fmt.Println(v.ImportPath, v.Name)
		}
		fmt.Println("=========== not generated ==========")
		for _, s := range reasons {
			fmt.Println(s)
		}
	}
	err = g.setOthers()
	if err != nil {
		return err
	}
	f := g.generateCode()

	err = g.output(f, fileName)
	if err != nil {
		return err
	}
	return nil
}

func (g *generator) getTargetFiles(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			path, err := g.getTargetFiles(filepath.Join(dir, file.Name()))
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

	var absPaths []string
	for _, path := range paths {
		abs, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		absPaths = append(absPaths, abs)
	}
	return absPaths, nil
}

func (g *generator) filter(sts []*structure.Structure, reasons []string) ([]*structure.Structure, []string) {
	newStructs := make([]*structure.Structure, 0)
	allOk := true
	for _, v := range sts {
		ok := true

		var rs []string
		if v.CanGen {
			for _, field := range v.Fields {
				if canGen, msgs := field.Node.CanGenerate(sts); !canGen {
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
			reasons = append(reasons, fmt.Sprintf("can not generate %s.%s", v.ImportPath, v.Name))
			reasons = append(reasons, "reason:")
			for _, s := range rs {
				reasons = append(reasons, "\t"+s)
			}
		}

		allOk = allOk && ok
	}
	if !allOk {
		return g.filter(newStructs, reasons)
	} else {
		return newStructs, reasons
	}
}

func (g *generator) setOthers() error {

	for i := range analyzedStructs {
		// todo : 整理
		analyzedStructs[i].Others = analyzedStructs
		continue

		others := make([]*structure.Structure, len(analyzedStructs)-1)
		index := 0
		for _, v := range analyzedStructs {
			if v.ImportPath != analyzedStructs[i].ImportPath || v.Name != analyzedStructs[i].Name {
				others[index] = v
				index++
			}
		}
		if index != len(others) {
			return fmt.Errorf("other package should be %d. but result is %d", len(others), index)
		}
		analyzedStructs[i].Others = others
	}
	return nil
}

func (g *generator) generateCode() *File {

	// todo : ソースコードが存在している場所だったら、そちらにパッケージ名をあわせる
	f := NewFilePath(g.outputPackageFullName())

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

	// todo : 名称修正
	for _, st := range analyzedStructs {
		st.CreateCode(f)
	}

	return f
}

func (g *generator) output(f *File, genFileName string) error {

	if err := os.MkdirAll(g.outputDir, 0777); err != nil {
		return err
	}

	fileName := g.outputDir + "/" + genFileName
	file, err := os.Create(fileName)
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

	if g.verbose {
		fmt.Println(fileName, "generated.")
	} else {
		fmt.Println(genFileName, "generated.")
	}
	return err
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
