package gen

import (
	"crypto/sha256"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/dave/jennifer/jen"
)

var analyzedStructs []analyzedStruct

const (
	pkTop = "github.com/shamaton/msgpackgen/msgpack"
	pkEnc = pkTop + "/enc"
	pkDec = pkTop + "/dec"

	idEncoder = "encoder"
	idDecoder = "decoder"

	outputFileName = "resolver.msgpackgen.go"
)

// todo : tagをmapのcaseに使いつつ、変数に代入するようにしないといけない

var funcIdMap = map[string]string{}

type Generator struct {
	fileSet              *token.FileSet
	targetPackages       map[string]bool
	parseFiles           []*ast.File
	fileNames            []string
	file2FullPackageName map[string]string
	file2PackageName     map[string]string
	noUserQualMap        map[string]bool

	outputDir           string
	outputPackageName   string
	outputPackagePrefix string

	pointer int
	verbose bool
	strict  bool
}

func (g *Generator) OutputPackageFullName() string {
	return fmt.Sprintf("%s/%s", g.outputPackagePrefix, g.outputPackageName)
}

type analyzedStruct struct {
	PackageName string
	Name        string
	Fields      []analyzedField
	NoUseQual   bool
}

type analyzedField struct {
	Name string
	Tag  string
	Type types.Type
	Ast  *analyzedASTFieldType
}

func NewGenerator() *Generator {
	return &Generator{
		targetPackages:       map[string]bool{},
		parseFiles:           []*ast.File{},
		fileNames:            []string{},
		file2FullPackageName: map[string]string{},
		file2PackageName:     map[string]string{},
		noUserQualMap:        map[string]bool{},
	}
}

func (g *Generator) Initialize(input, out string, pointer int, strict, verbose bool) {

	g.pointer = pointer
	g.strict = strict
	g.verbose = verbose

	outAbs, err := filepath.Abs(out)
	if err != nil {
		fmt.Println(err)
	}

	g.outputDir = outAbs
	paths := strings.SplitN(g.outputDir, "src/", 2)
	if len(paths) != 2 {
		fmt.Printf("%s get import path failed", out)
		return
	}
	fmt.Println("paths", paths)
	g.outputPackageName = paths[1]

	// todo : ファイル指定オプション

	fmt.Println(g.outputDir, outAbs)
	fmt.Println("output pacagek is 222... ", g.outputPackagePrefix, g.outputPackageName)
}

func (g *Generator) Dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, g.Dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		if filepath.Ext(file.Name()) == ".go" && !strings.HasSuffix(file.Name(), "_test.go") {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}

	var abss []string
	for _, path := range paths {
		abs, err := filepath.Abs(path)
		if err != nil {
			panic(err)
		}
		abss = append(abss, abs)
	}
	return abss
}
func privateFuncNamePattern(funcName string) string {
	return fmt.Sprintf("___%s", funcName)
}

