package msgpack

import (
	"errors"
	"strings"
	"testing"

	rawmsgpack "github.com/shamaton/msgpack/v3"
)

func preserveResolvers(t *testing.T) {
	t.Helper()
	encMap := encAsMapResolver
	encArray := encAsArrayResolver
	encMapTo := encAsMapToResolver
	encArrayTo := encAsArrayToResolver
	decMap := decAsMapResolver
	decArray := decAsArrayResolver
	structAsArray := StructAsArray()
	t.Cleanup(func() {
		encAsMapResolver = encMap
		encAsArrayResolver = encArray
		encAsMapToResolver = encMapTo
		encAsArrayToResolver = encArrayTo
		decAsMapResolver = decMap
		decAsArrayResolver = decArray
		SetStructAsArray(structAsArray)
	})
}

func TestInternalBufferEncodeAppendsFallback(t *testing.T) {
	preserveResolvers(t)
	SetResolver(noOpEncResolver, noOpEncResolver, noOpDecResolver, noOpDecResolver)

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

func TestInternalBufferEncodeResolverStates(t *testing.T) {
	preserveResolvers(t)
	oldErr := errors.New("old resolver must not be called")
	toErr := errors.New("to resolver error")
	SetResolver(
		func(any) ([]byte, error) { return nil, oldErr },
		noOpEncResolver,
		noOpDecResolver,
		noOpDecResolver,
	)
	SetToResolver(
		func(any, []byte) ([]byte, bool, error) {
			return nil, false, toErr
		},
		noOpEncToResolver,
	)
	if _, err := marshalAsMapTo(1, []byte{0x01}); !errors.Is(err, toErr) {
		t.Fatalf("marshalAsMapTo error = %v, want %v", err, toErr)
	}
	SetToResolver(
		func(any, []byte) ([]byte, bool, error) {
			return []byte{0xbb}, true, toErr
		},
		func(any, []byte) ([]byte, bool, error) {
			return []byte{0xcc}, true, toErr
		},
	)
	if _, err := marshalAsMapTo(1, []byte{0x01}); !errors.Is(err, toErr) {
		t.Fatalf("handled marshalAsMapTo error = %v, want %v", err, toErr)
	}
	if _, err := marshalAsArrayTo(1, []byte{0x01}); !errors.Is(err, toErr) {
		t.Fatalf("handled marshalAsArrayTo error = %v, want %v", err, toErr)
	}

	SetResolver(
		func(any) ([]byte, error) { return []byte{0xaa}, nil },
		noOpEncResolver,
		noOpDecResolver,
		noOpDecResolver,
	)
	SetToResolver(
		func(any, []byte) ([]byte, bool, error) {
			return []byte{0xbb}, true, nil
		},
		noOpEncToResolver,
	)
	got, err := marshalAsMapTo(1, []byte{0x01})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0xbb}) {
		t.Fatalf("handled marshalAsMapTo = %x, want bb", got)
	}

	SetToResolver(
		func(any, []byte) ([]byte, bool, error) {
			return []byte{0xcc}, false, nil
		},
		noOpEncToResolver,
	)
	got, err = marshalAsMapTo(1, []byte{0x01})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x01, 0xaa}) {
		t.Fatalf("unhandled marshalAsMapTo = %x, want 01aa", got)
	}
}

func TestInternalBufferEncodeFallbackIgnoresUnhandledResolverBuffer(t *testing.T) {
	preserveResolvers(t)
	SetResolver(noOpEncResolver, noOpEncResolver, noOpDecResolver, noOpDecResolver)
	SetToResolver(
		func(any, []byte) ([]byte, bool, error) {
			return []byte{0xee, 0xee}, false, nil
		},
		func(any, []byte) ([]byte, bool, error) {
			return []byte{0xdd, 0xdd}, false, nil
		},
	)

	mapPrefix := []byte{0x01}
	got, err := marshalAsMapTo(map[string]int{"a": 1}, mapPrefix)
	if err != nil {
		t.Fatal(err)
	}
	wantEncoded, err := rawmsgpack.MarshalAsMap(map[string]int{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	want := append([]byte{0x01}, wantEncoded...)
	if string(got) != string(want) {
		t.Fatalf("marshalAsMapTo fallback = %x, want %x", got, want)
	}

	arrayPrefix := []byte{0x02}
	got, err = marshalAsArrayTo([]int{3, 4}, arrayPrefix)
	if err != nil {
		t.Fatal(err)
	}
	wantEncoded, err = rawmsgpack.MarshalAsArray([]int{3, 4})
	if err != nil {
		t.Fatal(err)
	}
	want = append([]byte{0x02}, wantEncoded...)
	if string(got) != string(want) {
		t.Fatalf("marshalAsArrayTo fallback = %x, want %x", got, want)
	}
}

func TestInternalBufferEncodeFallbackUsesOriginalBufferLength(t *testing.T) {
	preserveResolvers(t)
	SetResolver(noOpEncResolver, noOpEncResolver, noOpDecResolver, noOpDecResolver)
	SetToResolver(
		func(_ any, buf []byte) ([]byte, bool, error) {
			b := buf[:cap(buf)]
			for i := len(buf); i < len(b); i++ {
				b[i] = 0xee
			}
			return b, false, nil
		},
		func(_ any, buf []byte) ([]byte, bool, error) {
			b := buf[:cap(buf)]
			for i := len(buf); i < len(b); i++ {
				b[i] = 0xdd
			}
			return b, false, nil
		},
	)

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

func TestLegacyResolverWorksWithInternalBufferEncode(t *testing.T) {
	preserveResolvers(t)
	SetResolver(
		func(any) ([]byte, error) { return []byte{0x81}, nil },
		func(any) ([]byte, error) { return []byte{0x91}, nil },
		noOpDecResolver,
		noOpDecResolver,
	)

	SetStructAsArray(false)
	got, err := Marshal(1)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x81}) {
		t.Fatalf("Marshal legacy map resolver = %x, want 81", got)
	}
	got, err = marshalWithBuffer(1, []byte{0x01})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x01, 0x81}) {
		t.Fatalf("marshalWithBuffer legacy map resolver = %x, want 0181", got)
	}

	SetStructAsArray(true)
	got, err = Marshal(1)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x91}) {
		t.Fatalf("Marshal legacy array resolver = %x, want 91", got)
	}
	got, err = marshalWithBuffer(1, []byte{0x02})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x02, 0x91}) {
		t.Fatalf("marshalWithBuffer legacy array resolver = %x, want 0291", got)
	}
}

