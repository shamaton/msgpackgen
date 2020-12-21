package main

import (
	"fmt"
	"go/types"
	"reflect"

	. "github.com/dave/jennifer/jen"
)

func (as *analyzedStruct) calcFunction(f *File) {
	v := "v"

	calcArraySizeCodes := make([]Code, 0)
	calcArraySizeCodes = append(calcArraySizeCodes, Id("size").Op(":=").Lit(0))
	calcArraySizeCodes = append(calcArraySizeCodes, Block(as.addSizePattern2("CalcStructHeader", Lit(len(as.Fields)))...))

	calcMapSizeCodes := make([]Code, 0)
	calcMapSizeCodes = append(calcMapSizeCodes, Id("size").Op(":=").Lit(0))
	calcMapSizeCodes = append(calcMapSizeCodes, Block(as.addSizePattern2("CalcStructHeader", Lit(len(as.Fields)))...))

	encArrayCodes := make([]Code, 0)
	encArrayCodes = append(encArrayCodes, Var().Err().Error())
	encArrayCodes = append(encArrayCodes, Id("offset").Op("=").Id(idEncoder).Dot("WriteStructHeader").Call(Lit(len(as.Fields)), Id("offset")))

	encMapCodes := make([]Code, 0)
	encMapCodes = append(encMapCodes, Var().Err().Error())
	encMapCodes = append(encMapCodes, Id("offset").Op("=").Id(idEncoder).Dot("WriteStructHeader").Call(Lit(len(as.Fields)), Id("offset")))

	decArrayCodes := make([]Code, 0)
	decArrayCodes = append(decArrayCodes, List(Id("offset"), Err()).Op(":=").Id(idDecoder).Dot("CheckStructHeader").Call(Lit(len(as.Fields)), Lit(0)))
	decArrayCodes = append(decArrayCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))

	decMapCodeSwitchCases := make([]Code, 0)

	for _, field := range as.Fields {
		calcMapSizeCodes = append(calcMapSizeCodes, as.addSizePattern1("CalcString", Lit(field.Name)))
		encMapCodes = append(encMapCodes, as.encPattern1("WriteString", Lit(field.Name), Id("offset")))

		cArray, cMap, eArray, eMap, dArray, dMap, _ := as.createFieldCode(field.Type, field.Name, true)
		calcArraySizeCodes = append(calcArraySizeCodes, cArray...)

		calcMapSizeCodes = append(calcMapSizeCodes, cMap...)

		encArrayCodes = append(encArrayCodes, eArray...)
		encMapCodes = append(encMapCodes, eMap...)

		decArrayCodes = append(decArrayCodes, dArray...)

		decMapCodeSwitchCases = append(decMapCodeSwitchCases, Case(Lit(field.Name)).Block(dMap...))

	}

	decMapCodeSwitchCases = append(decMapCodeSwitchCases, Default().Block(Id("offset").Op("=").Id(idDecoder).Dot("JumpOffset").Call(Id("offset"))))

	decMapCodes := make([]Code, 0)
	decMapCodes = append(decMapCodes, List(Id("offset"), Err()).Op(":=").Id(idDecoder).Dot("CheckStructHeader").Call(Lit(len(as.Fields)), Lit(0)))
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

	f.Func().Id("calcArraySize"+as.Name).Params(Id(v).Qual(as.PackageName, as.Name), Id(idEncoder).Op("*").Qual(pkEnc, "Encoder")).Params(Int(), Error()).Block(
		calcArraySizeCodes...,
	)

	f.Func().Id("calcMapSize"+as.Name).Params(Id(v).Qual(as.PackageName, as.Name), Id(idEncoder).Op("*").Qual(pkEnc, "Encoder")).Params(Int(), Error()).Block(
		calcMapSizeCodes...,
	)

	f.Func().Id("encodeArray"+as.Name).Params(Id(v).Qual(as.PackageName, as.Name), Id(idEncoder).Op("*").Qual(pkEnc, "Encoder"), Id("offset").Int()).Params(Index().Byte(), Int(), Error()).Block(
		encArrayCodes...,
	)

	f.Func().Id("encodeMap"+as.Name).Params(Id(v).Qual(as.PackageName, as.Name), Id(idEncoder).Op("*").Qual(pkEnc, "Encoder"), Id("offset").Int()).Params(Index().Byte(), Int(), Error()).Block(
		encMapCodes...,
	)

	f.Func().Id("decodeArray"+as.Name).Params(Id(v).Op("*").Qual(as.PackageName, as.Name), Id(idDecoder).Op("*").Qual(pkDec, "Decoder"), Id("offset").Int()).Params(Int(), Error()).Block(
		decArrayCodes...,
	)

	f.Func().Id("decodeMap"+as.Name).Params(Id(v).Op("*").Qual(as.PackageName, as.Name), Id(idDecoder).Op("*").Qual(pkDec, "Decoder"), Id("offset").Int()).Params(Int(), Error()).Block(
		decMapCodes...,
	)
}

