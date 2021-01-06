package t

import "github.com/shamaton/msgpackgen/internal/t/t"

//go:generate go run github.com/shamaton/msgpackgen

type A struct {
	Int int
	B   t.B
}
