package structure

import (
	"fmt"
	"go/ast"
	"math"

	. "github.com/dave/jennifer/jen"
	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpackgen/internal/generator/ptn"
)

// Structure has information needed for code generation
type Structure struct {
	ImportPath string
	Package    string
	Name       string
	Fields     []Field
	ZeroFields []Field
	NoUseQual  bool

	Others []*Structure
	File   *ast.File

	CanGen  bool
	Reasons []string
}

// Field has a field information in structure
type Field struct {
	Name      string
	Tag       string
	OmitEmpty bool
	Node      *Node
}

// CalcArraySizeFuncName gets the function name for each structure
func (st *Structure) CalcArraySizeFuncName() string {
	return st.createFuncName("calcArraySize")
}

// CalcArraySizeMaxFuncName gets the max-size function name for each structure.
func (st *Structure) CalcArraySizeMaxFuncName() string {
	return st.createFuncName("calcArraySizeMax")
}

// CalcArraySizeNoErrFuncName gets the no-error size function name for each structure.
func (st *Structure) CalcArraySizeNoErrFuncName() string {
	return st.createFuncName("calcArraySizeNoErr")
}

// CalcArraySizeMaxNoErrFuncName gets the no-error max-size function name for each structure.
func (st *Structure) CalcArraySizeMaxNoErrFuncName() string {
	return st.createFuncName("calcArraySizeMaxNoErr")
}

// CalcMapSizeFuncName gets the function name for each structure
func (st *Structure) CalcMapSizeFuncName() string {
	return st.createFuncName("calcMapSize")
}

// CalcMapSizeMaxFuncName gets the max-size function name for each structure.
func (st *Structure) CalcMapSizeMaxFuncName() string {
	return st.createFuncName("calcMapSizeMax")
}

// CalcMapSizeNoErrFuncName gets the no-error size function name for each structure.
func (st *Structure) CalcMapSizeNoErrFuncName() string {
	return st.createFuncName("calcMapSizeNoErr")
}

// CalcMapSizeMaxNoErrFuncName gets the no-error max-size function name for each structure.
func (st *Structure) CalcMapSizeMaxNoErrFuncName() string {
	return st.createFuncName("calcMapSizeMaxNoErr")
}

// EncodeArrayFuncName gets the function name for each structure
func (st *Structure) EncodeArrayFuncName() string {
	return st.createFuncName("encodeArray")
}

// EncodeMapFuncName gets the function name for each structure
func (st *Structure) EncodeMapFuncName() string {
	return st.createFuncName("encodeMap")
}

// DecodeArrayFuncName gets the function name for each structure
func (st *Structure) DecodeArrayFuncName() string {
	return st.createFuncName("decodeArray")
}

// DecodeMapFuncName gets the function name for each structure
func (st *Structure) DecodeMapFuncName() string {
	return st.createFuncName("decodeMap")
}

func (st *Structure) createFuncName(prefix string) string {
	return createFuncName(prefix, st.Name, st.ImportPath)
}

