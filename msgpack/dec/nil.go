package dec

import "github.com/shamaton/msgpack/v2/def"

// IsCodeNil returns true if the next data is def.Nil.
func (d *Decoder) IsCodeNil(offset int) bool {
	return def.Nil == d.data[offset]
}
