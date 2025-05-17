// Package encoding provides interfaces analogous to the io package, but oriented towards arbitrary type transmission, rather than raw bytes.
//
// The Encoder and Decoder base interfaces use [reflect.Value], rather than [any], because the reflect package is pretty much always involved in the marshaling of arbitrary types.
// It's better to directly impose it as standard, rather than having implementations switch back and forth between reflect and non-reflect code.
package encoding

import (
	"reflect"
)

type Decoder interface {
	Decode(reflect.Type) (reflect.Value, error)
}

type Encoder interface {
	Encode(reflect.Value) error
}
