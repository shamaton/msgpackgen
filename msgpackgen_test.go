package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shamaton/msgpack/def"
	define2 "github.com/shamaton/msgpackgen/internal/fortest/define"
	"github.com/shamaton/msgpackgen/internal/fortest/define/define"
	"github.com/shamaton/msgpackgen/msgpack"
)

var (
	iDir  = "."
	iFile = ""
	oDir  = "."
	oFile = defaultFileName
	ptr   = defaultPointerLevel
)

func TestMain(m *testing.M) {
	testBeforeRegister()
	RegisterGeneratedResolver()

	code := m.Run()

	os.Exit(code)
}

func testBeforeRegister() {
	{
		v := rand.Int()
		b, err := msgpack.Marshal(v)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var vv int
		err = msgpack.Unmarshal(b, &vv)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if v != vv {
			fmt.Println(v, vv, "different")
			os.Exit(1)
		}
	}
	msgpack.SetStructAsArray(true)
	{
		v := rand.Int()
		b, err := msgpack.Marshal(v)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var vv int
		err = msgpack.Unmarshal(b, &vv)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if v != vv {
			fmt.Println(v, vv, "different")
			os.Exit(1)
		}
	}
	msgpack.SetStructAsArray(false)
}

func TestGenerateCodeErrorInput(t *testing.T) {
	{
		d := "./noname"

		err := generate(d, iFile, oDir, oFile, ptr, true, false, false)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "input directory error") {
			t.Fatal(err)
		}
	}
	{
		d := "./noname"
		f := "foo.go"

		err := generate(d, f, oDir, oFile, ptr, true, false, false)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "at same time") {
			t.Fatal(err)
		}
	}
	{
		d := "main.go"

		err := generate(d, iFile, oDir, oFile, ptr, true, false, false)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "is not directory") {
			t.Fatal(err)
		}
	}
	{
		f := "foo.go"

		err := generate(iDir, f, oDir, oFile, ptr, true, false, false)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "input file error") {
			t.Fatal(err)
		}
	}
	{
		f := "internal"

		err := generate(iDir, f, oDir, oFile, ptr, true, false, false)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "is a directory") {
			t.Fatal(err)
		}
	}
	{
		f := "./testdata/test.sh"

		err := generate(iDir, f, oDir, oFile, ptr, true, false, false)
		if err == nil {
			t.Fatal("error has to return")
		}
		if !strings.Contains(err.Error(), "is not .go file") {
			t.Fatal(err)
		}
	}
}

