package structure

import (
	"crypto/sha256"
	"fmt"
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/shamaton/msgpackgen/internal/generator/ptn"
)

func isRootField(name string) bool {
	return strings.Contains(name, ".")
}

func createFuncName(prefix, name, importPath string) string {
	suffix := fmt.Sprintf("%x", sha256.Sum256([]byte(importPath)))
	return ptn.PrivateFuncName(fmt.Sprintf("%s%s_%s", prefix, name, suffix))
}

func createAddSizeCode(funcName string, params ...Code) Code {
	return Id("size").Op("+=").Id(ptn.IdEncoder).Dot(funcName).Call(params...)
}

func createAddSizeErrCheckCode(funcName string, params ...Code) []Code {
	return []Code{
		List(Id("s"), Err()).Op(":=").Id(ptn.IdEncoder).Dot(funcName).Call(params...),
		If(Err().Op("!=").Nil()).Block(
			Return(Lit(0), Err()),
		),
		Id("size").Op("+=").Id("s"),
	}

}

func createDecodeDefineVarCode(node *Node, structures []*Structure, varName string) ([]Code, string) {

	ptrCount, isParentTypeArrayOrMap := node.GetPointerInfo()

	codes := make([]Code, 0)
	receiverName := varName

	if ptrCount < 1 && !isParentTypeArrayOrMap {
		codes = append(codes, node.TypeJenChain(structures, Var().Id(receiverName)))
	} else if isParentTypeArrayOrMap {

		for i := 0; i < ptrCount; i++ {
			p := strings.Repeat("p", i+1)
			ptr := strings.Repeat("*", ptrCount-1-i)

			codes = append(codes, node.TypeJenChain(structures, Var().Id(varName+p).Op(ptr)))
		}
		receiverName = varName + strings.Repeat("p", ptrCount)
	} else {
		for i := 0; i < ptrCount; i++ {
			p := strings.Repeat("p", i)
			ptr := strings.Repeat("*", ptrCount-1-i)

			codes = append(codes, node.TypeJenChain(structures, Var().Id(varName+p).Op(ptr)))
		}
		receiverName = varName + strings.Repeat("p", ptrCount-1)
	}
	return codes, receiverName
}

func createDecodeSetValueCode(node *Node, varName, fieldName string) []Code {

	ptrCount, isParentTypeArrayOrMap := node.GetPointerInfo()

	var codes []Code
	if isParentTypeArrayOrMap {
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
				codes = append(codes, Id(fieldName).Op("=").Op("&").Id(tmp))
			}
		}
		if ptrCount < 1 {
			codes = append(codes, Id(fieldName).Op("=").Op("").Id(varName))
		}
	}

	return codes
}
