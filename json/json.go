package json

import (
	"encoding/json"
	"io"
	"reflect"
)

type Decoder struct {
	dec *json.Decoder
}

func DecoderMake(r io.Reader) Decoder {
	return Decoder{
		dec: json.NewDecoder(r),
	}
}

func (x Decoder) Decode(t reflect.Type) (reflect.Value, error) {
	o := reflect.New(t)
	err := x.dec.Decode(o.Interface())
	return o.Elem(), err
}

type Encoder struct {
	enc *json.Encoder
}

func EncoderMake(w io.Writer) Encoder {
	return Encoder{
		enc: json.NewEncoder(w),
	}
}

func (x Encoder) Encode(v reflect.Value) error {
	return x.enc.Encode(v.Interface())
}
