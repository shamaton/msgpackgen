package dec

import "github.com/shamaton/msgpack/v2/def"

func (d *Decoder) IsCodeNil(offset int) bool {
	return def.Nil == d.data[offset]
}
