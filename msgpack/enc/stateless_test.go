package enc

import (
	"bytes"
	"math"
	"strconv"
	"testing"
	"time"
)

func TestRequireAtExtendsBuffer(t *testing.T) {
	buf := make([]byte, 2, 8)
	buf[0], buf[1] = 0xaa, 0xbb
	base := &buf[:cap(buf)][0]

	got := RequireAt(buf, 5, 2)
	if len(got) != 7 {
		t.Fatalf("len = %d, want 7", len(got))
	}
	if &got[:cap(got)][0] != base {
		t.Fatal("RequireAt reallocated despite sufficient capacity")
	}
	if got[0] != 0xaa || got[1] != 0xbb {
		t.Fatalf("prefix = %x, want aabb", got[:2])
	}
}

func TestRequireAtAllocatesWhenCapacityIsInsufficient(t *testing.T) {
	buf := []byte{0xaa, 0xbb}
	got := RequireAt(buf, 2, 4)
	if len(got) != 6 {
		t.Fatalf("len = %d, want 6", len(got))
	}
	if got[0] != 0xaa || got[1] != 0xbb {
		t.Fatalf("prefix = %x, want aabb", got[:2])
	}
}

func TestStatelessWritersMatchLegacyEncoder(t *testing.T) {
	longString := stringsOfLen(1 << 16)
	instant := time.Unix(1<<35, 123).UTC()

	tests := []struct {
		name     string
		size     int
		writeOld func(*Encoder, int) int
		writeTo  func([]byte, int) int
	}{
		{
			name:     "int",
			size:     CalcInt(math.MinInt16),
			writeOld: func(e *Encoder, offset int) int { return e.WriteInt(math.MinInt16, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteIntTo(buf, math.MinInt16, offset) },
		},
		{
			name:     "int8",
			size:     CalcInt8(-5),
			writeOld: func(e *Encoder, offset int) int { return e.WriteInt8(-5, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteInt8To(buf, -5, offset) },
		},
		{
			name:     "int16",
			size:     CalcInt16(-129),
			writeOld: func(e *Encoder, offset int) int { return e.WriteInt16(-129, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteInt16To(buf, -129, offset) },
		},
		{
			name:     "int32",
			size:     CalcInt32(math.MinInt16 - 1),
			writeOld: func(e *Encoder, offset int) int { return e.WriteInt32(math.MinInt16-1, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteInt32To(buf, math.MinInt16-1, offset) },
		},
		{
			name:     "int64",
			size:     CalcInt64(math.MinInt32 - 1),
			writeOld: func(e *Encoder, offset int) int { return e.WriteInt64(math.MinInt32-1, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteInt64To(buf, math.MinInt32-1, offset) },
		},
		{
			name:     "uint",
			size:     CalcUint(uint(math.MaxUint16 + 1)),
			writeOld: func(e *Encoder, offset int) int { return e.WriteUint(uint(math.MaxUint16+1), offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteUintTo(buf, uint(math.MaxUint16+1), offset) },
		},
		{
			name:     "uint8",
			size:     CalcUint8(math.MaxInt8 + 1),
			writeOld: func(e *Encoder, offset int) int { return e.WriteUint8(math.MaxInt8+1, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteUint8To(buf, math.MaxInt8+1, offset) },
		},
		{
			name:     "uint16",
			size:     CalcUint16(math.MaxUint8 + 1),
			writeOld: func(e *Encoder, offset int) int { return e.WriteUint16(math.MaxUint8+1, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteUint16To(buf, math.MaxUint8+1, offset) },
		},
		{
			name:     "uint32",
			size:     CalcUint32(math.MaxUint16 + 1),
			writeOld: func(e *Encoder, offset int) int { return e.WriteUint32(math.MaxUint16+1, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteUint32To(buf, math.MaxUint16+1, offset) },
		},
		{
			name:     "uint64",
			size:     CalcUint64(math.MaxUint32 + 1),
			writeOld: func(e *Encoder, offset int) int { return e.WriteUint64(math.MaxUint32+1, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteUint64To(buf, math.MaxUint32+1, offset) },
		},
		{
			name:     "float32",
			size:     CalcFloat32(1.25),
			writeOld: func(e *Encoder, offset int) int { return e.WriteFloat32(1.25, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteFloat32To(buf, 1.25, offset) },
		},
		{
			name:     "float64",
			size:     CalcFloat64(1.25),
			writeOld: func(e *Encoder, offset int) int { return e.WriteFloat64(1.25, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteFloat64To(buf, 1.25, offset) },
		},
		{
			name:     "complex64",
			size:     CalcComplex64(1 + 2i),
			writeOld: func(e *Encoder, offset int) int { return e.WriteComplex64(1+2i, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteComplex64To(buf, 1+2i, offset) },
		},
		{
			name:     "complex128",
			size:     CalcComplex128(1 + 2i),
			writeOld: func(e *Encoder, offset int) int { return e.WriteComplex128(1+2i, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteComplex128To(buf, 1+2i, offset) },
		},
		{
			name:     "bool",
			size:     CalcBool(true),
			writeOld: func(e *Encoder, offset int) int { return e.WriteBool(true, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteBoolTo(buf, true, offset) },
		},
		{
			name:     "nil",
			size:     CalcNil(),
			writeOld: func(e *Encoder, offset int) int { return e.WriteNil(offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteNilTo(buf, offset) },
		},
		{
			name:     "byte",
			size:     CalcByte(0xcc),
			writeOld: func(e *Encoder, offset int) int { return e.WriteByte(0xcc, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteByteTo(buf, 0xcc, offset) },
		},
		{
			name:     "rune",
			size:     CalcRune('界'),
			writeOld: func(e *Encoder, offset int) int { return e.WriteRune('界', offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteRuneTo(buf, '界', offset) },
		},
		{
			name:     "string",
			size:     CalcString("hello"),
			writeOld: func(e *Encoder, offset int) int { return e.WriteString("hello", offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteStringTo(buf, "hello", offset) },
		},
		{
			name:     "string fix",
			size:     CalcStringFix(5),
			writeOld: func(e *Encoder, offset int) int { return e.WriteStringFix("hello", 5, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteStringFixTo(buf, "hello", 5, offset) },
		},
		{
			name:     "string8",
			size:     CalcString8(32),
			writeOld: func(e *Encoder, offset int) int { return e.WriteString8(stringsOfLen(32), 32, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteString8To(buf, stringsOfLen(32), 32, offset) },
		},
		{
			name:     "string16",
			size:     CalcString16(256),
			writeOld: func(e *Encoder, offset int) int { return e.WriteString16(stringsOfLen(256), 256, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteString16To(buf, stringsOfLen(256), 256, offset) },
		},
		{
			name:     "string32",
			size:     CalcString32(len(longString)),
			writeOld: func(e *Encoder, offset int) int { return e.WriteString32(longString, len(longString), offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteString32To(buf, longString, len(longString), offset) },
		},
		{
			name: "slice length",
			size: mustSize(func() (int, error) {
				return CalcSliceLength(16, false)
			}),
			writeOld: func(e *Encoder, offset int) int { return e.WriteSliceLength(16, offset, false) },
			writeTo:  func(buf []byte, offset int) int { return WriteSliceLengthTo(buf, 16, offset, false) },
		},
		{
			name: "byte slice length",
			size: mustSize(func() (int, error) {
				return CalcSliceLength(16, true)
			}),
			writeOld: func(e *Encoder, offset int) int { return e.WriteSliceLength(16, offset, true) },
			writeTo:  func(buf []byte, offset int) int { return WriteSliceLengthTo(buf, 16, offset, true) },
		},
		{
			name: "map length",
			size: mustSize(func() (int, error) {
				return CalcMapLength(16)
			}),
			writeOld: func(e *Encoder, offset int) int { return e.WriteMapLength(16, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteMapLengthTo(buf, 16, offset) },
		},
		{
			name:     "struct fix array",
			size:     CalcStructHeaderFix(2),
			writeOld: func(e *Encoder, offset int) int { return e.WriteStructHeaderFixAsArray(2, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteStructHeaderFixAsArrayTo(buf, 2, offset) },
		},
		{
			name:     "struct16 array",
			size:     CalcStructHeader16(16),
			writeOld: func(e *Encoder, offset int) int { return e.WriteStructHeader16AsArray(16, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteStructHeader16AsArrayTo(buf, 16, offset) },
		},
		{
			name:     "struct32 array",
			size:     CalcStructHeader32(1 << 16),
			writeOld: func(e *Encoder, offset int) int { return e.WriteStructHeader32AsArray(1<<16, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteStructHeader32AsArrayTo(buf, 1<<16, offset) },
		},
		{
			name:     "struct fix map",
			size:     CalcStructHeaderFix(2),
			writeOld: func(e *Encoder, offset int) int { return e.WriteStructHeaderFixAsMap(2, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteStructHeaderFixAsMapTo(buf, 2, offset) },
		},
		{
			name:     "struct16 map",
			size:     CalcStructHeader16(16),
			writeOld: func(e *Encoder, offset int) int { return e.WriteStructHeader16AsMap(16, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteStructHeader16AsMapTo(buf, 16, offset) },
		},
		{
			name:     "struct32 map",
			size:     CalcStructHeader32(1 << 16),
			writeOld: func(e *Encoder, offset int) int { return e.WriteStructHeader32AsMap(1<<16, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteStructHeader32AsMapTo(buf, 1<<16, offset) },
		},
		{
			name:     "time",
			size:     CalcTime(instant),
			writeOld: func(e *Encoder, offset int) int { return e.WriteTime(instant, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteTimeTo(buf, instant, offset) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefix := []byte{0xaa, 0xbb}
			old := NewEncoder()
			old.MakeBytes(len(prefix) + tt.size)
			copy(old.d, prefix)
			oldOffset := tt.writeOld(old, len(prefix))

			buf := append([]byte(nil), prefix...)
			buf = RequireAt(buf, len(prefix), tt.size)
			newOffset := tt.writeTo(buf, len(prefix))

			if oldOffset != len(prefix)+tt.size {
				t.Fatalf("legacy offset = %d, want %d", oldOffset, len(prefix)+tt.size)
			}
			if newOffset != oldOffset {
				t.Fatalf("offset = %d, want %d", newOffset, oldOffset)
			}
			if !bytes.Equal(buf[:newOffset], old.d[:oldOffset]) {
				t.Fatalf("bytes = %x, want %x", buf[:newOffset], old.d[:oldOffset])
			}
		})
	}
}

func TestStatelessWriterBoundariesMatchLegacyEncoder(t *testing.T) {
	string31 := stringsOfLen(31)
	string32 := stringsOfLen(32)
	string255 := stringsOfLen(math.MaxUint8)
	string256 := stringsOfLen(math.MaxUint8 + 1)
	string65535 := stringsOfLen(math.MaxUint16)
	string65536 := stringsOfLen(math.MaxUint16 + 1)
	timestamp32 := time.Unix(math.MaxUint32, 0).UTC()
	timestamp64 := time.Unix(1, 1).UTC()
	timestamp96 := time.Unix(1<<34, 0).UTC()

	tests := []struct {
		name     string
		size     int
		writeOld func(*Encoder, int) int
		writeTo  func([]byte, int) int
	}{
		{
			name:     "int positive fix max",
			size:     CalcInt(math.MaxInt8),
			writeOld: func(e *Encoder, offset int) int { return e.WriteInt(math.MaxInt8, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteIntTo(buf, math.MaxInt8, offset) },
		},
		{
			name:     "int uint8 min",
			size:     CalcInt(math.MaxInt8 + 1),
			writeOld: func(e *Encoder, offset int) int { return e.WriteInt(math.MaxInt8+1, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteIntTo(buf, math.MaxInt8+1, offset) },
		},
		{
			name:     "int negative fix min",
			size:     CalcInt(-32),
			writeOld: func(e *Encoder, offset int) int { return e.WriteInt(-32, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteIntTo(buf, -32, offset) },
		},
		{
			name:     "int int8 min",
			size:     CalcInt(-33),
			writeOld: func(e *Encoder, offset int) int { return e.WriteInt(-33, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteIntTo(buf, -33, offset) },
		},
		{
			name:     "uint positive fix max",
			size:     CalcUint(math.MaxInt8),
			writeOld: func(e *Encoder, offset int) int { return e.WriteUint(math.MaxInt8, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteUintTo(buf, math.MaxInt8, offset) },
		},
		{
			name:     "uint8 min",
			size:     CalcUint(math.MaxInt8 + 1),
			writeOld: func(e *Encoder, offset int) int { return e.WriteUint(math.MaxInt8+1, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteUintTo(buf, math.MaxInt8+1, offset) },
		},
		{
			name:     "uint16 min",
			size:     CalcUint(math.MaxUint8 + 1),
			writeOld: func(e *Encoder, offset int) int { return e.WriteUint(math.MaxUint8+1, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteUintTo(buf, math.MaxUint8+1, offset) },
		},
		{
			name:     "uint32 min",
			size:     CalcUint(math.MaxUint16 + 1),
			writeOld: func(e *Encoder, offset int) int { return e.WriteUint(math.MaxUint16+1, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteUintTo(buf, math.MaxUint16+1, offset) },
		},
		{
			name:     "uint64 min",
			size:     CalcUint64(math.MaxUint32 + 1),
			writeOld: func(e *Encoder, offset int) int { return e.WriteUint64(math.MaxUint32+1, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteUint64To(buf, math.MaxUint32+1, offset) },
		},
		{
			name:     "string fix max",
			size:     CalcString(string31),
			writeOld: func(e *Encoder, offset int) int { return e.WriteString(string31, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteStringTo(buf, string31, offset) },
		},
		{
			name:     "string8 min",
			size:     CalcString(string32),
			writeOld: func(e *Encoder, offset int) int { return e.WriteString(string32, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteStringTo(buf, string32, offset) },
		},
		{
			name:     "string8 max",
			size:     CalcString(string255),
			writeOld: func(e *Encoder, offset int) int { return e.WriteString(string255, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteStringTo(buf, string255, offset) },
		},
		{
			name:     "string16 min",
			size:     CalcString(string256),
			writeOld: func(e *Encoder, offset int) int { return e.WriteString(string256, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteStringTo(buf, string256, offset) },
		},
		{
			name:     "string16 max",
			size:     CalcString(string65535),
			writeOld: func(e *Encoder, offset int) int { return e.WriteString(string65535, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteStringTo(buf, string65535, offset) },
		},
		{
			name:     "string32 min",
			size:     CalcString(string65536),
			writeOld: func(e *Encoder, offset int) int { return e.WriteString(string65536, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteStringTo(buf, string65536, offset) },
		},
		{
			name: "slice fix max",
			size: mustSize(func() (int, error) {
				return CalcSliceLength(15, false)
			}),
			writeOld: func(e *Encoder, offset int) int { return e.WriteSliceLength(15, offset, false) },
			writeTo:  func(buf []byte, offset int) int { return WriteSliceLengthTo(buf, 15, offset, false) },
		},
		{
			name: "slice16 min",
			size: mustSize(func() (int, error) {
				return CalcSliceLength(16, false)
			}),
			writeOld: func(e *Encoder, offset int) int { return e.WriteSliceLength(16, offset, false) },
			writeTo:  func(buf []byte, offset int) int { return WriteSliceLengthTo(buf, 16, offset, false) },
		},
		{
			name: "slice16 max",
			size: mustSize(func() (int, error) {
				return CalcSliceLength(math.MaxUint16, false)
			}),
			writeOld: func(e *Encoder, offset int) int { return e.WriteSliceLength(math.MaxUint16, offset, false) },
			writeTo:  func(buf []byte, offset int) int { return WriteSliceLengthTo(buf, math.MaxUint16, offset, false) },
		},
		{
			name: "slice32 min",
			size: mustSize(func() (int, error) {
				return CalcSliceLength(math.MaxUint16+1, false)
			}),
			writeOld: func(e *Encoder, offset int) int { return e.WriteSliceLength(math.MaxUint16+1, offset, false) },
			writeTo:  func(buf []byte, offset int) int { return WriteSliceLengthTo(buf, math.MaxUint16+1, offset, false) },
		},
		{
			name: "map fix max",
			size: mustSize(func() (int, error) {
				return CalcMapLength(15)
			}),
			writeOld: func(e *Encoder, offset int) int { return e.WriteMapLength(15, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteMapLengthTo(buf, 15, offset) },
		},
		{
			name: "map16 min",
			size: mustSize(func() (int, error) {
				return CalcMapLength(16)
			}),
			writeOld: func(e *Encoder, offset int) int { return e.WriteMapLength(16, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteMapLengthTo(buf, 16, offset) },
		},
		{
			name: "map16 max",
			size: mustSize(func() (int, error) {
				return CalcMapLength(math.MaxUint16)
			}),
			writeOld: func(e *Encoder, offset int) int { return e.WriteMapLength(math.MaxUint16, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteMapLengthTo(buf, math.MaxUint16, offset) },
		},
		{
			name: "map32 min",
			size: mustSize(func() (int, error) {
				return CalcMapLength(math.MaxUint16 + 1)
			}),
			writeOld: func(e *Encoder, offset int) int { return e.WriteMapLength(math.MaxUint16+1, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteMapLengthTo(buf, math.MaxUint16+1, offset) },
		},
		{
			name:     "time timestamp32",
			size:     CalcTime(timestamp32),
			writeOld: func(e *Encoder, offset int) int { return e.WriteTime(timestamp32, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteTimeTo(buf, timestamp32, offset) },
		},
		{
			name:     "time timestamp64",
			size:     CalcTime(timestamp64),
			writeOld: func(e *Encoder, offset int) int { return e.WriteTime(timestamp64, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteTimeTo(buf, timestamp64, offset) },
		},
		{
			name:     "time timestamp96",
			size:     CalcTime(timestamp96),
			writeOld: func(e *Encoder, offset int) int { return e.WriteTime(timestamp96, offset) },
			writeTo:  func(buf []byte, offset int) int { return WriteTimeTo(buf, timestamp96, offset) },
		},
	}

	for _, tt := range tests {
		assertStatelessWriterMatchesLegacy(t, tt.name, tt.size, tt.writeOld, tt.writeTo)
	}
}

func TestStatelessLengthErrorsMatchLegacyEncoder(t *testing.T) {
	if strconv.IntSize < 64 {
		t.Skip("unsupported MaxUint32+1 length requires 64-bit int")
	}

	tooLong64 := uint64(math.MaxUint32) + 1
	tooLong := int(tooLong64)
	e := NewEncoder()
	if _, err := CalcSliceLength(tooLong, false); err == nil {
		t.Fatal("CalcSliceLength error = nil")
	}
	if _, err := e.CalcSliceLength(tooLong, false); err == nil {
		t.Fatal("legacy CalcSliceLength error = nil")
	}
	if _, err := CalcMapLength(tooLong); err == nil {
		t.Fatal("CalcMapLength error = nil")
	}
	if _, err := e.CalcMapLength(tooLong); err == nil {
		t.Fatal("legacy CalcMapLength error = nil")
	}
}

func stringsOfLen(n int) string {
	return string(bytes.Repeat([]byte{'a'}, n))
}

func mustSize(fn func() (int, error)) int {
	size, err := fn()
	if err != nil {
		panic(err)
	}
	return size
}

func assertStatelessWriterMatchesLegacy(
	t *testing.T,
	name string,
	size int,
	writeOld func(*Encoder, int) int,
	writeTo func([]byte, int) int,
) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		prefix := []byte{0xaa, 0xbb}
		old := NewEncoder()
		old.MakeBytes(len(prefix) + size)
		copy(old.d, prefix)
		oldOffset := writeOld(old, len(prefix))

		buf := append([]byte(nil), prefix...)
		buf = RequireAt(buf, len(prefix), size)
		newOffset := writeTo(buf, len(prefix))

		if oldOffset != len(prefix)+size {
			t.Fatalf("legacy offset = %d, want %d", oldOffset, len(prefix)+size)
		}
		if newOffset != oldOffset {
			t.Fatalf("offset = %d, want %d", newOffset, oldOffset)
		}
		if !bytes.Equal(buf[:newOffset], old.d[:oldOffset]) {
			t.Fatalf("bytes = %x, want %x", buf[:newOffset], old.d[:oldOffset])
		}
	})
}
