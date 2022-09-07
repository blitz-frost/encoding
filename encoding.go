// Package encoding provides interfaces analogous to the io package, but oriented towards arbitrary type transmission, rather than raw bytes.
//
// The Encoder and Decoder base interfaces use [reflect.Value], rather than [any], because the reflect package is pretty much always involved in the marshaling of arbitrary types.
// It's better to directly impose it as standard, rather than having implementations switch back and forth between reflect and non-reflect code.
package encoding

import (
	"reflect"

	"github.com/blitz-frost/io"
)

type Closer = io.Closer

type Codec interface {
	Decoder(io.Reader) (Decoder, error)
	Encoder(io.Writer) (Encoder, error)
}

type Decoder interface {
	Decode(reflect.Type) (reflect.Value, error)
	Closer
}

type Encoder interface {
	Encode(reflect.Value) error
	Closer
}
