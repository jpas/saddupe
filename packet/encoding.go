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
