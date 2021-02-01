package structure

import (
	"crypto/sha256"
	"fmt"
	"go/ast"
	"math"
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/shamaton/msgpackgen/internal/generator/ptn"
)

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

type Field struct {
	Name string
	Tag  string
	Node *Node
}

func (st *Structure) CalcArraySizeFuncName() string {
	return st.createFuncName("calcArraySize")
}

func (st *Structure) CalcMapSizeFuncName() string {
	return st.createFuncName("calcMapSize")
}

func (st *Structure) EncodeArrayFuncName() string {
	return st.createFuncName("encodeArray")
}

func (st *Structure) EncodeMapFuncName() string {
	return st.createFuncName("encodeMap")
}

func (st *Structure) DecodeArrayFuncName() string {
	return st.createFuncName("decodeArray")
}

func (st *Structure) DecodeMapFuncName() string {
	return st.createFuncName("decodeMap")
}

func (st *Structure) createFuncName(prefix string) string {
	return createFuncName(prefix, st.Name, st.ImportPath)
}

func createFuncName(prefix, name, importPath string) string {
	suffix := fmt.Sprintf("%x", sha256.Sum256([]byte(importPath)))
	return ptn.PrivateFuncName(fmt.Sprintf("%s%s_%s", prefix, name, suffix))
}

