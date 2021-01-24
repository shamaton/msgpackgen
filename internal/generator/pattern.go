package generator

import "fmt"

func privateFuncNamePattern(funcName string) string {
	return fmt.Sprintf("___%s", funcName)
}
