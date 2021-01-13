package dec

func (d *Decoder) AsRune(offset int) (rune, int, error) {
	v, offset, err := d.asInt(offset)
	return rune(v), offset, err
}