func TestInternalBufferEncodeStrictLegacyResolverDoesNotFallback(t *testing.T) {
	preserveResolvers(t)
	strictErr := errors.New("use strict option : undefined type")
	SetResolver(
		func(any) ([]byte, error) { return nil, strictErr },
		func(any) ([]byte, error) { return nil, strictErr },
		noOpDecResolver,
		noOpDecResolver,
	)

	for _, tt := range []struct {
		name string
		fn   func() ([]byte, error)
	}{
		{
			name: "MarshalAsMap",
			fn:   func() ([]byte, error) { return MarshalAsMap(map[string]int{"a": 1}) },
		},
		{
			name: "MarshalAsArray",
			fn:   func() ([]byte, error) { return MarshalAsArray([]int{1}) },
		},
		{
			name: "marshalAsMapTo",
			fn:   func() ([]byte, error) { return marshalAsMapTo(map[string]int{"a": 1}, []byte{0x01}) },
		},
		{
			name: "marshalAsArrayTo",
			fn:   func() ([]byte, error) { return marshalAsArrayTo([]int{1}, []byte{0x02}) },
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if b, err := tt.fn(); err == nil {
				t.Fatalf("error = nil, bytes = %x", b)
			} else if !strings.Contains(err.Error(), "use strict option") {
				t.Fatalf("error = %v, want strict option", err)
			}
		})
	}
}

func TestSetResolverResetsToResolver(t *testing.T) {
	preserveResolvers(t)
	SetToResolver(
		func(any, []byte) ([]byte, bool, error) { return []byte{0xbb}, true, nil },
		func(any, []byte) ([]byte, bool, error) { return []byte{0xcc}, true, nil },
	)
	SetResolver(
		func(any) ([]byte, error) { return []byte{0x11}, nil },
		func(any) ([]byte, error) { return []byte{0x22}, nil },
		noOpDecResolver,
		noOpDecResolver,
	)

	got, err := marshalAsMapTo(1, []byte{0x01})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x01, 0x11}) {
		t.Fatalf("marshalAsMapTo after SetResolver = %x, want 0111", got)
	}
	got, err = marshalAsArrayTo(1, []byte{0x02})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x02, 0x22}) {
		t.Fatalf("marshalAsArrayTo after SetResolver = %x, want 0222", got)
	}
}

func TestInternalBufferEncodeUsesStructAsArray(t *testing.T) {
	preserveResolvers(t)
	SetResolver(noOpEncResolver, noOpEncResolver, noOpDecResolver, noOpDecResolver)
	SetToResolver(
		func(any, []byte) ([]byte, bool, error) { return []byte{0x81}, true, nil },
		func(any, []byte) ([]byte, bool, error) { return []byte{0x91}, true, nil },
	)

	SetStructAsArray(false)
	got, err := marshalWithBuffer(1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x81}) {
		t.Fatalf("marshalWithBuffer map = %x, want 81", got)
	}

	SetStructAsArray(true)
	got, err = marshalWithBuffer(1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x91}) {
		t.Fatalf("marshalWithBuffer array = %x, want 91", got)
	}
}

