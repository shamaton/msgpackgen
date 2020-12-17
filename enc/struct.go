package encoding

import (
	"fmt"
	"math"
	"reflect"
	"sync"

	"github.com/shamaton/msgpack/def"
)

type structCache struct {
	indexes []int
	names   []string
	//common.Common
}

var cachemap = sync.Map{}

type structCalcFunc func(rv reflect.Value) (int, error)
type structWriteFunc func(rv reflect.Value, offset int) int

func (e *Encoder) getStructCalc(typ reflect.Type) structCalcFunc {

	for j := range extCoders {
		if extCoders[j].Type() == typ {
			return extCoders[j].CalcByteSize
		}
	}
	if e.asArray {
		return e.calcStructArray
	}
	return e.calcStructMap

}

func (e *Encoder) CalcStruct(num int) (int, error) {
	ret := def.Byte1
	if e.asArray {
		if num <= 0x0f {
			// format code only
		} else if num <= math.MaxUint16 {
			ret += def.Byte2
		} else if uint(num) <= math.MaxUint32 {
			ret += def.Byte4
		} else {
			// not supported error
			return 0, fmt.Errorf("not support this array length : %d", num)
		}
	} else {

		if num <= 0x0f {
			// format code only
		} else if num <= math.MaxUint16 {
			ret += def.Byte2
		} else if uint(num) <= math.MaxUint32 {
			ret += def.Byte4
		} else {
			// not supported error
			return 0, fmt.Errorf("not support this array length : %d", num)
		}
	}
	return ret, nil
}

func (e *Encoder) WriteStruct(num, offset int) int {
	if e.asArray {
		if num <= 0x0f {
			offset = e.setByte1Int(def.FixArray+num, offset)
		} else if num <= math.MaxUint16 {
			offset = e.setByte1Int(def.Array16, offset)
			offset = e.setByte2Int(num, offset)
		} else if uint(num) <= math.MaxUint32 {
			offset = e.setByte1Int(def.Array32, offset)
			offset = e.setByte4Int(num, offset)
		}
	} else {
		if num <= 0x0f {
			offset = e.setByte1Int(def.FixMap+num, offset)
		} else if num <= math.MaxUint16 {
			offset = e.setByte1Int(def.Map16, offset)
			offset = e.setByte2Int(num, offset)
		} else if uint(num) <= math.MaxUint32 {
			offset = e.setByte1Int(def.Map32, offset)
			offset = e.setByte4Int(num, offset)
		}
	}
	return offset
}

func (e *Encoder) calcStruct(rv reflect.Value) (int, error) {

	//if isTime, tm := e.isDateTime(rv); isTime {
	//	size := e.CalcTime(tm)
	//	return size, nil
	//}

	for i := range extCoders {
		if extCoders[i].Type() == rv.Type() {
			return extCoders[i].CalcByteSize(rv)
		}
	}

	if e.asArray {
		return e.calcStructArray(rv)
	}
	return e.calcStructMap(rv)
}

func (e *Encoder) calcStructArray(rv reflect.Value) (int, error) {
	return 0, nil

	/*
		ret := 0
		t := rv.Type()
		cache, find := cachemap.Load(t)
		var c *structCache
		if !find {
			c = &structCache{}
			for i := 0; i < rv.NumField(); i++ {
				field := t.Field(i)
				if ok, name := e.CheckField(field); ok {
					size, err := e.calcSize(rv.Field(i))
					if err != nil {
						return 0, err
					}
					ret += size
					c.indexes = append(c.indexes, i)
					c.names = append(c.names, name)
				}
			}
			cachemap.Store(t, c)
		} else {
			c = cache.(*structCache)
			for i := 0; i < len(c.indexes); i++ {
				size, err := e.calcSize(rv.Field(c.indexes[i]))
				if err != nil {
					return 0, err
				}
				ret += size
			}
		}

		// format size
		l := len(c.indexes)
		if l <= 0x0f {
			// format code only
		} else if l <= math.MaxUint16 {
			ret += def.Byte2
		} else if uint(l) <= math.MaxUint32 {
			ret += def.Byte4
		} else {
			// not supported error
			return 0, fmt.Errorf("not support this array length : %d", l)
		}
		return ret, nil

	*/
}

