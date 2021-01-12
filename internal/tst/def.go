package tst

// todo : ドットimportに未対応
import (
	tst2 "github.com/shamaton/msgpackgen/internal/tst/tst"
	. "github.com/shamaton/msgpackgen/internal/tst/tst/tst"
)

//go:generate go run github.com/shamaton/msgpackgen -strict

type A struct {
	Int  int
	Uint uint
	B    tst2.B
	BB   C
}

type NotGenStruct struct {
	Interface interface{}
	Int       int
}

type NotGeStruct2 struct {
	I interface{}
}

type NotGen struct {
	A  []float32
	M  map[float64]uint64
	N  NotGenStruct
	N2 NotGeStruct2
	D  Def2
	//NN tst.NotNotGen
}
