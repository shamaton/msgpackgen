package enc

import (
	"bytes"
	"encoding/binary"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/shamaton/msgpack/v3/def"
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

func TestWriters(t *testing.T) {
	longString := stringsOfLen(1 << 16)
	instant := time.Unix(1<<35, 123).UTC()

	tests := []writerCase{
		{name: "int", size: CalcInt(math.MinInt16), write: func(buf []byte, offset int) int { return WriteInt(buf, math.MinInt16, offset) }, want: bytesOf(def.Int16, 0x80, 0x00)},
		{name: "int8", size: CalcInt8(-5), write: func(buf []byte, offset int) int { return WriteInt8(buf, -5, offset) }, want: bytesOf(0xfb)},
		{name: "int16", size: CalcInt16(-129), write: func(buf []byte, offset int) int { return WriteInt16(buf, -129, offset) }, want: bytesOf(def.Int16, 0xff, 0x7f)},
		{name: "int32", size: CalcInt32(math.MinInt16 - 1), write: func(buf []byte, offset int) int { return WriteInt32(buf, math.MinInt16-1, offset) }, want: append(bytesOf(def.Int32), be32int(math.MinInt16-1)...)},
		{name: "int64", size: CalcInt64(math.MinInt32 - 1), write: func(buf []byte, offset int) int { return WriteInt64(buf, math.MinInt32-1, offset) }, want: append(bytesOf(def.Int64), be64int(math.MinInt32-1)...)},
		{name: "uint", size: CalcUint(uint(math.MaxUint16 + 1)), write: func(buf []byte, offset int) int { return WriteUint(buf, uint(math.MaxUint16+1), offset) }, want: append(bytesOf(def.Uint32), be32uint(math.MaxUint16+1)...)},
		{name: "uint8", size: CalcUint8(math.MaxInt8 + 1), write: func(buf []byte, offset int) int { return WriteUint8(buf, math.MaxInt8+1, offset) }, want: bytesOf(def.Uint8, 0x80)},
		{name: "uint16", size: CalcUint16(math.MaxUint8 + 1), write: func(buf []byte, offset int) int { return WriteUint16(buf, math.MaxUint8+1, offset) }, want: bytesOf(def.Uint16, 0x01, 0x00)},
		{name: "uint32", size: CalcUint32(math.MaxUint16 + 1), write: func(buf []byte, offset int) int { return WriteUint32(buf, math.MaxUint16+1, offset) }, want: append(bytesOf(def.Uint32), be32uint(math.MaxUint16+1)...)},
		{name: "uint64", size: CalcUint64(math.MaxUint32 + 1), write: func(buf []byte, offset int) int { return WriteUint64(buf, math.MaxUint32+1, offset) }, want: append(bytesOf(def.Uint64), be64uint(math.MaxUint32+1)...)},
		{name: "float32", size: CalcFloat32(1.25), write: func(buf []byte, offset int) int { return WriteFloat32(buf, 1.25, offset) }, want: append(bytesOf(def.Float32), be32uint(uint64(math.Float32bits(1.25)))...)},
		{name: "float64", size: CalcFloat64(1.25), write: func(buf []byte, offset int) int { return WriteFloat64(buf, 1.25, offset) }, want: append(bytesOf(def.Float64), be64uint(math.Float64bits(1.25))...)},
		{name: "complex64", size: CalcComplex64(1 + 2i), write: func(buf []byte, offset int) int { return WriteComplex64(buf, 1+2i, offset) }, want: concat(bytesOf(def.Fixext8, int(def.ComplexTypeCode())), be32uint(uint64(math.Float32bits(1))), be32uint(uint64(math.Float32bits(2))))},
		{name: "complex128", size: CalcComplex128(1 + 2i), write: func(buf []byte, offset int) int { return WriteComplex128(buf, 1+2i, offset) }, want: concat(bytesOf(def.Fixext16, int(def.ComplexTypeCode())), be64uint(math.Float64bits(1)), be64uint(math.Float64bits(2)))},
		{name: "bool", size: CalcBool(true), write: func(buf []byte, offset int) int { return WriteBool(buf, true, offset) }, want: bytesOf(def.True)},
		{name: "nil", size: CalcNil(), write: func(buf []byte, offset int) int { return WriteNil(buf, offset) }, want: bytesOf(def.Nil)},
		{name: "byte", size: CalcByte(0xcc), write: func(buf []byte, offset int) int { return WriteByte(buf, 0xcc, offset) }, want: bytesOf(0xcc)},
		{name: "rune", size: CalcRune('界'), write: func(buf []byte, offset int) int { return WriteRune(buf, '界', offset) }, want: bytesOf(def.Uint16, 0x75, 0x4c)},
		{name: "string", size: CalcString("hello"), write: func(buf []byte, offset int) int { return WriteString(buf, "hello", offset) }, want: append(bytesOf(def.FixStr+5), []byte("hello")...)},
		{name: "string fix", size: CalcStringFix(5), write: func(buf []byte, offset int) int { return WriteStringFix(buf, "hello", 5, offset) }, want: append(bytesOf(def.FixStr+5), []byte("hello")...)},
		{name: "string8", size: CalcString8(32), write: func(buf []byte, offset int) int { return WriteString8(buf, stringsOfLen(32), 32, offset) }, want: concat(bytesOf(def.Str8, 32), bytes.Repeat([]byte{'a'}, 32))},
		{name: "string16", size: CalcString16(256), write: func(buf []byte, offset int) int { return WriteString16(buf, stringsOfLen(256), 256, offset) }, want: concat(bytesOf(def.Str16, 0x01, 0x00), bytes.Repeat([]byte{'a'}, 256))},
		{name: "string32", size: CalcString32(len(longString)), write: func(buf []byte, offset int) int { return WriteString32(buf, longString, len(longString), offset) }, want: concat(bytesOf(def.Str32), be32uint(uint64(len(longString))), []byte(longString))},
		{name: "slice length", size: mustSize(func() (int, error) { return CalcSliceLength(16, false) }), write: func(buf []byte, offset int) int { return WriteSliceLength(buf, 16, offset, false) }, want: bytesOf(def.Array16, 0x00, 0x10)},
		{name: "byte slice length", size: mustSize(func() (int, error) { return CalcSliceLength(16, true) }), write: func(buf []byte, offset int) int { return WriteSliceLength(buf, 16, offset, true) }, want: bytesOf(def.Bin8, 0x10)},
		{name: "map length", size: mustSize(func() (int, error) { return CalcMapLength(16) }), write: func(buf []byte, offset int) int { return WriteMapLength(buf, 16, offset) }, want: bytesOf(def.Map16, 0x00, 0x10)},
		{name: "struct fix array", size: CalcStructHeaderFix(2), write: func(buf []byte, offset int) int { return WriteStructHeaderFixAsArray(buf, 2, offset) }, want: bytesOf(def.FixArray + 2)},
		{name: "struct16 array", size: CalcStructHeader16(16), write: func(buf []byte, offset int) int { return WriteStructHeader16AsArray(buf, 16, offset) }, want: bytesOf(def.Array16, 0x00, 0x10)},
		{name: "struct32 array", size: CalcStructHeader32(1 << 16), write: func(buf []byte, offset int) int { return WriteStructHeader32AsArray(buf, 1<<16, offset) }, want: append(bytesOf(def.Array32), be32uint(1<<16)...)},
		{name: "struct fix map", size: CalcStructHeaderFix(2), write: func(buf []byte, offset int) int { return WriteStructHeaderFixAsMap(buf, 2, offset) }, want: bytesOf(def.FixMap + 2)},
		{name: "struct16 map", size: CalcStructHeader16(16), write: func(buf []byte, offset int) int { return WriteStructHeader16AsMap(buf, 16, offset) }, want: bytesOf(def.Map16, 0x00, 0x10)},
		{name: "struct32 map", size: CalcStructHeader32(1 << 16), write: func(buf []byte, offset int) int { return WriteStructHeader32AsMap(buf, 1<<16, offset) }, want: append(bytesOf(def.Map32), be32uint(1<<16)...)},
		{name: "time", size: CalcTime(instant), write: func(buf []byte, offset int) int { return WriteTime(buf, instant, offset) }, want: concat(bytesOf(def.Ext8, 12, def.TimeStamp), be32uint(123), be64uint(1<<35))},
	}

	for _, tt := range tests {
		assertWriter(t, tt)
	}
}

func TestWriterBoundaries(t *testing.T) {
	string31 := stringsOfLen(31)
	string32 := stringsOfLen(32)
	string255 := stringsOfLen(math.MaxUint8)
	string256 := stringsOfLen(math.MaxUint8 + 1)
	string65535 := stringsOfLen(math.MaxUint16)
	string65536 := stringsOfLen(math.MaxUint16 + 1)
	timestamp32 := time.Unix(math.MaxUint32, 0).UTC()
	timestamp64 := time.Unix(1, 1).UTC()
	timestamp96 := time.Unix(1<<34, 0).UTC()

	tests := []writerCase{
		{name: "int positive fix max", size: CalcInt(math.MaxInt8), write: func(buf []byte, offset int) int { return WriteInt(buf, math.MaxInt8, offset) }, want: bytesOf(0x7f)},
		{name: "int uint8 min", size: CalcInt(math.MaxInt8 + 1), write: func(buf []byte, offset int) int { return WriteInt(buf, math.MaxInt8+1, offset) }, want: bytesOf(def.Uint8, 0x80)},
		{name: "int negative fix min", size: CalcInt(-32), write: func(buf []byte, offset int) int { return WriteInt(buf, -32, offset) }, want: bytesOf(0xe0)},
		{name: "int int8 min", size: CalcInt(-33), write: func(buf []byte, offset int) int { return WriteInt(buf, -33, offset) }, want: bytesOf(def.Int8, 0xdf)},
		{name: "uint positive fix max", size: CalcUint(math.MaxInt8), write: func(buf []byte, offset int) int { return WriteUint(buf, math.MaxInt8, offset) }, want: bytesOf(0x7f)},
		{name: "uint8 min", size: CalcUint(math.MaxInt8 + 1), write: func(buf []byte, offset int) int { return WriteUint(buf, math.MaxInt8+1, offset) }, want: bytesOf(def.Uint8, 0x80)},
		{name: "uint16 min", size: CalcUint(math.MaxUint8 + 1), write: func(buf []byte, offset int) int { return WriteUint(buf, math.MaxUint8+1, offset) }, want: bytesOf(def.Uint16, 0x01, 0x00)},
		{name: "uint32 min", size: CalcUint(math.MaxUint16 + 1), write: func(buf []byte, offset int) int { return WriteUint(buf, math.MaxUint16+1, offset) }, want: append(bytesOf(def.Uint32), be32uint(math.MaxUint16+1)...)},
		{name: "uint64 min", size: CalcUint64(math.MaxUint32 + 1), write: func(buf []byte, offset int) int { return WriteUint64(buf, math.MaxUint32+1, offset) }, want: append(bytesOf(def.Uint64), be64uint(math.MaxUint32+1)...)},
		{name: "string fix max", size: CalcString(string31), write: func(buf []byte, offset int) int { return WriteString(buf, string31, offset) }, want: concat(bytesOf(def.FixStr+31), []byte(string31))},
		{name: "string8 min", size: CalcString(string32), write: func(buf []byte, offset int) int { return WriteString(buf, string32, offset) }, want: concat(bytesOf(def.Str8, 32), []byte(string32))},
		{name: "string8 max", size: CalcString(string255), write: func(buf []byte, offset int) int { return WriteString(buf, string255, offset) }, want: concat(bytesOf(def.Str8, 0xff), []byte(string255))},
		{name: "string16 min", size: CalcString(string256), write: func(buf []byte, offset int) int { return WriteString(buf, string256, offset) }, want: concat(bytesOf(def.Str16, 0x01, 0x00), []byte(string256))},
		{name: "string16 max", size: CalcString(string65535), write: func(buf []byte, offset int) int { return WriteString(buf, string65535, offset) }, want: concat(bytesOf(def.Str16, 0xff, 0xff), []byte(string65535))},
		{name: "string32 min", size: CalcString(string65536), write: func(buf []byte, offset int) int { return WriteString(buf, string65536, offset) }, want: concat(bytesOf(def.Str32), be32uint(uint64(len(string65536))), []byte(string65536))},
		{name: "slice fix max", size: mustSize(func() (int, error) { return CalcSliceLength(15, false) }), write: func(buf []byte, offset int) int { return WriteSliceLength(buf, 15, offset, false) }, want: bytesOf(def.FixArray + 15)},
		{name: "slice16 min", size: mustSize(func() (int, error) { return CalcSliceLength(16, false) }), write: func(buf []byte, offset int) int { return WriteSliceLength(buf, 16, offset, false) }, want: bytesOf(def.Array16, 0x00, 0x10)},
		{name: "slice16 max", size: mustSize(func() (int, error) { return CalcSliceLength(math.MaxUint16, false) }), write: func(buf []byte, offset int) int { return WriteSliceLength(buf, math.MaxUint16, offset, false) }, want: bytesOf(def.Array16, 0xff, 0xff)},
		{name: "slice32 min", size: mustSize(func() (int, error) { return CalcSliceLength(math.MaxUint16+1, false) }), write: func(buf []byte, offset int) int { return WriteSliceLength(buf, math.MaxUint16+1, offset, false) }, want: append(bytesOf(def.Array32), be32uint(math.MaxUint16+1)...)},
		{name: "byte slice8 max", size: mustSize(func() (int, error) { return CalcSliceLength(math.MaxUint8, true) }), write: func(buf []byte, offset int) int { return WriteSliceLength(buf, math.MaxUint8, offset, true) }, want: bytesOf(def.Bin8, 0xff)},
		{name: "byte slice16 min", size: mustSize(func() (int, error) { return CalcSliceLength(math.MaxUint8+1, true) }), write: func(buf []byte, offset int) int { return WriteSliceLength(buf, math.MaxUint8+1, offset, true) }, want: bytesOf(def.Bin16, 0x01, 0x00)},
		{name: "byte slice16 max", size: mustSize(func() (int, error) { return CalcSliceLength(math.MaxUint16, true) }), write: func(buf []byte, offset int) int { return WriteSliceLength(buf, math.MaxUint16, offset, true) }, want: bytesOf(def.Bin16, 0xff, 0xff)},
		{name: "byte slice32 min", size: mustSize(func() (int, error) { return CalcSliceLength(math.MaxUint16+1, true) }), write: func(buf []byte, offset int) int { return WriteSliceLength(buf, math.MaxUint16+1, offset, true) }, want: append(bytesOf(def.Bin32), be32uint(math.MaxUint16+1)...)},
		{name: "map fix max", size: mustSize(func() (int, error) { return CalcMapLength(15) }), write: func(buf []byte, offset int) int { return WriteMapLength(buf, 15, offset) }, want: bytesOf(def.FixMap + 15)},
		{name: "map16 min", size: mustSize(func() (int, error) { return CalcMapLength(16) }), write: func(buf []byte, offset int) int { return WriteMapLength(buf, 16, offset) }, want: bytesOf(def.Map16, 0x00, 0x10)},
		{name: "map16 max", size: mustSize(func() (int, error) { return CalcMapLength(math.MaxUint16) }), write: func(buf []byte, offset int) int { return WriteMapLength(buf, math.MaxUint16, offset) }, want: bytesOf(def.Map16, 0xff, 0xff)},
		{name: "map32 min", size: mustSize(func() (int, error) { return CalcMapLength(math.MaxUint16 + 1) }), write: func(buf []byte, offset int) int { return WriteMapLength(buf, math.MaxUint16+1, offset) }, want: append(bytesOf(def.Map32), be32uint(math.MaxUint16+1)...)},
		{name: "time timestamp32", size: CalcTime(timestamp32), write: func(buf []byte, offset int) int { return WriteTime(buf, timestamp32, offset) }, want: concat(bytesOf(def.Fixext4, def.TimeStamp), be32uint(math.MaxUint32))},
		{name: "time timestamp64", size: CalcTime(timestamp64), write: func(buf []byte, offset int) int { return WriteTime(buf, timestamp64, offset) }, want: concat(bytesOf(def.Fixext8, def.TimeStamp), be64uint(1<<34|1))},
		{name: "time timestamp96", size: CalcTime(timestamp96), write: func(buf []byte, offset int) int { return WriteTime(buf, timestamp96, offset) }, want: concat(bytesOf(def.Ext8, 12, def.TimeStamp), be32uint(0), be64uint(1<<34))},
	}

	for _, tt := range tests {
		assertWriter(t, tt)
	}
}

func TestLengthErrors(t *testing.T) {
	if strconv.IntSize < 64 {
		t.Skip("unsupported MaxUint32+1 length requires 64-bit int")
	}

	tooLong := int(uint64(math.MaxUint32) + 1)
	if _, err := CalcSliceLength(tooLong, false); err == nil {
		t.Fatal("CalcSliceLength error = nil")
	}
	if _, err := CalcSliceLength(tooLong, true); err == nil {
		t.Fatal("CalcSliceLength byte slice error = nil")
	}
	if _, err := CalcMapLength(tooLong); err == nil {
		t.Fatal("CalcMapLength error = nil")
	}
}

func TestTimeEncodingDefaultsToUTC(t *testing.T) {
	local := time.FixedZone("JST", 9*60*60)
	localTime := time.Date(2026, 5, 22, 18, 45, 30, 987654321, local)
	utcTime := localTime.UTC()

	if CalcTime(localTime) != CalcTime(utcTime) {
		t.Fatalf("CalcTime(local) = %d, want %d", CalcTime(localTime), CalcTime(utcTime))
	}

	size := CalcTime(localTime)
	localEncoded := RequireAt(nil, 0, size)
	localEncodedOffset := WriteTime(localEncoded, localTime, 0)
	utcEncoded := RequireAt(nil, 0, size)
	utcEncodedOffset := WriteTime(utcEncoded, utcTime, 0)

	if localEncodedOffset != utcEncodedOffset {
		t.Fatalf("offsets = local:%d utc:%d", localEncodedOffset, utcEncodedOffset)
	}
	if !bytes.Equal(localEncoded[:localEncodedOffset], utcEncoded[:utcEncodedOffset]) {
		t.Fatalf("local bytes = %x, want utc %x", localEncoded[:localEncodedOffset], utcEncoded[:utcEncodedOffset])
	}
}

type writerCase struct {
	name  string
	size  int
	write func([]byte, int) int
	want  []byte
}

func assertWriter(t *testing.T, tt writerCase) {
	t.Helper()
	t.Run(tt.name, func(t *testing.T) {
		if tt.size != len(tt.want) {
			t.Fatalf("size = %d, want %d", tt.size, len(tt.want))
		}

		prefix := []byte{0xaa, 0xbb}
		buf := append([]byte(nil), prefix...)
		buf = RequireAt(buf, len(prefix), tt.size)
		offset := tt.write(buf, len(prefix))

		if offset != len(prefix)+tt.size {
			t.Fatalf("offset = %d, want %d", offset, len(prefix)+tt.size)
		}
		if !bytes.Equal(buf[:len(prefix)], prefix) {
			t.Fatalf("prefix = %x, want %x", buf[:len(prefix)], prefix)
		}
		if !bytes.Equal(buf[len(prefix):offset], tt.want) {
			t.Fatalf("bytes = %x, want %x", buf[len(prefix):offset], tt.want)
		}
	})
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

func bytesOf(values ...int) []byte {
	b := make([]byte, len(values))
	for i, v := range values {
		b[i] = byte(v)
	}
	return b
}

func concat(parts ...[]byte) []byte {
	var out []byte
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

func be32uint(v uint64) []byte {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(v))
	return b[:]
}

func be32int(v int64) []byte {
	return be32uint(uint64(uint32(int32(v))))
}

func be64uint(v uint64) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], v)
	return b[:]
}

func be64int(v int64) []byte {
	return be64uint(uint64(v))
}
