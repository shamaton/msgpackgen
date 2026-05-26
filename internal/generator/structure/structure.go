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
	NoUseQual  bool

	Others []*Structure
	File   *ast.File

	CanGen  bool
	Reasons []string
}

// Field has a field information in structure
type Field struct {
	Name string
	Tag  string
	Node *Node
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

	calcArraySizeCodes := make([]Code, 0)
	calcArraySizeCodes = append(calcArraySizeCodes, Id("size").Op(":=").Lit(0))
	calcArraySizeCodes = append(calcArraySizeCodes, calcStruct)

	calcArraySizeMaxCodes := make([]Code, 0)
	calcArraySizeMaxCodes = append(calcArraySizeMaxCodes, Id("size").Op(":=").Lit(0))
	calcArraySizeMaxCodes = append(calcArraySizeMaxCodes, calcStruct)

	calcMapSizeCodes := make([]Code, 0)
	calcMapSizeCodes = append(calcMapSizeCodes, Id("size").Op(":=").Lit(0))
	calcMapSizeCodes = append(calcMapSizeCodes, calcStruct)

	calcMapSizeMaxCodes := make([]Code, 0)
	calcMapSizeMaxCodes = append(calcMapSizeMaxCodes, Id("size").Op(":=").Lit(0))
	calcMapSizeMaxCodes = append(calcMapSizeMaxCodes, calcStruct)

	canCalcSizeNoErr := st.CanCalcSizeNoErr()
	calcArraySizeNoErrCodes := make([]Code, 0)
	calcArraySizeMaxNoErrCodes := make([]Code, 0)
	calcMapSizeNoErrCodes := make([]Code, 0)
	calcMapSizeMaxNoErrCodes := make([]Code, 0)
	if canCalcSizeNoErr {
		calcArraySizeNoErrCodes = append(calcArraySizeNoErrCodes, Id("size").Op(":=").Lit(0), calcStruct)
		calcArraySizeMaxNoErrCodes = append(calcArraySizeMaxNoErrCodes, Id("size").Op(":=").Lit(0), calcStruct)
		calcMapSizeNoErrCodes = append(calcMapSizeNoErrCodes, Id("size").Op(":=").Lit(0), calcStruct)
		calcMapSizeMaxNoErrCodes = append(calcMapSizeMaxNoErrCodes, Id("size").Op(":=").Lit(0), calcStruct)
	}

	encArrayCodes := make([]Code, 0)
	encArrayCodes = append(encArrayCodes, Var().Err().Error())
	encArrayCodes = append(encArrayCodes, encStructArray)

	encMapCodes := make([]Code, 0)
	encMapCodes = append(encMapCodes, Var().Err().Error())
	encMapCodes = append(encMapCodes, encStructMap)

	decArrayCodes := make([]Code, 0)
	decArrayCodes = append(decArrayCodes, List(Id("offset"), Err()).Op(":=").Id(ptn.IdDecoder).Dot("CheckStructHeader").Call(Lit(len(st.Fields)), Id("offset")))
	decArrayCodes = append(decArrayCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))

	decMapCodeSwitchCases := make([]Code, 0)
	decKeySliceVar := make([]Code, 0)

	for i, field := range st.Fields {
		fieldName := "v." + field.Name

		calcKeyStringCode, writeKeyStringCode := st.createKeyStringCode(field.Tag)
		calcMapSizeCodes = append(calcMapSizeCodes, calcKeyStringCode)
		calcMapSizeMaxCodes = append(calcMapSizeMaxCodes, calcKeyStringCode)
		encMapCodes = append(encMapCodes, writeKeyStringCode)
		if canCalcSizeNoErr {
			calcMapSizeNoErrCodes = append(calcMapSizeNoErrCodes, calcKeyStringCode)
			calcMapSizeMaxNoErrCodes = append(calcMapSizeMaxNoErrCodes, calcKeyStringCode)
		}

		cArray, cMap, eArray, eMap, dArray, dMap := st.createFieldCode(field.Node, fieldName, fieldName)
		cArrayMax, cMapMax := st.createFieldMaxCode(field.Node, fieldName)
		calcArraySizeCodes = append(calcArraySizeCodes, cArray...)
		calcArraySizeMaxCodes = append(calcArraySizeMaxCodes, cArrayMax...)

		calcMapSizeCodes = append(calcMapSizeCodes, cMap...)
		calcMapSizeMaxCodes = append(calcMapSizeMaxCodes, cMapMax...)

		if canCalcSizeNoErr {
			cArrayNoErr, cMapNoErr := st.createFieldSizeNoErrCode(field.Node, fieldName, false)
			cArrayMaxNoErr, cMapMaxNoErr := st.createFieldSizeNoErrCode(field.Node, fieldName, true)
			calcArraySizeNoErrCodes = append(calcArraySizeNoErrCodes, cArrayNoErr...)
			calcArraySizeMaxNoErrCodes = append(calcArraySizeMaxNoErrCodes, cArrayMaxNoErr...)
			calcMapSizeNoErrCodes = append(calcMapSizeNoErrCodes, cMapNoErr...)
			calcMapSizeMaxNoErrCodes = append(calcMapSizeMaxNoErrCodes, cMapMaxNoErr...)
		}

		encArrayCodes = append(encArrayCodes, eArray...)
		encMapCodes = append(encMapCodes, eMap...)

		decArrayCodes = append(decArrayCodes, dArray...)

		tagBytes := []byte(field.Tag)
		lits := make([]Code, len(tagBytes))
		for i := range tagBytes {
			lits = append(lits, Lit(tagBytes[i]))
		}
		decKeySliceVar = append(decKeySliceVar, Values(lits...).Id(",").Commentf("%s", field.Tag))

		decMapCodeSwitchCases = append(decMapCodeSwitchCases, Case(Lit(i)).Block(
			append(dMap, Id("count").Op("++"))...,
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

	decMapCodes = append(decMapCodes, Id("keys").Op(":=").Index().Index().Byte().Block(decKeySliceVar...))

	decMapCodes = append(decMapCodes, List(Id("offset"), Err()).Op(":=").Id(ptn.IdDecoder).Dot("CheckStructHeader").Call(Lit(len(st.Fields)), Id("offset")))
	decMapCodes = append(decMapCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	//decMapCodes = append(decMapCodes, Id("dataLen").Op(":=").Id(ptn.IdDecoder).Dot("Len").Call())
	//decMapCodes = append(decMapCodes, For(Id("count").Op("<").Id("dataLen").Block(
	decMapCodes = append(decMapCodes, Id("count").Op(":=").Lit(0))
	decMapCodes = append(decMapCodes, For(Id("count").Op("<").Lit(len(st.Fields)).Block(
		Var().Id("dataKey").Index().Byte(),
		List(Id("dataKey"), Id("offset"), Err()).Op("=").Id(ptn.IdDecoder).Dot("AsStringBytes").Call(Id("offset")),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),

		Id("fieldIndex").Op(":=").Lit(-1),
		For(List(Id("i"), Id("key"))).Op(":=").Range().Id("keys").Block(
			If(Len(Id("dataKey")).Op("!=").Len(Id("key"))).Block(
				Continue(),
			),

			Id("fieldIndex").Op("=").Id("i"),
			For(Id("dataKeyIndex")).Op(":=").Range().Id("dataKey").Block(
				If(Id("dataKey").Index(Id("dataKeyIndex")).Op("!=").Id("key").Index(Id("dataKeyIndex"))).Block(
					Id("fieldIndex").Op("=").Lit(-1),
					Break(),
				),
			),
			If(Id("fieldIndex").Op(">=").Lit(0)).Block(
				Break(),
			),
		),

		Switch(Id("fieldIndex")).Block(
			decMapCodeSwitchCases...,
		),
	)))

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
