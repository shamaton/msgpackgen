package generator

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

func (as *Structure) CalcArraySizeFuncName() string {
	return as.createFuncName("calcArraySize")
}

func (as *Structure) CalcMapSizeFuncName() string {
	return as.createFuncName("calcMapSize")
}

func (as *Structure) EncodeArrayFuncName() string {
	return as.createFuncName("encodeArray")
}

func (as *Structure) EncodeMapFuncName() string {
	return as.createFuncName("encodeMap")
}

func (as *Structure) DecodeArrayFuncName() string {
	return as.createFuncName("decodeArray")
}

func (as *Structure) DecodeMapFuncName() string {
	return as.createFuncName("decodeMap")
}

func (as *Structure) createFuncName(prefix string) string {
	return createFuncName(prefix, as.Name, as.ImportPath)
}

func createFuncName(prefix, name, importPath string) string {
	suffix := fmt.Sprintf("%x", sha256.Sum256([]byte(importPath)))
	return ptn.PrivateFuncName(fmt.Sprintf("%s%s_%s", prefix, name, suffix))
}

func (as *Structure) CreateCode(f *File) {
	v := "v"

	calcStruct, encStructArray, encStructMap := as.CreateStructCode(len(as.Fields))

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
	decArrayCodes = append(decArrayCodes, List(Id("offset"), Err()).Op(":=").Id(IdDecoder).Dot("CheckStructHeader").Call(Lit(len(as.Fields)), Id("offset")))
	decArrayCodes = append(decArrayCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))

	decMapCodeSwitchCases := make([]Code, 0)

	for _, field := range as.Fields {
		fieldName := "v." + field.Name

		calcKeyStringCode, writeKeyStringCode := as.CreateKeyStringCode(field.Tag)
		calcMapSizeCodes = append(calcMapSizeCodes, calcKeyStringCode)
		encMapCodes = append(encMapCodes, writeKeyStringCode)

		cArray, cMap, eArray, eMap, dArray, dMap, _ := as.createFieldCode(field.Node, fieldName, fieldName)
		calcArraySizeCodes = append(calcArraySizeCodes, cArray...)

		calcMapSizeCodes = append(calcMapSizeCodes, cMap...)

		encArrayCodes = append(encArrayCodes, eArray...)
		encMapCodes = append(encMapCodes, eMap...)

		decArrayCodes = append(decArrayCodes, dArray...)

		decMapCodeSwitchCases = append(decMapCodeSwitchCases, Case(Lit(field.Tag)).Block(dMap...))

	}

	decMapCodeSwitchCases = append(decMapCodeSwitchCases, Default().Block(Id("offset").Op("=").Id(IdDecoder).Dot("JumpOffset").Call(Id("offset"))))

	decMapCodes := make([]Code, 0)
	decMapCodes = append(decMapCodes, List(Id("offset"), Err()).Op(":=").Id(IdDecoder).Dot("CheckStructHeader").Call(Lit(len(as.Fields)), Id("offset")))
	decMapCodes = append(decMapCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	decMapCodes = append(decMapCodes, Id("dataLen").Op(":=").Id(IdDecoder).Dot("Len").Call())
	decMapCodes = append(decMapCodes, For(Id("offset").Op("<").Id("dataLen").Block(
		Var().Id("s").String(),
		List(Id("s"), Id("offset"), Err()).Op("=").Id(IdDecoder).Dot("AsString").Call(Id("offset")),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Switch(Id("s")).Block(
			decMapCodeSwitchCases...,
		),
	)))

	var firstEncParam, firstDecParam *Statement
	if as.NoUseQual {
		firstEncParam = Id(v).Id(as.Name)
		firstDecParam = Id(v).Op("*").Id(as.Name)
	} else {
		firstEncParam = Id(v).Qual(as.ImportPath, as.Name)
		firstDecParam = Id(v).Op("*").Qual(as.ImportPath, as.Name)
	}

	f.Comment(fmt.Sprintf("// calculate size from %s.%s\n", as.ImportPath, as.Name)).
		Func().Id(as.CalcArraySizeFuncName()).Params(firstEncParam, Id(IdEncoder).Op("*").Qual(PkEnc, "Encoder")).Params(Int(), Error()).Block(
		append(calcArraySizeCodes, Return(Id("size"), Nil()))...,
	)

	f.Comment(fmt.Sprintf("// calculate size from %s.%s\n", as.ImportPath, as.Name)).
		Func().Id(as.CalcMapSizeFuncName()).Params(firstEncParam, Id(IdEncoder).Op("*").Qual(PkEnc, "Encoder")).Params(Int(), Error()).Block(
		append(calcMapSizeCodes, Return(Id("size"), Nil()))...,
	)

	f.Comment(fmt.Sprintf("// encode from %s.%s\n", as.ImportPath, as.Name)).
		Func().Id(as.EncodeArrayFuncName()).Params(firstEncParam, Id(IdEncoder).Op("*").Qual(PkEnc, "Encoder"), Id("offset").Int()).Params(Index().Byte(), Int(), Error()).Block(
		append(encArrayCodes, Return(Id(IdEncoder).Dot("EncodedBytes").Call(), Id("offset"), Err()))...,
	)

	f.Comment(fmt.Sprintf("// encode from %s.%s\n", as.ImportPath, as.Name)).
		Func().Id(as.EncodeMapFuncName()).Params(firstEncParam, Id(IdEncoder).Op("*").Qual(PkEnc, "Encoder"), Id("offset").Int()).Params(Index().Byte(), Int(), Error()).Block(
		append(encMapCodes, Return(Id(IdEncoder).Dot("EncodedBytes").Call(), Id("offset"), Err()))...,
	)

	f.Comment(fmt.Sprintf("// decode to %s.%s\n", as.ImportPath, as.Name)).
		Func().Id(as.DecodeArrayFuncName()).Params(firstDecParam, Id(IdDecoder).Op("*").Qual(PkDec, "Decoder"), Id("offset").Int()).Params(Int(), Error()).Block(
		append(decArrayCodes, Return(Id("offset"), Err()))...,
	)

	f.Comment(fmt.Sprintf("// decode to %s.%s\n", as.ImportPath, as.Name)).
		Func().Id(as.DecodeMapFuncName()).Params(firstDecParam, Id(IdDecoder).Op("*").Qual(PkDec, "Decoder"), Id("offset").Int()).Params(Int(), Error()).Block(

		append(decMapCodes, Return(Id("offset"), Err()))...,
	)
}

