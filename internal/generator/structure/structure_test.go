package structure

import (
	"bytes"
	"math"
	"strings"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
)

func TestEncodedStringBytes(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		prefix []byte
	}{
		{
			name:   "fixstr",
			value:  strings.Repeat("a", 31),
			prefix: []byte{byte(def.FixStr + 31)},
		},
		{
			name:   "str8",
			value:  strings.Repeat("a", 32),
			prefix: []byte{byte(def.Str8), 32},
		},
		{
			name:   "str16",
			value:  strings.Repeat("a", math.MaxUint8+1),
			prefix: []byte{byte(def.Str16), 0x01, 0x00},
		},
		{
			name:   "str32",
			value:  strings.Repeat("a", math.MaxUint16+1),
			prefix: []byte{byte(def.Str32), 0x00, 0x01, 0x00, 0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodedStringBytes(tt.value)
			if !bytes.HasPrefix(got, tt.prefix) {
				t.Fatalf("prefix = %x, want %x", got[:len(tt.prefix)], tt.prefix)
			}
			if string(got[len(tt.prefix):]) != tt.value {
				t.Fatal("payload mismatch")
			}
			if len(got) != len(tt.prefix)+len(tt.value) {
				t.Fatalf("len = %d, want %d", len(got), len(tt.prefix)+len(tt.value))
			}
		})
	}
}
