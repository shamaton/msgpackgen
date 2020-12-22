package main

import (
	"fmt"

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

		cArray, cMap, eArray, eMap, dArray, dMap, _ := as.createFieldCode(field.Ast, field.Name, true)
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
		append(calcArraySizeCodes, Return(Id("size"), Nil()))...,
	)

	f.Func().Id("calcMapSize"+as.Name).Params(Id(v).Qual(as.PackageName, as.Name), Id(idEncoder).Op("*").Qual(pkEnc, "Encoder")).Params(Int(), Error()).Block(
		append(calcMapSizeCodes, Return(Id("size"), Nil()))...,
	)

	f.Func().Id("encodeArray"+as.Name).Params(Id(v).Qual(as.PackageName, as.Name), Id(idEncoder).Op("*").Qual(pkEnc, "Encoder"), Id("offset").Int()).Params(Index().Byte(), Int(), Error()).Block(
		append(encArrayCodes, Return(Id(idEncoder).Dot("EncodedBytes").Call(), Id("offset"), Err()))...,
	)

	f.Func().Id("encodeMap"+as.Name).Params(Id(v).Qual(as.PackageName, as.Name), Id(idEncoder).Op("*").Qual(pkEnc, "Encoder"), Id("offset").Int()).Params(Index().Byte(), Int(), Error()).Block(
		append(encMapCodes, Return(Id(idEncoder).Dot("EncodedBytes").Call(), Id("offset"), Err()))...,
	)

	f.Func().Id("decodeArray"+as.Name).Params(Id(v).Op("*").Qual(as.PackageName, as.Name), Id(idDecoder).Op("*").Qual(pkDec, "Decoder"), Id("offset").Int()).Params(Int(), Error()).Block(
		append(decArrayCodes, Return(Id("offset"), Err()))...,
	)

	f.Func().Id("decodeMap"+as.Name).Params(Id(v).Op("*").Qual(as.PackageName, as.Name), Id(idDecoder).Op("*").Qual(pkDec, "Decoder"), Id("offset").Int()).Params(Int(), Error()).Block(

		append(decMapCodes, Return(Id("offset"), Err()))...,
	)
}

func (as *analyzedStruct) createFieldCode(ast *analyzedASTFieldType, fieldName string, isRoot bool) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {

	switch {
	case ast.IsIdentical():
		fmt.Println("basic", fieldName, ast)
		return as.createBasicCode(ast, fieldName, isRoot)

	case ast.IsArray():
		fmt.Println("slice", fieldName, ast)
		fmt.Println("type string.................................. ", ast.TypeString())
		return as.createSliceCode(ast, fieldName, isRoot)

	case ast.IsMap():
		return as.createMapCode(ast, fieldName, isRoot)

	case ast.IsStruct():
		fieldValue := Id(fieldName)
		if isRoot {
			fieldValue = Id("v").Dot(fieldName)
		}

		if ast.ImportPath == "time" {

			cArray = append(cArray, as.addSizePattern1("CalcTime", fieldValue))
			eArray = append(eArray, as.encPattern1("WriteTime", fieldValue, Id("offset")))

			cMap = append(cMap, as.addSizePattern1("CalcTime", fieldValue))
			eMap = append(eMap, as.encPattern1("WriteTime", fieldValue, Id("offset")))

			dArray = append(dArray, as.decodeBasicPattern(ast, fieldName, "offset", "time.Time", "AsDateTime", isRoot)...)
			dMap = append(dMap, as.decodeBasicPattern(ast, fieldName, "offset", "time.Time", "AsDateTime", isRoot)...)
		} else {
			// todo : 対象のパッケージかどうかをちゃんと判断する
			cArray, cMap, eArray, eMap, dArray, dMap = as.createNamedCode(fieldName, ast, fieldValue, isRoot)
		}
		//if (types.Identical(types.Struct{}, ast))
		fmt.Println("named", fieldName, ast, as.PackageName)

	default:
		// todo : error

	}

	return cArray, cMap, eArray, eMap, dArray, dMap, err
}

