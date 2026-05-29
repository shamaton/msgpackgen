package fallback

import (
	"errors"
	"testing"
	"time"

	"github.com/shamaton/msgpack/v3/def"
)

func TestDecodeReturnsErrorForTruncatedInput(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		target func() any
		fn     func([]byte, any) error
	}{
		{
			name: "array time truncated fixext4",
			data: []byte{def.Fixext4, 0xff},
			target: func() any {
				var v time.Time
				return &v
			},
			fn: UnmarshalAsArray,
		},
		{
			name: "array string truncated payload",
			data: []byte{def.Str8, 2, 'a'},
			target: func() any {
				var v string
				return &v
			},
			fn: UnmarshalAsArray,
		},
		{
			name: "array slice truncated header",
			data: []byte{def.Array16, 0},
			target: func() any {
				var v []int
				return &v
			},
			fn: UnmarshalAsArray,
		},
		{
			name: "map truncated header",
			data: []byte{def.Map16, 0},
			target: func() any {
				var v map[string]int
				return &v
			},
			fn: UnmarshalAsMap,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn(tt.data, tt.target())
			if !errors.Is(err, def.ErrTooShortBytes) {
				t.Fatalf("error = %v, want %v", err, def.ErrTooShortBytes)
			}
		})
	}
}
