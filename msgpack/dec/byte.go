package dec

func (d *Decoder) AsByte(offset int) (byte, int, error) {
	b, offset := d.readSize1(offset)
	return b, offset, nil
}
