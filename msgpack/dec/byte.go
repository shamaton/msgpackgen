package dec

// AsByte checks codes and returns the got bytes as byte
func (d *Decoder) AsByte(offset int) (byte, int, error) {
	b, offset := d.readSize1(offset)
	return b, offset, nil
}
