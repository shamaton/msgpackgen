package main

import (
	"fmt"
	"go/types"
	"path/filepath"
	"reflect"
	"strings"

	. "github.com/dave/jennifer/jen"
)

var analyzedStructs []analyzedStruct

const (
	pkTop = "github.com/shamaton/msgpackgen"
	pkEnc = "github.com/shamaton/msgpackgen/enc"
	pkDec = "github.com/shamaton/msgpackgen/dec"

	idEncoder = "encoder"
	idDecoder = "decoder"
)

func main() {
	fileName := "msgpackgen_struct.go"
	path, err := filepath.Abs(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(findStructs(path))
	GenerateCode()
}

func GenerateCode() {

	f := NewFilePath("msgpackgen/resolver")

	f.Func().Id("init").Params().Block(
		Qual(pkTop, "SetResolver").Call(Id("encode"), Id("decode")),
	)

	decodeTopTemplate("decode", f).Block(
		// todo : Qual
		If(Qual(pkTop, "StructAsArray").Call()).Block(
			Return(Id("decodeAsArray").Call(Id("data"), Id("i"))),
		).Else().Block(
			Return(Id("decodeAsMap").Call(Id("data"), Id("i"))),
		),
	)

	decodeTopTemplate("decodeAsArray", f).Block(
		Switch(Id("v").Op(":=").Id("i").Assert(Type())).Block(
			cases()...,
		),
		Return(False(), Nil()),
	)

	decodeTopTemplate("decodeAsMap", f).Block(
		Return(False(), Nil()),
	)

	encodeTopTemplate("encode", f).Block(
		If(Qual(pkTop, "StructAsArray").Call()).Block(
			Return(Id("encodeAsArray").Call(Id("i"))),
		).Else().Block(
			Return(Id("encodeAsMap").Call(Id("i"))),
		),
	)

	encodeTopTemplate("encodeAsArray", f).Block(
		Return(Nil(), Nil()),
	)

	encodeTopTemplate("encodeAsMap", f).Block(
		Return(Nil(), Nil()),
	)

	for _, st := range analyzedStructs {
		calcFunction(st, f)
	}

	fmt.Printf("%#v", f)
}

func decodeTopTemplate(name string, f *File) *Statement {
	return f.Func().Id(name).Params(Id("data").Index().Byte(), Id("i").Interface()).Params(Bool(), Error())
}

func encodeTopTemplate(name string, f *File) *Statement {
	return f.Func().Id(name).Params(Id("i").Interface()).Params(Index().Byte(), Error())
}

func cases() []Code {
	var states []Code
	for _, v := range analyzedStructs {
		states = append(states, Case(Id("*"+v.Name)).Block(
			// 			_, err := decodestructTest(v, dec.NewDecoder(data), 0)
			List(Id("_"), Err()).Op(":=").Id("decode"+v.Name).Call(Id("v"), Qual(pkDec, "NewDecoder").Call(Id("data")), Id("0")),
			Return(True(), Err())))
	}
	return states
}

func calcFunction(st analyzedStruct, f *File) {
	v := "v"

	calcArraySizeCodes := make([]Code, 0)
	calcArraySizeCodes = append(calcArraySizeCodes, Id("size").Op(":=").Lit(0))
	calcArraySizeCodes = append(calcArraySizeCodes, Block(addSizePattern2("CalcStructHeader", Lit(len(st.Fields)))...))

	calcMapSizeCodes := make([]Code, 0)
	calcMapSizeCodes = append(calcMapSizeCodes, Id("size").Op(":=").Lit(0))
	calcMapSizeCodes = append(calcMapSizeCodes, Block(addSizePattern2("CalcStructHeader", Lit(len(st.Fields)))...))

	encArrayCodes := make([]Code, 0)
	encArrayCodes = append(encArrayCodes, Var().Err().Error())
	encArrayCodes = append(encArrayCodes, Id("offset").Op("=").Id(idEncoder).Dot("WriteStructHeader").Call(Lit(len(st.Fields)), Id("offset")))

	encMapCodes := make([]Code, 0)
	encMapCodes = append(encMapCodes, Var().Err().Error())
	encMapCodes = append(encMapCodes, Id("offset").Op("=").Id(idEncoder).Dot("WriteStructHeader").Call(Lit(len(st.Fields)), Id("offset")))

	decArrayCodes := make([]Code, 0)
	decArrayCodes = append(decArrayCodes, List(Id("offset"), Err()).Op(":=").Id(idDecoder).Dot("CheckStructHeader").Call(Lit(len(st.Fields)), Lit(0)))
	decArrayCodes = append(decArrayCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))

	decMapCodeSwitchCases := make([]Code, 0)

	for _, field := range st.Fields {
		calcMapSizeCodes = append(calcMapSizeCodes, addSizePattern1("CalcString", Lit(field.Name)))
		encMapCodes = append(encMapCodes, encPattern1("WriteString", Lit(field.Name), Id("offset")))

		cArray, cMap, eArray, eMap, dArray, dMap, _ := createFieldCode(field.Type, field.Name, true)
		calcArraySizeCodes = append(calcArraySizeCodes, cArray...)

		calcMapSizeCodes = append(calcMapSizeCodes, cMap...)

		encArrayCodes = append(encArrayCodes, eArray...)
		encMapCodes = append(encMapCodes, eMap...)

		decArrayCodes = append(decArrayCodes, dArray...)

		decMapCodeSwitchCases = append(decMapCodeSwitchCases, Case(Lit(field.Name)).Block(dMap...))

	}

	decMapCodeSwitchCases = append(decMapCodeSwitchCases, Default().Block(Id("offset").Op("=").Id(idDecoder).Dot("JumpOffset").Call(Id("offset"))))

	decMapCodes := make([]Code, 0)
	decMapCodes = append(decMapCodes, List(Id("offset"), Err()).Op(":=").Id(idDecoder).Dot("CheckStructHeader").Call(Lit(len(st.Fields)), Lit(0)))
	decMapCodes = append(decMapCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	decMapCodes = append(decMapCodes, Id("dataLen").Op(":=").Id(idDecoder).Dot("Len").Call())
	decMapCodes = append(decMapCodes, For(Id("offset").Op("<").Id("dataLen").Block(
		Var().Id("s").String(),
		List(Id("s"), Id("offset"), Err()).Op("=").Id(idDecoder).Dot("AsString").Call(Id("offset")),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Switch(Id("s")).Block(
			decMapCodeSwitchCases...,
		),
	)))

	f.Func().Id("calcArraySize"+st.Name).Params(Id(v).Qual(st.PackageName, st.Name), Id(idEncoder).Op("*").Qual(pkEnc, "Encoder")).Params(Int(), Error()).Block(
		calcArraySizeCodes...,
	)

	f.Func().Id("calcMapSize"+st.Name).Params(Id(v).Qual(st.PackageName, st.Name), Id(idEncoder).Op("*").Qual(pkEnc, "Encoder")).Params(Int(), Error()).Block(
		calcMapSizeCodes...,
	)

	f.Func().Id("encodeArray"+st.Name).Params(Id(v).Qual(st.PackageName, st.Name), Id(idEncoder).Op("*").Qual(pkEnc, "Encoder"), Id("offset").Int()).Params(Index().Byte(), Int(), Error()).Block(
		encArrayCodes...,
	)

	f.Func().Id("encodeMap"+st.Name).Params(Id(v).Qual(st.PackageName, st.Name), Id(idEncoder).Op("*").Qual(pkEnc, "Encoder"), Id("offset").Int()).Params(Index().Byte(), Int(), Error()).Block(
		encMapCodes...,
	)

	f.Func().Id("decodeArray"+st.Name).Params(Id(v).Op("*").Qual(st.PackageName, st.Name), Id(idDecoder).Op("*").Qual(pkDec, "Decoder"), Id("offset").Int()).Params(Int(), Error()).Block(
		decArrayCodes...,
	)

	f.Func().Id("decodeMap"+st.Name).Params(Id(v).Op("*").Qual(st.PackageName, st.Name), Id(idDecoder).Op("*").Qual(pkDec, "Decoder"), Id("offset").Int()).Params(Int(), Error()).Block(
		decMapCodes...,
	)
}

func createFieldCode(fieldType types.Type, fieldName string, isRoot bool) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {

	switch reflect.TypeOf(fieldType) {
	case reflect.TypeOf(&types.Basic{}):
		fmt.Println("basic", fieldName, fieldType)
		return createBasicCode(fieldType, fieldName, isRoot)

	case reflect.TypeOf(&types.Slice{}):
		fmt.Println("slice", fieldName, fieldType)
		child := fieldType.(*types.Slice).Elem()

		name, childName := "", ""
		if isRoot {
			name = "v." + fieldName
			childName = "vv"
		} else {
			childName = fieldName + "v"
		}

		ca, _, _, _, _, _, _ := createFieldCode(child, childName, false)

		statements := addSizePattern2("CalcSliceLength", Id(fmt.Sprintf("len(%s)", name)))
		statements = append(statements, For(List(Id("_"), Id(childName)).Op(":=").Range().Id(name)).Block(
			ca...,
		))

		cArray = append(cArray, If(Id(name).Op("!=").Nil()).Block(
			statements...,
		).Else().Block(
			addSizePattern1("CalcNil"),
		))

	case reflect.TypeOf(&types.Map{}):
		mp := fieldType.(*types.Map)
		key := mp.Key()
		value := mp.Elem()
		fmt.Println("map", fieldName, fieldType)
		fmt.Println(key, value)

		name, childKey, childValue := "", "", ""
		if isRoot {
			name = "v." + fieldName
			childKey = "kk"
			childValue = "vv"
		} else {
			childKey = fieldName + "k"
			childValue = fieldName + "v"
		}

		caKey, _, _, _, _, _, _ := createFieldCode(key, childKey, false)
		caValue, _, _, _, _, _, _ := createFieldCode(value, childValue, false)

		statements := addSizePattern2("CalcMapLength", Id(fmt.Sprintf("len(%s)", name)))
		statements = append(statements, For(List(Id(childKey), Id(childValue)).Op(":=").Range().Id(name)).Block(
			append(caKey, caValue...)...,
		))

		cArray = append(cArray, If(Id(name).Op("!=").Nil()).Block(
			statements...,
		).Else().Block(
			addSizePattern1("CalcNil"),
		))

	case reflect.TypeOf(&types.Named{}):
		if fieldType.String() == "time.Time" {

			fieldValue := Id(fieldName)
			if isRoot {
				fieldValue = Id("v").Dot(fieldName)
			}
			cArray = append(cArray, addSizePattern1("CalcTime", fieldValue))
			eArray = append(eArray, encPattern1("WriteTime", fieldValue, Id("offset")))

			cMap = append(cMap, addSizePattern1("CalcTime", fieldValue))
			eMap = append(eMap, encPattern1("WriteTime", fieldValue, Id("offset")))

			dArray = append(dArray, decodeBasicPattern(fieldType, fieldName, "offset", "time.Time", "AsDateTime", isRoot)...)
			dMap = append(dMap, decodeBasicPattern(fieldType, fieldName, "offset", "time.Time", "AsDateTime", isRoot)...)
		} else if strings.HasPrefix(fieldType.String(), "github") {
			// todo : 対象のパッケージかどうかをちゃんと判断する
		}
		//if (types.Identical(types.Struct{}, fieldType))
		fmt.Println("named", fieldName, fieldType)

	default:
		// todo : error

	}

	return cArray, cMap, eArray, eMap, dArray, dMap, err
}

func createBasicCode(fieldType types.Type, fieldName string, isRoot bool) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {
	isType := func(kind types.BasicKind) bool {
		return types.Identical(types.Typ[kind], fieldType)
	}

	offset := "offset"

	fieldValue := Id(fieldName)
	if isRoot {
		fieldValue = Id("v").Dot(fieldName)
	}

	switch {
	case isType(types.Int):
		cArray = append(cArray, addSizePattern1("CalcInt", Id("int64").Call(fieldValue)))
		eArray = append(eArray, encPattern1("WriteInt", Id("int64").Call(fieldValue), Id(offset)))

		cMap = append(cMap, addSizePattern1("CalcInt", Id("int64").Call(fieldValue)))
		eMap = append(eMap, encPattern1("WriteInt", Id("int64").Call(fieldValue), Id(offset)))

		dArray = append(dArray, decodeBasicPattern(fieldType, fieldName, offset, "int64", "AsInt", isRoot)...)
		dMap = append(dMap, decodeBasicPattern(fieldType, fieldName, offset, "int64", "AsInt", isRoot)...)

	case isType(types.Uint):
		cArray = append(cArray, addSizePattern1("CalcUint", Id("uint64").Call(fieldValue)))

	case isType(types.String):
		cArray = append(cArray, addSizePattern1("CalcString", fieldValue))

	case isType(types.Float64):
		cArray = append(cArray, addSizePattern1("CalcFloat64", fieldValue))
	default:
		// todo error

	}
	return cArray, cMap, eArray, eMap, dArray, dMap, err
}

func addSizePattern1(funcName string, params ...Code) Code {
	return Id("size").Op("+=").Id(idEncoder).Dot(funcName).Call(params...)
}

func addSizePattern2(funcName string, params ...Code) []Code {
	return []Code{
		List(Id("s"), Err()).Op(":=").Id(idEncoder).Dot(funcName).Call(params...),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id("size").Op("+=").Id("s"),
	}

}

func encPattern1(funcName string, params ...Code) Code {
	return Id("offset").Op("=").Id(idEncoder).Dot(funcName).Call(params...)
}

func decodeBasicPattern(fieldType types.Type, fieldName, offsetName, varTypeName, decoderFuncName string, isRoot bool) []Code {

	vName := "vv"
	setName := "v." + fieldName
	if !isRoot {
		vName = fieldName + "v"
		setName = fieldName
	}

	return []Code{Block(
		Var().Id(vName).Id(varTypeName),
		List(Id(vName), Id(offsetName), Err()).Op("=").Id(idDecoder).Dot(decoderFuncName).Call(Id(offsetName)),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id(setName).Op("=").Id(fieldType.String()).Call(Id(vName)),
	)}
}
