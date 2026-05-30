package dec

// AsByte checks codes and returns the got bytes as byte
func (d *Decoder) AsByte(offset int) (byte, int, error) {
	b, offset, err := d.readSize1(offset)
	if err != nil {
		return 0, 0, err
	}
	return b, offset, nil
}
