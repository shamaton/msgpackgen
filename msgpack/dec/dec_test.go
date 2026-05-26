package dec

import (
	"errors"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
)

func TestDecoderReturnsErrorForTruncatedInput(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*Decoder) error
	}{
		{
			name: "AsBool empty",
			fn: func(d *Decoder) error {
				_, _, err := d.AsBool(0)
				return err
			},
		},
		{
			name: "AsByte empty",
			fn: func(d *Decoder) error {
				_, _, err := d.AsByte(0)
				return err
			},
		},
		{
			name: "AsInt truncated uint16",
			fn: func(d *Decoder) error {
				_, _, err := d.AsInt(0)
				return err
			},
		},
		{
			name: "AsUint truncated uint16",
			fn: func(d *Decoder) error {
				_, _, err := d.AsUint(0)
				return err
			},
		},
		{
			name: "AsFloat32 truncated float32",
			fn: func(d *Decoder) error {
				_, _, err := d.AsFloat32(0)
				return err
			},
		},
		{
			name: "AsFloat64 truncated float64",
			fn: func(d *Decoder) error {
				_, _, err := d.AsFloat64(0)
				return err
			},
		},
		{
			name: "AsComplex64 truncated fixext8",
			fn: func(d *Decoder) error {
				_, _, err := d.AsComplex64(0)
				return err
			},
		},
		{
			name: "AsComplex128 truncated fixext16",
			fn: func(d *Decoder) error {
				_, _, err := d.AsComplex128(0)
				return err
			},
		},
		{
			name: "AsStringBytes truncated str16",
			fn: func(d *Decoder) error {
				_, _, err := d.AsStringBytes(0)
				return err
			},
		},
		{
			name: "AsStringBytes truncated payload",
			fn: func(d *Decoder) error {
				_, _, err := d.AsStringBytes(0)
				return err
			},
		},
		{
			name: "SliceLength truncated array16",
			fn: func(d *Decoder) error {
				_, _, err := d.SliceLength(0)
				return err
			},
		},
		{
			name: "MapLength truncated map16",
			fn: func(d *Decoder) error {
				_, _, err := d.MapLength(0)
				return err
			},
		},
		{
			name: "CheckStructHeader truncated array16",
			fn: func(d *Decoder) error {
				_, err := d.CheckStructHeader(1, 0)
				return err
			},
		},
		{
			name: "AsDateTime truncated fixext4",
			fn: func(d *Decoder) error {
				_, _, err := d.AsDateTime(0)
				return err
			},
		},
	}

	inputs := map[string][]byte{
		"AsBool empty":                        {},
		"AsByte empty":                        {},
		"AsInt truncated uint16":              {def.Uint16},
		"AsUint truncated uint16":             {def.Uint16},
		"AsFloat32 truncated float32":         {def.Float32},
		"AsFloat64 truncated float64":         {def.Float64},
		"AsComplex64 truncated fixext8":       {def.Fixext8, byte(def.ComplexTypeCode()), 0},
		"AsComplex128 truncated fixext16":     {def.Fixext16, byte(def.ComplexTypeCode()), 0},
		"AsStringBytes truncated str16":       {def.Str16, 0},
		"AsStringBytes truncated payload":     {def.Str8, 2, 'a'},
		"SliceLength truncated array16":       {def.Array16, 0},
		"MapLength truncated map16":           {def.Map16, 0},
		"CheckStructHeader truncated array16": {def.Array16, 0},
		"AsDateTime truncated fixext4":        {def.Fixext4, 0xff, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn(NewDecoder(inputs[tt.name]))
			if !errors.Is(err, def.ErrTooShortBytes) {
				t.Fatalf("error = %v, want %v", err, def.ErrTooShortBytes)
			}
		})
	}
}

func TestIsCodeNilCheckedBoundsSafe(t *testing.T) {
	d := NewDecoder(nil)
	if _, err := d.IsCodeNilChecked(0); !errors.Is(err, def.ErrTooShortBytes) {
		t.Fatalf("IsCodeNilChecked error = %v, want %v", err, def.ErrTooShortBytes)
	}
}
