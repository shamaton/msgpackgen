package main

import (
	"fmt"
	"go/types"
	"path/filepath"
	"reflect"

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

	//for _, field := range st.Fields {
	//	switch {
	//	case types.Identical(types.Typ[types.Int], field.Type):
	//		calcArraySizeCodes = append(calcArraySizeCodes, addSizePattern1("CalcInt", Id("int64").Call(Id("v").Dot(field.Name))))
	//
	//		//case types.Identical(types.Typ[types.], field.Type):
	//
	//	}
	//}

	for _, field := range st.Fields {
		cArray, cMap, _, _, _, _, _ := createFieldCode(field.Type, field.Name, true)
		calcArraySizeCodes = append(calcArraySizeCodes, cArray...)
		calcMapSizeCodes = append(calcMapSizeCodes, cMap...)

		//switch reflect.TypeOf(field.Type) {
		//case reflect.TypeOf(&types.Basic{}):
		//	fmt.Println("basic", field.Name, field.Type)
		//	cArray, cMap, _, _, _, _, _ := createBasicCode(field)
		//	calcArraySizeCodes = append(calcArraySizeCodes, cArray)
		//	calcMapSizeCodes = append(calcMapSizeCodes, cMap...)
		//
		//case reflect.TypeOf(&types.Slice{}):
		//	fmt.Println("slice", field.Name, field.Type)
		//	fmt.Println(field.Type.(*types.Slice).Elem().)
		//
		//case reflect.TypeOf(&types.Map{}):
		//	fmt.Println("map", field.Name, field.Type)
		//
		//case reflect.TypeOf(&types.Named{}):
		//	fmt.Println("named", field.Name, field.Type)
		//
		//default:
		//	// todo : error
		//
		//}
	}

	f.Func().Id("calcArraySize"+st.Name).Params(Id(v).Qual(st.PackageName, st.Name), Id(idEncoder).Op("*").Qual(pkEnc, "Encoder")).Params(Int(), Error()).Block(
		calcArraySizeCodes...,
	)

	f.Func().Id("calcMapSize"+st.Name).Params(Id(v).Qual(st.PackageName, st.Name), Id(idEncoder).Op("*").Qual(pkEnc, "Encoder")).Params(Int(), Error()).Block(
		calcMapSizeCodes...,
	)
}

func createFieldCode(fieldType types.Type, fieldName string, isRoot bool) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {

	switch reflect.TypeOf(fieldType) {
	case reflect.TypeOf(&types.Basic{}):
		fmt.Println("basic", fieldName, fieldType)
		return createBasicCode(fieldType, fieldName)

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
		fmt.Println("named", fieldName, fieldType)

	default:
		// todo : error

	}

	return cArray, cMap, eArray, eMap, dArray, dMap, err
}

func createBasicCode(fieldType types.Type, fieldName string) (cArray []Code, cMap []Code, eArray []Code, eMap []Code, dArray []Code, dMap []Code, err error) {
	isType := func(kind types.BasicKind) bool {
		return types.Identical(types.Typ[kind], fieldType)
	}
	offset := "offset"

	switch {
	case isType(types.Int):
		cArray = append(cArray, addSizePattern1("CalcInt", Id("int64").Call(Id("v").Dot(fieldName))))
		eArray = append(eArray, addSizePattern1("WriteInt", Id("int64").Call(Id("v").Dot(fieldName)), Id(offset)))

		cMap = append(cMap, addSizePattern1("CalcString", Lit(fieldName)))
		cMap = append(cMap, addSizePattern1("CalcInt", Id("int64").Call(Id("v").Dot(fieldName))))
		eMap = append(eMap, addSizePattern1("WriteString", Lit(fieldName), Id(offset)))
		eMap = append(eMap, addSizePattern1("WriteInt", Id("int64").Call(Id("v").Dot(fieldName)), Id(offset)))

	case isType(types.Uint):
		cArray = append(cArray, addSizePattern1("CalcUint", Id("uint64").Call(Id("v").Dot(fieldName))))

	case isType(types.String):
		cArray = append(cArray, addSizePattern1("CalcString", Id("uint64").Call(Id("v").Dot(fieldName))))

	case isType(types.Float64):
		cArray = append(cArray, addSizePattern1("CalcFloat64", Id("float64").Call(Id("v").Dot(fieldName))))
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