func (g *Generator) Generate() {

	for _, st := range analyzedStructs {
		funcIdMap[st.PackageName] = fmt.Sprintf("%x", sha256.Sum256([]byte(st.PackageName)))
	}

	fmt.Println("output pacagek is ... ", g.outputPackagePrefix, g.outputPackageName)

	// todo : ソースコードが存在している場所だったら、そちらにパッケージ名をあわせる
	f := NewFilePath(g.OutputPackageFullName())

	registerName := "RegisterGeneratedResolver"
	f.HeaderComment("// Code generated by msgpackgen. DO NOT EDIT.\n// Thank you for using and generating.")
	f.Comment(fmt.Sprintf("// %s registers generated resolver.\n", registerName)).
		Func().Id(registerName).Params().Block(
		Qual(pkTop, "SetResolver").Call(
			Id(privateFuncNamePattern("encodeAsMap")),
			Id(privateFuncNamePattern("encodeAsArray")),
			Id(privateFuncNamePattern("decodeAsMap")),
			Id(privateFuncNamePattern("decodeAsArray")),
		),
	)

	g.decodeTopTemplate("decode", f).Block(
		If(Qual(pkTop, "StructAsArray").Call()).Block(
			Return(Id(privateFuncNamePattern("decodeAsArray")).Call(Id("data"), Id("i"))),
		).Else().Block(
			Return(Id(privateFuncNamePattern("decodeAsMap")).Call(Id("data"), Id("i"))),
		),
	)

	encReturn := Return(Nil(), Nil())
	decReturn := Return(False(), Nil())
	if g.strict {
		encReturn = Return(Nil(), Qual("fmt", "Errorf").Call(Lit("use strict option : undefined type")))
		decReturn = Return(False(), Qual("fmt", "Errorf").Call(Lit("use strict option : undefined type")))
	}

	g.decodeTopTemplate("decodeAsArray", f).Block(
		Switch(Id("v").Op(":=").Id("i").Assert(Type())).Block(
			g.decodeAsArrayCases()...,
		),
		decReturn,
	)

	g.decodeTopTemplate("decodeAsMap", f).Block(
		Switch(Id("v").Op(":=").Id("i").Assert(Type())).Block(
			g.decodeAsMapCases()...,
		),
		decReturn,
	)

	g.encodeTopTemplate("encode", f).Block(
		If(Qual(pkTop, "StructAsArray").Call()).Block(
			Return(Id(privateFuncNamePattern("encodeAsArray")).Call(Id("i"))),
		).Else().Block(
			Return(Id(privateFuncNamePattern("encodeAsMap")).Call(Id("i"))),
		),
	)

	g.encodeTopTemplate("encodeAsArray", f).Block(
		Switch(Id("v").Op(":=").Id("i").Assert(Type())).Block(
			g.encodeAsArrayCases()...,
		),
		encReturn,
	)

	g.encodeTopTemplate("encodeAsMap", f).Block(
		Switch(Id("v").Op(":=").Id("i").Assert(Type())).Block(
			g.encodeAsMapCases()...,
		),
		encReturn,
	)

	// todo : 名称修正
	for _, st := range analyzedStructs {
		st.calcFunction(f)
	}

	if err := os.MkdirAll(g.outputDir, 0777); err != nil {
		fmt.Println(err)
	}

	fileName := g.outputDir + "/" + outputFileName
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%#v", f)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fileName, "generated.")

}

func (g *Generator) decodeTopTemplate(name string, f *File) *Statement {
	return f.Comment(fmt.Sprintf("// %s\n", name)).
		Func().Id(privateFuncNamePattern(name)).Params(Id("data").Index().Byte(), Id("i").Interface()).Params(Bool(), Error())
}

func (g *Generator) encodeTopTemplate(name string, f *File) *Statement {
	return f.Comment(fmt.Sprintf("// %s\n", name)).
		Func().Id(privateFuncNamePattern(name)).Params(Id("i").Interface()).Params(Index().Byte(), Error())
}

func (g *Generator) encodeAsArrayCases() []Code {
	var states []Code
	for _, v := range analyzedStructs {
		for _, vv := range []string{"", "*"} {
			var caseStatement *Statement
			var errID *Statement
			if v.NoUseQual {
				caseStatement = Op(vv).Id(v.Name)
				errID = Lit(v.Name)
			} else {
				caseStatement = Op(vv).Qual(v.PackageName, v.Name)
				errID = Lit(v.PackageName + "." + v.Name)
				//errID = Id("\"").Qual(v.PackageName, v.Name).Id("\"")
			}

			states = append(states, Case(caseStatement).Block(
				Id(idEncoder).Op(":=").Qual(pkEnc, "NewEncoder").Call(),
				List(Id("size"), Err()).Op(":=").Id(v.calcArraySizeFuncName()).Call(Id(vv+"v"), Id(idEncoder)),
				If(Err().Op("!=").Nil()).Block(
					Return(Nil(), Err()),
				),
				Id(idEncoder).Dot("MakeBytes").Call(Id("size")),
				List(Id("b"), Id("offset"), Err()).Op(":=").Id(v.encodeArrayFuncName()).Call(Id(vv+"v"), Id(idEncoder), Lit(0)),
				If(Err().Op("!=").Nil()).Block(
					Id(idEncoder).Dot("ReleaseBytes").Call(),
					Return(Nil(), Err()),
				),
				If(Id("size").Op("!=").Id("offset")).Block(
					Id(idEncoder).Dot("ReleaseBytes").Call(),
					Return(Nil(), Qual("fmt", "Errorf").Call(Lit("%s size / offset different %d : %d"), errID, Id("size"), Id("offset"))),
				),
				Id(idEncoder).Dot("ReleaseBytes").Call(),
				Return(Id("b"), Err()),
			))
		}
	}
	return states
}

