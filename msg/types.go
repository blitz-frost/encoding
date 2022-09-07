package msg

import (
	"reflect"
)

type Void struct{}

func (x Void) Close() error {
	return nil
}

func (x Void) Encode(v reflect.Value) error {
	return nil
}

func (x Void) ReaderTake(r Reader) error {
	return r.Close()
}
