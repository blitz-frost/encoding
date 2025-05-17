package msg

import "github.com/blitz-frost/msg"

type ExchangeLoop struct {
	v msg.Loop

	in ExchangeInput
}

func (x *ExchangeLoop) Close() {
	if x.v.Final == nil {
		x.v.Final = func() error {
			return x.in.Cancel()
		}
	}
	x.v.Close()
}

func (x *ExchangeLoop) Final(taker ExchangeInputTaker) {
	x.v.Final = func() error {
		return taker(x.in)
	}
}

func (x *ExchangeLoop) Run(inlet ExchangeInlet, taker ExchangeInputTaker) error {
	f1 := func() error {
		var err error
		x.in, err = inlet()
		return err
	}

	f2 := func() error {
		return taker(x.in)
	}

	return x.v.Run(f1, f2)
}

type Loop struct {
	v msg.Loop

	in Input
}

func (x *Loop) Close() {
	if x.v.Final == nil {
		x.v.Final = func() error {
			return x.in.Close()
		}
	}
	x.v.Close()
}

func (x *Loop) Final(taker InputTaker) {
	x.v.Final = func() error {
		return taker(x.in)
	}
}

func (x *Loop) Run(inlet Inlet, taker InputTaker) error {
	f1 := func() error {
		var err error
		x.in, err = inlet()
		return err
	}

	f2 := func() error {
		return taker(x.in)
	}

	return x.v.Run(f1, f2)
}
