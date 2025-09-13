package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/shamaton/msgpack/v2/def"
	define2 "github.com/shamaton/msgpackgen/internal/fortest/define"
	"github.com/shamaton/msgpackgen/internal/fortest/define/define"
	"github.com/shamaton/msgpackgen/msgpack"
)

func TestSwitchDefaultBehaviour(t *testing.T) {
	msgpack.SetStructAsArray(false)

	v := inside{Int: 1}
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
	v := testingValue{Int: -8, Int8: math.MinInt8, Int16: math.MinInt16}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{Int: -108, Int8: math.MaxInt8, Int16: math.MaxInt16}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{Int: -30108, Int32: math.MinInt32, Int64: math.MinInt64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{Int: -1030108, Int32: math.MaxInt32, Int64: math.MaxInt64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{
		Int:                                      math.MinInt64 + 12345,
		Abcdefghijabcdefghijabcdefghijabcdefghij: rand.Int(),
		AbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghijAbcdefghij: rand.Int(),
	}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	if v.Function() != v.Int*2 {
		t.Errorf("value different %d, %d", v.Function(), v.Int*2)
	}

	{
		var r testingInt
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
	{
		var r testingInt
		f32 := testingFloat32{F: 2.345}
		b, err := msgpack.MarshalAsArray(f32)
		if err != nil {
			t.Error(err)
		}

		err = msgpack.UnmarshalAsArray(b, &r)
		if err != nil {
			t.Error(err)
		}

		if r.I != 2 {
			t.Error("different value", r.I)
		}
	}

	{
		var r testingInt
		f64 := testingFloat64{F: 6.789}
		b, err := msgpack.MarshalAsArray(f64)
		if err != nil {
			t.Error(err)
		}

		err = msgpack.UnmarshalAsArray(b, &r)
		if err != nil {
			t.Error(err)
		}

		if r.I != 6 {
			t.Error("different value", r.I)
		}
	}
}

func TestUint(t *testing.T) {
	v := testingValue{Uint: 8, Uint8: math.MaxUint8}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{Uint: 130, Uint16: math.MaxUint16}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{Uint: 30130, Uint32: math.MaxUint32}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{Uint: 1030130, Uint64: math.MaxUint64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{Uint: math.MaxUint64 - 12345}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

	{
		vv := inside{Int: -1}
		var r testingUint
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
		vv := inside{Int: math.MinInt8}
		var r testingUint
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
		vv := inside{Int: math.MinInt16}
		var r testingUint
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
		vv := inside{Int: math.MinInt32}
		var r testingUint
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
		vv := inside{Int: math.MinInt32 - 1}
		var r testingUint
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
		var r testingUint
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
	v := testingValue{Float32: 0, Float64: 0}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{Float32: -1, Float64: -1}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{Float32: math.SmallestNonzeroFloat32, Float64: math.SmallestNonzeroFloat64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{Float32: math.MaxFloat32, Float64: math.MaxFloat64}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

	{
		vv := inside{Int: 1}
		var r testingFloat32
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
		vv := inside{Int: -1}
		var r testingFloat32
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
		vv := inside{Int: 1}
		var r testingFloat64
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
		vv := inside{Int: -1}
		var r testingFloat64
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
		vv := testingFloat32{F: 1.23}
		var r testingFloat64
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
		var r testingFloat32
		err := msgpack.UnmarshalAsArray([]byte{0x91, def.Nil}, &r)
		if err != nil {
			t.Error(err)
		}
		if r.F != 0 {
			t.Errorf("not equal %v", r)
		}
	}
	{
		var r testingFloat64
		err := msgpack.UnmarshalAsArray([]byte{0x91, def.Nil}, &r)
		if err != nil {
			t.Error(err)
		}
		if r.F != 0 {
			t.Errorf("not equal %v", r)
		}
	}
	{
		var r testingFloat32
		err := msgpack.UnmarshalAsArray([]byte{0x91, def.True}, &r)
		if err == nil {
			t.Errorf("error must occur")
		}
		if !strings.Contains(err.Error(), "AsFloat32") {
			t.Error(err)
		}
	}
	{
		var r testingFloat64
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
	v := testingValue{String: ""}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{String: strings.Repeat(base, 1)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{String: strings.Repeat(base, 8)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{String: strings.Repeat(base, (math.MaxUint16/len(base))-1)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{String: strings.Repeat(base, (math.MaxUint16/len(base))+1)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

	{
		var r testingString
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
	v := testingValue{Bool: true}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{Bool: false}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

}

func TestByteRune(t *testing.T) {
	v := testingValue{Byte: 127, Rune: 'a'}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

	{
		var r testingBool
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
	v := testingValue{Complex64: complex(1, 2), Complex128: complex(3, 4)}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{
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

	b64, err := msgpack.MarshalAsArray(testingComplex64{})
	if err != nil {
		t.Error(err)
	}
	b128, err := msgpack.MarshalAsArray(testingComplex128{})
	if err != nil {
		t.Error(err)
	}

	msgpack.SetComplexTypeCode(-122)

	{
		var r testingComplex64
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
		var r testingComplex128
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
	v := testingValue{
		MapIntInt: map[string]int{"1": 2, "3": 4, "5": 6, "7": 8, "9": 10},
	}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{
		MapIntInt: map[string]int{},
	}
	for i := 0; i < 1000; i++ {
		v.MapIntInt[fmt.Sprint(i)] = i + 1
	}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

	v = testingValue{
		MapIntInt: map[string]int{},
	}
	for i := 0; i < math.MaxUint16+1; i++ {
		v.MapIntInt[fmt.Sprint(i)] = i + 1
	}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}
	v = testingValue{
		MapIntInt: nil,
	}
	if err := checkValue(v); err != nil {
		t.Error(err)
	}

	{
		var r testingMap
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
		vv := 123
		v := testingValue{Pint: &vv}
		if err := checkValue(v); err != nil {
			t.Error(err)
		}
	}
	{
		vv := "this is pointer"
		vvv := &vv
		v := testingValue{P2string: &vvv}
		if err := checkValue(v); err != nil {
			t.Error(err)
		}
	}
	{
		vv := float32(1.234)
		vvv := &vv
		vvvv := &vvv
		vvvvv := &vvvv
		v := testingValue{P3float32: vvvvv}
		if err := checkValue(v); err != nil {
			t.Error(err)
		}
	}

	check := func(v testingValue) error {
		var v1, v2 testingValue
		f1 := func() (bool, interface{}, interface{}) {
			return v1.P3float32 == nil, v1.P3float32, nil
		}
		f2 := func() (bool, interface{}, interface{}) {
			return v2.P3float32 == nil, v1.P3float32, nil
		}

		return _checkValue(v, &v1, &v2, f1, f2)
	}
	{

		vv := float32(1.234)
		vvv := &vv
		vvvv := &vvv
		vvvvv := &vvvv

		vvv = nil
		v := testingValue{P3float32: vvvvv}
		if err := check(v); err != nil {
			t.Error(err)
		}
		vvvv = nil
		v = testingValue{P3float32: vvvvv}
		if err := check(v); err != nil {
			t.Error(err)
		}
		vvvvv = nil
		v = testingValue{P3float32: vvvvv}
		if err := check(v); err != nil {
			t.Error(err)
		}
	}
	{
		vv := make([]*int, 100)
		for i := range vv {
			if i%2 == 0 {
				vvv := i
				vv[i] = &vvv
			} else {
				vv[i] = nil
			}
		}

		v := testingValue{IntPointers: vv}
		if err := checkValue(v); err != nil {
			t.Error(err)
		}
	}

	check2 := func(v testingValue) error {
		var v1, v2 testingValue
		f1 := func() (bool, interface{}, interface{}) {
			mp := map[uint]string{}
			for kk, vv := range v1.MapPointers {
				vvv := *vv
				mp[*kk] = *vvv
			}
			for kk, vv := range v.MapPointers {
				vvv := *vv
				if str, ok := mp[*kk]; !ok || str != *vvv {
					return false, str, *vvv
				}
			}
			return true, nil, nil
		}

		f2 := func() (bool, interface{}, interface{}) {
			mp := map[uint]string{}
			for kk, vv := range v2.MapPointers {
				vvv := *vv
				mp[*kk] = *vvv
			}
			for kk, vv := range v.MapPointers {
				vvv := *vv
				if str, ok := mp[*kk]; !ok || str != *vvv {
					return false, str, *vvv
				}
			}
			return true, nil, nil
		}

		return _checkValue(v, &v1, &v2, f1, f2)
	}
	{
		vvv := make(map[*uint]**string, 10)
		for i := 0; i < 10; i++ {
			k := uint(i)
			v := fmt.Sprint(i, i, i)
			vv := &v
			vvv[&k] = &vv
		}

		v := testingValue{MapPointers: vvv}
		if err := check2(v); err != nil {
			t.Error(err)
		}
	}
}

func TestTime(t *testing.T) {
	{
		v := testingTime{Time: time.Now()}
		b, err := msgpack.Marshal(v)
		if err != nil {
			t.Error(err)
		}
		var vv testingTime
		err = msgpack.Unmarshal(b, &vv)
		if err != nil {
			t.Error(err)
		}
		if v.Time.UnixNano() != vv.Time.UnixNano() {
			t.Errorf("time different %v, %v", v.Time, vv.Time)
		}
	}
	{
		vv := testingTime{}
		b, err := msgpack.MarshalAsArray(vv)
		if err != nil {
			t.Error(err)
		}
		var vvv testingTime
		err = msgpack.UnmarshalAsArray(b, &vvv)
		if err != nil {
			t.Error(err)
		}
		if vv.Time.UnixNano() != vvv.Time.UnixNano() {
			t.Errorf("time different %v, %v", vv.Time, vvv.Time)
		}
	}
	{
		v := define2.AA{Time: time.Now()}
		b, err := msgpack.Marshal(v)
		if err != nil {
			t.Error(err)
		}
		var vv define2.AA
		err = msgpack.Unmarshal(b, &vv)
		if err != nil {
			t.Error(err)
		}
		if v.UnixNano() != vv.UnixNano() {
			t.Errorf("time different %v, %v", v.Time, vv.Time)
		}
	}
	{
		now := time.Now().Unix()
		nowByte := make([]byte, 4)
		binary.BigEndian.PutUint32(nowByte, uint32(now))

		var r testingTime
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
	v := testingTimePointer{Time: &add}
	b, err := msgpack.Marshal(v)
	if err != nil {
		t.Error(err)
	}
	var r testingTimePointer
	err = msgpack.Unmarshal(b, &r)
	if err != nil {
		t.Error(err)
	}
	if v.Time.UnixNano() != r.Time.UnixNano() {
		t.Errorf("time different %v, %v", v.Time, r.Time)
	}

	vv := testingTimePointer{}
	b, err = msgpack.MarshalAsArray(vv)
	if err != nil {
		t.Error(err)
	}
	var rr testingTimePointer
	err = msgpack.UnmarshalAsArray(b, &rr)
	if err != nil {
		t.Error(err)
	}
	if vv.Time != nil || rr.Time != nil {
		t.Errorf("time different %v, %v", vv.Time, rr.Time)
	}
}

func checkValue(v testingValue, eqs ...func() (bool, interface{}, interface{})) error {
	var v1, v2 testingValue
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

	check := func(v testingValue) error {
		var v1, v2 testingValue
		return _checkValue(v, &v1, &v2)
	}

	v := testingValue{Slice: f(15)}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = testingValue{Slice: f(150)}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = testingValue{Slice: f(math.MaxUint16 + 1)}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = testingValue{Slice: nil}
	if err := check(v); err != nil {
		t.Error(err)
	}

	v = testingValue{}
	v.Bytes = make([]byte, 100)
	for i := range v.Bytes {
		v.Bytes[i] = byte(rand.Intn(255))
	}
	if err := check(v); err != nil {
		t.Error(err)
	}

	v = testingValue{}
	v.Bytes = make([]byte, math.MaxUint8+1)
	for i := range v.Bytes {
		v.Bytes[i] = byte(rand.Intn(255))
	}
	if err := check(v); err != nil {
		t.Error(err)
	}

	v = testingValue{}
	v.Bytes = make([]byte, math.MaxUint16+1)
	for i := range v.Bytes {
		v.Bytes[i] = byte(rand.Intn(255))
	}
	if err := check(v); err != nil {
		t.Error(err)
	}

	v = testingValue{}
	v.DoubleSlice = make([][]int16, 3)
	for i := range v.DoubleSlice {
		v.DoubleSlice[i] = make([]int16, i+1)
		for j := range v.DoubleSlice[i] {
			v.DoubleSlice[i][j] = int16(i*2 + j)
		}
	}
	if err := check(v); err != nil {
		t.Error(err)
	}

	v = testingValue{}
	for i := range v.DoubleArray {
		for j := range v.DoubleArray[i] {
			v.DoubleArray[i][j] = int16(i*4 + j)
		}
	}
	if err := check(v); err != nil {
		t.Error(err)
	}

	v = testingValue{}
	v.TripleBytes = make([][][]byte, 5)
	for i := range v.TripleBytes {
		v.TripleBytes[i] = make([][]byte, i+1)
		for j := range v.TripleBytes[i] {
			v.TripleBytes[i][j] = make([]byte, j+1)
			for k := range v.TripleBytes[i][j] {
				v.TripleBytes[i][j][k] = byte(i*2 + j*1 + k)
			}
		}
	}
	if err := check(v); err != nil {
		t.Error(err)
	}

	v = testingValue{}
	v.DoubleSlicePointerMap = make([][]**map[string]int, 2)
	for i := range v.DoubleSlicePointerMap {
		v.DoubleSlicePointerMap[i] = make([]**map[string]int, 4)
		for j := range v.DoubleSlicePointerMap[i] {
			m := map[string]int{fmt.Sprint(i*100 + j): i*50 + j}
			mp := &m
			mpp := &mp
			v.DoubleSlicePointerMap[i][j] = mpp
		}
	}
	if err := check(v); err != nil {
		t.Error(err)
	}

	v = testingValue{}
	v.MapDoubleSlicePointerInt = make(map[string][][]**int)
	for i := 0; i < 3; i++ {
		key := fmt.Sprint(i)
		v.MapDoubleSlicePointerInt[key] = make([][]**int, i+3)
		for j := range v.MapDoubleSlicePointerInt[key] {
			v.MapDoubleSlicePointerInt[key][j] = make([]**int, j+2)
			for k := range v.MapDoubleSlicePointerInt[key][j] {
				a := rand.Int()
				ap := &a
				app := &ap
				v.MapDoubleSlicePointerInt[key][j][k] = app
			}
		}
	}
	if err := check(v); err != nil {
		t.Error(err)
	}

	v = testingValue{}
	v.Bytes = make([]byte, math.MaxUint32+1)
	_, err := msgpack.MarshalAsArray(v)
	if err == nil || !strings.Contains(err.Error(), "not support this array length") {
		t.Error("something wrong", err)
	}
	runtime.GC()

	{
		var r testingSlice
		err = msgpack.UnmarshalAsArray([]byte{0x91, def.True}, &r)
		if err == nil {
			t.Errorf("error must occur")
		}
		if !strings.Contains(err.Error(), "SliceLength") {
			t.Error(err)
		}
	}
	//{ windows panic...
	//	vv := testingSlice{}
	//	vv.Slice = make([]int8, math.MaxUint32+1)
	//	_, err := msgpack.MarshalAsArray(v)
	//	if err == nil || !strings.Contains(err.Error(), "not support this array length") {
	//		t.Error("something wrong", err)
	//	}
	//}

	runtime.GC()
}

func TestArray(t *testing.T) {
	check := func(v testingArrays) error {
		var v1, v2 testingArrays
		return _checkValue(v, &v1, &v2)
	}

	var v testingArrays

	v = testingArrays{}
	for i := range v.Array1 {
		v.Array1[i] = float32(rand.Intn(0xff))
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = testingArrays{}
	for i := range v.Array2 {
		v.Array2[i] = "a"
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = testingArrays{}
	for i := range v.Array3 {
		v.Array3[i] = rand.Intn(0xff) > 0x7f
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = testingArrays{}
	for i := range v.Array4 {
		v.Array4[i] = rand.Intn(math.MaxInt32)
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = testingArrays{}
	for i := range v.Array5 {
		v.Array5[i] = rand.Intn(math.MaxInt32)
	}
	if err := check(v); err != nil {
		t.Error(err)
	}
	v = testingArrays{}
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
		if b, v1, v2 := eqs[0](); !b {
			return fmt.Errorf("not equal u1 %v, %v", v1, v2)
		}
		if b, v1, v2 := eqs[1](); !b {
			return fmt.Errorf("not equal u2 %v, %v", v1, v2)
		}
	}
	return nil
}

func TestStruct(t *testing.T) {
	check := func(v testingStruct) (testingStruct, testingStruct, error) {
		f := func() (bool, interface{}, interface{}) {
			return true, nil, nil
		}
		var r1, r2 testingStruct
		err := _checkValue(v, &r1, &r2, f, f)
		return r1, r2, err
	}

	v := testingStruct{}
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

	v = testingStruct{}
	v.Inside = inside{Int: rand.Int()}
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

	v = testingStruct{}
	v.Outside = outside{Int: rand.Int()}
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

	v = testingStruct{}
	v.A = define2.A{Int: rand.Int(), B: define.B{Int: rand.Int()}}
	v1, v2, err = check(v)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(v.A, v1.A) || !reflect.DeepEqual(v.A, v2.A) {
		t.Error("value different", v.A, v1.A, v2.A)
	}
	if v.A.B.Int != v1.B.Int || v.B.Int != v2.A.B.Int { //nolint:staticcheck // keeping "A" explicit for clarity
		t.Error("value something wrong", v.A.Int, v1.Int, v2.Int)
	}
	if v.A.Int == v1.Int || v.A.Int == v2.Int {
		t.Error("value something wrong", v.A.Int, v1.Int, v2.Int)
	}
	if v.B.Int == v1.Int || v.B.Int == v2.Int {
		t.Error("value something wrong", v.A.Int, v1.Int, v2.Int)
	}

	v = testingStruct{}
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

	v = testingStruct{}
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

	v = testingStruct{}
	v.R = &recursive{Int: rand.Int()}
	v.R.R = &recursive{Int: rand.Int()}
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
		var r inside
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
	v := testingTag{Tag: 1, Ignore: rand.Int(), Omit: rand.Int()}
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

	var v1, v2 testingTag
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

	v := testingValue{Int: -1, Uint: 1}

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
	v1, v2 := new(testingValue), new(testingValue)
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
	v3p, v4p := new(testingValue), new(testingValue)
	v3, v4 := &v3p, &v4p
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
	if v3v := *v3; !reflect.DeepEqual(v, *v3v) {
		t.Error("not equal v3", v, v3v)
	}
	if v4v := *v4; !reflect.DeepEqual(v, *v4v) {
		t.Error("not equal v4", v, v4v)
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
	v5pp, v6pp := new(testingValue), new(testingValue)
	v5p, v6p := &v5pp, &v6pp
	v5, v6 := &v5p, &v6p
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
	err := checkUndefined(notGenerated1{}, notGenerated2{}, &notGenerated1{}, &notGenerated2{})
	if err != nil {
		t.Error(err)
	}
	err = checkUndefined(notGenerated3{}, notGenerated4{}, &notGenerated3{}, &notGenerated4{})
	if err != nil {
		t.Error(err)
	}
	err = checkUndefined(notGenerated5{}, notGenerated6{}, &notGenerated5{}, &notGenerated6{})
	if err != nil {
		t.Error(err)
	}
	err = checkUndefined(notGenerated7{}, notGenerated8{}, &notGenerated7{}, &notGenerated8{})
	if err != nil {
		t.Error(err)
	}
	err = checkUndefined(notGenerated10{}, notGenerated10{}, &notGenerated10{}, &notGenerated10{})
	if err != nil {
		t.Error(err)
	}
}

func TestPrivate(t *testing.T) {
	v := private{}
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