func TestGenerateCodeGoPathOutside(t *testing.T) {

	g := os.Getenv("GOPATH")
	path := os.Getenv("PATH")
	err := os.Setenv("GOPATH", path)
	if err != nil {
		t.Fatal(err)
	}

	err = generate(iDir, iFile, oDir, oFile, ptr, true, false, false)
	if err == nil {
		t.Fatal("error has to return")
	}
	if !strings.Contains(err.Error(), "outside gopath") {
		t.Fatal(err)
	}

	err = os.Setenv("GOPATH", g)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateCodeDuplicateTag(t *testing.T) {

	f := "./testdata/def.go"

	err := generate(iDir, f, oDir, oFile, ptr, true, false, false)
	if err == nil {
		t.Fatal("error has to return")
	}
	if !strings.Contains(err.Error(), "duplicate tags") {
		t.Fatal(err)
	}
}

func TestGenerateCodeDryRun(t *testing.T) {

	err := generate(iDir, iFile, "", oFile, -1, true, false, false)
	if err != nil {
		t.Fatal("error has to return")
	}
}

func TestGenerateCodeOK(t *testing.T) {
	var err error
	err = flag.CommandLine.Set("strict", "true")
	if err != nil {
		t.Fatal(err)
	}
	err = flag.CommandLine.Set("v", "true")
	if err != nil {
		t.Fatal(err)
	}
	err = flag.CommandLine.Set("pointer", "2")
	if err != nil {
		t.Fatal(err)
	}
	err = flag.CommandLine.Set("output-file", "resolver_test.go")
	if err != nil {
		t.Fatal(err)
	}

	// diff resolver_test.go main.go | wc -l
	main()
}

func TestSwitchDefaultBehaviour(t *testing.T) {
	msgpack.SetStructAsArray(false)

	v := Inside{Int: 1}
	b1, err := msgpack.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	msgpack.SetStructAsArray(true)

	b2, err := msgpack.Marshal(v)
	if err != nil {
		t.Error(err)
	}
	msgpack.SetStructAsArray(false)

	if b1[0] != 0x81 || b2[0] != 0x91 {
		t.Fatalf("format may be different 0x%x, 0x%x", b1[0], b2[0])
	}

}

func TestInt(t *testing.T) {
	v := TestingValue{Int: -8, Int8: math.MinInt8, Int16: math.MinInt16}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{Int: -108, Int8: math.MaxInt8, Int16: math.MaxInt16}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{Int: -30108, Int32: math.MinInt32, Int64: math.MinInt64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{Int: -1030108, Int32: math.MaxInt32, Int64: math.MaxInt64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{
		Int:                  math.MinInt64 + 12345,
		Abcdefghijabcdefghij: rand.Int(),
		AbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghij: rand.Int(),
	}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	if v.Function() != v.Int*2 {
		t.Errorf("value diffrent %d, %d", v.Function(), v.Int*2)
	}

	{
		var r TestingInt
		r.I = 1234
		err := msgpack.UnmarshalAsArray([]byte{0x91, def.Nil}, &r)
		if err != nil {
			t.Error(err)
		}
		if r.I != 0 {
			t.Errorf("not equal %v", r)
		}

		err = msgpack.UnmarshalAsArray([]byte{0x91, def.True}, &r)
		if err == nil || !strings.Contains(err.Error(), "AsInt") {
			t.Error("something wrong", err)
		}
	}
}

func TestUint(t *testing.T) {
	v := TestingValue{Uint: 8, Uint8: math.MaxUint8}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{Uint: 130, Uint16: math.MaxUint16}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{Uint: 30130, Uint32: math.MaxUint32}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{Uint: 1030130, Uint64: math.MaxUint64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{Uint: math.MaxUint64 - 12345}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

	{
		vv := Inside{Int: -1}
		var r TestingUint
		b, err := msgpack.MarshalAsArray(vv)
		if err != nil {
			t.Error(err)
		}
		err = msgpack.UnmarshalAsArray(b, &r)
		if err != nil {
			t.Error(err)
		}
		if r.U != math.MaxUint64 {
			t.Errorf("not equal %v", r)
		}
	}

	{
		vv := Inside{Int: math.MinInt8}
		var r TestingUint
		b, err := msgpack.MarshalAsArray(vv)
		if err != nil {
			t.Error(err)
		}
		err = msgpack.UnmarshalAsArray(b, &r)
		if err != nil {
			t.Error(err)
		}
		if r.U != math.MaxUint64+math.MinInt8+1 {
			t.Errorf("not equal %v", r)
		}
	}

	{
		vv := Inside{Int: math.MinInt16}
		var r TestingUint
		b, err := msgpack.MarshalAsArray(vv)
		if err != nil {
			t.Error(err)
		}
		err = msgpack.UnmarshalAsArray(b, &r)
		if err != nil {
			t.Error(err)
		}
		if r.U != math.MaxUint64+math.MinInt16+1 {
			t.Errorf("not equal %v", r)
		}
	}

	{
		vv := Inside{Int: math.MinInt32}
		var r TestingUint
		b, err := msgpack.MarshalAsArray(vv)
		if err != nil {
			t.Error(err)
		}
		err = msgpack.UnmarshalAsArray(b, &r)
		if err != nil {
			t.Error(err)
		}
		if r.U != math.MaxUint64+math.MinInt32+1 {
			t.Errorf("not equal %v", r)
		}
	}
	{
		vv := Inside{Int: math.MinInt32 - 1}
		var r TestingUint
		b, err := msgpack.MarshalAsArray(vv)
		if err != nil {
			t.Error(err)
		}
		err = msgpack.UnmarshalAsArray(b, &r)
		if err != nil {
			t.Error(err)
		}
		if r.U != math.MaxUint64+math.MinInt32 {
			t.Errorf("not equal %v", r)
		}
	}
	{
		var r TestingUint
		r.U = 1234
		err := msgpack.UnmarshalAsArray([]byte{0x91, def.Nil}, &r)
		if err != nil {
			t.Error(err)
		}
		if r.U != 0 {
			t.Errorf("not equal %v", r)
		}

		err = msgpack.UnmarshalAsArray([]byte{0x91, def.True}, &r)
		if err == nil {
			t.Errorf("error must occur")
		}
		if !strings.Contains(err.Error(), "AsUint") {
			t.Error(err)
		}
	}
}

func TestFloat(t *testing.T) {
	v := TestingValue{Float32: 0, Float64: 0}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{Float32: -1, Float64: -1}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{Float32: math.SmallestNonzeroFloat32, Float64: math.SmallestNonzeroFloat64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{Float32: math.MaxFloat32, Float64: math.MaxFloat64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

	{
		vv := Inside{Int: 1}
		var r TestingFloat32
		b, err := msgpack.MarshalAsArray(vv)
		if err != nil {
			t.Error(err)
		}
		err = msgpack.UnmarshalAsArray(b, &r)
		if err != nil {
			t.Error(err)
		}
		if r.F != 1 {
			t.Errorf("not equal %v", r)
		}
	}
	{
		vv := Inside{Int: -1}
		var r TestingFloat32
		b, err := msgpack.MarshalAsArray(vv)
		if err != nil {
			t.Error(err)
		}
		err = msgpack.UnmarshalAsArray(b, &r)
		if err != nil {
			t.Error(err)
		}
		if r.F != -1 {
			t.Errorf("not equal %v", r)
		}
	}

	{
		vv := Inside{Int: 1}
		var r TestingFloat64
		b, err := msgpack.MarshalAsArray(vv)
		if err != nil {
			t.Error(err)
		}
		err = msgpack.UnmarshalAsArray(b, &r)
		if err != nil {
			t.Error(err)
		}
		if r.F != 1 {
			t.Errorf("not equal %v", r)
		}
	}
	{
		vv := Inside{Int: -1}
		var r TestingFloat64
		b, err := msgpack.MarshalAsArray(vv)
		if err != nil {
			t.Error(err)
		}
		err = msgpack.UnmarshalAsArray(b, &r)
		if err != nil {
			t.Error(err)
		}
		if r.F != -1 {
			t.Errorf("not equal %v", r)
		}
	}
	{
		vv := TestingFloat32{F: 1.23}
		var r TestingFloat64
		b, err := msgpack.MarshalAsArray(vv)
		if err != nil {
			t.Error(err)
		}
		err = msgpack.UnmarshalAsArray(b, &r)
		if err != nil {
			t.Error(err)
		}
		if float32(r.F) != 1.23 {
			t.Errorf("not equal %v", r)
		}
	}
	{
		var r TestingFloat32
		err := msgpack.UnmarshalAsArray([]byte{0x91, def.Nil}, &r)
		if err != nil {
			t.Error(err)
		}
		if r.F != 0 {
			t.Errorf("not equal %v", r)
		}
	}
	{
		var r TestingFloat64
		err := msgpack.UnmarshalAsArray([]byte{0x91, def.Nil}, &r)
		if err != nil {
			t.Error(err)
		}
		if r.F != 0 {
			t.Errorf("not equal %v", r)
		}
	}
	{
		var r TestingFloat32
		err := msgpack.UnmarshalAsArray([]byte{0x91, def.True}, &r)
		if err == nil {
			t.Errorf("error must occur")
		}
		if !strings.Contains(err.Error(), "AsFloat32") {
			t.Error(err)
		}
	}
	{
		var r TestingFloat64
		err := msgpack.UnmarshalAsArray([]byte{0x91, def.True}, &r)
		if err == nil {
			t.Errorf("error must occur")
		}
		if !strings.Contains(err.Error(), "AsFloat64") {
			t.Error(err)
		}
	}
}

func TestString(t *testing.T) {
	base := "abcdefghijklmnopqrstuvwxyz12345"
	v := TestingValue{String: ""}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{String: strings.Repeat(base, 1)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{String: strings.Repeat(base, 8)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{String: strings.Repeat(base, (math.MaxUint16/len(base))-1)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{String: strings.Repeat(base, (math.MaxUint16/len(base))+1)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

	{
		var r TestingString
		r.S = "setset"
		err := msgpack.UnmarshalAsArray([]byte{0x91, def.Nil}, &r)
		if err != nil {
			t.Error(err)
		}
		if r.S != "" {
			t.Errorf("not equal %v", r)
		}

		err = msgpack.UnmarshalAsArray([]byte{0x91, def.True}, &r)
		if err == nil || !strings.Contains(err.Error(), "StringByteLength") {
			t.Error("something wrong", err)
		}
	}
}

func TestBool(t *testing.T) {
	v := TestingValue{Bool: true}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{Bool: false}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

}

func TestByteRune(t *testing.T) {
	v := TestingValue{Byte: 127, Rune: 'a'}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

	{
		var r TestingBool
		r.B = true
		err := msgpack.UnmarshalAsArray([]byte{0x91, def.Nil}, &r)
		if err != nil {
			t.Error(err)
		}
		if r.B != false {
			t.Errorf("not equal %v", r)
		}
		err = msgpack.UnmarshalAsArray([]byte{0x91, def.Uint8, 0x01}, &r)
		if err == nil {
			t.Errorf("error must occur")
		}
		if !strings.Contains(err.Error(), "AsBool") {
			t.Error(err)
		}
	}
}

func TestComplex(t *testing.T) {
	v := TestingValue{Complex64: complex(1, 2), Complex128: complex(3, 4)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{
		Complex64:  complex(math.MaxFloat32, math.SmallestNonzeroFloat32),
		Complex128: complex(math.MaxFloat64, math.SmallestNonzeroFloat64),
	}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	msgpack.SetComplexTypeCode(-123)
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

	b64, err := msgpack.MarshalAsArray(TestingComplex64{})
	if err != nil {
		t.Error(err)
	}
	b128, err := msgpack.MarshalAsArray(TestingComplex128{})
	if err != nil {
		t.Error(err)
	}

	msgpack.SetComplexTypeCode(-122)

	{
		var r TestingComplex64
		err = msgpack.UnmarshalAsArray(b64, &r)
		if err == nil || !strings.Contains(err.Error(), "fixext8") {
			t.Error(err)
		}
		err := msgpack.UnmarshalAsArray([]byte{0x91, def.True}, &r)
		if err == nil {
			t.Errorf("error must occur")
		}
		if !strings.Contains(err.Error(), "AsComplex64") {
			t.Error(err)
		}
	}
	{
		var r TestingComplex128
		err = msgpack.UnmarshalAsArray(b128, &r)
		if err == nil || !strings.Contains(err.Error(), "fixext16") {
			t.Error(err)
		}
		err := msgpack.UnmarshalAsArray([]byte{0x91, def.True}, &r)
		if err == nil {
			t.Errorf("error must occur")
		}
		if !strings.Contains(err.Error(), "AsComplex128") {
			t.Error(err)
		}
	}
}

func TestMap(t *testing.T) {
	v := TestingValue{
		MapIntInt: map[string]int{"1": 2, "3": 4, "5": 6, "7": 8, "9": 10},
	}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{
		MapIntInt: map[string]int{},
	}
	for i := 0; i < 1000; i++ {
		v.MapIntInt[fmt.Sprint(i)] = i + 1
	}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

	v = TestingValue{
		MapIntInt: map[string]int{},
	}
	for i := 0; i < math.MaxUint16+1; i++ {
		v.MapIntInt[fmt.Sprint(i)] = i + 1
	}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{
		MapIntInt: nil,
	}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

	{
		var r TestingMap
		err := msgpack.UnmarshalAsArray([]byte{0x91, def.True}, &r)
		if err == nil {
			t.Errorf("error must occur")
		}
		if !strings.Contains(err.Error(), "MapLength") {
			t.Error(err)
		}
	}
}

func TestPointerValue(t *testing.T) {
	{
		_v := 123
		v := TestingValue{Pint: &_v}
		if err := checkValue(v); err != nil {
			t.Error(err)
		}
	}
	{
		_v := "this is pointer"
		__v := &_v
		v := TestingValue{P2string: &__v}
		if err := checkValue(v); err != nil {
			t.Error(err)
		}
	}
	{
		_v := float32(1.234)
		__v := &_v
		___v := &__v
		____v := &___v
		v := TestingValue{P3float32: ____v}
		if err := checkValue(v); err != nil {
			t.Error(err)
		}
	}

	check := func(v TestingValue) error {
		var v1, v2 TestingValue
		f1 := func() (bool, interface{}, interface{}) {
			return v1.P3float32 == nil, v1.P3float32, nil
		}
		f2 := func() (bool, interface{}, interface{}) {
			return v2.P3float32 == nil, v1.P3float32, nil
		}

		return _checkValue(v, &v1, &v2, f1, f2)
	}
	{

		_v := float32(1.234)
		__v := &_v
		___v := &__v
		____v := &___v

		__v = nil
		v := TestingValue{P3float32: ____v}
		if err := check(v); err != nil {
			t.Error(err)
		}
		___v = nil
		v = TestingValue{P3float32: ____v}
		if err := check(v); err != nil {
			t.Error(err)
		}
		____v = nil
		v = TestingValue{P3float32: ____v}
		if err := check(v); err != nil {
			t.Error(err)
		}
	}
	{
		_v := make([]*int, 100)
		for i := range _v {
			if i%2 == 0 {
				__v := i
				_v[i] = &__v
			} else {
				_v[i] = nil
			}
		}

		v := TestingValue{IntPointers: _v}
		if err := checkValue(v); err != nil {
			t.Error(err)
		}
	}

	check2 := func(v TestingValue) error {
		var v1, v2 TestingValue
		f1 := func() (bool, interface{}, interface{}) {
			mp := map[uint]string{}
			for _k, _v := range v1.MapPointers {
				__v := *_v
				mp[*_k] = *__v
			}
			for _k, _v := range v.MapPointers {
				__v := *_v
				if str, ok := mp[*_k]; !ok || str != *__v {
					return false, str, *__v
				}
			}
			return true, nil, nil
		}

		f2 := func() (bool, interface{}, interface{}) {
			mp := map[uint]string{}
			for _k, _v := range v2.MapPointers {
				__v := *_v
				mp[*_k] = *__v
			}
			for _k, _v := range v.MapPointers {
				__v := *_v
				if str, ok := mp[*_k]; !ok || str != *__v {
					return false, str, *__v
				}
			}
			return true, nil, nil
		}

		return _checkValue(v, &v1, &v2, f1, f2)
	}
	{
		_v := make(map[*uint]**string, 10)
		for i := 0; i < 10; i++ {
			k := uint(i)
			v := fmt.Sprint(i, i, i)
			vv := &v
			_v[&k] = &vv
		}

		v := TestingValue{MapPointers: _v}
		if err := check2(v); err != nil {
			t.Error(err)
		}
	}
}

func TestTime(t *testing.T) {
	{
		v := TestingTime{Time: time.Now()}
		b, err := msgpack.Marshal(v)
		if err != nil {
			t.Error(err)
		}
		var _v TestingTime
		err = msgpack.Unmarshal(b, &_v)
		if err != nil {
			t.Error(err)
		}
		if v.Time.UnixNano() != _v.Time.UnixNano() {
			t.Errorf("time different %v, %v", v.Time, _v.Time)
		}
	}
	{
		vv := TestingTime{}
		b, err := msgpack.MarshalAsArray(vv)
		if err != nil {
			t.Error(err)
		}
		var _vv TestingTime
		err = msgpack.UnmarshalAsArray(b, &_vv)
		if err != nil {
			t.Error(err)
		}
		if vv.Time.UnixNano() != _vv.Time.UnixNano() {
			t.Errorf("time different %v, %v", vv.Time, _vv.Time)
		}
	}
	{
		v := define2.AA{Time: time.Now()}
		b, err := msgpack.Marshal(v)
		if err != nil {
			t.Error(err)
		}
		var _v define2.AA
		err = msgpack.Unmarshal(b, &_v)
		if err != nil {
			t.Error(err)
		}
		if v.UnixNano() != _v.UnixNano() {
			t.Errorf("time different %v, %v", v.Time, _v.Time)
		}
	}
	{
		now := time.Now().Unix()
		nowByte := make([]byte, 4)
		binary.BigEndian.PutUint32(nowByte, uint32(now))

		var r TestingTime
		c := def.TimeStamp
		err := msgpack.UnmarshalAsArray(append([]byte{def.FixArray + 1, def.Fixext4, byte(c)}, nowByte...), &r)
		if err != nil {
			t.Error(err)
		}
		if r.Time.Unix() != now {
			t.Error("different time", r.Time.Unix(), now, nowByte)
		}

		_, err = msgpack.MarshalAsArray(r)
		if err != nil {
			t.Error(err)
		}

		err = msgpack.UnmarshalAsArray(append([]byte{def.FixArray + 1, def.Fixext4, 3}, nowByte...), &r)
		if err == nil || !strings.Contains(err.Error(), "fixext4. time type is different") {
			t.Error("something wrong", err)
		}

		err = msgpack.UnmarshalAsArray([]byte{def.FixArray + 1, def.Fixext8, 3}, &r)
		if err == nil || !strings.Contains(err.Error(), "fixext8. time type is different") {
			t.Error("something wrong", err)
		}

		err = msgpack.UnmarshalAsArray([]byte{def.FixArray + 1, def.Ext8, 11}, &r)
		if err == nil || !strings.Contains(err.Error(), "ext8. time ext length is different") {
			t.Error("something wrong", err)
		}

		err = msgpack.UnmarshalAsArray([]byte{def.FixArray + 1, def.Ext8, 12, 3}, &r)
		if err == nil || !strings.Contains(err.Error(), "ext8. time type is different") {
			t.Error("something wrong", err)
		}

		nanoByte := make([]byte, 64)
		for i := range nanoByte[:30] {
			nanoByte[i] = 0xff
		}
		b := append([]byte{def.FixArray + 1, def.Fixext8, byte(c)}, nanoByte...)
		err = msgpack.UnmarshalAsArray(b, &r)
		if err == nil || !strings.Contains(err.Error(), "in timestamp 64 formats") {
			t.Error(err)
		}

		nanoByte = make([]byte, 96)
		for i := range nanoByte[:32] {
			nanoByte[i] = 0xff
		}
		b = append([]byte{def.FixArray + 1, def.Ext8, byte(12), byte(c)}, nanoByte...)
		err = msgpack.UnmarshalAsArray(b, &r)
		if err == nil || !strings.Contains(err.Error(), "in timestamp 96 formats") {
			t.Error(err)
		}

		err = msgpack.UnmarshalAsArray([]byte{def.FixArray + 1, def.Fixext1}, &r)
		if err == nil || !strings.Contains(err.Error(), "AsDateTime") {
			t.Error("something wrong", err)
		}
	}
}

func TestTimePointer(t *testing.T) {
	now := time.Now()
	add := now.Add(1 * time.Minute)
	v := TestingTimePointer{Time: &add}
	b, err := msgpack.Marshal(v)
	if err != nil {
		t.Error(err)
	}
	var _v TestingTimePointer
	err = msgpack.Unmarshal(b, &_v)
	if err != nil {
		t.Error(err)
	}
	if v.Time.UnixNano() != _v.Time.UnixNano() {
		t.Errorf("time different %v, %v", v.Time, _v.Time)
	}

	vv := TestingTimePointer{}
	b, err = msgpack.MarshalAsArray(vv)
	if err != nil {
		t.Error(err)
	}
	var _vv TestingTimePointer
	err = msgpack.UnmarshalAsArray(b, &_vv)
	if err != nil {
		t.Error(err)
	}
	if vv.Time != nil || _vv.Time != nil {
		t.Errorf("time different %v, %v", vv.Time, _vv.Time)
	}
}

func checkValue(v TestingValue, eqs ...func() (bool, interface{}, interface{})) error {
	var v1, v2 TestingValue
	return _checkValue(v, &v1, &v2, eqs...)
}

func TestSlice(t *testing.T) {

	f := func(l int) []int8 {
		slice := make([]int8, l)
		for i := range slice {
			slice[i] = int8(rand.Intn(math.MaxUint8) + math.MinInt8)
		}
		return slice
	}

	check := func(v TestingValue) error {
		var v1, v2 TestingValue
		return _checkValue(v, &v1, &v2)
	}

	v := TestingValue{Slice: f(15)}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{Slice: f(150)}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{Slice: f(math.MaxUint16 + 1)}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = TestingValue{Slice: nil}
	if err := check(v); err != nil {
		t.Error(err)
	}

	v = TestingValue{}
	v.Bytes = make([]byte, 100)
	for i := range v.Bytes {
		v.Bytes[i] = byte(rand.Intn(255))
	}
	if err := check(v); err != nil {
		t.Error(err)
	}

	v = TestingValue{}
	v.Bytes = make([]byte, math.MaxUint8+1)
	for i := range v.Bytes {
		v.Bytes[i] = byte(rand.Intn(255))
	}
	if err := check(v); err != nil {
		t.Error(err)
	}

	v = TestingValue{}
	v.Bytes = make([]byte, math.MaxUint16+1)
	for i := range v.Bytes {
		v.Bytes[i] = byte(rand.Intn(255))
	}
	if err := check(v); err != nil {
		t.Error(err)
	}

	v = TestingValue{}
	v.Bytes = make([]byte, math.MaxUint32+1)
	_, err := msgpack.MarshalAsArray(v)
	if err == nil || !strings.Contains(err.Error(), "not support this array length") {
		t.Error("something wrong", err)
	}

	{
		var r TestingSlice
		err = msgpack.UnmarshalAsArray([]byte{0x91, def.True}, &r)
		if err == nil {
			t.Errorf("error must occur")
		}
		if !strings.Contains(err.Error(), "SliceLength") {
			t.Error(err)
		}
	}
	{
		vv := TestingSlice{}
		vv.Slice = make([]int8, math.MaxUint32+1)
		_, err := msgpack.MarshalAsArray(v)
		if err == nil || !strings.Contains(err.Error(), "not support this array length") {
			t.Error("something wrong", err)
		}
	}

}

func TestArray(t *testing.T) {
	check := func(v TestingArrays) error {
		var v1, v2 TestingArrays
		return _checkValue(v, &v1, &v2)
	}

	var v TestingArrays

	v = TestingArrays{}
	for i := range v.Array1 {
		v.Array1[i] = float32(rand.Intn(0xff))
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = TestingArrays{}
	for i := range v.Array2 {
		v.Array2[i] = "a"
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = TestingArrays{}
	for i := range v.Array3 {
		v.Array3[i] = rand.Intn(0xff) > 0x7f
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = TestingArrays{}
	for i := range v.Array4 {
		v.Array4[i] = rand.Intn(math.MaxInt32)
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = TestingArrays{}
	for i := range v.Array5 {
		v.Array5[i] = rand.Intn(math.MaxInt32)
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = TestingArrays{}
	for i := range v.Array6 {
		v.Array6[i] = rand.Intn(math.MaxInt32)
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
}

func _checkValue(v interface{}, u1, u2 interface{}, eqs ...func() (bool, interface{}, interface{})) error {
	b1, b2, err1, err2 := marshal(v, v)
	if err1 != nil {
		return fmt.Errorf("marshal to b1 failed %v", err1)
	}
	if err2 != nil {
		return fmt.Errorf("marshal to b2 failed %v", err2)
	}

	err1, err2 = unmarshal(b1, b2, u1, u2)
	if err1 != nil {
		return fmt.Errorf("unmarshal to u1 failed %v", err1)
	}
	if err2 != nil {
		return fmt.Errorf("unmarshal to u2 failed %v", err2)
	}

	if len(eqs) < 2 {
		if !reflect.DeepEqual(v, reflect.ValueOf(u1).Elem().Interface()) {
			return fmt.Errorf("not equal u1 %v, %v", v, u1)
		}
		if !reflect.DeepEqual(v, reflect.ValueOf(u2).Elem().Interface()) {
			return fmt.Errorf("not equal u2 %v, %v", v, u2)
		}
	} else {
		if b, _v1, _v2 := eqs[0](); !b {
			return fmt.Errorf("not equal u1 %v, %v", _v1, _v2)
		}
		if b, _v1, _v2 := eqs[1](); !b {
			return fmt.Errorf("not equal u2 %v, %v", _v1, _v2)
		}
	}
	return nil
}

func TestStruct(t *testing.T) {
	check := func(v TestingStruct) (TestingStruct, TestingStruct, error) {
		f := func() (bool, interface{}, interface{}) {
			return true, nil, nil
		}
		var r1, r2 TestingStruct
		err := _checkValue(v, &r1, &r2, f, f)
		return r1, r2, err
	}

	v := TestingStruct{}
	v.Int = rand.Int()
	v1, v2, err := check(v)
	if err != nil {
		t.Error(err)
	}
	if v.Int != v1.Int || v.Int != v2.Int {
		t.Error("value different", v.Int, v1.Int, v2.Int)
	}
	if v.Int == v1.Inside.Int || v.Int == v2.Outside.Int {
		t.Error("value something wrong", v.Int, v1.Inside.Int, v2.Outside.Int)
	}

	v = TestingStruct{}
	v.Inside = Inside{Int: rand.Int()}
	v1, v2, err = check(v)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(v.Inside, v1.Inside) || !reflect.DeepEqual(v.Inside, v2.Inside) {
		t.Error("value different", v.Inside, v1.Inside, v2.Inside)
	}
	if v.Inside.Int == v1.Int || v.Inside.Int == v2.Outside.Int {
		t.Error("value something wrong", v.Inside.Int, v1.Int, v2.Outside.Int)
	}

	v = TestingStruct{}
	v.Outside = Outside{Int: rand.Int()}
	v1, v2, err = check(v)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(v.Outside, v1.Outside) || !reflect.DeepEqual(v.Outside, v2.Outside) {
		t.Error("value different", v.Outside, v1.Outside, v2.Outside)
	}
	if v.Outside.Int == v1.Int || v.Outside.Int == v2.Inside.Int {
		t.Error("value something wrong", v.Outside.Int, v1.Int, v2.Inside.Int)
	}

	v = TestingStruct{}
	v.A = define2.A{Int: rand.Int(), B: define.B{Int: rand.Int()}}
	v1, v2, err = check(v)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(v.A, v1.A) || !reflect.DeepEqual(v.A, v2.A) {
		t.Error("value different", v.A, v1.A, v2.A)
	}
	if v.A.B.Int != v1.B.Int || v.B.Int != v2.A.B.Int {
		t.Error("value something wrong", v.A.Int, v1.Int, v2.Int)
	}
	if v.A.Int == v1.Int || v.A.Int == v2.Int {
		t.Error("value something wrong", v.A.Int, v1.Int, v2.Int)
	}
	if v.B.Int == v1.Int || v.B.Int == v2.Int {
		t.Error("value something wrong", v.A.Int, v1.Int, v2.Int)
	}

	v = TestingStruct{}
	v.BB = define.DotImport{Int: rand.Int()}
	v1, v2, err = check(v)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(v.BB, v1.BB) || !reflect.DeepEqual(v.BB, v2.BB) {
		t.Error("value different", v.BB, v1.BB, v2.BB)
	}
	if v.BB.Int == v1.Int || v.BB.Int == v2.Int {
		t.Error("value something wrong", v.BB.Int, v1.Int, v2.Int)
	}

	v = TestingStruct{}
	v.Time = define.Time{Int: rand.Int()}
	v1, v2, err = check(v)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(v.Time, v1.Time) || !reflect.DeepEqual(v.Time, v2.Time) {
		t.Error("value different", v.Time, v1.Time, v2.Time)
	}
	if v.Time.Int == v1.Int || v.Time.Int == v2.Int {
		t.Error("value something wrong", v.Time.Int, v1.Int, v2.Int)
	}

	v = TestingStruct{}
	v.R = &Recursive{Int: rand.Int()}
	v.R.R = &Recursive{Int: rand.Int()}
	v1, v2, err = check(v)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(v.R, v1.R) || !reflect.DeepEqual(v.R, v2.R) {
		t.Error("value different", v.R, v1.R, v2.R)
	}
	if v.R.Int != v1.R.Int || v.R.Int != v2.R.Int {
		t.Error("value different", v.R.Int, v1.R.Int, v2.R.Int)
	}
	if v.R.R.Int != v1.R.R.Int || v.R.R.Int != v2.R.R.Int {
		t.Error("value different", v.R.R.Int, v1.R.R.Int, v2.R.R.Int)
	}
	if v.R.R.R != nil || v1.R.R.R != nil || v2.R.R.R != nil {
		t.Error("value different", v.R.R.R, v1.R.R.R, v2.R.R.R)
	}

	{
		b1 := []byte{def.Array32, 0x00, 0x00, 0x00, 0x02}
		var r Inside
		err := msgpack.UnmarshalAsArray(b1, &r)
		if err == nil || !strings.Contains(err.Error(), "data length wrong") {
			t.Error("something wrong", err)
		}
		b2 := []byte{def.Map32, 0x00, 0x00, 0x00, 0x02}
		err = msgpack.UnmarshalAsArray(b2, &r)
		if err == nil || !strings.Contains(err.Error(), "data length wrong") {
			t.Error("something wrong", err)
		}
	}
}

func TestTag(t *testing.T) {
	v := TestingTag{Tag: 1, Ignore: rand.Int(), Omit: rand.Int()}
	b1, b2, e1, e2 := marshal(v, v)
	if e1 != nil || e2 != nil {
		t.Error(e1, e2)
	}
	if len(b1) <= 6 {
		t.Errorf("tag does not recognize %v % x", v, b1)
	}
	if len(b2) != 2 {
		t.Errorf("something wrong %v % x", v, b2)
	}

	var v1, v2 TestingTag
	e1, e2 = unmarshal(b1, b2, &v1, &v2)
	if e1 != nil || e2 != nil {
		t.Error(e1, e2)
	}
	if v.Tag != v1.Tag || v.Tag != v2.Tag {
		t.Errorf("not equal value %d, %d, %d", v.Tag, v1.Tag, v2.Tag)
	}
	if v.Ignore == v1.Ignore || v.Ignore == v2.Ignore {
		t.Errorf("equal value %d, %d, %d", v.Ignore, v1.Ignore, v2.Ignore)
	}
	if v.Omit == v1.Omit || v.Omit == v2.Omit {
		t.Errorf("equal value %d, %d, %d", v.Omit, v1.Omit, v2.Omit)
	}
}

func TestPointer(t *testing.T) {

	v := TestingValue{Int: -1, Uint: 1}

	//// OK
	// encode single pointer
	b1, b2, err1, err2 := marshal(&v, &v)
	if err1 != nil {
		t.Error(err1)
	}
	if err2 != nil {
		t.Error(err2)
	}

	// decode double pointer
	v1, v2 := new(TestingValue), new(TestingValue)
	err1, err2 = unmarshal(b1, b2, &v1, &v2)
	if err1 != nil {
		t.Error(err1)
	}
	if err2 != nil {
		t.Error(err2)
	}

	// encode double pointer
	b3, b4, err1, err2 := marshal(&v1, &v2)
	if err1 != nil {
		t.Error(err1)
	}
	if err2 != nil {
		t.Error(err2)
	}

	// decode triple pointer
	_v3, _v4 := new(TestingValue), new(TestingValue)
	v3, v4 := &_v3, &_v4
	err1, err2 = unmarshal(b3, b4, &v3, &v4)
	if err1 != nil {
		t.Error(err1)
	}
	if err2 != nil {
		t.Error(err2)
	}

	if !reflect.DeepEqual(v, *v1) {
		t.Error("not equal v1", v, *v1)
	}
	if !reflect.DeepEqual(v, *v2) {
		t.Error("not equal v2", v, *v2)
	}
	if _v := *v3; !reflect.DeepEqual(v, *_v) {
		t.Error("not equal v3", v, _v)
	}
	if _v := *v4; !reflect.DeepEqual(v, *_v) {
		t.Error("not equal v4", v, _v)
	}

	//// NG
	// encode triple pointer
	b5, b6, err1, err2 := marshal(&v3, &v4)
	if err1 != nil && !strings.Contains(err1.Error(), "strict") {
		t.Error(err1)
	}
	if err1 == nil {
		t.Error("error should occur at marshalling v3 pointer")
	}
	if err2 != nil && !strings.Contains(err2.Error(), "strict") {
		t.Error(err2)
	}
	if err2 == nil {
		t.Error("error should occur at marshalling v4 pointer")
	}

	// decode quad pointer
	__v5, __v6 := new(TestingValue), new(TestingValue)
	_v5, _v6 := &__v5, &__v6
	v5, v6 := &_v5, &_v6
	err1, err2 = unmarshal(b5, b6, &v5, &v6)
	if err1 != nil && !strings.Contains(err1.Error(), "strict") {
		t.Error(err1)
	}
	if err1 == nil {
		t.Error("error should occur at unmarshalling b5 pointer")
	}
	if err2 != nil && !strings.Contains(err2.Error(), "strict") {
		t.Error(err2)
	}
	if err2 == nil {
		t.Error("error should occur at unmarshalling b6 pointer")
	}
}

func TestNotGenerated(t *testing.T) {
	err := checkUndefined(NotGenerated1{}, NotGenerated2{}, &NotGenerated1{}, &NotGenerated2{})
	if err != nil {
		t.Error(err)
	}
	err = checkUndefined(NotGenerated3{}, NotGenerated4{}, &NotGenerated3{}, &NotGenerated4{})
	if err != nil {
		t.Error(err)
	}
	err = checkUndefined(NotGenerated5{}, NotGenerated6{}, &NotGenerated5{}, &NotGenerated6{})
	if err != nil {
		t.Error(err)
	}
	err = checkUndefined(NotGenerated7{}, NotGenerated8{}, &NotGenerated7{}, &NotGenerated8{})
	if err != nil {
		t.Error(err)
	}
}

func TestPrivate(t *testing.T) {
	v := Private{}
	v.SetInt()
	b1, b2, err1, err2 := marshal(v, v)
	if err1 != nil || err2 != nil {
		t.Errorf("somthing wrong %v, %v", err1, err2)
	}
	if len(b1) != 1 || b1[0] != 0x80 {
		t.Errorf("data is wrong % x", b1)
	}
	if len(b2) != 1 || b2[0] != 0x90 {
		t.Errorf("data is wrong % x", b2)
	}

	vc1 := &v
	vc2 := &vc1
	vc3 := &vc2
	vc4 := &vc3
	err := forCoverage(vc1, vc2, vc3, vc4)
	if err != nil {
		t.Error(err)
	}
}

func forCoverage(v1, v2, v3, v4 interface{}) error {

	// encode single pointer
	b1, b2, err1, err2 := marshal(v1, v1)
	if err1 != nil {
		return fmt.Errorf("marshal to b1 error %v", err1)
	}
	if err2 != nil {
		return fmt.Errorf("marshal to b2 error %v", err2)
	}

	// decode double pointer
	err1, err2 = unmarshal(b1, b2, v2, v2)
	if err1 != nil {
		return fmt.Errorf("unmarshal from b1 error %v", err1)
	}
	if err2 != nil {
		return fmt.Errorf("unmarshal from b2 error %v", err2)
	}

	// encode double pointer
	b3, b4, err1, err2 := marshal(v2, v2)
	if err1 != nil {
		return fmt.Errorf("marshal to b3 error %v", err1)
	}
	if err2 != nil {
		return fmt.Errorf("marshal to b4 error %v", err2)
	}

	// decode triple pointer
	err1, err2 = unmarshal(b3, b4, v3, v3)
	if err1 != nil {
		return fmt.Errorf("unmarshal from b3 error %v", err1)
	}
	if err2 != nil {
		return fmt.Errorf("unmarshal from b4 error %v", err2)
	}

	//// NG
	// encode triple pointer / decode quad pointer
	err := checkUndefined(v3, v3, v4, v4)
	if err != nil {
		return fmt.Errorf("for coverage; %v", err)
	}
	return nil
}

func checkUndefined(m1, m2, u1, u2 interface{}) error {
	b1, b2, err1, err2 := marshal(m1, m2)
	if err1 != nil && !strings.Contains(err1.Error(), "use strict option") {
		return fmt.Errorf("check undefined: marshal to b1 error %v", err1)
	}
	if err1 == nil {
		return fmt.Errorf("check undefined: error should occur at marshalling m1")
	}
	if err2 != nil && !strings.Contains(err2.Error(), "use strict option") {
		return fmt.Errorf("check undefined: marshal to b2 error %v", err2)
	}
	if err2 == nil {
		return fmt.Errorf("check undefined: error should occur at marshalling m2")
	}

	err1, err2 = unmarshal(b1, b2, u1, u2)
	if err1 != nil && !strings.Contains(err1.Error(), "use strict option") {
		return fmt.Errorf("check undefined: unmarshal to u1 error %v", err1)
	}
	if err1 == nil {
		return fmt.Errorf("check undefined: error should occur at unmarshalling u1")
	}
	if err2 != nil && !strings.Contains(err2.Error(), "use strict option") {
		return fmt.Errorf("ucheck undefined: unmarshal to u2 error %v", err2)
	}
	if err2 == nil {
		return fmt.Errorf("check undefined: error should occur at unmarshalling u2")
	}
	return nil
}

func marshal(v1, v2 interface{}) ([]byte, []byte, error, error) {
	b1, e1 := msgpack.Marshal(v1)
	b2, e2 := msgpack.MarshalAsArray(v2)
	return b1, b2, e1, e2
}

func unmarshal(b1, b2 []byte, v1, v2 interface{}) (error, error) {
	return msgpack.Unmarshal(b1, v1), msgpack.UnmarshalAsArray(b2, v2)
}
