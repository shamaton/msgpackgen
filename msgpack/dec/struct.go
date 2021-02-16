package dec

import (
	"encoding/binary"
	"fmt"

	"github.com/shamaton/msgpack/v2/def"
)

// CheckStructHeader checks if fieldNum matches the number of fields on the data.
func (d *Decoder) CheckStructHeader(fieldNum, offset int) (int, error) {
	code, offset := d.readSize1(offset)
	var l int
	switch {
	case d.isFixSlice(code):
		l = int(code - def.FixArray)

	case code == def.Array16:
		bs, o := d.readSize2(offset)
		l = int(binary.BigEndian.Uint16(bs))
		offset = o
	case code == def.Array32:
		bs, o := d.readSize4(offset)
		l = int(binary.BigEndian.Uint32(bs))
		offset = o

	case d.isFixMap(code):
		l = int(code - def.FixMap)
	case code == def.Map16:
		bs, o := d.readSize2(offset)
		l = int(binary.BigEndian.Uint16(bs))
		offset = o
	case code == def.Map32:
		bs, o := d.readSize4(offset)
		l = int(binary.BigEndian.Uint32(bs))
		offset = o
	}

	if fieldNum != l {
		return 0, fmt.Errorf("data length wrong %d : %d", fieldNum, l)
	}
	return offset, nil
}

//func (d *Decoder) JumpOffset(offset int) int {
//	code, offset := d.readSize1(offset)
//	switch {
//	case code == def.True, code == def.False, code == def.Nil:
//		// do nothing
//
//	case d.isPositiveFixNum(code) || d.isNegativeFixNum(code):
//		// do nothing
//	case code == def.Uint8, code == def.Int8:
//		offset += def.Byte1
//	case code == def.Uint16, code == def.Int16:
//		offset += def.Byte2
//	case code == def.Uint32, code == def.Int32, code == def.Float32:
//		offset += def.Byte4
//	case code == def.Uint64, code == def.Int64, code == def.Float64:
//		offset += def.Byte8
//
//	case d.isFixString(code):
//		offset += int(code - def.FixStr)
//	case code == def.Str8, code == def.Bin8:
//		b, o := d.readSize1(offset)
//		o += int(b)
//		offset = o
//	case code == def.Str16, code == def.Bin16:
//		bs, o := d.readSize2(offset)
//		o += int(binary.BigEndian.Uint16(bs))
//		offset = o
//	case code == def.Str32, code == def.Bin32:
//		bs, o := d.readSize4(offset)
//		o += int(binary.BigEndian.Uint32(bs))
//		offset = o
//
//	case d.isFixSlice(code):
//		l := int(code - def.FixArray)
//		for i := 0; i < l; i++ {
//			offset = d.JumpOffset(offset)
//		}
//	case code == def.Array16:
//		bs, o := d.readSize2(offset)
//		l := int(binary.BigEndian.Uint16(bs))
//		for i := 0; i < l; i++ {
//			o = d.JumpOffset(o)
//		}
//		offset = o
//	case code == def.Array32:
//		bs, o := d.readSize4(offset)
//		l := int(binary.BigEndian.Uint32(bs))
//		for i := 0; i < l; i++ {
//			o = d.JumpOffset(o)
//		}
//		offset = o
//
//	case d.isFixMap(code):
//		l := int(code - def.FixMap)
//		for i := 0; i < l*2; i++ {
//			offset = d.JumpOffset(offset)
//		}
//	case code == def.Map16:
//		bs, o := d.readSize2(offset)
//		l := int(binary.BigEndian.Uint16(bs))
//		for i := 0; i < l*2; i++ {
//			o = d.JumpOffset(o)
//		}
//		offset = o
//	case code == def.Map32:
//		bs, o := d.readSize4(offset)
//		l := int(binary.BigEndian.Uint32(bs))
//		for i := 0; i < l*2; i++ {
//			o = d.JumpOffset(o)
//		}
//		offset = o
//
//	case code == def.Fixext1:
//		offset += def.Byte1 + def.Byte1
//	case code == def.Fixext2:
//		offset += def.Byte1 + def.Byte2
//	case code == def.Fixext4:
//		offset += def.Byte1 + def.Byte4
//	case code == def.Fixext8:
//		offset += def.Byte1 + def.Byte8
//	case code == def.Fixext16:
//		offset += def.Byte1 + def.Byte16
//
//	case code == def.Ext8:
//		b, o := d.readSize1(offset)
//		o += def.Byte1 + int(b)
//		offset = o
//	case code == def.Ext16:
//		bs, o := d.readSize2(offset)
//		o += def.Byte1 + int(binary.BigEndian.Uint16(bs))
//		offset = o
//	case code == def.Ext32:
//		bs, o := d.readSize4(offset)
//		o += def.Byte1 + int(binary.BigEndian.Uint32(bs))
//		offset = o
//
//	}
//	return offset
//}
