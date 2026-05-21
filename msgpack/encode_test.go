package msgpack

import (
	"errors"
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

func TestMarshalToAppendsFallback(t *testing.T) {
	preserveResolvers(t)
	SetResolver(noOpEncResolver, noOpEncResolver, noOpDecResolver, noOpDecResolver)

	prefix := []byte{0x01, 0x02}
	input := []int{3, 4}
	got, err := MarshalAsArrayTo(input, prefix[:1])
	if err != nil {
		t.Fatal(err)
	}
	wantEncoded, err := rawmsgpack.MarshalAsArray(input)
	if err != nil {
		t.Fatal(err)
	}
	want := append([]byte{0x01}, wantEncoded...)
	if string(got) != string(want) {
		t.Fatalf("MarshalAsArrayTo = %x, want %x", got, want)
	}
	if prefix[0] != 0x01 {
		t.Fatalf("prefix mutated: %x", prefix)
	}

	got, err = MarshalAsMapTo(input, nil)
	if err != nil {
		t.Fatal(err)
	}
	want, err = rawmsgpack.MarshalAsMap(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("MarshalAsMapTo nil buf = %x, want %x", got, want)
	}
}

func TestMarshalToResolverStates(t *testing.T) {
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
	if _, err := MarshalAsMapTo(1, []byte{0x01}); !errors.Is(err, toErr) {
		t.Fatalf("MarshalAsMapTo error = %v, want %v", err, toErr)
	}
	SetToResolver(
		func(any, []byte) ([]byte, bool, error) {
			return []byte{0xbb}, true, toErr
		},
		func(any, []byte) ([]byte, bool, error) {
			return []byte{0xcc}, true, toErr
		},
	)
	if _, err := MarshalAsMapTo(1, []byte{0x01}); !errors.Is(err, toErr) {
		t.Fatalf("handled MarshalAsMapTo error = %v, want %v", err, toErr)
	}
	if _, err := MarshalAsArrayTo(1, []byte{0x01}); !errors.Is(err, toErr) {
		t.Fatalf("handled MarshalAsArrayTo error = %v, want %v", err, toErr)
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
	got, err := MarshalAsMapTo(1, []byte{0x01})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0xbb}) {
		t.Fatalf("handled MarshalAsMapTo = %x, want bb", got)
	}

	SetToResolver(
		func(any, []byte) ([]byte, bool, error) {
			return []byte{0xcc}, false, nil
		},
		noOpEncToResolver,
	)
	got, err = MarshalAsMapTo(1, []byte{0x01})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x01, 0xaa}) {
		t.Fatalf("unhandled MarshalAsMapTo = %x, want 01aa", got)
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

	got, err := MarshalAsMapTo(1, []byte{0x01})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x01, 0x11}) {
		t.Fatalf("MarshalAsMapTo after SetResolver = %x, want 0111", got)
	}
	got, err = MarshalAsArrayTo(1, []byte{0x02})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x02, 0x22}) {
		t.Fatalf("MarshalAsArrayTo after SetResolver = %x, want 0222", got)
	}
}

func TestMarshalToUsesStructAsArray(t *testing.T) {
	preserveResolvers(t)
	SetResolver(noOpEncResolver, noOpEncResolver, noOpDecResolver, noOpDecResolver)
	SetToResolver(
		func(any, []byte) ([]byte, bool, error) { return []byte{0x81}, true, nil },
		func(any, []byte) ([]byte, bool, error) { return []byte{0x91}, true, nil },
	)

	SetStructAsArray(false)
	got, err := MarshalTo(1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x81}) {
		t.Fatalf("MarshalTo map = %x, want 81", got)
	}

	SetStructAsArray(true)
	got, err = MarshalTo(1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x91}) {
		t.Fatalf("MarshalTo array = %x, want 91", got)
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

	got, err := MarshalAsMapTo(1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x11}) {
		t.Fatalf("MarshalAsMapTo with nil To resolver = %x, want 11", got)
	}
	got, err = MarshalAsArrayTo(1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string([]byte{0x22}) {
		t.Fatalf("MarshalAsArrayTo with nil To resolver = %x, want 22", got)
	}
}
