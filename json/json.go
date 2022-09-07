package json

import (
	"encoding/json"
	"reflect"

	"github.com/blitz-frost/encoding"
	"github.com/blitz-frost/io"
)

var Codec encoding.Codec = codec{}

type codec struct{}

func (x codec) Decoder(r io.Reader) (encoding.Decoder, error) {
	return MakeDecoder(r), nil
}

func (x codec) Encoder(w io.Writer) (encoding.Encoder, error) {
	return MakeEncoder(w), nil
}

type Decoder struct {
	c   encoding.Closer
	dec *json.Decoder
}

func MakeDecoder(r io.Reader) Decoder {
	return Decoder{
		c:   r,
		dec: json.NewDecoder(r),
	}
}

func (x Decoder) Close() error {
	return x.c.Close()
}

func (x Decoder) Decode(t reflect.Type) (reflect.Value, error) {
	o := reflect.New(t)
	err := x.dec.Decode(o.Interface())
	return o.Elem(), err
}

type Encoder struct {
	c   encoding.Closer
	enc *json.Encoder
}

func MakeEncoder(w io.Writer) Encoder {
	return Encoder{
		c:   w,
		enc: json.NewEncoder(w),
	}
}

func (x Encoder) Close() error {
	return x.c.Close()
}

func (x Encoder) Encode(v reflect.Value) error {
	return x.enc.Encode(v.Interface())
}
