package dec

import (
	"github.com/shamaton/msgpack/v3/def"
)

// AsBool checks codes and returns the got bytes as bool
func (d *Decoder) AsBool(offset int) (bool, int, error) {
	code, offset, err := d.readSize1(offset)
	if err != nil {
		return false, 0, err
	}

	switch code {
	case def.True:
		return true, offset, nil
	case def.False, def.Nil:
		return false, offset, nil
	}
	return false, 0, d.errorTemplate(code, "AsBool")
}