// CreateCode creates codes to serialize structure
func (st *Structure) CreateCode(f *File) {
	v := "v"

	calcStruct, encStructArray, encStructMap := st.createStructCode(len(st.Fields))
	hasOmitEmpty := st.hasOmitEmptyField()

	calcArraySizeCodes := make([]Code, 0)
	calcArraySizeCodes = append(calcArraySizeCodes, Id("size").Op(":=").Lit(0))
	calcArraySizeCodes = append(calcArraySizeCodes, calcStruct)

	calcArraySizeMaxCodes := make([]Code, 0)
	calcArraySizeMaxCodes = append(calcArraySizeMaxCodes, Id("size").Op(":=").Lit(0))
	calcArraySizeMaxCodes = append(calcArraySizeMaxCodes, calcStruct)

	calcMapSizeCodes := make([]Code, 0)
	calcMapSizeCodes = append(calcMapSizeCodes, Id("size").Op(":=").Lit(0))
	if !hasOmitEmpty {
		calcMapSizeCodes = append(calcMapSizeCodes, calcStruct)
	}

	calcMapSizeMaxCodes := make([]Code, 0)
	calcMapSizeMaxCodes = append(calcMapSizeMaxCodes, Id("size").Op(":=").Lit(0))
	if !hasOmitEmpty {
		calcMapSizeMaxCodes = append(calcMapSizeMaxCodes, calcStruct)
	}

	canCalcSizeNoErr := st.CanCalcSizeNoErr()
	calcArraySizeNoErrCodes := make([]Code, 0)
	calcArraySizeMaxNoErrCodes := make([]Code, 0)
	calcMapSizeNoErrCodes := make([]Code, 0)
	calcMapSizeMaxNoErrCodes := make([]Code, 0)
	if canCalcSizeNoErr {
		calcArraySizeNoErrCodes = append(calcArraySizeNoErrCodes, Id("size").Op(":=").Lit(0), calcStruct)
		calcArraySizeMaxNoErrCodes = append(calcArraySizeMaxNoErrCodes, Id("size").Op(":=").Lit(0), calcStruct)
		calcMapSizeNoErrCodes = append(calcMapSizeNoErrCodes, Id("size").Op(":=").Lit(0))
		calcMapSizeMaxNoErrCodes = append(calcMapSizeMaxNoErrCodes, Id("size").Op(":=").Lit(0))
		if !hasOmitEmpty {
			calcMapSizeNoErrCodes = append(calcMapSizeNoErrCodes, calcStruct)
			calcMapSizeMaxNoErrCodes = append(calcMapSizeMaxNoErrCodes, calcStruct)
		}
	}

	encArrayCodes := make([]Code, 0)
	encArrayCodes = append(encArrayCodes, Var().Err().Error())
	encArrayCodes = append(encArrayCodes, encStructArray)

	encMapCodes := make([]Code, 0)
	encMapCodes = append(encMapCodes, Var().Err().Error())
	if !hasOmitEmpty {
		encMapCodes = append(encMapCodes, encStructMap)
	}

	decArrayCodes := make([]Code, 0)
	decArrayCodes = append(decArrayCodes, List(Id("offset"), Err()).Op(":=").Id(ptn.IdDecoder).Dot("CheckStructHeader").Call(Lit(len(st.Fields)), Id("offset")))
	decArrayCodes = append(decArrayCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))

	decMapCodeSwitchCases := make([]Code, 0)

	mapFieldCountVar := "fieldNum"
	if hasOmitEmpty {
		calcMapSizeCodes = append(calcMapSizeCodes, Id(mapFieldCountVar).Op(":=").Lit(0))
		calcMapSizeMaxCodes = append(calcMapSizeMaxCodes, Id(mapFieldCountVar).Op(":=").Lit(0))
		encMapCodes = append(encMapCodes, Id(mapFieldCountVar).Op(":=").Lit(0))
		if canCalcSizeNoErr {
			calcMapSizeNoErrCodes = append(calcMapSizeNoErrCodes, Id(mapFieldCountVar).Op(":=").Lit(0))
			calcMapSizeMaxNoErrCodes = append(calcMapSizeMaxNoErrCodes, Id(mapFieldCountVar).Op(":=").Lit(0))
		}
		for _, field := range st.Fields {
			if field.OmitEmpty {
				fieldName := "v." + field.Name
				calcMapSizeCodes = append(calcMapSizeCodes, If(st.createOmitEmptyCondition(field, fieldName)).Block(Id(mapFieldCountVar).Op("++")))
				calcMapSizeMaxCodes = append(calcMapSizeMaxCodes, If(st.createOmitEmptyCondition(field, fieldName)).Block(Id(mapFieldCountVar).Op("++")))
				encMapCodes = append(encMapCodes, If(st.createOmitEmptyCondition(field, fieldName)).Block(Id(mapFieldCountVar).Op("++")))
				if canCalcSizeNoErr {
					calcMapSizeNoErrCodes = append(calcMapSizeNoErrCodes, If(st.createOmitEmptyCondition(field, fieldName)).Block(Id(mapFieldCountVar).Op("++")))
					calcMapSizeMaxNoErrCodes = append(calcMapSizeMaxNoErrCodes, If(st.createOmitEmptyCondition(field, fieldName)).Block(Id(mapFieldCountVar).Op("++")))
				}
			} else {
				calcMapSizeCodes = append(calcMapSizeCodes, Id(mapFieldCountVar).Op("++"))
				calcMapSizeMaxCodes = append(calcMapSizeMaxCodes, Id(mapFieldCountVar).Op("++"))
				encMapCodes = append(encMapCodes, Id(mapFieldCountVar).Op("++"))
				if canCalcSizeNoErr {
					calcMapSizeNoErrCodes = append(calcMapSizeNoErrCodes, Id(mapFieldCountVar).Op("++"))
					calcMapSizeMaxNoErrCodes = append(calcMapSizeMaxNoErrCodes, Id(mapFieldCountVar).Op("++"))
				}
			}
		}

		calcMapSizeCodes = append(calcMapSizeCodes, st.createStructHeaderCalcCode(mapFieldCountVar))
		calcMapSizeMaxCodes = append(calcMapSizeMaxCodes, st.createStructHeaderCalcCode(mapFieldCountVar))
		encMapCodes = append(encMapCodes, st.createStructHeaderEncMapCode(mapFieldCountVar))
		if canCalcSizeNoErr {
			calcMapSizeNoErrCodes = append(calcMapSizeNoErrCodes, st.createStructHeaderCalcCode(mapFieldCountVar))
			calcMapSizeMaxNoErrCodes = append(calcMapSizeMaxNoErrCodes, st.createStructHeaderCalcCode(mapFieldCountVar))
		}
	}

	for _, field := range st.Fields {
		fieldName := "v." + field.Name

		calcKeyStringCode, writeKeyStringCode := st.createKeyStringCode(field.Tag)
		calcMapSizeCodes = append(calcMapSizeCodes, st.wrapOmitEmptyMapCode(field, fieldName, calcKeyStringCode)...)
		calcMapSizeMaxCodes = append(calcMapSizeMaxCodes, st.wrapOmitEmptyMapCode(field, fieldName, calcKeyStringCode)...)
		encMapCodes = append(encMapCodes, st.wrapOmitEmptyMapCode(field, fieldName, writeKeyStringCode)...)
		if canCalcSizeNoErr {
			calcMapSizeNoErrCodes = append(calcMapSizeNoErrCodes, st.wrapOmitEmptyMapCode(field, fieldName, calcKeyStringCode)...)
			calcMapSizeMaxNoErrCodes = append(calcMapSizeMaxNoErrCodes, st.wrapOmitEmptyMapCode(field, fieldName, calcKeyStringCode)...)
		}

		cArray, cMap, eArray, eMap, dArray, dMap := st.createFieldCode(field.Node, fieldName, fieldName)
		cArrayMax, cMapMax := st.createFieldMaxCode(field.Node, fieldName)
		calcArraySizeCodes = append(calcArraySizeCodes, cArray...)
		calcArraySizeMaxCodes = append(calcArraySizeMaxCodes, cArrayMax...)

		calcMapSizeCodes = append(calcMapSizeCodes, st.wrapOmitEmptyMapCode(field, fieldName, cMap...)...)
		calcMapSizeMaxCodes = append(calcMapSizeMaxCodes, st.wrapOmitEmptyMapCode(field, fieldName, cMapMax...)...)

		if canCalcSizeNoErr {
			cArrayNoErr, cMapNoErr := st.createFieldSizeNoErrCode(field.Node, fieldName, false)
			cArrayMaxNoErr, cMapMaxNoErr := st.createFieldSizeNoErrCode(field.Node, fieldName, true)
			calcArraySizeNoErrCodes = append(calcArraySizeNoErrCodes, cArrayNoErr...)
			calcArraySizeMaxNoErrCodes = append(calcArraySizeMaxNoErrCodes, cArrayMaxNoErr...)
			calcMapSizeNoErrCodes = append(calcMapSizeNoErrCodes, st.wrapOmitEmptyMapCode(field, fieldName, cMapNoErr...)...)
			calcMapSizeMaxNoErrCodes = append(calcMapSizeMaxNoErrCodes, st.wrapOmitEmptyMapCode(field, fieldName, cMapMaxNoErr...)...)
		}

		encArrayCodes = append(encArrayCodes, eArray...)
		encMapCodes = append(encMapCodes, st.wrapOmitEmptyMapCode(field, fieldName, eMap...)...)

		decArrayCodes = append(decArrayCodes, dArray...)

		decMapCaseCodes := append(dMap, Id("count").Op("++"))
		if hasOmitEmpty && !field.OmitEmpty {
			decMapCaseCodes = append(decMapCaseCodes, Id(st.requiredFieldFoundVar(field)).Op("=").True())
		}
		decMapCodeSwitchCases = append(decMapCodeSwitchCases, Case(Lit(field.Tag)).Block(
			decMapCaseCodes...,
		// dMap...,
		),
		)
	}

	// not use jump offset
	//decMapCodeSwitchCases = append(decMapCodeSwitchCases, Default().Block(
	//	Id("offset").Op("=").Id(ptn.IdDecoder).Dot("JumpOffset").Call(Id("offset")),
	//	),
	//)
	decMapCodeSwitchCases = append(decMapCodeSwitchCases, Default().Block(
		Return(Lit(0), Qual("fmt", "Errorf").Call(Lit("unknown key[%s] found"), String().Id("(dataKey)"))),
	),
	)

	decMapCodes := make([]Code, 0)

	if hasOmitEmpty {
		decMapCodes = append(decMapCodes, List(Id("dataLen"), Id("offset"), Err()).Op(":=").Id(ptn.IdDecoder).Dot("MapLength").Call(Id("offset")))
	} else {
		decMapCodes = append(decMapCodes, List(Id("offset"), Err()).Op(":=").Id(ptn.IdDecoder).Dot("CheckStructHeader").Call(Lit(len(st.Fields)), Id("offset")))
	}
	decMapCodes = append(decMapCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	if hasOmitEmpty {
		decMapCodes = append(decMapCodes, If(Id("dataLen").Op(">").Lit(len(st.Fields))).Block(
			Return(Lit(0), Qual("fmt", "Errorf").Call(Lit("data length wrong %d : %d"), Lit(len(st.Fields)), Id("dataLen"))),
		))
		for _, field := range st.Fields {
			if !field.OmitEmpty {
				decMapCodes = append(decMapCodes, Id(st.requiredFieldFoundVar(field)).Op(":=").False())
			}
		}
	}
	//decMapCodes = append(decMapCodes, Id("dataLen").Op(":=").Id(ptn.IdDecoder).Dot("Len").Call())
	//decMapCodes = append(decMapCodes, For(Id("count").Op("<").Id("dataLen").Block(
	decMapCodes = append(decMapCodes, Id("count").Op(":=").Lit(0))
	mapDecodeLoopCodes := []Code{
		Var().Id("dataKey").Index().Byte(),
		List(Id("dataKey"), Id("offset"), Err()).Op("=").Id(ptn.IdDecoder).Dot("AsStringBytes").Call(Id("offset")),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),

		Switch(String().Call(Id("dataKey"))).Block(
			decMapCodeSwitchCases...,
		),
	}
	if hasOmitEmpty {
		decMapCodes = append(decMapCodes, For(Id("count").Op("<").Id("dataLen")).Block(mapDecodeLoopCodes...))
		for _, field := range st.Fields {
			if !field.OmitEmpty {
				decMapCodes = append(decMapCodes, If(Op("!").Id(st.requiredFieldFoundVar(field))).Block(
					Return(Lit(0), Qual("fmt", "Errorf").Call(Lit("required key[%s] not found"), Lit(field.Tag))),
				))
			}
		}
	} else {
		decMapCodes = append(decMapCodes, For(Id("count").Op("<").Lit(len(st.Fields))).Block(mapDecodeLoopCodes...))
	}

	var firstEncParam, firstDecParam *Statement
	if st.NoUseQual {
		firstEncParam = Id(v).Op("*").Id(st.Name)
		firstDecParam = Id(v).Op("*").Id(st.Name)
	} else {
		firstEncParam = Id(v).Op("*").Qual(st.ImportPath, st.Name)
		firstDecParam = Id(v).Op("*").Qual(st.ImportPath, st.Name)
	}

	f.Comment(fmt.Sprintf("// calculate size from %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.CalcArraySizeFuncName()).Params(firstEncParam).Params(Int(), Error()).Block(
		append(calcArraySizeCodes, Return(Id("size"), Nil()))...,
	)

	f.Comment(fmt.Sprintf("// calculate max size from %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.CalcArraySizeMaxFuncName()).Params(firstEncParam).Params(Int(), Error()).Block(
		append(calcArraySizeMaxCodes, Return(Id("size"), Nil()))...,
	)

	f.Comment(fmt.Sprintf("// calculate size from %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.CalcMapSizeFuncName()).Params(firstEncParam).Params(Int(), Error()).Block(
		append(calcMapSizeCodes, Return(Id("size"), Nil()))...,
	)

	f.Comment(fmt.Sprintf("// calculate max size from %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.CalcMapSizeMaxFuncName()).Params(firstEncParam).Params(Int(), Error()).Block(
		append(calcMapSizeMaxCodes, Return(Id("size"), Nil()))...,
	)

	if canCalcSizeNoErr {
		f.Comment(fmt.Sprintf("// calculate no-error size from %s.%s\n", st.ImportPath, st.Name)).
			Func().Id(st.CalcArraySizeNoErrFuncName()).Params(firstEncParam).Params(Int()).Block(
			append(calcArraySizeNoErrCodes, Return(Id("size")))...,
		)

		f.Comment(fmt.Sprintf("// calculate no-error max size from %s.%s\n", st.ImportPath, st.Name)).
			Func().Id(st.CalcArraySizeMaxNoErrFuncName()).Params(firstEncParam).Params(Int()).Block(
			append(calcArraySizeMaxNoErrCodes, Return(Id("size")))...,
		)

		f.Comment(fmt.Sprintf("// calculate no-error size from %s.%s\n", st.ImportPath, st.Name)).
			Func().Id(st.CalcMapSizeNoErrFuncName()).Params(firstEncParam).Params(Int()).Block(
			append(calcMapSizeNoErrCodes, Return(Id("size")))...,
		)

		f.Comment(fmt.Sprintf("// calculate no-error max size from %s.%s\n", st.ImportPath, st.Name)).
			Func().Id(st.CalcMapSizeMaxNoErrFuncName()).Params(firstEncParam).Params(Int()).Block(
			append(calcMapSizeMaxNoErrCodes, Return(Id("size")))...,
		)
	}

	f.Comment(fmt.Sprintf("// encode from %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.EncodeArrayFuncName()).Params(firstEncParam, Id("buf").Index().Byte(), Id("offset").Int()).Params(Int(), Error()).Block(
		append(encArrayCodes, Return(Id("offset"), Err()))...,
	)

	f.Comment(fmt.Sprintf("// encode from %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.EncodeMapFuncName()).Params(firstEncParam, Id("buf").Index().Byte(), Id("offset").Int()).Params(Int(), Error()).Block(
		append(encMapCodes, Return(Id("offset"), Err()))...,
	)

	f.Comment(fmt.Sprintf("// decode to %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.DecodeArrayFuncName()).Params(firstDecParam, Id(ptn.IdDecoder).Op("*").Qual(ptn.PkDec, "Decoder"), Id("offset").Int()).Params(Int(), Error()).Block(
		append(decArrayCodes, Return(Id("offset"), Err()))...,
	)

	f.Comment(fmt.Sprintf("// decode to %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.DecodeMapFuncName()).Params(firstDecParam, Id(ptn.IdDecoder).Op("*").Qual(ptn.PkDec, "Decoder"), Id("offset").Int()).Params(Int(), Error()).Block(

		append(decMapCodes, Return(Id("offset"), Err()))...,
	)

	st.createNotEmptyFunc(f)
	for _, field := range st.Fields {
		if field.OmitEmpty {
			st.createOmitEmptyFieldFunc(f, field)
		}
	}
}

func (st *Structure) hasOmitEmptyField() bool {
	for _, field := range st.Fields {
		if field.OmitEmpty {
			return true
		}
	}
	return false
}

func (st *Structure) createOmitEmptyCondition(field Field, fieldName string) Code {
	return Id(st.omitEmptyFieldFuncName(field)).Call(Id(fieldName))
}

func (st *Structure) wrapOmitEmptyMapCode(field Field, fieldName string, codes ...Code) []Code {
	if !field.OmitEmpty {
		return codes
	}
	return []Code{If(st.createOmitEmptyCondition(field, fieldName)).Block(codes...)}
}

func (st *Structure) requiredFieldFoundVar(field Field) string {
	return "found" + field.Name
}

func (st *Structure) notEmptyFuncName() string {
	return st.createFuncName("isNotEmpty")
}

func (st *Structure) omitEmptyFieldFuncName(field Field) string {
	return st.createFuncName("isNotEmpty" + field.Name)
}

func (st *Structure) createNotEmptyFunc(f *File) {
	var param *Statement
	if st.NoUseQual {
		param = Id("v").Id(st.Name)
	} else {
		param = Id("v").Qual(st.ImportPath, st.Name)
	}

	codes := make([]Code, 0)
	for _, field := range st.ZeroFields {
		codes = append(codes, st.createReturnIfNotEmptyCode(field.Node, "v."+field.Name, 0)...)
	}
	codes = append(codes, Return(False()))

	f.Func().Id(st.notEmptyFuncName()).Params(param).Bool().Block(codes...)
}

func (st *Structure) createOmitEmptyFieldFunc(f *File, field Field) {
	codes := st.createReturnIfNotEmptyCode(field.Node, "v", 0)
	codes = append(codes, Return(False()))

	f.Func().Id(st.omitEmptyFieldFuncName(field)).Params(
		field.Node.TypeJenChain(st.Others, Id("v")),
	).Bool().Block(codes...)
}

func (st *Structure) createReturnIfNotEmptyCode(node *Node, valueName string, depth int) []Code {
	switch {
	case node.IsIdentical():
		return []Code{
			If(st.createIdentNotEmptyCondition(node, valueName)).Block(Return(True())),
		}

	case node.IsSlice(), node.IsMap(), node.IsPointer():
		return []Code{
			If(Id(valueName).Op("!=").Nil()).Block(Return(True())),
		}

	case node.IsArray():
		childName := fmt.Sprintf("vv%d", depth)
		return []Code{
			For(List(Id("_"), Id(childName)).Op(":=").Range().Id(valueName)).Block(
				st.createReturnIfNotEmptyCode(node.Elm(), childName, depth+1)...,
			),
		}

	case node.IsStruct():
		if node.ImportPath == "time" {
			return []Code{
				If(Op("!").Id(valueName).Dot("IsZero").Call()).Block(Return(True())),
			}
		}
		return []Code{
			If(Id(createFuncName("isNotEmpty", node.StructName, node.ImportPath)).Call(Id(valueName))).Block(Return(True())),
		}
	}

	return nil
}

func (st *Structure) createIdentNotEmptyCondition(node *Node, valueName string) Code {
	switch node.IdenticalName {
	case "bool":
		return Id(valueName)
	case "string":
		return Id(valueName).Op("!=").Lit("")
	default:
		return Id(valueName).Op("!=").Lit(0)
	}
}

func (st *Structure) createStructHeaderCalcCode(fieldNum string) Code {
	return If(Id(fieldNum).Op("<=").Lit(0x0f)).Block(
		Id("size").Op("+=").Qual(ptn.PkEnc, "CalcStructHeaderFix").Call(Id(fieldNum)),
	).Else().If(Id(fieldNum).Op("<=").Qual("math", "MaxUint16")).Block(
		Id("size").Op("+=").Qual(ptn.PkEnc, "CalcStructHeader16").Call(Id(fieldNum)),
	).Else().Block(
		Id("size").Op("+=").Qual(ptn.PkEnc, "CalcStructHeader32").Call(Id(fieldNum)),
	)
}

func (st *Structure) createStructHeaderEncMapCode(fieldNum string) Code {
	return If(Id(fieldNum).Op("<=").Lit(0x0f)).Block(
		Id("offset").Op("=").Qual(ptn.PkEnc, "WriteStructHeaderFixAsMap").Call(Id("buf"), Id(fieldNum), Id("offset")),
	).Else().If(Id(fieldNum).Op("<=").Qual("math", "MaxUint16")).Block(
		Id("offset").Op("=").Qual(ptn.PkEnc, "WriteStructHeader16AsMap").Call(Id("buf"), Id(fieldNum), Id("offset")),
	).Else().Block(
		Id("offset").Op("=").Qual(ptn.PkEnc, "WriteStructHeader32AsMap").Call(Id("buf"), Id(fieldNum), Id("offset")),
	)
}

func (st *Structure) createStructCode(fieldNum int) (Code, Code, Code) {

	suffix := ""
	if fieldNum <= 0x0f {
		suffix = "Fix"
	} else if fieldNum <= math.MaxUint16 {
		suffix = "16"
	} else if uint(fieldNum) <= math.MaxUint32 {
		suffix = "32"
	}

	return Id("size").Op("+=").Qual(ptn.PkEnc, "CalcStructHeader"+suffix).Call(Lit(fieldNum)),
		Id("offset").Op("=").Qual(ptn.PkEnc, "WriteStructHeader"+suffix+"AsArray").Call(Id("buf"), Lit(fieldNum), Id("offset")),
		Id("offset").Op("=").Qual(ptn.PkEnc, "WriteStructHeader"+suffix+"AsMap").Call(Id("buf"), Lit(fieldNum), Id("offset"))
}

func (st *Structure) createKeyStringCode(v string) (Code, Code) {
	keyBytes := encodedStringBytes(v)
	keyLen := len(keyBytes)

	return Id("size").Op("+=").Lit(len(keyBytes)),
		Id("offset").Op("+=").Id("copy").Call(Id("buf").Index(Id("offset").Op(":").Id("offset").Op("+").Lit(keyLen)), Lit(string(keyBytes)))
}

func encodedStringBytes(v string) []byte {
	l := len(v)
	b := make([]byte, 0, def.Byte1+def.Byte4+l)
	if l < 32 {
		b = append(b, byte(def.FixStr+l))
	} else if l <= math.MaxUint8 {
		b = append(b, byte(def.Str8), byte(l))
	} else if l <= math.MaxUint16 {
		b = append(b, byte(def.Str16), byte(l>>8), byte(l))
	} else {
		b = append(b, byte(def.Str32), byte(l>>24), byte(l>>16), byte(l>>8), byte(l))
	}
	return append(b, v...)
}

func (st *Structure) createFieldCode(node *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	switch {
	case node.IsIdentical():
		cArray, cMap, eArray, eMap, dArray, dMap = st.createIdentCode(node, encodeFieldName, decodeFieldName)

	case node.IsSlice():
		cArray, cMap, eArray, eMap, dArray, dMap = st.createSliceCode(node, encodeFieldName, decodeFieldName)

	case node.IsArray():
		cArray, cMap, eArray, eMap, dArray, dMap = st.createArrayCode(node, encodeFieldName, decodeFieldName)

	case node.IsMap():
		cArray, cMap, eArray, eMap, dArray, dMap = st.createMapCode(node, encodeFieldName, decodeFieldName)

	case node.IsPointer():
		cArray, cMap, eArray, eMap, dArray, dMap = st.createPointerCode(node, encodeFieldName, decodeFieldName)

	case node.IsStruct():

		if node.ImportPath == "time" {
			cArray, cMap, eArray, eMap, dArray, dMap = st.createTimeCode(encodeFieldName, decodeFieldName, node)
		} else {
			cArray, cMap, eArray, eMap, dArray, dMap = st.createNamedCode(encodeFieldName, decodeFieldName, node)
		}
	}

	return
}
