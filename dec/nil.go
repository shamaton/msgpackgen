package dec

import "github.com/shamaton/msgpack/def"

func (d *Decoder) isCodeNil(v byte) bool {
	return def.Nil == v
}
