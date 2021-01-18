package tst_test

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shamaton/msgpackgen/internal/tst"
	"github.com/shamaton/msgpackgen/msgpack"
)

func TestMain(m *testing.M) {
	tst.RegisterGeneratedResolver()

	code := m.Run()

	// resetGeneratedCode()

	os.Exit(code)
}

func resetGeneratedCode() {
	tpl := template.Must(template.New("").Parse(`package tst

import "fmt"

func RegisterGeneratedResolver() {
	fmt.Println("this is dummy.")
}
`))

	file, err := os.Create("./resolver.msgpackgen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = tpl.Execute(file, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func TestInt(t *testing.T) {
	v := tst.ValueChecking{Int: -8, Int8: math.MinInt8, Int16: math.MinInt16}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = tst.ValueChecking{Int: -108, Int8: math.MaxInt8, Int16: math.MaxInt16}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = tst.ValueChecking{Int: -30108, Int32: math.MinInt32, Int64: math.MinInt64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = tst.ValueChecking{Int: -1030108, Int32: math.MaxInt32, Int64: math.MaxInt64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = tst.ValueChecking{Int: math.MinInt64 + 12345}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

}

func TestUint(t *testing.T) {
	v := tst.ValueChecking{Uint: 8, Uint8: math.MaxUint8}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = tst.ValueChecking{Uint: 130, Uint16: math.MaxUint16}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = tst.ValueChecking{Uint: 30130, Uint32: math.MaxUint32}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = tst.ValueChecking{Uint: 1030130, Uint64: math.MaxUint64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = tst.ValueChecking{Uint: math.MaxUint64 - 12345}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

}

func TestFloat(t *testing.T) {
	v := tst.ValueChecking{Float32: 0, Float64: 0}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = tst.ValueChecking{Float32: -1, Float64: -1}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = tst.ValueChecking{Float32: math.SmallestNonzeroFloat32, Float64: math.SmallestNonzeroFloat64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = tst.ValueChecking{Float32: math.MaxFloat32, Float64: math.MaxFloat64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

}

func TestString(t *testing.T) {
	base := "abcdefghijklmnopqrstuvwxyz12345"
	v := tst.ValueChecking{String: ""}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = tst.ValueChecking{String: strings.Repeat(base, 1)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = tst.ValueChecking{String: strings.Repeat(base, 8)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = tst.ValueChecking{String: strings.Repeat(base, (math.MaxUint16/len(base))-1)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = tst.ValueChecking{String: strings.Repeat(base, (math.MaxUint16/len(base))+1)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
}

func TestBool(t *testing.T) {
	v := tst.ValueChecking{Bool: true}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = tst.ValueChecking{Bool: false}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

}

func TestByteRune(t *testing.T) {
	v := tst.ValueChecking{Byte: 127, Rune: 'a'}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
}

func TestComplex(t *testing.T) {
	v := tst.ValueChecking{Complex64: complex(1, 2), Complex128: complex(3, 4)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = tst.ValueChecking{
		Complex64:  complex(math.MaxFloat32, math.SmallestNonzeroFloat32),
		Complex128: complex(math.MaxFloat64, math.SmallestNonzeroFloat64),
	}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
}

func TestTime(t *testing.T) {
	v := tst.TimeChecking{Time: time.Now()}
	b, err := msgpack.Encode(v)
	if err != nil {
		t.Error(err)
	}
	var _v tst.TimeChecking
	err = msgpack.Decode(b, &_v)
	if err != nil {
		t.Error(err)
	}
	if v.Time.UnixNano() != _v.Time.UnixNano() {
		t.Errorf("time different %v, %v", v.Time, _v.Time)
	}

	vv := tst.TimeChecking{}
	b, err = msgpack.EncodeAsArray(vv)
	if err != nil {
		t.Error(err)
	}
	var _vv tst.TimeChecking
	err = msgpack.DecodeAsArray(b, &_vv)
	if err != nil {
		t.Error(err)
	}
	if vv.Time.UnixNano() != _vv.Time.UnixNano() {
		t.Errorf("time different %v, %v", vv.Time, _vv.Time)
	}
}

func checkValue(v tst.ValueChecking) error {
	var v1, v2 tst.ValueChecking
	return _checkValue(v, &v1, &v2)
}

func TestSliceArray(t *testing.T) {

	f := func(l int) []int {
		slice := make([]int, l)
		for i := range slice {
			slice[i] = rand.Intn(math.MaxInt32)
		}
		return slice
	}

	check := func(v tst.SliceArray) error {
		var v1, v2 tst.SliceArray
		return _checkValue(v, &v1, &v2)
	}

	v := tst.SliceArray{Slice: f(15)}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = tst.SliceArray{Slice: f(30015)}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = tst.SliceArray{Slice: f(1030015)}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = tst.SliceArray{}
	for i := range v.Array1 {
		v.Array1[i] = float32(rand.Intn(0xff))
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = tst.SliceArray{}
	for i := range v.Array2 {
		v.Array2[i] = "a"
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = tst.SliceArray{}
	for i := range v.Array3 {
		v.Array3[i] = rand.Intn(0xff) > 0x7f
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = tst.SliceArray{}
	for i := range v.Array4 {
		v.Array4[i] = rand.Intn(math.MaxInt32)
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = tst.SliceArray{}
	for i := range v.Array5 {
		v.Array5[i] = rand.Intn(math.MaxInt32)
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = tst.SliceArray{}
	for i := range v.Array6 {
		v.Array6[i] = rand.Intn(math.MaxInt32)
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
}

func _checkValue(v interface{}, u1, u2 interface{}) error {
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

	if !reflect.DeepEqual(v, reflect.ValueOf(u1).Elem().Interface()) {
		return fmt.Errorf("not equal u1 %v, %v", v, u1)
	}
	if !reflect.DeepEqual(v, reflect.ValueOf(u2).Elem().Interface()) {
		return fmt.Errorf("not equal u2 %v, %v", v, u2)
	}
	return nil
}

func TestPointer(t *testing.T) {

	v := tst.ValueChecking{Int: -1, Uint: 1}

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
	v1, v2 := new(tst.ValueChecking), new(tst.ValueChecking)
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
	_v3, _v4 := new(tst.ValueChecking), new(tst.ValueChecking)
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
	__v5, __v6 := new(tst.ValueChecking), new(tst.ValueChecking)
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
	err := checkUndefined(tst.NotGenerated1{}, tst.NotGenerated2{}, &tst.NotGenerated1{}, &tst.NotGenerated2{})
	if err != nil {
		t.Error(err)
	}
	err = checkUndefined(tst.NotGenerated3{}, tst.NotGenerated4{}, &tst.NotGenerated3{}, &tst.NotGenerated4{})
	if err != nil {
		t.Error(err)
	}
}

func TestPrivate(t *testing.T) {
	v := tst.Private{}
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
	b1, e1 := msgpack.Encode(v1)
	b2, e2 := msgpack.EncodeAsArray(v2)
	return b1, b2, e1, e2
}

func unmarshal(b1, b2 []byte, v1, v2 interface{}) (error, error) {
	return msgpack.Decode(b1, v1), msgpack.DecodeAsArray(b2, v2)
}