func (as *Structure) CreateKeyStringCode(v string) (Code, Code) {
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

	return Id("size").Op("+=").Id(IdEncoder).Dot("CalcString" + suffix).Call(Lit(l)),
		Id("offset").Op("=").Id(IdEncoder).Dot("WriteString"+suffix).Call(Lit(v), Lit(l), Id("offset"))
}

func (as *Structure) CreateStructCode(fieldNum int) (Code, Code, Code) {

	suffix := ""
	if fieldNum <= 0x0f {
		suffix = "Fix"
	} else if fieldNum <= math.MaxUint16 {
		suffix = "16"
	} else if uint(fieldNum) <= math.MaxUint32 {
		suffix = "32"
	}

	return Id("size").Op("+=").Id(IdEncoder).Dot("CalcStructHeader" + suffix).Call(Lit(fieldNum)),
		Id("offset").Op("=").Id(IdEncoder).Dot(" WriteStructHeader"+suffix+"AsArray").Call(Lit(fieldNum), Id("offset")),
		Id("offset").Op("=").Id(IdEncoder).Dot(" WriteStructHeader"+suffix+"AsMap").Call(Lit(fieldNum), Id("offset"))
}

func (as *Structure) createFieldCode(ast *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {

	switch {
	case ast.IsIdentical():
		return as.createBasicCode(ast, encodeFieldName, decodeFieldName)

	case ast.IsSlice():
		return as.createSliceCode(ast, encodeFieldName, decodeFieldName)

	case ast.IsArray():
		return as.createArrayCode(ast, encodeFieldName, decodeFieldName)

	case ast.IsMap():
		return as.createMapCode(ast, encodeFieldName, decodeFieldName)

	case ast.IsPointer():
		return as.createPointerCode(ast, encodeFieldName, decodeFieldName)

	case ast.IsStruct():

		ptrOp := ""
		node := ast
		for {
			if node.HasParent() && node.Parent.IsPointer() {
				ptrOp += "*"
				node = node.Parent
			} else {
				break
			}
		}

		fieldValue := Id(encodeFieldName)

		// todo : ポインタでの動作検証
		if ast.ImportPath == "time" {
			cArray = append(cArray, as.addSizePattern1("CalcTime", fieldValue))
			eArray = append(eArray, as.encPattern1("WriteTime", fieldValue, Id("offset")))

			cMap = append(cMap, as.addSizePattern1("CalcTime", fieldValue))
			eMap = append(eMap, as.encPattern1("WriteTime", fieldValue, Id("offset")))

			dArray = append(dArray, as.decodeBasicPattern(ast, decodeFieldName, "offset", "AsDateTime")...)
			dMap = append(dMap, as.decodeBasicPattern(ast, decodeFieldName, "offset", "AsDateTime")...)
		} else {
			// todo : 対象のパッケージかどうかをちゃんと判断する
			cArray, cMap, eArray, eMap, dArray, dMap = as.createNamedCode(encodeFieldName, decodeFieldName, ast)
		}

	default:
		// todo : error

	}

	return cArray, cMap, eArray, eMap, dArray, dMap, err
}

func (as *Structure) createPointerCode(ast *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {

	encodeChildName := encodeFieldName + "p"
	if isRootField(encodeFieldName) {
		encodeChildName = "vp"
	}

	ca, _, ea, _, da, _, _ := as.createFieldCode(ast.Elm(), encodeChildName, decodeFieldName)

	cArray = make([]Code, 0)
	cArray = append(cArray, If(Id(encodeFieldName).Op("!=").Nil()).Block(
		append([]Code{
			Id(encodeChildName).Op(":=").Op("*").Id(encodeFieldName),
		}, ca...)...,
	).Else().Block(
		Id("size").Op("+=").Id(IdEncoder).Dot("CalcNil").Call(),
	))

	eArray = make([]Code, 0)
	eArray = append(eArray, If(Id(encodeFieldName).Op("!=").Nil()).Block(
		append([]Code{
			Id(encodeChildName).Op(":=").Op("*").Id(encodeFieldName),
		}, ea...)...,
	).Else().Block(
		Id("offset").Op("=").Id(IdEncoder).Dot("WriteNil").Call(Id("offset")),
	))

	// todo : ようかくにん、重複コードをスキップ
	isParentPointer := ast.HasParent() && ast.Parent.IsPointer()
	if isParentPointer {
		dArray = da
	} else {
		dArray = make([]Code, 0)
		dArray = append(dArray, If(Op("!").Id(IdDecoder).Dot("IsCodeNil").Call(Id("offset"))).Block(
			da...,
		).Else().Block(
			Id("offset").Op("++"),
		))
	}

	return cArray, cArray, eArray, eArray, dArray, dArray, err
}

func (as *Structure) createMapCode(ast *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {

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

	caKey, _, eaKey, _, daKey, _, _ := as.createFieldCode(key, encodeChildKey, decodeChildKey)
	caValue, _, eaValue, _, daValue, _, _ := as.createFieldCode(value, encodeChildValue, decodeChildValue)

	calcCodes := as.addSizePattern2("CalcMapLength", Len(Id(encodeFieldName)))
	calcCodes = append(calcCodes, For(List(Id(encodeChildKey), Id(encodeChildValue)).Op(":=").Range().Id(encodeFieldName)).Block(
		append(caKey, caValue...)...,
	))

	cArray = append(cArray, If(Id(encodeFieldName).Op("!=").Nil()).Block(
		calcCodes...,
	).Else().Block(
		as.addSizePattern1("CalcNil"),
	))

	encCodes := make([]Code, 0)
	encCodes = append(encCodes, Id("offset").Op("=").Id(IdEncoder).Dot("WriteMapLength").Call(Len(Id(encodeFieldName)), Id("offset")))
	encCodes = append(encCodes, For(List(Id(encodeChildKey), Id(encodeChildValue)).Op(":=").Range().Id(encodeFieldName)).Block(
		append(eaKey, eaValue...)...,
	))

	eArray = append(eArray, If(Id(encodeFieldName).Op("!=").Nil()).Block(
		encCodes...,
	).Else().Block(
		Id("offset").Op("=").Id(IdEncoder).Dot("WriteNil").Call(Id("offset")),
	))

	decCodes := make([]Code, 0)
	decCodes = append(decCodes, ast.TypeJenChain(as.Others, Var().Id(decodeChildValue)))
	decCodes = append(decCodes, Var().Id(decodeChildValue+"l").Int())
	decCodes = append(decCodes, List(Id(decodeChildValue+"l"), Id("offset"), Err()).Op("=").Id(IdDecoder).Dot("MapLength").Call(Id("offset")))
	decCodes = append(decCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	decCodes = append(decCodes, Id(decodeChildValue).Op("=").Make(ast.TypeJenChain(as.Others), Id(decodeChildValue+"l")))

	da := []Code{ast.Key.TypeJenChain(as.Others, Var().Id(decodeChildKey+"v"))}
	da = append(da, daKey...)
	da = append(da, ast.Value.TypeJenChain(as.Others, Var().Id(decodeChildValue+"v")))
	da = append(da, daValue...)
	da = append(da, Id(decodeChildValue).Index(Id(decodeChildKey+"v")).Op("=").Id(decodeChildValue+"v"))

	decCodes = append(decCodes, For(Id(decodeChildValue+"i").Op(":=").Lit(0).Op(";").Id(decodeChildValue+"i").Op("<").Id(decodeChildValue+"l").Op(";").Id(decodeChildValue+"i").Op("++")).Block(
		da...,
	))

	// todo : 不要なコードがあるはず
	ptrOp := ""
	andOp := ""
	prtCount := 0
	node := ast
	for {
		if node.HasParent() && node.Parent.IsPointer() {
			ptrOp += "*"
			andOp += "&"
			prtCount++
			node = node.Parent
		} else {
			break
		}
	}

	name := decodeChildValue
	if prtCount > 0 {
		andOp = "&"
	}
	for i := 0; i < prtCount-1; i++ {
		n := "_" + name
		decCodes = append(decCodes, Id(n).Op(":=").Op("&").Id(name))
		name = n
	}

	decCodes = append(decCodes, Id(decodeFieldName).Op("=").Op(andOp).Id(name))

	// todo : ようかくにん、重複コードをスキップ
	if ast.HasParent() && ast.Parent.IsPointer() {
		dArray = decCodes
	} else {

		dArray = append(dArray, If(Op("!").Id(IdDecoder).Dot("IsCodeNil").Call(Id("offset"))).Block(
			decCodes...,
		).Else().Block(
			Id("offset").Op("++"),
		))
	}

	return cArray, cArray, eArray, eArray, dArray, dArray, nil
}

func (as *Structure) createSliceCode(ast *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {

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

	ca, _, ea, _, da, _, _ := as.createFieldCode(ast.Elm(), encodeChildName, decodeChildName)
	isChildByte := ast.Elm().IsIdentical() && ast.Elm().IdenticalName == "byte"

	calcCodes := as.addSizePattern2("CalcSliceLength", Len( /*Op(ptrOp).*/ Id(encodeFieldName)), Lit(isChildByte))
	calcCodes = append(calcCodes, For(List(Id("_"), Id(encodeChildName)).Op(":=").Range(). /*Op(ptrOp).*/ Id(encodeFieldName)).Block(
		ca...,
	))

	cArray = append(cArray, If( /*Op(ptrOp).*/ Id(encodeFieldName).Op("!=").Nil()).Block(
		calcCodes...,
	).Else().Block(
		as.addSizePattern1("CalcNil"),
	))

	encCodes := make([]Code, 0)
	encCodes = append(encCodes, Id("offset").Op("=").Id(IdEncoder).Dot("WriteSliceLength").Call(Len( /*Op(ptrOp).*/ Id(encodeFieldName)), Id("offset"), Lit(isChildByte)))
	encCodes = append(encCodes, For(List(Id("_"), Id(encodeChildName)).Op(":=").Range(). /*Op(ptrOp).*/ Id(encodeFieldName)).Block(
		ea...,
	))

	eArray = append(eArray, If( /*Op(ptrOp).*/ Id(encodeFieldName).Op("!=").Nil()).Block(
		encCodes...,
	).Else().Block(
		Id("offset").Op("=").Id(IdEncoder).Dot("WriteNil").Call(Id("offset")),
	))

	decCodes := make([]Code, 0)
	decCodes = append(decCodes, ast.TypeJenChain(as.Others, Var().Id(decodeChildName)))
	decCodes = append(decCodes, Var().Id(decodeChildLengthName).Int())
	decCodes = append(decCodes, List(Id(decodeChildLengthName), Id("offset"), Err()).Op("=").Id(IdDecoder).Dot("SliceLength").Call(Id("offset")))
	decCodes = append(decCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	decCodes = append(decCodes, Id(decodeChildName).Op("=").Make(ast.TypeJenChain(as.Others), Id(decodeChildLengthName)))

	da = append([]Code{ast.Elm().TypeJenChain(as.Others, Var().Id(decodeChildChildName))}, da...)
	da = append(da, Id(decodeChildName).Index(Id(decodeChildIndexName)).Op("=").Id(decodeChildChildName))

	decCodes = append(decCodes, For(Id(decodeChildIndexName).Op(":=").Range().Id(decodeChildName)).Block(
		da...,
	))

	// todo : 不要なコードがあるはず
	ptrOp := ""
	andOp := ""
	prtCount := 0
	node := ast
	for {
		if node.HasParent() && node.Parent.IsPointer() {
			ptrOp += "*"
			andOp += "&"
			prtCount++
			node = node.Parent
		} else {
			break
		}
	}

	name := decodeChildName
	if prtCount > 0 {
		andOp = "&"
	}
	for i := 0; i < prtCount-1; i++ {
		n := "_" + name
		decCodes = append(decCodes, Id(n).Op(":=").Op("&").Id(name))
		name = n
	}

	decCodes = append(decCodes, Id(decodeFieldName).Op("=").Op(andOp).Id(name))

	// todo : ようかくにん、重複コードをスキップ
	if ast.HasParent() && ast.Parent.IsPointer() {
		dArray = decCodes
	} else {

		dArray = append(dArray, If(Op("!").Id(IdDecoder).Dot("IsCodeNil").Call(Id("offset"))).Block(
			decCodes...,
		).Else().Block(
			Id("offset").Op("++"),
		))
	}

	return cArray, cArray, eArray, eArray, dArray, dArray, nil
}

func (as *Structure) createArrayCode(ast *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {

	encodeChildName := encodeFieldName + "v"
	if isRootField(encodeFieldName) {
		encodeChildName = "vv"
	}

	decodeChildName := decodeFieldName + "v"
	if isRootField(decodeFieldName) {
		decodeChildName = "vv"
	}

	ca, _, ea, _, da, _, _ := as.createFieldCode(ast.Elm(), encodeChildName, decodeChildName)
	isChildByte := ast.Elm().IsIdentical() && ast.Elm().IdenticalName == "byte"

	calcCodes := as.addSizePattern2("CalcSliceLength", Len( /*Op(ptrOp).*/ Id(encodeFieldName)), Lit(isChildByte))
	calcCodes = append(calcCodes, For(List(Id("_"), Id(encodeChildName)).Op(":=").Range(). /*Op(ptrOp).*/ Id(encodeFieldName)).Block(
		ca...,
	))

	cArray = append(cArray /* If(Op(ptrOp).Id(name).Op("!=").Nil()).*/, Block(
		calcCodes...,
	), /*.Else().Block(
		as.addSizePattern1("CalcNil"),
	)*/)

	encCodes := make([]Code, 0)
	encCodes = append(encCodes, Id("offset").Op("=").Id(IdEncoder).Dot("WriteSliceLength").Call(Len( /*Op(ptrOp).*/ Id(encodeFieldName)), Id("offset"), Lit(isChildByte)))
	encCodes = append(encCodes, For(List(Id("_"), Id(encodeChildName)).Op(":=").Range(). /*Op(ptrOp).*/ Id(encodeFieldName)).Block(
		ea...,
	))

	eArray = append(eArray /*If(Op(ptrOp).Id(name).Op("!=").Nil()).*/, Block(
		encCodes...,
	), /*.Else().Block(
		Id("offset").Op("=").Id(IdEncoder).Dot("WriteNil").Call(Id("offset")),
	)*/)

	decCodes := make([]Code, 0)
	decCodes = append(decCodes, ast.TypeJenChain(as.Others, Var().Id(decodeChildName)))
	decCodes = append(decCodes, Var().Id(decodeChildName+"l").Int())
	decCodes = append(decCodes, List(Id(decodeChildName+"l"), Id("offset"), Err()).Op("=").Id(IdDecoder).Dot("SliceLength").Call(Id("offset")))
	decCodes = append(decCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	decCodes = append(decCodes, If(Id(decodeChildName+"l").Op(">").Id(fmt.Sprint(ast.ArrayLen))).Block(
		Return(Lit(0), Qual("fmt", "Errorf").Call(Lit("length size(%d) is over array size(%d)"), Id(decodeChildName+"l"), Id(fmt.Sprint(ast.ArrayLen)))),
	))

	da = append([]Code{ast.Elm().TypeJenChain(as.Others, Var().Id(decodeChildName+"v"))}, da...)
	da = append(da, Id(decodeChildName).Index(Id(decodeChildName+"i")).Op("=").Id(decodeChildName+"v"))

	decCodes = append(decCodes, For(Id(decodeChildName+"i").Op(":=").Range().Id(decodeChildName).Index(Id(":"+decodeChildName+"l"))).Block(
		da...,
	))

	// todo : 不要なコードがあるはず
	ptrOp := ""
	andOp := ""
	prtCount := 0
	node := ast
	for {
		if node.HasParent() && node.Parent.IsPointer() {
			ptrOp += "*"
			andOp += "&"
			prtCount++
			node = node.Parent
		} else {
			break
		}
	}

	name := decodeChildName
	if prtCount > 0 {
		andOp = "&"
	}
	for i := 0; i < prtCount-1; i++ {
		n := "_" + name
		decCodes = append(decCodes, Id(n).Op(":=").Op("&").Id(name))
		name = n
	}

	decCodes = append(decCodes, Id(decodeFieldName).Op("=").Op(andOp).Id(name))

	// todo : ようかくにん、重複コードをスキップ
	if ast.HasParent() && ast.Parent.IsPointer() {
		dArray = decCodes
	} else {

		dArray = append(dArray, If(Op("!").Id(IdDecoder).Dot("IsCodeNil").Call(Id("offset"))).Block(
			decCodes...,
		).Else().Block(
			Id("offset").Op("++"),
		))
	}

	return cArray, cArray, eArray, eArray, dArray, dArray, nil
}

func (as *Structure) createBasicCode(ast *Node, encodeFieldName, decodeFieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {

	funcSuffix := strings.Title(ast.IdenticalName)

	cArray = append(cArray, as.addSizePattern1("Calc"+funcSuffix, Id(encodeFieldName)))
	eArray = append(eArray, as.encPattern1("Write"+funcSuffix, Id(encodeFieldName), Id("offset")))

	cMap = append(cMap, as.addSizePattern1("Calc"+funcSuffix, Id(encodeFieldName)))
	eMap = append(eMap, as.encPattern1("Write"+funcSuffix, Id(encodeFieldName), Id("offset")))

	dArray = append(dArray, as.decodeBasicPattern(ast, decodeFieldName, "offset", "As"+funcSuffix)...)
	dMap = append(dMap, as.decodeBasicPattern(ast, decodeFieldName, "offset", "As"+funcSuffix)...)

	return cArray, cMap, eArray, eMap, dArray, dMap, err
}

func (as *Structure) addSizePattern1(funcName string, params ...Code) Code {
	return Id("size").Op("+=").Id(IdEncoder).Dot(funcName).Call(params...)
}

func (as *Structure) addSizePattern2(funcName string, params ...Code) []Code {
	return []Code{
		List(Id("s"), Err()).Op(":=").Id(IdEncoder).Dot(funcName).Call(params...),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id("size").Op("+=").Id("s"),
	}

}

func (as *Structure) encPattern1(funcName string, params ...Code) Code {
	return Id("offset").Op("=").Id(IdEncoder).Dot(funcName).Call(params...)
}

func isRootField(name string) bool {
	return strings.Contains(name, ".")
}

func (as *Structure) decodeBasicPattern(ast *Node, fieldName, offsetName, decoderFuncName string) []Code {

	varName := fieldName + "v"
	if isRootField(fieldName) {
		varName = "vv"
	}

	node := ast
	ptrCount := 0
	isParentTypeArrayOrMap := false

	for {
		if node.HasParent() {
			node = node.Parent
			if node.IsPointer() {
				ptrCount++
			} else if node.IsSlice() || node.IsArray() || node.IsMap() {
				isParentTypeArrayOrMap = true
				break
			} else {
				// todo : error or empty
			}
		} else {
			break
		}
	}

	codes := make([]Code, 0)
	recieverName := varName

	if ptrCount < 1 && !isParentTypeArrayOrMap {
		codes = append(codes, ast.TypeJenChain(as.Others, Var().Id(recieverName)))
	} else if isParentTypeArrayOrMap {

		for i := 0; i < ptrCount; i++ {
			p := strings.Repeat("p", i+1)
			kome := strings.Repeat("*", ptrCount-1-i)
			codes = append(codes, ast.TypeJenChain(as.Others, Var().Id(varName+p).Op(kome)))
		}
		recieverName = varName + strings.Repeat("p", ptrCount)
	} else {
		for i := 0; i < ptrCount; i++ {
			p := strings.Repeat("p", i)
			kome := strings.Repeat("*", ptrCount-1-i)
			codes = append(codes, ast.TypeJenChain(as.Others, Var().Id(varName+p).Op(kome)))
		}
		recieverName = varName + strings.Repeat("p", ptrCount-1)
	}

	codes = append(codes,
		List(Id(recieverName), Id(offsetName), Err()).Op("=").Id(IdDecoder).Dot(decoderFuncName).Call(Id(offsetName)),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
	)

	codes = as.createDecodeSetVarPattern(ptrCount, varName, fieldName /*setVarName*/, isParentTypeArrayOrMap, codes)

	// array or map
	if isParentTypeArrayOrMap {
		return codes
	}

	return []Code{Block(codes...)}
}

func (as *Structure) createDecodeSetVarPattern(ptrCount int, varName, setVarName string, isLastSkip bool, codes []Code) []Code {

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

func (as *Structure) createNamedCode(encodeFieldName, decodeFieldName string, ast *Node) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	sizeName := "size_" + encodeFieldName
	if isRootField(encodeFieldName) {
		sizeName = strings.ReplaceAll(sizeName, ".", "_")
	}

	cArray = []Code{
		List(Id(sizeName), Err()).
			Op(":=").
			Id(createFuncName("calcArraySize", ast.StructName, ast.ImportPath)).Call(Id(encodeFieldName), Id(IdEncoder)),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id("size").Op("+=").Id(sizeName),
	}

	cMap = []Code{
		List(Id(sizeName), Err()).
			Op(":=").
			Id(createFuncName("calcMapSize", ast.StructName, ast.ImportPath)).Call(Id(encodeFieldName), Id(IdEncoder)),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id("size").Op("+=").Id(sizeName),
	}

	eArray = []Code{
		List(Id("_"), Id("offset"), Err()).
			Op("=").
			Id(createFuncName("encodeArray", ast.StructName, ast.ImportPath)).Call(Id(encodeFieldName), Id(IdEncoder), Id("offset")),
		If(Err().Op("!=").Nil()).Block(
			Return(Nil(), Lit(0), Err()),
		),
	}

	eMap = []Code{
		List(Id("_"), Id("offset"), Err()).
			Op("=").
			Id(createFuncName("encodeMap", ast.StructName, ast.ImportPath)).Call(Id(encodeFieldName), Id(IdEncoder), Id("offset")),
		If(Err().Op("!=").Nil()).Block(
			Return(Nil(), Lit(0), Err()),
		),
	}

	dArray = append(dArray, as.decodeNamedPattern(ast, decodeFieldName, "decodeArray")...)
	dMap = append(dMap, as.decodeNamedPattern(ast, decodeFieldName, "decodeMap")...)

	return
}

func (as *Structure) decodeNamedPattern(ast *Node, fieldName, decodeFuncName string) []Code {

	varName := fieldName + "v"
	if isRootField(fieldName) {
		varName = "vv"
	}

	node := ast
	ptrCount := 0
	isParentTypeArrayOrMap := false

	for {
		if node.HasParent() {
			node = node.Parent
			if node.IsPointer() {
				ptrCount++
			} else if node.IsSlice() || node.IsArray() || node.IsMap() {
				isParentTypeArrayOrMap = true
				break
			} else {
				// todo : error or empty
			}
		} else {
			break
		}
	}

	codes := make([]Code, 0)
	receiverName := varName

	if ptrCount < 1 && !isParentTypeArrayOrMap {
		codes = append(codes, ast.TypeJenChain(as.Others, Var().Id(receiverName)))
	} else if isParentTypeArrayOrMap {

		for i := 0; i < ptrCount; i++ {
			p := strings.Repeat("p", i+1)
			kome := strings.Repeat("*", ptrCount-1-i)

			codes = append(codes, ast.TypeJenChain(as.Others, Var().Id(varName+p).Op(kome)))
		}
		receiverName = varName + strings.Repeat("p", ptrCount)
	} else {
		for i := 0; i < ptrCount; i++ {
			p := strings.Repeat("p", i)
			kome := strings.Repeat("*", ptrCount-1-i)

			codes = append(codes, ast.TypeJenChain(as.Others, Var().Id(varName+p).Op(kome)))
		}
		receiverName = varName + strings.Repeat("p", ptrCount-1)
	}

	codes = append(codes,
		List(Id("offset"), Err()).Op("=").Id(createFuncName(decodeFuncName, ast.StructName, ast.ImportPath)).Call(Op("&").Id(receiverName), Id(IdDecoder), Id("offset")),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
	)

	codes = as.createDecodeSetVarPattern(ptrCount, varName, fieldName /*setVarName*/, isParentTypeArrayOrMap, codes)

	// array or map
	if isParentTypeArrayOrMap {
		return codes
	}

	return []Code{Block(codes...)}
}
