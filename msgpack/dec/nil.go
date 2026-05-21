package dec

import "github.com/shamaton/msgpack/v3/def"

// IsCodeNil returns true if the next data is def.Nil.
func (d *Decoder) IsCodeNil(offset int) bool {
	if !d.canRead(offset, def.Byte1) {
		return false
	}
	return def.Nil == d.data[offset]
}

// IsCodeNilChecked returns true if the next data is def.Nil.
func (d *Decoder) IsCodeNilChecked(offset int) (bool, error) {
	if !d.canRead(offset, def.Byte1) {
		return false, def.ErrTooShortBytes
	}
	return def.Nil == d.data[offset], nil
}
