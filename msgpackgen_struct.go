package msgpackgen

import "time"

type StructTest struct {
	A      int
	B      float32
	String string
	Bool   bool
	Uint64 uint64
	Now    time.Time
	Slice  []uint
	Map    map[string]float64
	//ItemData Item
	//Items    []Item
	Interface interface{}
	Piyo      Hoge
}

type Hoge struct {
	Fuga int
}