func TestPublicMarshalUsesInternalToResolver(t *testing.T) {
	preserveResolvers(t)
	SetResolver(
		func(any) ([]byte, error) { return []byte{0x11}, nil },
		func(any) ([]byte, error) { return []byte{0x22}, nil },
		noOpDecResolver,
		noOpDecResolver,
	)
	SetToResolver(
		func(_ any, buf []byte) ([]byte, bool, error) {
			return append(buf, 0x81), true, nil
		},
		func(_ any, buf []byte) ([]byte, bool, error) {
			return append(buf, 0x91), true, nil
		},
	)

	got, err := MarshalAsMap(1)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x81}) {
		t.Fatalf("MarshalAsMap = %x, want 81", got)
	}

	got, err = MarshalAsArray(1)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x91}) {
		t.Fatalf("MarshalAsArray = %x, want 91", got)
	}

	SetStructAsArray(false)
	got, err = Marshal(1)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x81}) {
		t.Fatalf("Marshal map = %x, want 81", got)
	}

	SetStructAsArray(true)
	got, err = Marshal(1)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x91}) {
		t.Fatalf("Marshal array = %x, want 91", got)
	}
}

func TestSetToResolverAcceptsNil(t *testing.T) {
	preserveResolvers(t)
	SetResolver(
		func(any) ([]byte, error) { return []byte{0x11}, nil },
		func(any) ([]byte, error) { return []byte{0x22}, nil },
		noOpDecResolver,
		noOpDecResolver,
	)
	SetToResolver(nil, nil)

	got, err := marshalAsMapTo(1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x11}) {
		t.Fatalf("marshalAsMapTo with nil To resolver = %x, want 11", got)
	}
	got, err = marshalAsArrayTo(1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x22}) {
		t.Fatalf("marshalAsArrayTo with nil To resolver = %x, want 22", got)
	}
}

func TestResolverRegistrationOrderCompatibility(t *testing.T) {
	preserveResolvers(t)

	registerOld := func(mapByte, arrayByte byte) {
		SetResolver(
			func(any) ([]byte, error) { return []byte{mapByte}, nil },
			func(any) ([]byte, error) { return []byte{arrayByte}, nil },
			noOpDecResolver,
			noOpDecResolver,
		)
	}
	registerNew := func(mapByte, arrayByte byte) {
		SetResolver(
			func(any) ([]byte, error) { return []byte{mapByte + 0x10}, nil },
			func(any) ([]byte, error) { return []byte{arrayByte + 0x10}, nil },
			noOpDecResolver,
			noOpDecResolver,
		)
		SetToResolver(
			func(_ any, buf []byte) ([]byte, bool, error) {
				return append(buf, mapByte), true, nil
			},
			func(_ any, buf []byte) ([]byte, bool, error) {
				return append(buf, arrayByte), true, nil
			},
		)
	}

	registerOld(0x11, 0x12)
	registerNew(0x21, 0x22)
	if got, err := marshalAsMapTo(1, []byte{0x01}); err != nil {
		t.Fatal(err)
	} else if string(got) != string([]byte{0x01, 0x21}) {
		t.Fatalf("old then new marshalAsMapTo = %x, want 0121", got)
	}
	if got, err := marshalAsArrayTo(1, []byte{0x02}); err != nil {
		t.Fatal(err)
	} else if string(got) != string([]byte{0x02, 0x22}) {
		t.Fatalf("old then new marshalAsArrayTo = %x, want 0222", got)
	}

	registerNew(0x31, 0x32)
	registerOld(0x41, 0x42)
	if got, err := marshalAsMapTo(1, []byte{0x03}); err != nil {
		t.Fatal(err)
	} else if string(got) != string([]byte{0x03, 0x41}) {
		t.Fatalf("new then old marshalAsMapTo = %x, want 0341", got)
	}
	if got, err := marshalAsArrayTo(1, []byte{0x04}); err != nil {
		t.Fatal(err)
	} else if string(got) != string([]byte{0x04, 0x42}) {
		t.Fatalf("new then old marshalAsArrayTo = %x, want 0442", got)
	}

	registerNew(0x51, 0x52)
	registerNew(0x61, 0x62)
	if got, err := marshalAsMapTo(1, nil); err != nil {
		t.Fatal(err)
	} else if string(got) != string([]byte{0x61}) {
		t.Fatalf("new then new marshalAsMapTo = %x, want 61", got)
	}
	if got, err := marshalAsArrayTo(1, nil); err != nil {
		t.Fatal(err)
	} else if string(got) != string([]byte{0x62}) {
		t.Fatalf("new then new marshalAsArrayTo = %x, want 62", got)
	}

	SetResolver(noOpEncResolver, noOpEncResolver, noOpDecResolver, noOpDecResolver)
	SetToResolver(nil, nil)
	input := []int{7, 8}
	got, err := marshalAsArrayTo(input, []byte{0x05})
	if err != nil {
		t.Fatal(err)
	}
	wantEncoded, err := rawmsgpack.MarshalAsArray(input)
	if err != nil {
		t.Fatal(err)
	}
	want := append([]byte{0x05}, wantEncoded...)
	if string(got) != string(want) {
		t.Fatalf("nil/default resolvers fallback = %x, want %x", got, want)
	}
}
