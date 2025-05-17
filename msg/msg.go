package msg

import (
	"reflect"

	"github.com/blitz-frost/io/msg"
)

type Decoder func(reflect.Type) (reflect.Value, error)

type Encoder func(reflect.Value) error

type ExchangeInlet func() (ExchangeInput, error)

type ExchangeInput struct {
	Decode Decoder
	Close  msg.Closer
	Output Outlet
	Cancel msg.Canceler
}

type ExchangeInputTaker func(ExchangeInput) error

type ExchangeOutlet func() (ExchangeOutput, error)

type ExchangeOutput struct {
	Encode Encoder
	Close  msg.Closer
	Input  Inlet
	Cancel msg.Canceler
}

type Inlet func() (Input, error)

type Input struct {
	Decode Decoder
	Close  msg.Closer
}

type InputTaker func(Input) error

type Outlet func() (Output, error)

type Output struct {
	Encode Encoder
	Close  msg.Closer

	// optional
	Cancel msg.Canceler
}
