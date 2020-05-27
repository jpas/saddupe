package packet

import "reflect"

type Encoder interface {
	Encode() ([]byte, error)
}

type Decoder interface {
	Decode([]byte) error
}

type EncodeDecoder interface {
	Encoder
	Decoder
}

func decode(b []byte, target Decoder) (Decoder, error) {
	typ := reflect.TypeOf(target)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	d := reflect.New(typ).Interface().(Decoder)
	err := d.Decode(b)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func boolToBit(b bool, offset int) uint64 {
	if b {
		return 1 << offset
	}
	return 0
}
func boolsToByte(b0, b1, b2, b3, b4, b5, b6, b7 bool) byte {
	hi := boolToBit(b7, 7) | boolToBit(b6, 6) | boolToBit(b5, 5) | boolToBit(b4, 4)
	lo := boolToBit(b3, 3) | boolToBit(b2, 2) | boolToBit(b1, 1) | boolToBit(b0, 0)
	return byte(hi | lo)
}

func bitIsSet(b byte, offset int) bool {
	mask := byte(1 << offset)
	return b&mask == mask
}

func byteToBools(b byte) (bool, bool, bool, bool, bool, bool, bool, bool) {
	b0 := bitIsSet(b, 0)
	b1 := bitIsSet(b, 1)
	b2 := bitIsSet(b, 2)
	b3 := bitIsSet(b, 3)
	b4 := bitIsSet(b, 4)
	b5 := bitIsSet(b, 5)
	b6 := bitIsSet(b, 6)
	b7 := bitIsSet(b, 7)
	return b0, b1, b2, b3, b4, b5, b6, b7
}

func bit(offset int, b bool) byte {
	if b {
		return 1 << offset
	}
	return 0
}
