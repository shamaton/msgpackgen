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
			a := as.createNamedCode(fieldName, fieldValue, isRoot)
			cArray = append(cArray, a...)
		}
		//if (types.Identical(types.Struct{}, fieldType))
		fmt.Println("named", fieldName, fieldType, as.PackageName)

	default:
		// todo : error

	}

	return cArray, cMap, eArray, eMap, dArray, dMap, err
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

	switch {
	case isType(types.Int):
		cArray = append(cArray, as.addSizePattern1("CalcInt", Id("int64").Call(fieldValue)))
		eArray = append(eArray, as.encPattern1("WriteInt", Id("int64").Call(fieldValue), Id(offset)))

		cMap = append(cMap, as.addSizePattern1("CalcInt", Id("int64").Call(fieldValue)))
		eMap = append(eMap, as.encPattern1("WriteInt", Id("int64").Call(fieldValue), Id(offset)))

		dArray = append(dArray, as.decodeBasicPattern(fieldType, fieldName, offset, "int64", "AsInt", isRoot)...)
		dMap = append(dMap, as.decodeBasicPattern(fieldType, fieldName, offset, "int64", "AsInt", isRoot)...)

	case isType(types.Uint):
		cArray = append(cArray, as.addSizePattern1("CalcUint", Id("uint64").Call(fieldValue)))

	case isType(types.String):
		cArray = append(cArray, as.addSizePattern1("CalcString", fieldValue))

	case isType(types.Float64):
		cArray = append(cArray, as.addSizePattern1("CalcFloat64", fieldValue))
	default:
		// todo error

	}
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

func (as *analyzedStruct) createNamedCode(fieldName string, fieldValue Code, isRoot bool) []Code {

	return []Code{
		List(Id(fieldName+"Size"), Err()).
			Op(":=").
			Id(as.calcArraySizeFuncName()).Call(fieldValue, Id(idEncoder)),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id("size").Op("+=").Id(fieldName + "Size"),
	}
}