func (as *analyzedStruct) createFieldCode(fieldType types.Type, fieldName string, isRoot bool) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {

	switch reflect.TypeOf(fieldType) {
	case reflect.TypeOf(&types.Basic{}):
		fmt.Println("basic", fieldName, fieldType)
		return as.createBasicCode(fieldType, fieldName, isRoot)

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

		ca, _, _, _, _, _, _ := as.createFieldCode(child, childName, false)

		statements := as.addSizePattern2("CalcSliceLength", Id(fmt.Sprintf("len(%s)", name)))
		statements = append(statements, For(List(Id("_"), Id(childName)).Op(":=").Range().Id(name)).Block(
			ca...,
		))

		cArray = append(cArray, If(Id(name).Op("!=").Nil()).Block(
			statements...,
		).Else().Block(
			as.addSizePattern1("CalcNil"),
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

		caKey, _, _, _, _, _, _ := as.createFieldCode(key, childKey, false)
		caValue, _, _, _, _, _, _ := as.createFieldCode(value, childValue, false)

		statements := as.addSizePattern2("CalcMapLength", Id(fmt.Sprintf("len(%s)", name)))
		statements = append(statements, For(List(Id(childKey), Id(childValue)).Op(":=").Range().Id(name)).Block(
			append(caKey, caValue...)...,
		))

		cArray = append(cArray, If(Id(name).Op("!=").Nil()).Block(
			statements...,
		).Else().Block(
			as.addSizePattern1("CalcNil"),
		))

	case reflect.TypeOf(&types.Named{}):
		fieldValue := Id(fieldName)
		if isRoot {
			fieldValue = Id("v").Dot(fieldName)
		}

		if fieldType.String() == "time.Time" {

			cArray = append(cArray, as.addSizePattern1("CalcTime", fieldValue))
			eArray = append(eArray, as.encPattern1("WriteTime", fieldValue, Id("offset")))

			cMap = append(cMap, as.addSizePattern1("CalcTime", fieldValue))
			eMap = append(eMap, as.encPattern1("WriteTime", fieldValue, Id("offset")))

			dArray = append(dArray, as.decodeBasicPattern(fieldType, fieldName, "offset", "time.Time", "AsDateTime", isRoot)...)
			dMap = append(dMap, as.decodeBasicPattern(fieldType, fieldName, "offset", "time.Time", "AsDateTime", isRoot)...)
		} else {
			// todo : 対象のパッケージかどうかをちゃんと判断する
			cArray, cMap, eArray, eMap, dArray, dMap = as.createNamedCode(fieldName, fieldValue, isRoot)
		}
		//if (types.Identical(types.Struct{}, fieldType))
		fmt.Println("named", fieldName, fieldType, as.PackageName)

	default:
		// todo : error

	}

	return cArray, cMap, eArray, eMap, dArray, dMap, err
}

func (as *analyzedStruct) createSliceCode(fieldType types.Type, fieldName string, isRoot bool) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {
	child := fieldType.(*types.Slice).Elem()

	name, childName := "", ""
	if isRoot {
		name = "v." + fieldName
		childName = "vv"
	} else {
		childName = fieldName + "v"
	}

	ca, _, ea, _, _, _, _ := as.createFieldCode(child, childName, false)

	calcCodes := as.addSizePattern2("CalcSliceLength", Id(fmt.Sprintf("len(%s)", name)))
	calcCodes = append(calcCodes, For(List(Id("_"), Id(childName)).Op(":=").Range().Id(name)).Block(
		ca...,
	))

	cArray = append(cArray, If(Id(name).Op("!=").Nil()).Block(
		calcCodes...,
	).Else().Block(
		as.addSizePattern1("CalcNil"),
	))

	encCodes := as.addSizePattern2("WriteSliceLength", Id(fmt.Sprintf("len(%s)", name)), Id("offset"))
	encCodes = append(encCodes, For(List(Id("_"), Id(childName)).Op(":=").Range().Id(name)).Block(
		ea...,
	))

	eArray = append(eArray, If(Id(name).Op("!=").Nil()).Block(
		encCodes...,
	).Else().Block(
		Id("offset").Op("=").Id(idEncoder).Dot("WriteNil").Call(Id("offset")),
	))

	decCodes := make([]Code, 0)

	// todo : 構造体かBasicで別処理する必要がありそう
	// todo : tetest側でgenするようにして試さないといけない
	// todo : field.Pkgで構造体のQual参照できるかも
	checkChild := child
	for {
		switch reflect.TypeOf(child) {
		case reflect.TypeOf(&types.Basic{}):
			decCodes = append(decCodes, Var().Id(childName).Id(fieldType.String()))
			break

		case reflect.TypeOf(&types.Named{}):

		case reflect.TypeOf(&types.Slice{}):
			checkChild = checkChild.(*types.Slice).Elem()
			// continue

		case reflect.TypeOf(&types.Pointer{}):
			// todo

		default:

		}
	}

	switch reflect.TypeOf(child) {
	case reflect.TypeOf(&types.Basic{}):
		decCodes = append(decCodes, Var().Id(childName).Id(fieldType.String()))
	}

	decCodes = append(decCodes, For(List(Id("_"), Id(childName)).Op(":=").Range().Id(name)).Block(
		ea...,
	))

	dArray = append(dArray, If(Op("!").Id(idDecoder).Dot("IsCodeNil").Call(Id("offset"))).Block(
		decCodes...,
	).Else().Block(
		Id("offset").Op("++"),
	))

	return cArray, cArray, eArray, eArray, dArray, dArray, nil
}

func (as *analyzedStruct) createBasicCode(fieldType types.Type, fieldName string, isRoot bool) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {
	isType := func(kind types.BasicKind) bool {
		return types.Identical(types.Typ[kind], fieldType)
	}

	offset := "offset"

	fieldValue := Id(fieldName)
	if isRoot {
		fieldValue = Id("v").Dot(fieldName)
	}

	var (
		castName   = ""
		funcSuffix = ""
	)

	// todo : byte

	switch {
	case isType(types.Int), isType(types.Int8), isType(types.Int16), isType(types.Int32), isType(types.Int64):
		castName = "int64"
		funcSuffix = "Int"

	case isType(types.Uint), isType(types.Uint8), isType(types.Uint16), isType(types.Uint32), isType(types.Uint64):
		castName = "uint64"
		funcSuffix = "Uint"

	case isType(types.String):
		castName = "string"
		funcSuffix = "String"

	case isType(types.Float32):
		castName = "float32"
		funcSuffix = "Float32"

	case isType(types.Float64):
		castName = "float64"
		funcSuffix = "Float64"

	case isType(types.Bool):
		castName = "bool"
		funcSuffix = "Bool"
	default:
		// todo error

	}

	cArray = append(cArray, as.addSizePattern1("Calc"+funcSuffix, Id(castName).Call(fieldValue)))
	eArray = append(eArray, as.encPattern1("Write"+funcSuffix, Id(castName).Call(fieldValue), Id(offset)))

	cMap = append(cMap, as.addSizePattern1("Calc"+funcSuffix, Id(castName).Call(fieldValue)))
	eMap = append(eMap, as.encPattern1("Write"+funcSuffix, Id(castName).Call(fieldValue), Id(offset)))

	dArray = append(dArray, as.decodeBasicPattern(fieldType, fieldName, offset, castName, "As"+funcSuffix, isRoot)...)
	dMap = append(dMap, as.decodeBasicPattern(fieldType, fieldName, offset, castName, "As"+funcSuffix, isRoot)...)

	return cArray, cMap, eArray, eMap, dArray, dMap, err
}

func (as *analyzedStruct) addSizePattern1(funcName string, params ...Code) Code {
	return Id("size").Op("+=").Id(idEncoder).Dot(funcName).Call(params...)
}

func (as *analyzedStruct) addSizePattern2(funcName string, params ...Code) []Code {
	return []Code{
		List(Id("s"), Err()).Op(":=").Id(idEncoder).Dot(funcName).Call(params...),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id("size").Op("+=").Id("s"),
	}

}

func (as *analyzedStruct) encPattern1(funcName string, params ...Code) Code {
	return Id("offset").Op("=").Id(idEncoder).Dot(funcName).Call(params...)
}

func (as *analyzedStruct) decodeBasicPattern(fieldType types.Type, fieldName, offsetName, varTypeName, decoderFuncName string, isRoot bool) []Code {

	varName, setVarName := as.decodeVarPattern(fieldName, isRoot)

	return []Code{Block(
		Var().Id(varName).Id(varTypeName),
		List(Id(varName), Id(offsetName), Err()).Op("=").Id(idDecoder).Dot(decoderFuncName).Call(Id(offsetName)),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id(setVarName).Op("=").Id(fieldType.String()).Call(Id(varName)),
	)}
}

func (as *analyzedStruct) createNamedCode(fieldName string, fieldValue Code, isRoot bool) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	cArray = []Code{
		List(Id("size"+fieldName), Err()).
			Op(":=").
			Id(as.calcArraySizeFuncName()).Call(fieldValue, Id(idEncoder)),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id("size").Op("+=").Id(fieldName + "Size"),
	}

	cMap = []Code{
		List(Id("size"+fieldName), Err()).
			Op(":=").
			Id(as.calcMapSizeFuncName()).Call(fieldValue, Id(idEncoder)),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id("size").Op("+=").Id(fieldName + "Size"),
	}

	eArray = []Code{
		List(Id("_"), Id("offset"), Err()).
			Op("=").
			Id(as.encodeArrayFuncName()).Call(fieldValue, Id(idEncoder)),
		If(Err().Op("!=").Nil()).Block(
			Return(Nil(), Lit(0), Err()),
		),
	}

	eMap = []Code{
		List(Id("_"), Id("offset"), Err()).
			Op("=").
			Id(as.encodeMapFuncName()).Call(fieldValue, Id(idEncoder)),
		If(Err().Op("!=").Nil()).Block(
			Return(Nil(), Lit(0), Err()),
		),
	}

	varName, setVarName := as.decodeVarPattern(fieldName, isRoot)

	dArray = []Code{
		Block(
			Var().Id(varName).Qual(as.PackageName, as.Name),
			List(Id("offset"), Err()).Op("=").Id(as.decodeArrayFuncName()).Call(Op("&").Id(varName), Id(idDecoder), Id("offset")),
			If(Err().Op("!=").Nil()).Block(
				Return(Lit(0), Err()),
			),
			Id(setVarName).Op("=").Id(varName),
		),
	}

	// dArrayと一緒
	dMap = []Code{
		Block(
			Var().Id(varName).Qual(as.PackageName, as.Name),
			List(Id("offset"), Err()).Op("=").Id(as.decodeArrayFuncName()).Call(Op("&").Id(varName), Id(idDecoder), Id("offset")),
			If(Err().Op("!=").Nil()).Block(
				Return(Lit(0), Err()),
			),
			Id(setVarName).Op("=").Id(varName),
		),
	}
	return
}

func (as *analyzedStruct) decodeVarPattern(fieldName string, isRoot bool) (varName string, setVarName string) {

	varName = "vv"
	setVarName = "v." + fieldName
	if !isRoot {
		varName = fieldName + "v"
		setVarName = fieldName
	}
	return
}