func (as *analyzedStruct) createMapCode(ast *analyzedASTFieldType, fieldName string, isRoot bool) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {

	key, value := ast.KeyValue()
	fmt.Println("map", fieldName, ast)
	fmt.Println(key, value)
	fmt.Println("type string.................................. ", ast.TypeString())

	name, childKey, childValue := "", "", ""
	if isRoot {
		name = "v." + fieldName
		childKey = "kk"
		childValue = "vv"
	} else {
		childKey = fieldName + "k"
		childValue = fieldName + "v"
	}

	caKey, _, eaKey, _, daKey, _, _ := as.createFieldCode(key, childKey, false)
	caValue, _, eaValue, _, daValue, _, _ := as.createFieldCode(value, childValue, false)

	calcCodes := as.addSizePattern2("CalcMapLength", Id(fmt.Sprintf("len(%s)", name)))
	calcCodes = append(calcCodes, For(List(Id(childKey), Id(childValue)).Op(":=").Range().Id(name)).Block(
		append(caKey, caValue...)...,
	))

	cArray = append(cArray, If(Id(name).Op("!=").Nil()).Block(
		calcCodes...,
	).Else().Block(
		as.addSizePattern1("CalcNil"),
	))

	encCodes := make([]Code, 0)
	encCodes = append(encCodes, Id("offset").Op("=").Id(idEncoder).Dot("WriteMapLength").Call(Id(fmt.Sprintf("len(%s)", name)), Id("offset")))
	encCodes = append(encCodes, For(List(Id(childKey), Id(childValue)).Op(":=").Range().Id(name)).Block(
		append(eaKey, eaValue...)...,
	))

	eArray = append(eArray, If(Id(name).Op("!=").Nil()).Block(
		encCodes...,
	).Else().Block(
		Id("offset").Op("=").Id(idEncoder).Dot("WriteNil").Call(Id("offset")),
	))

	decCodes := make([]Code, 0)
	decCodes = append(decCodes, ast.TypeJenChain(Var().Id(childValue)))
	decCodes = append(decCodes, Var().Id(childValue+"l").Int())
	decCodes = append(decCodes, List(Id(childValue+"l"), Id("offset"), Err()).Op("=").Id(idDecoder).Dot("MapLength").Call(Id("offset")))
	decCodes = append(decCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	decCodes = append(decCodes, Id(childValue).Op("=").Make(ast.TypeJenChain(), Id(childValue+"l")))
	decCodes = append(decCodes, For(Id(childValue+"i").Op(":=").Lit(0).Op(";").Id(childValue+"i").Op("<").Id(childValue+"l").Op(";").Id(childValue+"i").Op("++")).Block(
		append(daKey, daValue...)...,
	))
	decCodes = append(decCodes, Id(name).Op("=").Id(childValue))

	dArray = append(dArray, If(Op("!").Id(idDecoder).Dot("IsCodeNil").Call(Id("offset"))).Block(
		decCodes...,
	).Else().Block(
		Id("offset").Op("++"),
	))

	return cArray, cArray, eArray, eArray, dArray, dArray, nil
}

func (as *analyzedStruct) createSliceCode(ast *analyzedASTFieldType, fieldName string, isRoot bool) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {

	name, childName := "", ""
	if isRoot {
		name = "v." + fieldName
		childName = "vv"
	} else {
		childName = fieldName + "v"
	}

	ca, _, ea, _, da, _, _ := as.createFieldCode(ast.Elm(), childName, false)

	calcCodes := as.addSizePattern2("CalcSliceLength", Id(fmt.Sprintf("len(%s)", name)))
	calcCodes = append(calcCodes, For(List(Id("_"), Id(childName)).Op(":=").Range().Id(name)).Block(
		ca...,
	))

	cArray = append(cArray, If(Id(name).Op("!=").Nil()).Block(
		calcCodes...,
	).Else().Block(
		as.addSizePattern1("CalcNil"),
	))

	encCodes := make([]Code, 0)
	encCodes = append(encCodes, Id("offset").Op("=").Id(idEncoder).Dot("WriteSliceLength").Call(Id(fmt.Sprintf("len(%s)", name)), Id("offset")))
	encCodes = append(encCodes, For(List(Id("_"), Id(childName)).Op(":=").Range().Id(name)).Block(
		ea...,
	))

	eArray = append(eArray, If(Id(name).Op("!=").Nil()).Block(
		encCodes...,
	).Else().Block(
		Id("offset").Op("=").Id(idEncoder).Dot("WriteNil").Call(Id("offset")),
	))

	decCodes := make([]Code, 0)
	decCodes = append(decCodes, ast.TypeJenChain(Var().Id(childName)))
	decCodes = append(decCodes, Var().Id(childName+"l").Int())
	decCodes = append(decCodes, List(Id(childName+"l"), Id("offset"), Err()).Op("=").Id(idDecoder).Dot("SliceLength").Call(Id("offset")))
	decCodes = append(decCodes, If(Err().Op("!=").Nil()).Block(
		Return(Lit(0), Err()),
	))
	decCodes = append(decCodes, Id(childName).Op("=").Make(ast.TypeJenChain(), Id(childName+"l")))
	decCodes = append(decCodes, For(Id(childName+"i").Op(":=").Range().Id(childName)).Block(
		da...,
	))
	decCodes = append(decCodes, Id(name).Op("=").Id(childName))

	dArray = append(dArray, If(Op("!").Id(idDecoder).Dot("IsCodeNil").Call(Id("offset"))).Block(
		decCodes...,
	).Else().Block(
		Id("offset").Op("++"),
	))

	return cArray, cArray, eArray, eArray, dArray, dArray, nil
}

func (as *analyzedStruct) createBasicCode(ast *analyzedASTFieldType, fieldName string, isRoot bool) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {

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

	switch ast.IdenticalName {
	case "int", "int8", "int16", "int32", "int64":
		castName = "int64"
		funcSuffix = "Int"

	case "uint", "uint8", "uint16", "uint32", "uint64":
		castName = "uint64"
		funcSuffix = "Uint"

	case "string":
		castName = "string"
		funcSuffix = "String"

	case "float32":
		castName = "float32"
		funcSuffix = "Float32"

	case "float64":
		castName = "float64"
		funcSuffix = "Float64"

	case "bool":
		castName = "bool"
		funcSuffix = "Bool"
	default:
		// todo error

	}

	cArray = append(cArray, as.addSizePattern1("Calc"+funcSuffix, Id(castName).Call(fieldValue)))
	eArray = append(eArray, as.encPattern1("Write"+funcSuffix, Id(castName).Call(fieldValue), Id(offset)))

	cMap = append(cMap, as.addSizePattern1("Calc"+funcSuffix, Id(castName).Call(fieldValue)))
	eMap = append(eMap, as.encPattern1("Write"+funcSuffix, Id(castName).Call(fieldValue), Id(offset)))

	dArray = append(dArray, as.decodeBasicPattern(ast, fieldName, offset, castName, "As"+funcSuffix, isRoot)...)
	dMap = append(dMap, as.decodeBasicPattern(ast, fieldName, offset, castName, "As"+funcSuffix, isRoot)...)

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

func (as *analyzedStruct) decodeBasicPattern(ast *analyzedASTFieldType, fieldName, offsetName, varTypeName, decoderFuncName string, isRoot bool) []Code {

	varName, setVarName := as.decodeVarPattern(fieldName, isRoot)

	return []Code{Block(
		Var().Id(varName).Id(varTypeName),
		List(Id(varName), Id(offsetName), Err()).Op("=").Id(idDecoder).Dot(decoderFuncName).Call(Id(offsetName)),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id(setVarName).Op("=").Id(ast.IdenticalName).Call(Id(varName)),
	)}
}

func (as *analyzedStruct) createNamedCode(fieldName string, ast *analyzedASTFieldType, fieldValue Code, isRoot bool) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

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
			Var().Id(varName).Qual(ast.ImportPath, ast.StructName),
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
			Var().Id(varName).Qual(ast.ImportPath, ast.StructName),
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
