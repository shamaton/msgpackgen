package main

import (
	"fmt"
	"strings"

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

	var firstEncParam, firstDecParam *Statement
	if as.NoUseQual {
		firstEncParam = Id(v).Id(as.Name)
		firstDecParam = Id(v).Op("*").Id(as.Name)
	} else {
		firstEncParam = Id(v).Qual(as.PackageName, as.Name)
		firstDecParam = Id(v).Op("*").Qual(as.PackageName, as.Name)
	}

	f.Comment(fmt.Sprintf("// calculate size from %s.%s\n", as.PackageName, as.Name)).
		Func().Id(as.calcArraySizeFuncName()).Params(firstEncParam, Id(idEncoder).Op("*").Qual(pkEnc, "Encoder")).Params(Int(), Error()).Block(
		append(calcArraySizeCodes, Return(Id("size"), Nil()))...,
	)

	f.Comment(fmt.Sprintf("// calculate size from %s.%s\n", as.PackageName, as.Name)).
		Func().Id(as.calcMapSizeFuncName()).Params(firstEncParam, Id(idEncoder).Op("*").Qual(pkEnc, "Encoder")).Params(Int(), Error()).Block(
		append(calcMapSizeCodes, Return(Id("size"), Nil()))...,
	)

	f.Comment(fmt.Sprintf("// encode from %s.%s\n", as.PackageName, as.Name)).
		Func().Id(as.encodeArrayFuncName()).Params(firstEncParam, Id(idEncoder).Op("*").Qual(pkEnc, "Encoder"), Id("offset").Int()).Params(Index().Byte(), Int(), Error()).Block(
		append(encArrayCodes, Return(Id(idEncoder).Dot("EncodedBytes").Call(), Id("offset"), Err()))...,
	)

	f.Comment(fmt.Sprintf("// encode from %s.%s\n", as.PackageName, as.Name)).
		Func().Id(as.encodeMapFuncName()).Params(firstEncParam, Id(idEncoder).Op("*").Qual(pkEnc, "Encoder"), Id("offset").Int()).Params(Index().Byte(), Int(), Error()).Block(
		append(encMapCodes, Return(Id(idEncoder).Dot("EncodedBytes").Call(), Id("offset"), Err()))...,
	)

	f.Comment(fmt.Sprintf("// decode to %s.%s\n", as.PackageName, as.Name)).
		Func().Id(as.decodeArrayFuncName()).Params(firstDecParam, Id(idDecoder).Op("*").Qual(pkDec, "Decoder"), Id("offset").Int()).Params(Int(), Error()).Block(
		append(decArrayCodes, Return(Id("offset"), Err()))...,
	)

	f.Comment(fmt.Sprintf("// decode to %s.%s\n", as.PackageName, as.Name)).
		Func().Id(as.decodeMapFuncName()).Params(firstDecParam, Id(idDecoder).Op("*").Qual(pkDec, "Decoder"), Id("offset").Int()).Params(Int(), Error()).Block(

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

			dArray = append(dArray, as.decodeBasicPattern(ast, fieldName, "offset", "AsDateTime", isRoot)...)
			dMap = append(dMap, as.decodeBasicPattern(ast, fieldName, "offset", "AsDateTime", isRoot)...)
		} else {
			// todo : 対象のパッケージかどうかをちゃんと判断する
			cArray, cMap, eArray, eMap, dArray, dMap = as.createNamedCode(fieldName, ast, fieldValue, isRoot)
		}
		fmt.Println("named", fieldName, ast, as.PackageName)

	case ast.IsPointer():
		return as.createPointerCode(ast, fieldName, isRoot)

	default:
		// todo : error

	}

	return cArray, cMap, eArray, eMap, dArray, dMap, err
}

