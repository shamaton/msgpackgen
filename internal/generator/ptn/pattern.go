package ptn

import "fmt"

// PrivateFuncName gets a function name that adapts funcName to the pattern
func PrivateFuncName(funcName string) string {
	return fmt.Sprintf("___%s", funcName)
}
