package tst

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

	tst2 "github.com/shamaton/msgpackgen/internal/tst/tst/tst"

	"github.com/shamaton/msgpackgen/internal/tst/tst"

	"github.com/shamaton/msgpackgen/msgpack"
)

func TestMain(m *testing.M) {
	RegisterGeneratedResolver()

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
	v := ValueChecking{Int: -8, Int8: math.MinInt8, Int16: math.MinInt16}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{Int: -108, Int8: math.MaxInt8, Int16: math.MaxInt16}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{Int: -30108, Int32: math.MinInt32, Int64: math.MinInt64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{Int: -1030108, Int32: math.MaxInt32, Int64: math.MaxInt64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{Int: math.MinInt64 + 12345}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

}

func TestUint(t *testing.T) {
	v := ValueChecking{Uint: 8, Uint8: math.MaxUint8}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{Uint: 130, Uint16: math.MaxUint16}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{Uint: 30130, Uint32: math.MaxUint32}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{Uint: 1030130, Uint64: math.MaxUint64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{Uint: math.MaxUint64 - 12345}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

}

func TestFloat(t *testing.T) {
	v := ValueChecking{Float32: 0, Float64: 0}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{Float32: -1, Float64: -1}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{Float32: math.SmallestNonzeroFloat32, Float64: math.SmallestNonzeroFloat64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{Float32: math.MaxFloat32, Float64: math.MaxFloat64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

}

func TestString(t *testing.T) {
	base := "abcdefghijklmnopqrstuvwxyz12345"
	v := ValueChecking{String: ""}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{String: strings.Repeat(base, 1)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{String: strings.Repeat(base, 8)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{String: strings.Repeat(base, (math.MaxUint16/len(base))-1)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{String: strings.Repeat(base, (math.MaxUint16/len(base))+1)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
}

func TestBool(t *testing.T) {
	v := ValueChecking{Bool: true}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{Bool: false}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

}

func TestByteRune(t *testing.T) {
	v := ValueChecking{Byte: 127, Rune: 'a'}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
}

func TestComplex(t *testing.T) {
	v := ValueChecking{Complex64: complex(1, 2), Complex128: complex(3, 4)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{
		Complex64:  complex(math.MaxFloat32, math.SmallestNonzeroFloat32),
		Complex128: complex(math.MaxFloat64, math.SmallestNonzeroFloat64),
	}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
}

func TestMap(t *testing.T) {
	v := ValueChecking{
		MapIntInt: map[string]int{"1": 2, "3": 4, "5": 6, "7": 8, "9": 10},
	}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{
		MapIntInt: make(map[string]int, 1000),
	}
	for i := 0; i < len(v.MapIntInt); i++ {
		v.MapIntInt[fmt.Sprint(i)] = i + 1
	}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

	v = ValueChecking{
		MapIntInt: make(map[string]int, math.MaxUint16+1),
	}
	for i := 0; i < len(v.MapIntInt); i++ {
		v.MapIntInt[fmt.Sprint(i)] = i + 1
	}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{
		MapIntInt: nil,
	}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
}

func TestPointerValue(t *testing.T) {
	{
		_v := 123
		v := ValueChecking{Pint: &_v}
		if err := checkValue(v); err != nil {
			t.Error(err)
		}
	}
	{
		_v := "this is pointer"
		__v := &_v
		v := ValueChecking{P2string: &__v}
		if err := checkValue(v); err != nil {
			t.Error(err)
		}
	}
	{
		_v := float32(1.234)
		__v := &_v
		___v := &__v
		____v := &___v
		v := ValueChecking{P3float32: ____v}
		if err := checkValue(v); err != nil {
			t.Error(err)
		}
	}

	check := func(v ValueChecking) error {
		var v1, v2 ValueChecking
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
		v := ValueChecking{P3float32: ____v}
		if err := check(v); err != nil {
			t.Error(err)
		}
		___v = nil
		v = ValueChecking{P3float32: ____v}
		if err := check(v); err != nil {
			t.Error(err)
		}
		____v = nil
		v = ValueChecking{P3float32: ____v}
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

		v := ValueChecking{IntPointers: _v}
		if err := checkValue(v); err != nil {
			t.Error(err)
		}
	}

	check2 := func(v ValueChecking) error {
		var v1, v2 ValueChecking
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

		v := ValueChecking{MapPointers: _v}
		if err := check2(v); err != nil {
			t.Error(err)
		}
	}
}

func TestTime(t *testing.T) {
	now := time.Now()
	add := now.Add(1 * time.Minute)
	v := TimeChecking{Time: time.Now(), TimePointer: &add}
	b, err := msgpack.Encode(v)
	if err != nil {
		t.Error(err)
	}
	var _v TimeChecking
	err = msgpack.Decode(b, &_v)
	if err != nil {
		t.Error(err)
	}
	if v.Time.UnixNano() != _v.Time.UnixNano() {
		t.Errorf("time different %v, %v", v.Time, _v.Time)
	}
	if v.TimePointer.UnixNano() != _v.TimePointer.UnixNano() {
		t.Errorf("time different %v, %v", v.Time, _v.Time)
	}

	vv := TimeChecking{}
	b, err = msgpack.EncodeAsArray(vv)
	if err != nil {
		t.Error(err)
	}
	var _vv TimeChecking
	err = msgpack.DecodeAsArray(b, &_vv)
	if err != nil {
		t.Error(err)
	}
	if vv.Time.UnixNano() != _vv.Time.UnixNano() {
		t.Errorf("time different %v, %v", vv.Time, _vv.Time)
	}
	if vv.TimePointer != nil || _vv.TimePointer != nil {
		t.Errorf("time different %v, %v", vv.Time, _vv.Time)
	}
}

func checkValue(v ValueChecking, eqs ...func() (bool, interface{}, interface{})) error {
	var v1, v2 ValueChecking
	return _checkValue(v, &v1, &v2, eqs...)
}

func TestSliceArray(t *testing.T) {

	f := func(l int) []int {
		slice := make([]int, l)
		for i := range slice {
			slice[i] = rand.Intn(math.MaxInt32)
		}
		return slice
	}

	check := func(v ValueChecking) error {
		var v1, v2 ValueChecking
		return _checkValue(v, &v1, &v2)
	}

	v := ValueChecking{Slice: f(15)}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{Slice: f(30015)}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{Slice: f(1030015)}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{Slice: nil}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{}
	for i := range v.Array1 {
		v.Array1[i] = float32(rand.Intn(0xff))
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{}
	for i := range v.Array2 {
		v.Array2[i] = "a"
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{}
	for i := range v.Array3 {
		v.Array3[i] = rand.Intn(0xff) > 0x7f
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{}
	for i := range v.Array4 {
		v.Array4[i] = rand.Intn(math.MaxInt32)
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{}
	for i := range v.Array5 {
		v.Array5[i] = rand.Intn(math.MaxInt32)
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = ValueChecking{}
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
	if v.Int == v1.Emb.Int || v.Int == v2.Embedded.Int {
		t.Error("value something wrong", v.Int, v1.Emb.Int, v2.Embedded.Int)
	}

	v = TestingStruct{}
	v.Embedded = Embedded{Int: rand.Int()}
	v1, v2, err = check(v)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(v.Embedded, v1.Embedded) || !reflect.DeepEqual(v.Embedded, v2.Embedded) {
		t.Error("value different", v.Embedded, v1.Embedded, v2.Embedded)
	}
	if v.Embedded.Int == v1.Int || v.Embedded.Int == v2.Emb.Int {
		t.Error("value something wrong", v.Embedded.Int, v1.Int, v2.Emb.Int)
	}

	v = TestingStruct{}
	v.Emb = embedded{Int: rand.Int()}
	v1, v2, err = check(v)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(v.Emb, v1.Emb) || !reflect.DeepEqual(v.Emb, v2.Emb) {
		t.Error("value different", v.Emb, v1.Emb, v2.Emb)
	}
	if v.Emb.Int == v1.Int || v.Emb.Int == v2.Embedded.Int {
		t.Error("value something wrong", v.Emb.Int, v1.Int, v2.Embedded.Int)
	}

	v = TestingStruct{}
	v.A = tst.A{Int: rand.Int(), B: tst2.B{Int: rand.Int()}}
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
	v.BB = tst2.DotImport{Int: rand.Int()}
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
	v.Time = tst2.Time{Int: rand.Int()}
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
}

func TestPointer(t *testing.T) {

	v := ValueChecking{Int: -1, Uint: 1}

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
	v1, v2 := new(ValueChecking), new(ValueChecking)
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
	_v3, _v4 := new(ValueChecking), new(ValueChecking)
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
	__v5, __v6 := new(ValueChecking), new(ValueChecking)
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
}

func TestPrivate(t *testing.T) {
	v := Private{}
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