func (st *Structure) CreateCode(f *File) {
	v := "v"

	calcStruct, encStructArray, encStructMap := st.CreateStructCode(len(st.Fields))

	calcArraySizeCodes := make([]Code, 0)
	calcArraySizeCodes = append(calcArraySizeCodes, Id("size").Op(":=").Lit(0))
	calcArraySizeCodes = append(calcArraySizeCodes, calcStruct)

	calcMapSizeCodes := make([]Code, 0)
	calcMapSizeCodes = append(calcMapSizeCodes, Id("size").Op(":=").Lit(0))
	calcMapSizeCodes = append(calcMapSizeCodes, calcStruct)

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

	for _, field := range st.Fields {
		fieldName := "v." + field.Name

		calcKeyStringCode, writeKeyStringCode := st.CreateKeyStringCode(field.Tag)
		calcMapSizeCodes = append(calcMapSizeCodes, calcKeyStringCode)
		encMapCodes = append(encMapCodes, writeKeyStringCode)

		cArray, cMap, eArray, eMap, dArray, dMap := st.createFieldCode(field.Node, fieldName, fieldName)
		calcArraySizeCodes = append(calcArraySizeCodes, cArray...)

		calcMapSizeCodes = append(calcMapSizeCodes, cMap...)

		encArrayCodes = append(encArrayCodes, eArray...)
		encMapCodes = append(encMapCodes, eMap...)

		decArrayCodes = append(decArrayCodes, dArray...)

		decMapCodeSwitchCases = append(decMapCodeSwitchCases, Case(Lit(field.Tag)).Block(
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
		Return(Lit(0), Qual("fmt", "Errorf").Call(Lit("unknown key[%s] found"), Id("s"))),
	),
	)

	decMapCodes := make([]Code, 0)
	decMapCodes = append(decMapCodes, List(Id("offset"), Err()).Op(":=").Id(ptn.IdDecoder).Dot("CheckStructHeader").Call(Lit(len(st.Fields)), Id("offset")))
	decMapCodes = append(decMapCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	//decMapCodes = append(decMapCodes, Id("dataLen").Op(":=").Id(ptn.IdDecoder).Dot("Len").Call())
	//decMapCodes = append(decMapCodes, For(Id("count").Op("<").Id("dataLen").Block(
	decMapCodes = append(decMapCodes, Id("count").Op(":=").Lit(0))
	decMapCodes = append(decMapCodes, For(Id("count").Op("<").Lit(len(st.Fields)).Block(
		Var().Id("s").String(),
		List(Id("s"), Id("offset"), Err()).Op("=").Id(ptn.IdDecoder).Dot("AsString").Call(Id("offset")),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Switch(Id("s")).Block(
			decMapCodeSwitchCases...,
		),
	)))

	var firstEncParam, firstDecParam *Statement
	if st.NoUseQual {
		firstEncParam = Id(v).Id(st.Name)
		firstDecParam = Id(v).Op("*").Id(st.Name)
	} else {
		firstEncParam = Id(v).Qual(st.ImportPath, st.Name)
		firstDecParam = Id(v).Op("*").Qual(st.ImportPath, st.Name)
	}

	f.Comment(fmt.Sprintf("// calculate size from %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.CalcArraySizeFuncName()).Params(firstEncParam, Id(ptn.IdEncoder).Op("*").Qual(ptn.PkEnc, "Encoder")).Params(Int(), Error()).Block(
		append(calcArraySizeCodes, Return(Id("size"), Nil()))...,
	)

	f.Comment(fmt.Sprintf("// calculate size from %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.CalcMapSizeFuncName()).Params(firstEncParam, Id(ptn.IdEncoder).Op("*").Qual(ptn.PkEnc, "Encoder")).Params(Int(), Error()).Block(
		append(calcMapSizeCodes, Return(Id("size"), Nil()))...,
	)

	f.Comment(fmt.Sprintf("// encode from %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.EncodeArrayFuncName()).Params(firstEncParam, Id(ptn.IdEncoder).Op("*").Qual(ptn.PkEnc, "Encoder"), Id("offset").Int()).Params(Index().Byte(), Int(), Error()).Block(
		append(encArrayCodes, Return(Id(ptn.IdEncoder).Dot("EncodedBytes").Call(), Id("offset"), Err()))...,
	)

	f.Comment(fmt.Sprintf("// encode from %s.%s\n", st.ImportPath, st.Name)).
		Func().Id(st.EncodeMapFuncName()).Params(firstEncParam, Id(ptn.IdEncoder).Op("*").Qual(ptn.PkEnc, "Encoder"), Id("offset").Int()).Params(Index().Byte(), Int(), Error()).Block(
		append(encMapCodes, Return(Id(ptn.IdEncoder).Dot("EncodedBytes").Call(), Id("offset"), Err()))...,
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

func (st *Structure) CreateKeyStringCode(v string) (Code, Code) {
	l := len(v)
	suffix := ""
	if l < 32 {
		suffix = "Fix"
	} else if l <= math.MaxUint8 {
		suffix = "8"
	} else if l <= math.MaxUint16 {
		suffix = "16"
	} else {
		suffix = "32"
	}

	return Id("size").Op("+=").Id(ptn.IdEncoder).Dot("CalcString" + suffix).Call(Lit(l)),
		Id("offset").Op("=").Id(ptn.IdEncoder).Dot("WriteString"+suffix).Call(Lit(v), Lit(l), Id("offset"))
}

func (st *Structure) CreateStructCode(fieldNum int) (Code, Code, Code) {

	suffix := ""
	if fieldNum <= 0x0f {
		suffix = "Fix"
	} else if fieldNum <= math.MaxUint16 {
		suffix = "16"
	} else if uint(fieldNum) <= math.MaxUint32 {
		suffix = "32"
	}

	return Id("size").Op("+=").Id(ptn.IdEncoder).Dot("CalcStructHeader" + suffix).Call(Lit(fieldNum)),
		Id("offset").Op("=").Id(ptn.IdEncoder).Dot(" WriteStructHeader"+suffix+"AsArray").Call(Lit(fieldNum), Id("offset")),
		Id("offset").Op("=").Id(ptn.IdEncoder).Dot(" WriteStructHeader"+suffix+"AsMap").Call(Lit(fieldNum), Id("offset"))
}

func (st *Structure) createFieldCode(ast *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	switch {
	case ast.IsIdentical():
		return st.createBasicCode(ast, encodeFieldName, decodeFieldName)

	case ast.IsSlice():
		return st.createSliceCode(ast, encodeFieldName, decodeFieldName)

	case ast.IsArray():
		return st.createArrayCode(ast, encodeFieldName, decodeFieldName)

	case ast.IsMap():
		return st.createMapCode(ast, encodeFieldName, decodeFieldName)

	case ast.IsPointer():
		return st.createPointerCode(ast, encodeFieldName, decodeFieldName)

	case ast.IsStruct():

		if ast.ImportPath == "time" {
			fieldValue := Id(encodeFieldName)
			cArray = append(cArray, st.addSizePattern1("CalcTime", fieldValue))
			eArray = append(eArray, st.encPattern1("WriteTime", fieldValue, Id("offset")))

			cMap = append(cMap, st.addSizePattern1("CalcTime", fieldValue))
			eMap = append(eMap, st.encPattern1("WriteTime", fieldValue, Id("offset")))

			dArray = append(dArray, st.decodeBasicPattern(ast, decodeFieldName, "offset", "AsDateTime")...)
			dMap = append(dMap, st.decodeBasicPattern(ast, decodeFieldName, "offset", "AsDateTime")...)
		} else {
			return st.createNamedCode(encodeFieldName, decodeFieldName, ast)
		}
	}

	return cArray, cMap, eArray, eMap, dArray, dMap
}

func (st *Structure) createPointerCode(ast *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	encodeChildName := encodeFieldName + "p"
	if isRootField(encodeFieldName) {
		encodeChildName = "vp"
	}

	ca, _, ea, _, da, _ := st.createFieldCode(ast.Elm(), encodeChildName, decodeFieldName)

	cArray = make([]Code, 0)
	cArray = append(cArray, If(Id(encodeFieldName).Op("!=").Nil()).Block(
		append([]Code{
			Id(encodeChildName).Op(":=").Op("*").Id(encodeFieldName),
		}, ca...)...,
	).Else().Block(
		Id("size").Op("+=").Id(ptn.IdEncoder).Dot("CalcNil").Call(),
	))

	eArray = make([]Code, 0)
	eArray = append(eArray, If(Id(encodeFieldName).Op("!=").Nil()).Block(
		append([]Code{
			Id(encodeChildName).Op(":=").Op("*").Id(encodeFieldName),
		}, ea...)...,
	).Else().Block(
		Id("offset").Op("=").Id(ptn.IdEncoder).Dot("WriteNil").Call(Id("offset")),
	))

	if ast.IsParentPointer() {
		dArray = da
	} else {
		dArray = make([]Code, 0)
		dArray = append(dArray, If(Op("!").Id(ptn.IdDecoder).Dot("IsCodeNil").Call(Id("offset"))).Block(
			da...,
		).Else().Block(
			Id("offset").Op("++"),
		))
	}

	return cArray, cArray, eArray, eArray, dArray, dArray
}

func (st *Structure) createMapCode(ast *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	key, value := ast.KeyValue()

	encodeChildKey, encodeChildValue := encodeFieldName+"k", encodeFieldName+"v"
	if isRootField(encodeFieldName) {
		encodeChildKey = "kk"
		encodeChildValue = "vv"
	}

	decodeChildKey, decodeChildValue := decodeFieldName+"k", decodeFieldName+"v"
	if isRootField(decodeFieldName) {
		decodeChildKey = "kk"
		decodeChildValue = "vv"
	}

	caKey, _, eaKey, _, daKey, _ := st.createFieldCode(key, encodeChildKey, decodeChildKey)
	caValue, _, eaValue, _, daValue, _ := st.createFieldCode(value, encodeChildValue, decodeChildValue)

	calcCodes := st.addSizePattern2("CalcMapLength", Len(Id(encodeFieldName)))
	calcCodes = append(calcCodes, For(List(Id(encodeChildKey), Id(encodeChildValue)).Op(":=").Range().Id(encodeFieldName)).Block(
		append(caKey, caValue...)...,
	))

	cArray = append(cArray, If(Id(encodeFieldName).Op("!=").Nil()).Block(
		calcCodes...,
	).Else().Block(
		st.addSizePattern1("CalcNil"),
	))

	encCodes := make([]Code, 0)
	encCodes = append(encCodes, Id("offset").Op("=").Id(ptn.IdEncoder).Dot("WriteMapLength").Call(Len(Id(encodeFieldName)), Id("offset")))
	encCodes = append(encCodes, For(List(Id(encodeChildKey), Id(encodeChildValue)).Op(":=").Range().Id(encodeFieldName)).Block(
		append(eaKey, eaValue...)...,
	))

	eArray = append(eArray, If(Id(encodeFieldName).Op("!=").Nil()).Block(
		encCodes...,
	).Else().Block(
		Id("offset").Op("=").Id(ptn.IdEncoder).Dot("WriteNil").Call(Id("offset")),
	))

	decCodes := make([]Code, 0)
	decCodes = append(decCodes, ast.TypeJenChain(st.Others, Var().Id(decodeChildValue)))
	decCodes = append(decCodes, Var().Id(decodeChildValue+"l").Int())
	decCodes = append(decCodes, List(Id(decodeChildValue+"l"), Id("offset"), Err()).Op("=").Id(ptn.IdDecoder).Dot("MapLength").Call(Id("offset")))
	decCodes = append(decCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	decCodes = append(decCodes, Id(decodeChildValue).Op("=").Make(ast.TypeJenChain(st.Others), Id(decodeChildValue+"l")))

	da := []Code{ast.Key.TypeJenChain(st.Others, Var().Id(decodeChildKey+"v"))}
	da = append(da, daKey...)
	da = append(da, ast.Value.TypeJenChain(st.Others, Var().Id(decodeChildValue+"v")))
	da = append(da, daValue...)
	da = append(da, Id(decodeChildValue).Index(Id(decodeChildKey+"v")).Op("=").Id(decodeChildValue+"v"))

	decCodes = append(decCodes, For(Id(decodeChildValue+"i").Op(":=").Lit(0).Op(";").Id(decodeChildValue+"i").Op("<").Id(decodeChildValue+"l").Op(";").Id(decodeChildValue+"i").Op("++")).Block(
		da...,
	))

	name := decodeChildValue
	andOp := ""
	prtCount, _ := ast.GetPointerInfo()

	if prtCount > 0 {
		andOp = "&"
	}

	for i := 0; i < prtCount-1; i++ {
		n := "_" + name
		decCodes = append(decCodes, Id(n).Op(":=").Op("&").Id(name))
		name = n
	}

	decCodes = append(decCodes, Id(decodeFieldName).Op("=").Op(andOp).Id(name))

	if ast.HasParent() && ast.Parent.IsPointer() {
		dArray = decCodes
	} else {

		dArray = append(dArray, If(Op("!").Id(ptn.IdDecoder).Dot("IsCodeNil").Call(Id("offset"))).Block(
			decCodes...,
		).Else().Block(
			Id("offset").Op("++"),
		))
	}

	return cArray, cArray, eArray, eArray, dArray, dArray
}

func (st *Structure) createSliceCode(ast *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	encodeChildName, decodeChildName := encodeFieldName+"v", decodeFieldName+""
	if isRootField(encodeFieldName) {
		encodeChildName = "vv"
	}
	if isRootField(decodeFieldName) {
		decodeChildName = "vv"
	}

	decodeChildLengthName := decodeChildName + "l"
	decodeChildIndexName := decodeChildName + "i"
	decodeChildChildName := decodeChildName + "v"

	ca, cm, ea, em, da, dm := st.createFieldCode(ast.Elm(), encodeChildName, decodeChildName)
	isChildByte := ast.Elm().IsIdentical() && ast.Elm().IdenticalName == "byte"

	// calc array
	caCodes := st.addSizePattern2("CalcSliceLength", Len( /*Op(ptrOp).*/ Id(encodeFieldName)), Lit(isChildByte))
	caCodes = append(caCodes, For(List(Id("_"), Id(encodeChildName)).Op(":=").Range(). /*Op(ptrOp).*/ Id(encodeFieldName)).Block(
		ca...,
	))

	cArray = append(cArray, If( /*Op(ptrOp).*/ Id(encodeFieldName).Op("!=").Nil()).Block(
		caCodes...,
	).Else().Block(
		st.addSizePattern1("CalcNil"),
	))

	// calc map
	cmCodes := st.addSizePattern2("CalcSliceLength", Len( /*Op(ptrOp).*/ Id(encodeFieldName)), Lit(isChildByte))
	cmCodes = append(cmCodes, For(List(Id("_"), Id(encodeChildName)).Op(":=").Range(). /*Op(ptrOp).*/ Id(encodeFieldName)).Block(
		cm...,
	))

	cMap = append(cMap, If( /*Op(ptrOp).*/ Id(encodeFieldName).Op("!=").Nil()).Block(
		cmCodes...,
	).Else().Block(
		st.addSizePattern1("CalcNil"),
	))

	// encode array
	eaCodes := make([]Code, 0)
	eaCodes = append(eaCodes, Id("offset").Op("=").Id(ptn.IdEncoder).Dot("WriteSliceLength").Call(Len( /*Op(ptrOp).*/ Id(encodeFieldName)), Id("offset"), Lit(isChildByte)))
	eaCodes = append(eaCodes, For(List(Id("_"), Id(encodeChildName)).Op(":=").Range(). /*Op(ptrOp).*/ Id(encodeFieldName)).Block(
		ea...,
	))

	eArray = append(eArray, If( /*Op(ptrOp).*/ Id(encodeFieldName).Op("!=").Nil()).Block(
		eaCodes...,
	).Else().Block(
		Id("offset").Op("=").Id(ptn.IdEncoder).Dot("WriteNil").Call(Id("offset")),
	))

	// encode map
	emCodes := make([]Code, 0)
	emCodes = append(emCodes, Id("offset").Op("=").Id(ptn.IdEncoder).Dot("WriteSliceLength").Call(Len( /*Op(ptrOp).*/ Id(encodeFieldName)), Id("offset"), Lit(isChildByte)))
	emCodes = append(emCodes, For(List(Id("_"), Id(encodeChildName)).Op(":=").Range(). /*Op(ptrOp).*/ Id(encodeFieldName)).Block(
		em...,
	))

	eMap = append(eMap, If( /*Op(ptrOp).*/ Id(encodeFieldName).Op("!=").Nil()).Block(
		emCodes...,
	).Else().Block(
		Id("offset").Op("=").Id(ptn.IdEncoder).Dot("WriteNil").Call(Id("offset")),
	))

	daCodes := make([]Code, 0)
	daCodes = append(daCodes, ast.TypeJenChain(st.Others, Var().Id(decodeChildName)))
	daCodes = append(daCodes, Var().Id(decodeChildLengthName).Int())
	daCodes = append(daCodes, List(Id(decodeChildLengthName), Id("offset"), Err()).Op("=").Id(ptn.IdDecoder).Dot("SliceLength").Call(Id("offset")))
	daCodes = append(daCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	daCodes = append(daCodes, Id(decodeChildName).Op("=").Make(ast.TypeJenChain(st.Others), Id(decodeChildLengthName)))

	da = append([]Code{ast.Elm().TypeJenChain(st.Others, Var().Id(decodeChildChildName))}, da...)
	da = append(da, Id(decodeChildName).Index(Id(decodeChildIndexName)).Op("=").Id(decodeChildChildName))

	daCodes = append(daCodes, For(Id(decodeChildIndexName).Op(":=").Range().Id(decodeChildName)).Block(
		da...,
	))

	dmCodes := make([]Code, 0)
	dmCodes = append(dmCodes, ast.TypeJenChain(st.Others, Var().Id(decodeChildName)))
	dmCodes = append(dmCodes, Var().Id(decodeChildLengthName).Int())
	dmCodes = append(dmCodes, List(Id(decodeChildLengthName), Id("offset"), Err()).Op("=").Id(ptn.IdDecoder).Dot("SliceLength").Call(Id("offset")))
	dmCodes = append(dmCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	dmCodes = append(dmCodes, Id(decodeChildName).Op("=").Make(ast.TypeJenChain(st.Others), Id(decodeChildLengthName)))

	dm = append([]Code{ast.Elm().TypeJenChain(st.Others, Var().Id(decodeChildChildName))}, dm...)
	dm = append(dm, Id(decodeChildName).Index(Id(decodeChildIndexName)).Op("=").Id(decodeChildChildName))

	dmCodes = append(dmCodes, For(Id(decodeChildIndexName).Op(":=").Range().Id(decodeChildName)).Block(
		dm...,
	))

	name := decodeChildName
	andOp := ""
	prtCount, _ := ast.GetPointerInfo()
	if prtCount > 0 {
		andOp = "&"
	}
	for i := 0; i < prtCount-1; i++ {
		n := "_" + name
		daCodes = append(daCodes, Id(n).Op(":=").Op("&").Id(name))
		dmCodes = append(dmCodes, Id(n).Op(":=").Op("&").Id(name))
		name = n
	}

	daCodes = append(daCodes, Id(decodeFieldName).Op("=").Op(andOp).Id(name))
	dmCodes = append(dmCodes, Id(decodeFieldName).Op("=").Op(andOp).Id(name))

	if ast.HasParent() && ast.Parent.IsPointer() {
		dArray = daCodes
		dMap = dmCodes
	} else {
		dArray = append(dArray, If(Op("!").Id(ptn.IdDecoder).Dot("IsCodeNil").Call(Id("offset"))).Block(
			daCodes...,
		).Else().Block(
			Id("offset").Op("++"),
		))
		dMap = append(dMap, If(Op("!").Id(ptn.IdDecoder).Dot("IsCodeNil").Call(Id("offset"))).Block(
			dmCodes...,
		).Else().Block(
			Id("offset").Op("++"),
		))
	}
	return
}

func (st *Structure) createArrayCode(ast *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	encodeChildName := encodeFieldName + "v"
	if isRootField(encodeFieldName) {
		encodeChildName = "vv"
	}

	decodeChildName := decodeFieldName + "v"
	if isRootField(decodeFieldName) {
		decodeChildName = "vv"
	}

	ca, _, ea, _, da, _ := st.createFieldCode(ast.Elm(), encodeChildName, decodeChildName)
	isChildByte := ast.Elm().IsIdentical() && ast.Elm().IdenticalName == "byte"

	calcCodes := st.addSizePattern2("CalcSliceLength", Len( /*Op(ptrOp).*/ Id(encodeFieldName)), Lit(isChildByte))
	calcCodes = append(calcCodes, For(List(Id("_"), Id(encodeChildName)).Op(":=").Range(). /*Op(ptrOp).*/ Id(encodeFieldName)).Block(
		ca...,
	))

	cArray = append(cArray /* If(Op(ptrOp).Id(name).Op("!=").Nil()).*/, Block(
		calcCodes...,
	), /*.Else().Block(
		st.addSizePattern1("CalcNil"),
	)*/)

	encCodes := make([]Code, 0)
	encCodes = append(encCodes, Id("offset").Op("=").Id(ptn.IdEncoder).Dot("WriteSliceLength").Call(Len( /*Op(ptrOp).*/ Id(encodeFieldName)), Id("offset"), Lit(isChildByte)))
	encCodes = append(encCodes, For(List(Id("_"), Id(encodeChildName)).Op(":=").Range(). /*Op(ptrOp).*/ Id(encodeFieldName)).Block(
		ea...,
	))

	eArray = append(eArray /*If(Op(ptrOp).Id(name).Op("!=").Nil()).*/, Block(
		encCodes...,
	), /*.Else().Block(
		Id("offset").Op("=").Id(IdEncoder).Dot("WriteNil").Call(Id("offset")),
	)*/)

	decCodes := make([]Code, 0)
	decCodes = append(decCodes, ast.TypeJenChain(st.Others, Var().Id(decodeChildName)))
	decCodes = append(decCodes, Var().Id(decodeChildName+"l").Int())
	decCodes = append(decCodes, List(Id(decodeChildName+"l"), Id("offset"), Err()).Op("=").Id(ptn.IdDecoder).Dot("SliceLength").Call(Id("offset")))
	decCodes = append(decCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	decCodes = append(decCodes, If(Id(decodeChildName+"l").Op(">").Id(fmt.Sprint(ast.ArrayLen))).Block(
		Return(Lit(0), Qual("fmt", "Errorf").Call(Lit("length size(%d) is over array size(%d)"), Id(decodeChildName+"l"), Id(fmt.Sprint(ast.ArrayLen)))),
	))

	da = append([]Code{ast.Elm().TypeJenChain(st.Others, Var().Id(decodeChildName+"v"))}, da...)
	da = append(da, Id(decodeChildName).Index(Id(decodeChildName+"i")).Op("=").Id(decodeChildName+"v"))

	decCodes = append(decCodes, For(Id(decodeChildName+"i").Op(":=").Range().Id(decodeChildName).Index(Id(":"+decodeChildName+"l"))).Block(
		da...,
	))

	name := decodeChildName
	andOp := ""
	prtCount, _ := ast.GetPointerInfo()
	if prtCount > 0 {
		andOp = "&"
	}
	for i := 0; i < prtCount-1; i++ {
		n := "_" + name
		decCodes = append(decCodes, Id(n).Op(":=").Op("&").Id(name))
		name = n
	}

	decCodes = append(decCodes, Id(decodeFieldName).Op("=").Op(andOp).Id(name))

	if ast.HasParent() && ast.Parent.IsPointer() {
		dArray = decCodes
	} else {

		dArray = append(dArray, If(Op("!").Id(ptn.IdDecoder).Dot("IsCodeNil").Call(Id("offset"))).Block(
			decCodes...,
		).Else().Block(
			Id("offset").Op("++"),
		))
	}

	return cArray, cArray, eArray, eArray, dArray, dArray
}

func (st *Structure) createBasicCode(ast *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	funcSuffix := strings.Title(ast.IdenticalName)

	cArray = append(cArray, st.addSizePattern1("Calc"+funcSuffix, Id(encodeFieldName)))
	eArray = append(eArray, st.encPattern1("Write"+funcSuffix, Id(encodeFieldName), Id("offset")))

	cMap = append(cMap, st.addSizePattern1("Calc"+funcSuffix, Id(encodeFieldName)))
	eMap = append(eMap, st.encPattern1("Write"+funcSuffix, Id(encodeFieldName), Id("offset")))

	dArray = append(dArray, st.decodeBasicPattern(ast, decodeFieldName, "offset", "As"+funcSuffix)...)
	dMap = append(dMap, st.decodeBasicPattern(ast, decodeFieldName, "offset", "As"+funcSuffix)...)

	return
}

func (st *Structure) addSizePattern1(funcName string, params ...Code) Code {
	return Id("size").Op("+=").Id(ptn.IdEncoder).Dot(funcName).Call(params...)
}

func (st *Structure) addSizePattern2(funcName string, params ...Code) []Code {
	return []Code{
		List(Id("s"), Err()).Op(":=").Id(ptn.IdEncoder).Dot(funcName).Call(params...),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id("size").Op("+=").Id("s"),
	}

}

func (st *Structure) encPattern1(funcName string, params ...Code) Code {
	return Id("offset").Op("=").Id(ptn.IdEncoder).Dot(funcName).Call(params...)
}

func (st *Structure) decodeBasicPattern(ast *Node, fieldName, offsetName, decoderFuncName string) []Code {

	varName := fieldName + "v"
	if isRootField(fieldName) {
		varName = "vv"
	}

	ptrCount, isParentTypeArrayOrMap := ast.GetPointerInfo()

	codes := make([]Code, 0)
	receiverName := varName

	if ptrCount < 1 && !isParentTypeArrayOrMap {
		codes = append(codes, ast.TypeJenChain(st.Others, Var().Id(receiverName)))
	} else if isParentTypeArrayOrMap {

		for i := 0; i < ptrCount; i++ {
			p := strings.Repeat("p", i+1)
			ptr := strings.Repeat("*", ptrCount-1-i)
			codes = append(codes, ast.TypeJenChain(st.Others, Var().Id(varName+p).Op(ptr)))
		}
		receiverName = varName + strings.Repeat("p", ptrCount)
	} else {
		for i := 0; i < ptrCount; i++ {
			p := strings.Repeat("p", i)
			ptr := strings.Repeat("*", ptrCount-1-i)
			codes = append(codes, ast.TypeJenChain(st.Others, Var().Id(varName+p).Op(ptr)))
		}
		receiverName = varName + strings.Repeat("p", ptrCount-1)
	}

	codes = append(codes,
		List(Id(receiverName), Id(offsetName), Err()).Op("=").Id(ptn.IdDecoder).Dot(decoderFuncName).Call(Id(offsetName)),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
	)

	codes = st.createDecodeSetVarPattern(ptrCount, varName, fieldName /*setVarName*/, isParentTypeArrayOrMap, codes)

	// array or map
	if isParentTypeArrayOrMap {
		return codes
	}

	return []Code{Block(codes...)}
}

func (st *Structure) createDecodeSetVarPattern(ptrCount int, varName, setVarName string, isLastSkip bool, codes []Code) []Code {

	if isLastSkip {
		for i := 0; i < ptrCount; i++ {
			tmp1 := varName + strings.Repeat("p", ptrCount-1-i)
			tmp2 := varName + strings.Repeat("p", ptrCount-i)
			codes = append(codes, Id(tmp1).Op("=").Op("&").Id(tmp2))
		}
	} else {

		for i := 0; i < ptrCount; i++ {
			if i != ptrCount-1 {
				tmp1 := varName + strings.Repeat("p", ptrCount-2-i)
				tmp2 := varName + strings.Repeat("p", ptrCount-1-i)
				codes = append(codes, Id(tmp1).Op("=").Op("&").Id(tmp2))
			} else {
				// last
				tmp := varName + strings.Repeat("p", 0)
				codes = append(codes, Id(setVarName).Op("=").Op("&").Id(tmp))
			}
		}
		if ptrCount < 1 {
			codes = append(codes, Id(setVarName).Op("=").Op("").Id(varName))
		}
	}

	return codes
}

func (st *Structure) createNamedCode(encodeFieldName, decodeFieldName string, ast *Node) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	sizeName := "size_" + encodeFieldName
	if isRootField(encodeFieldName) {
		sizeName = strings.ReplaceAll(sizeName, ".", "_")
	}

	cArray = []Code{
		List(Id(sizeName), Err()).
			Op(":=").
			Id(createFuncName("calcArraySize", ast.StructName, ast.ImportPath)).Call(Id(encodeFieldName), Id(ptn.IdEncoder)),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id("size").Op("+=").Id(sizeName),
	}

	cMap = []Code{
		List(Id(sizeName), Err()).
			Op(":=").
			Id(createFuncName("calcMapSize", ast.StructName, ast.ImportPath)).Call(Id(encodeFieldName), Id(ptn.IdEncoder)),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id("size").Op("+=").Id(sizeName),
	}

	eArray = []Code{
		List(Id("_"), Id("offset"), Err()).
			Op("=").
			Id(createFuncName("encodeArray", ast.StructName, ast.ImportPath)).Call(Id(encodeFieldName), Id(ptn.IdEncoder), Id("offset")),
		If(Err().Op("!=").Nil()).Block(
			Return(Nil(), Lit(0), Err()),
		),
	}

	eMap = []Code{
		List(Id("_"), Id("offset"), Err()).
			Op("=").
			Id(createFuncName("encodeMap", ast.StructName, ast.ImportPath)).Call(Id(encodeFieldName), Id(ptn.IdEncoder), Id("offset")),
		If(Err().Op("!=").Nil()).Block(
			Return(Nil(), Lit(0), Err()),
		),
	}

	dArray = append(dArray, st.decodeNamedPattern(ast, decodeFieldName, "decodeArray")...)
	dMap = append(dMap, st.decodeNamedPattern(ast, decodeFieldName, "decodeMap")...)

	return
}

func (st *Structure) decodeNamedPattern(ast *Node, fieldName, decodeFuncName string) []Code {

	varName := fieldName + "v"
	if isRootField(fieldName) {
		varName = "vv"
	}

	ptrCount, isParentTypeArrayOrMap := ast.GetPointerInfo()

	codes := make([]Code, 0)
	receiverName := varName

	if ptrCount < 1 && !isParentTypeArrayOrMap {
		codes = append(codes, ast.TypeJenChain(st.Others, Var().Id(receiverName)))
	} else if isParentTypeArrayOrMap {

		for i := 0; i < ptrCount; i++ {
			p := strings.Repeat("p", i+1)
			ptr := strings.Repeat("*", ptrCount-1-i)

			codes = append(codes, ast.TypeJenChain(st.Others, Var().Id(varName+p).Op(ptr)))
		}
		receiverName = varName + strings.Repeat("p", ptrCount)
	} else {
		for i := 0; i < ptrCount; i++ {
			p := strings.Repeat("p", i)
			ptr := strings.Repeat("*", ptrCount-1-i)

			codes = append(codes, ast.TypeJenChain(st.Others, Var().Id(varName+p).Op(ptr)))
		}
		receiverName = varName + strings.Repeat("p", ptrCount-1)
	}

	codes = append(codes,
		List(Id("offset"), Err()).Op("=").Id(createFuncName(decodeFuncName, ast.StructName, ast.ImportPath)).Call(Op("&").Id(receiverName), Id(ptn.IdDecoder), Id("offset")),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
	)

	codes = st.createDecodeSetVarPattern(ptrCount, varName, fieldName /*setVarName*/, isParentTypeArrayOrMap, codes)

	// array or map
	if isParentTypeArrayOrMap {
		return codes
	}

	return []Code{Block(codes...)}
}

func isRootField(name string) bool {
	return strings.Contains(name, ".")
}