func (e *Encoder) calcStructMap(rv reflect.Value) (int, error) {
	return 0, nil

	/*
		ret := 0
		t := rv.Type()
		cache, find := cachemap.Load(t)
		var c *structCache
		if !find {
			c = &structCache{}
			for i := 0; i < rv.NumField(); i++ {
				if ok, name := e.CheckField(rv.Type().Field(i)); ok {
					keySize := def.Byte1 + e.CalcString(name)
					valueSize, err := e.calcSize(rv.Field(i))
					if err != nil {
						return 0, err
					}
					ret += keySize + valueSize
					c.indexes = append(c.indexes, i)
					c.names = append(c.names, name)
				}
			}
			cachemap.Store(t, c)
		} else {
			c = cache.(*structCache)
			for i := 0; i < len(c.indexes); i++ {
				keySize := def.Byte1 + e.CalcString(c.names[i])
				valueSize, err := e.calcSize(rv.Field(c.indexes[i]))
				if err != nil {
					return 0, err
				}
				ret += keySize + valueSize
			}
		}

		// format size
		l := len(c.indexes)
		if l <= 0x0f {
			// format code only
		} else if l <= math.MaxUint16 {
			ret += def.Byte2
		} else if uint(l) <= math.MaxUint32 {
			ret += def.Byte4
		} else {
			// not supported error
			return 0, fmt.Errorf("not support this array length : %d", l)
		}
		return ret, nil

	*/
}

func (e *Encoder) getStructWriter(typ reflect.Type) structWriteFunc {

	for i := range extCoders {
		if extCoders[i].Type() == typ {
			return func(rv reflect.Value, offset int) int {
				return extCoders[i].WriteToBytes(rv, offset, &e.d)
			}
		}
	}

	if e.asArray {
		return e.writeStructArray
	}
	return e.writeStructMap
}

func (e *Encoder) writeStruct(rv reflect.Value, offset int) int {
	/*
		if isTime, tm := e.isDateTime(rv); isTime {
			return e.WriteTime(tm, offset)
		}
	*/

	for i := range extCoders {
		if extCoders[i].Type() == rv.Type() {
			return extCoders[i].WriteToBytes(rv, offset, &e.d)
		}
	}

	if e.asArray {
		return e.writeStructArray(rv, offset)
	}
	return e.writeStructMap(rv, offset)
}

func (e *Encoder) writeStructArray(rv reflect.Value, offset int) int {

	cache, _ := cachemap.Load(rv.Type())
	c := cache.(*structCache)

	// write format
	num := len(c.indexes)
	if num <= 0x0f {
		offset = e.setByte1Int(def.FixArray+num, offset)
	} else if num <= math.MaxUint16 {
		offset = e.setByte1Int(def.Array16, offset)
		offset = e.setByte2Int(num, offset)
	} else if uint(num) <= math.MaxUint32 {
		offset = e.setByte1Int(def.Array32, offset)
		offset = e.setByte4Int(num, offset)
	}

	for i := 0; i < num; i++ {
		offset = e.create(rv.Field(c.indexes[i]), offset)
	}
	return offset
}

func (e *Encoder) writeStructMap(rv reflect.Value, offset int) int {

	cache, _ := cachemap.Load(rv.Type())
	c := cache.(*structCache)

	// format size
	num := len(c.indexes)
	if num <= 0x0f {
		offset = e.setByte1Int(def.FixMap+num, offset)
	} else if num <= math.MaxUint16 {
		offset = e.setByte1Int(def.Map16, offset)
		offset = e.setByte2Int(num, offset)
	} else if uint(num) <= math.MaxUint32 {
		offset = e.setByte1Int(def.Map32, offset)
		offset = e.setByte4Int(num, offset)
	}

	for i := 0; i < num; i++ {
		offset = e.WriteString(c.names[i], offset)
		offset = e.create(rv.Field(c.indexes[i]), offset)
	}
	return offset
}