func (as *analyzedStruct) createPointerCode(ast *analyzedASTFieldType, fieldName string, isRoot bool) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {

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

	name := fieldName
	if isRoot {
		name = "v." + fieldName
	}
	ca, _, ea, _, da, _, _ := as.createFieldCode(ast.Elm(), fieldName, isRoot)

	cArray = make([]Code, 0)
	cArray = append(cArray, If(Op(ptrOp).Id(name).Op("!=").Nil()).Block(
		ca...,
	).Else().Block(
		Id("size").Op("+=").Id(idEncoder).Dot("CalcNil").Call(),
	))

	eArray = make([]Code, 0)
	eArray = append(eArray, If(Op(ptrOp).Id(name).Op("!=").Nil()).Block(
		ea...,
	).Else().Block(
		Id("offset").Op("=").Id(idEncoder).Dot("WriteNil").Call(Id("offset")),
	))

	// todo : ようかくにん、重複コードをスキップ
	if len(ptrOp) < 1 {
		dArray = make([]Code, 0)
		dArray = append(dArray, If(Op("!").Id(idDecoder).Dot("IsCodeNil").Call(Id("offset"))).Block(
			da...,
		).Else().Block(
			Id("offset").Op("++"),
		))
	} else {
		dArray = da
	}

	return cArray, cArray, eArray, eArray, dArray, dArray, err
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

	ptrOp := ""
	andOp := ""
	node := ast
	for {
		if node.HasParent() && node.Parent.IsPointer() {
			ptrOp += "*"
			andOp += "&"
			node = node.Parent
		} else {
			break
		}
	}

	caKey, _, eaKey, _, daKey, _, _ := as.createFieldCode(key, childKey, false)
	caValue, _, eaValue, _, daValue, _, _ := as.createFieldCode(value, childValue, false)

	calcCodes := as.addSizePattern2("CalcMapLength", Len(Op(ptrOp).Id(name)))
	calcCodes = append(calcCodes, For(List(Id(childKey), Id(childValue)).Op(":=").Range().Op(ptrOp).Id(name)).Block(
		append(caKey, caValue...)...,
	))

	cArray = append(cArray, If(Id(name).Op("!=").Nil()).Block(
		calcCodes...,
	).Else().Block(
		as.addSizePattern1("CalcNil"),
	))

	encCodes := make([]Code, 0)
	encCodes = append(encCodes, Id("offset").Op("=").Id(idEncoder).Dot("WriteMapLength").Call(Len(Op(ptrOp).Id(name)), Id("offset")))
	encCodes = append(encCodes, For(List(Id(childKey), Id(childValue)).Op(":=").Range().Op(ptrOp).Id(name)).Block(
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

	da := []Code{ast.Key.TypeJenChain(Var().Id(childKey + "v"))}
	da = append(da, daKey...)
	da = append(da, ast.Value.TypeJenChain(Var().Id(childValue+"v")))
	da = append(da, daValue...)
	da = append(da, Id(childValue).Index(Id(childKey+"v")).Op("=").Id(childValue+"v"))

	decCodes = append(decCodes, For(Id(childValue+"i").Op(":=").Lit(0).Op(";").Id(childValue+"i").Op("<").Id(childValue+"l").Op(";").Id(childValue+"i").Op("++")).Block(
		da..., //append(append(daKey, daValue...), Id(childValue).Index(Id(childKey+"v")).Op("=").Id(childValue+"v"))...,
	))
	decCodes = append(decCodes, Id(name).Op("=").Op(andOp).Id(childValue))

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

	ptrOp := ""
	andOp := ""
	node := ast
	for {
		if node.HasParent() && node.Parent.IsPointer() {
			ptrOp += "*"
			andOp += "&"
			node = node.Parent
		} else {
			break
		}
	}

	ca, _, ea, _, da, _, _ := as.createFieldCode(ast.Elm(), childName, false)

	calcCodes := as.addSizePattern2("CalcSliceLength", Len(Op(ptrOp).Id(name)))
	calcCodes = append(calcCodes, For(List(Id("_"), Id(childName)).Op(":=").Range().Op(ptrOp).Id(name)).Block(
		ca...,
	))

	cArray = append(cArray, If(Op(ptrOp).Id(name).Op("!=").Nil()).Block(
		calcCodes...,
	).Else().Block(
		as.addSizePattern1("CalcNil"),
	))

	encCodes := make([]Code, 0)
	encCodes = append(encCodes, Id("offset").Op("=").Id(idEncoder).Dot("WriteSliceLength").Call(Len(Op(ptrOp).Id(name)), Id("offset")))
	encCodes = append(encCodes, For(List(Id("_"), Id(childName)).Op(":=").Range().Op(ptrOp).Id(name)).Block(
		ea...,
	))

	eArray = append(eArray, If(Op(ptrOp).Id(name).Op("!=").Nil()).Block(
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

	da = append([]Code{ast.Elm().TypeJenChain(Var().Id(childName + "v"))}, da...)
	da = append(da, Id(childName).Index(Id(childName+"i")).Op("=").Id(childName+"v"))

	decCodes = append(decCodes, For(Id(childName+"i").Op(":=").Range().Id(childName)).Block(
		da...,
	))
	decCodes = append(decCodes, Id(name).Op("=").Op(andOp).Id(childName))

	// todo : ようかくにん、重複コードをスキップ
	if ast.HasParent() && ast.Parent.IsPointer() {
		dArray = decCodes
	} else {

		dArray = append(dArray, If(Op("!").Id(idDecoder).Dot("IsCodeNil").Call(Id("offset"))).Block(
			decCodes...,
		).Else().Block(
			Id("offset").Op("++"),
		))
	}

	return cArray, cArray, eArray, eArray, dArray, dArray, nil
}

func (as *analyzedStruct) createBasicCode(ast *analyzedASTFieldType, fieldName string, isRoot bool) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {

	offset := "offset"

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

	fieldValue := Op(ptrOp).Id(fieldName)
	if isRoot {
		fieldValue = Op(ptrOp).Id("v").Dot(fieldName)
	}

	var (
		//castName   = ""
		funcSuffix = ""
	)

	// todo : byte

	//switch ast.IdenticalName {
	//case "int", "int8", "int16", "int32", "int64":
	//	castName = "int64"
	//	funcSuffix = "Int"
	//
	//case "uint", "uint8", "uint16", "uint32", "uint64":
	//	castName = "uint64"
	//	funcSuffix = "Uint"
	//
	//case "string":
	//	castName = "string"
	//	funcSuffix = "String"
	//
	//case "float32":
	//	castName = "float32"
	//	funcSuffix = "Float32"
	//
	//case "float64":
	//	castName = "float64"
	//	funcSuffix = "Float64"
	//
	//case "bool":
	//	castName = "bool"
	//	funcSuffix = "Bool"
	//default:
	//	// todo error
	//
	//}
	funcSuffix = strings.Title(ast.IdenticalName)

	cArray = append(cArray, as.addSizePattern1("Calc"+funcSuffix, fieldValue))
	eArray = append(eArray, as.encPattern1("Write"+funcSuffix, fieldValue, Id(offset)))

	cMap = append(cMap, as.addSizePattern1("Calc"+funcSuffix, fieldValue))
	eMap = append(eMap, as.encPattern1("Write"+funcSuffix, fieldValue, Id(offset)))

	dArray = append(dArray, as.decodeBasicPattern(ast, fieldName, offset, "As"+funcSuffix, isRoot)...)
	dMap = append(dMap, as.decodeBasicPattern(ast, fieldName, offset, "As"+funcSuffix, isRoot)...)

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

func (as *analyzedStruct) decodeBasicPattern(ast *analyzedASTFieldType, fieldName, offsetName, decoderFuncName string, isRoot bool) []Code {

	// todo : ポインタの場合, vvp / vvvを使う必要
	varName, setVarName := as.decodeVarPattern(fieldName, isRoot)

	node := ast
	ptrCount := 0
	isParentTypeArrayOrMap := false

	for {
		if node.HasParent() {
			node = node.Parent
			if node.IsPointer() {
				ptrCount++
			} else if node.IsArray() || node.IsMap() {
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
		codes = append(codes, ast.TypeJenChain(Var().Id(recieverName)))
	} else if isParentTypeArrayOrMap {

		for i := 0; i < ptrCount; i++ {
			p := strings.Repeat("p", i+1)
			kome := strings.Repeat("*", ptrCount-1-i)
			codes = append(codes, ast.TypeJenChain(Var().Id(varName+p).Op(kome)))
		}
		recieverName = varName + strings.Repeat("p", ptrCount)
	} else {
		for i := 0; i < ptrCount; i++ {
			p := strings.Repeat("p", i)
			kome := strings.Repeat("*", ptrCount-1-i)
			codes = append(codes, ast.TypeJenChain(Var().Id(varName+p).Op(kome)))
		}
		recieverName = varName + strings.Repeat("p", ptrCount-1)
	}

	codes = append(codes,
		//ast.TypeJenChain(Var().Id(varName)), // todo : これは外だし
		List(Id(recieverName), Id(offsetName), Err()).Op("=").Id(idDecoder).Dot(decoderFuncName).Call(Id(offsetName)),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
	)

	codes = as.createDecodeSetVarPattern(ptrCount, varName, setVarName, isParentTypeArrayOrMap, codes)

	// array or map
	if isParentTypeArrayOrMap {
		return codes
	}

	//for i := 0; i < ptrCount; i++ {
	//	if i != ptrCount-1 {
	//		tmp1 := varName + strings.Repeat("p", i)
	//		tmp2 := varName + strings.Repeat("p", i+1)
	//		commons = append(commons, Id(tmp2).Op("=").Op("&").Id(tmp1))
	//	} else {
	//		// last
	//		tmp := varName + strings.Repeat("p", i)
	//		commons = append(commons, Id(setVarName).Op("=").Op("&").Id(tmp))
	//	}
	//}
	//if ptrCount < 1 {
	//	commons = append(commons, Id(setVarName).Op("=").Op("").Id(varName))
	//}
	return []Code{Block(codes...)}
}

func (as *analyzedStruct) createDecodeSetVarPattern(ptrCount int, varName, setVarName string, isLastSkip bool, codes []Code) []Code {

	if isLastSkip {
		for i := 0; i < ptrCount; i++ {
			tmp1 := varName + strings.Repeat("p", ptrCount-1-i)
			tmp2 := varName + strings.Repeat("p", ptrCount-i)
			codes = append(codes, Id(tmp1).Op("=").Op("&").Id(tmp2))
		}
	} else {

		for i := 0; i < ptrCount; i++ {
			if i != ptrCount-1 {
				tmp1 := varName + strings.Repeat("p", ptrCount-2+i)
				tmp2 := varName + strings.Repeat("p", ptrCount-1+i)
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

func (as *analyzedStruct) createNamedCode(fieldName string, ast *analyzedASTFieldType, fieldValue Code, isRoot bool) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code) {

	cArray = []Code{
		List(Id("size"+fieldName), Err()).
			Op(":=").
			Id(createFuncName("calcArraySize", ast.StructName, ast.ImportPath)).Call(fieldValue, Id(idEncoder)),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id("size").Op("+=").Id("size" + fieldName),
	}

	cMap = []Code{
		List(Id("size"+fieldName), Err()).
			Op(":=").
			Id(createFuncName("calcMapSize", ast.StructName, ast.ImportPath)).Call(fieldValue, Id(idEncoder)),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id("size").Op("+=").Id("size" + fieldName),
	}

	eArray = []Code{
		List(Id("_"), Id("offset"), Err()).
			Op("=").
			Id(createFuncName("encodeArray", ast.StructName, ast.ImportPath)).Call(fieldValue, Id(idEncoder), Id("offset")),
		If(Err().Op("!=").Nil()).Block(
			Return(Nil(), Lit(0), Err()),
		),
	}

	eMap = []Code{
		List(Id("_"), Id("offset"), Err()).
			Op("=").
			Id(createFuncName("encodeMap", ast.StructName, ast.ImportPath)).Call(fieldValue, Id(idEncoder), Id("offset")),
		If(Err().Op("!=").Nil()).Block(
			Return(Nil(), Lit(0), Err()),
		),
	}

	//varName, setVarName := as.decodeVarPattern(fieldName, isRoot)

	dArray = append(dArray, as.decodeNamedPattern(ast, fieldName, "decodeArray", isRoot)...)
	dMap = append(dMap, as.decodeNamedPattern(ast, fieldName, "decodeMap", isRoot)...)

	//dArray = []Code{
	//	Block(
	//		Var().Id(varName).Qual(ast.ImportPath, ast.StructName),
	//		List(Id("offset"), Err()).Op("=").Id(createFuncName("decodeArray", ast.StructName, ast.ImportPath)).Call(Op("&").Id(varName), Id(idDecoder), Id("offset")),
	//		If(Err().Op("!=").Nil()).Block(
	//			Return(Lit(0), Err()),
	//		),
	//		Id(setVarName).Op("=").Id(varName),
	//	),
	//}
	//
	//// dArrayと一緒
	//dMap = []Code{
	//	Block(
	//		Var().Id(varName).Qual(ast.ImportPath, ast.StructName),
	//		List(Id("offset"), Err()).Op("=").Id(createFuncName("decodeMap", ast.StructName, ast.ImportPath)).Call(Op("&").Id(varName), Id(idDecoder), Id("offset")),
	//		If(Err().Op("!=").Nil()).Block(
	//			Return(Lit(0), Err()),
	//		),
	//		Id(setVarName).Op("=").Id(varName),
	//	),
	//}
	return
}

func (as *analyzedStruct) decodeNamedPattern(ast *analyzedASTFieldType, fieldName, decodeFuncName string, isRoot bool) []Code {

	// todo : ポインタの場合, vvp / vvvを使う必要
	varName, setVarName := as.decodeVarPattern(fieldName, isRoot)

	node := ast
	ptrCount := 0
	isParentTypeArrayOrMap := false

	for {
		if node.HasParent() {
			node = node.Parent
			if node.IsPointer() {
				ptrCount++
			} else if node.IsArray() || node.IsMap() {
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
		codes = append(codes, ast.TypeJenChain(Var().Id(recieverName)))
	} else if isParentTypeArrayOrMap {

		for i := 0; i < ptrCount; i++ {
			p := strings.Repeat("p", i+1)
			kome := strings.Repeat("*", ptrCount-1-i)
			codes = append(codes, ast.TypeJenChain(Var().Id(varName+p).Op(kome)))
		}
		recieverName = varName + strings.Repeat("p", ptrCount)
	} else {
		for i := 0; i < ptrCount; i++ {
			p := strings.Repeat("p", i)
			kome := strings.Repeat("*", ptrCount-1-i)
			codes = append(codes, ast.TypeJenChain(Var().Id(varName+p).Op(kome)))
		}
		recieverName = varName + strings.Repeat("p", ptrCount-1)
	}

	codes = append(codes,
		//ast.TypeJenChain(Var().Id(varName)), // todo : これは外だし
		List(Id("offset"), Err()).Op("=").Id(createFuncName(decodeFuncName, ast.StructName, ast.ImportPath)).Call(Op("&").Id(varName), Id(idDecoder), Id("offset")),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
	)

	codes = as.createDecodeSetVarPattern(ptrCount, varName, setVarName, isParentTypeArrayOrMap, codes)

	// array or map
	if isParentTypeArrayOrMap {
		return codes
	}

	//for i := 0; i < ptrCount; i++ {
	//	if i != ptrCount-1 {
	//		tmp1 := varName + strings.Repeat("p", i)
	//		tmp2 := varName + strings.Repeat("p", i+1)
	//		commons = append(commons, Id(tmp2).Op("=").Op("&").Id(tmp1))
	//	} else {
	//		// last
	//		tmp := varName + strings.Repeat("p", i)
	//		commons = append(commons, Id(setVarName).Op("=").Op("&").Id(tmp))
	//	}
	//}
	//if ptrCount < 1 {
	//	commons = append(commons, Id(setVarName).Op("=").Op("").Id(varName))
	//}
	return []Code{Block(codes...)}
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
