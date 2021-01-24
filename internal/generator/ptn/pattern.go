package ptn

import "fmt"

func PrivateFuncName(funcName string) string {
	return fmt.Sprintf("___%s", funcName)
}