func (g *Generator) encodeAsMapCases() []Code {
	var states []Code
	for _, v := range analyzedStructs {
		for _, vv := range []string{"", "*"} {
			var caseStatement *Statement
			var errID *Statement
			if v.NoUseQual {
				caseStatement = Op(vv).Id(v.Name)
				errID = Lit(v.Name)
			} else {
				caseStatement = Op(vv).Qual(v.PackageName, v.Name)
				errID = Lit(v.PackageName + "." + v.Name)
				//errID = Id("\"").Qual(v.PackageName, v.Name).Id("\"")
			}

			states = append(states, Case(caseStatement).Block(
				Id(idEncoder).Op(":=").Qual(pkEnc, "NewEncoder").Call(),
				List(Id("size"), Err()).Op(":=").Id(v.calcMapSizeFuncName()).Call(Id(vv+"v"), Id(idEncoder)),
				If(Err().Op("!=").Nil()).Block(
					Return(Nil(), Err()),
				),
				Id(idEncoder).Dot("MakeBytes").Call(Id("size")),
				List(Id("b"), Id("offset"), Err()).Op(":=").Id(v.encodeMapFuncName()).Call(Id(vv+"v"), Id(idEncoder), Lit(0)),
				If(Err().Op("!=").Nil()).Block(
					Id(idEncoder).Dot("ReleaseBytes").Call(),
					Return(Nil(), Err()),
				),
				If(Id("size").Op("!=").Id("offset")).Block(
					Id(idEncoder).Dot("ReleaseBytes").Call(),
					Return(Nil(), Qual("fmt", "Errorf").Call(Lit("%s size / offset different %d : %d"), errID, Id("size"), Id("offset"))),
				),
				Id(idEncoder).Dot("ReleaseBytes").Call(),
				Return(Id("b"), Err()),
			))
		}
	}
	return states
}

func (g *Generator) decodeAsArrayCases() []Code {
	var states []Code
	for _, v := range analyzedStructs {
		for _, vv := range []string{"", "*"} {

			var caseStatement *Statement
			if v.NoUseQual {
				caseStatement = Op("*" + vv).Id(v.Name)
			} else {
				caseStatement = Op("*"+vv).Qual(v.PackageName, v.Name)
			}

			states = append(states, Case(caseStatement).Block(
				List(Id("_"), Err()).Op(":=").Id(v.decodeArrayFuncName()).Call(Id(vv+"v"), Qual(pkDec, "NewDecoder").Call(Id("data")), Id("0")),
				Return(True(), Err())))
		}
	}
	return states
}

func (g *Generator) decodeAsMapCases() []Code {
	var states []Code
	for _, v := range analyzedStructs {
		for _, vv := range []string{"", "*"} {
			var caseStatement *Statement
			if v.NoUseQual {
				caseStatement = Op("*" + vv).Id(v.Name)
			} else {
				caseStatement = Op("*"+vv).Qual(v.PackageName, v.Name)
			}

			states = append(states, Case(caseStatement).Block(
				List(Id("_"), Err()).Op(":=").Id(v.decodeMapFuncName()).Call(Id(vv+"v"), Qual(pkDec, "NewDecoder").Call(Id("data")), Id("0")),
				Return(True(), Err())))
		}
	}
	return states
}

func (as *analyzedStruct) calcArraySizeFuncName() string {
	return createFuncName("calcArraySize", as.Name, as.PackageName)
}

func (as *analyzedStruct) calcMapSizeFuncName() string {
	return createFuncName("calcMapSize", as.Name, as.PackageName)
}

func (as *analyzedStruct) encodeArrayFuncName() string {
	return createFuncName("encodeArray", as.Name, as.PackageName)
}

func (as *analyzedStruct) encodeMapFuncName() string {
	return createFuncName("encodeMap", as.Name, as.PackageName)
}

func (as *analyzedStruct) decodeArrayFuncName() string {
	return createFuncName("decodeArray", as.Name, as.PackageName)
}

func (as *analyzedStruct) decodeMapFuncName() string {
	return createFuncName("decodeMap", as.Name, as.PackageName)
}

func createFuncName(prefix, name, packageName string) string {
	return privateFuncNamePattern(fmt.Sprintf("%s%s_%s", prefix, name, funcIdMap[packageName]))
}
