package msgpack

import (
	"testing"

	rawmsgpack "github.com/shamaton/msgpack/v3"
)

func preserveStructAsArray(t *testing.T) {
	t.Helper()
	structAsArray := StructAsArray()
	t.Cleanup(func() {
		SetStructAsArray(structAsArray)
	})
}

func TestInternalBufferEncodeAppendsFallback(t *testing.T) {
	prefix := []byte{0x01, 0x02}
	input := []int{3, 4}
	got, err := marshalAsArrayTo(input, prefix[:1])
	if err != nil {
		t.Fatal(err)
	}
	wantEncoded, err := rawmsgpack.MarshalAsArray(input)
	if err != nil {
		t.Fatal(err)
	}
	want := append([]byte{0x01}, wantEncoded...)
	if string(got) != string(want) {
		t.Fatalf("marshalAsArrayTo = %x, want %x", got, want)
	}
	if prefix[0] != 0x01 {
		t.Fatalf("prefix mutated: %x", prefix)
	}

	got, err = marshalAsMapTo(input, nil)
	if err != nil {
		t.Fatal(err)
	}
	want, err = rawmsgpack.MarshalAsMap(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("marshalAsMapTo nil buf = %x, want %x", got, want)
	}
}

func TestInternalBufferEncodeFallbackUsesOriginalBufferLength(t *testing.T) {
	mapPrefix := make([]byte, 2, 64)
	mapPrefix[0], mapPrefix[1] = 0x01, 0x02
	got, err := marshalAsMapTo(map[string]int{"a": 1}, mapPrefix)
	if err != nil {
		t.Fatal(err)
	}
	wantEncoded, err := rawmsgpack.MarshalAsMap(map[string]int{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	want := append([]byte{0x01, 0x02}, wantEncoded...)
	if string(got) != string(want) {
		t.Fatalf("marshalAsMapTo fallback = %x, want %x", got, want)
	}
	if string(got[:2]) != string([]byte{0x01, 0x02}) {
		t.Fatalf("marshalAsMapTo prefix = %x, want 0102", got[:2])
	}

	arrayPrefix := make([]byte, 2, 64)
	arrayPrefix[0], arrayPrefix[1] = 0x03, 0x04
	got, err = marshalAsArrayTo([]int{3, 4}, arrayPrefix)
	if err != nil {
		t.Fatal(err)
	}
	wantEncoded, err = rawmsgpack.MarshalAsArray([]int{3, 4})
	if err != nil {
		t.Fatal(err)
	}
	want = append([]byte{0x03, 0x04}, wantEncoded...)
	if string(got) != string(want) {
		t.Fatalf("marshalAsArrayTo fallback = %x, want %x", got, want)
	}
	if string(got[:2]) != string([]byte{0x03, 0x04}) {
		t.Fatalf("marshalAsArrayTo prefix = %x, want 0304", got[:2])
	}
}

func TestInternalBufferEncodeUsesStructAsArray(t *testing.T) {
	preserveStructAsArray(t)

	input := []int{7, 8}
	SetStructAsArray(false)
	got, err := marshalWithBuffer(input, []byte{0x01})
	if err != nil {
		t.Fatal(err)
	}
	wantEncoded, err := rawmsgpack.MarshalAsMap(input)
	if err != nil {
		t.Fatal(err)
	}
	want := append([]byte{0x01}, wantEncoded...)
	if string(got) != string(want) {
		t.Fatalf("marshalWithBuffer map = %x, want %x", got, want)
	}

	SetStructAsArray(true)
	got, err = marshalWithBuffer(input, []byte{0x02})
	if err != nil {
		t.Fatal(err)
	}
	wantEncoded, err = rawmsgpack.MarshalAsArray(input)
	if err != nil {
		t.Fatal(err)
	}
	want = append([]byte{0x02}, wantEncoded...)
	if string(got) != string(want) {
		t.Fatalf("marshalWithBuffer array = %x, want %x", got, want)
	}
}

func TestPublicMarshalUsesFallbackRuntime(t *testing.T) {
	preserveStructAsArray(t)

	input := map[string]int{"a": 1}
	got, err := MarshalAsMap(input)
	if err != nil {
		t.Fatal(err)
	}
	want, err := rawmsgpack.MarshalAsMap(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("MarshalAsMap = %x, want %x", got, want)
	}

	got, err = MarshalAsArray(input)
	if err != nil {
		t.Fatal(err)
	}
	want, err = rawmsgpack.MarshalAsArray(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("MarshalAsArray = %x, want %x", got, want)
	}

	SetStructAsArray(false)
	got, err = Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	want, err = rawmsgpack.MarshalAsMap(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("Marshal map = %x, want %x", got, want)
	}

	SetStructAsArray(true)
	got, err = Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	want, err = rawmsgpack.MarshalAsArray(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("Marshal array = %x, want %x", got, want)
	}
}
