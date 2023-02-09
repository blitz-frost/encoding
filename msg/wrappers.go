package msg

import (
	"errors"
	"io"
	"reflect"

	"github.com/blitz-frost/encoding"
	msgio "github.com/blitz-frost/io/msg"
	"github.com/blitz-frost/msg"
)

type exchangeReader struct {
	Reader
	writerGiver
}

func makeExchangeReader(er msgio.ExchangeReader, codec encoding.Codec) (exchangeReader, error) {
	dec, err := codec.Decoder(er)
	if err != nil {
		return exchangeReader{}, err
	}
	return exchangeReader{
		Reader: dec,
		writerGiver: writerGiver{
			wg:    er,
			codec: codec,
		},
	}, nil
}

type exchangeReaderChainer struct {
	ert   ExchangeReaderTaker
	codec encoding.Codec
}

func (x *exchangeReaderChainer) ReaderChain(ert ExchangeReaderTaker) error {
	x.ert = ert
	return nil
}

func (x *exchangeReaderChainer) ReaderTake(er msgio.ExchangeReader) error {
	r, err := makeExchangeReader(er, x.codec)
	if err != nil {
		return err
	}
	return x.ert.ReaderTake(r)
}

type exchangeWriter struct {
	w  Writer
	c  msg.Closer // underlying msgio.ExchangeWriter Close method
	rg readerGiver
}

func makeExchangeWriter(ew msgio.ExchangeWriter, codec encoding.Codec) (exchangeWriter, error) {
	// the encoder sees a regular writer that it might try to close when finished
	// but for an ExchangeWriter that means canceling the message, so we mask it
	enc, err := codec.Encoder(noCloseWriter{ew})
	if err != nil {
		return exchangeWriter{}, err
	}
	return exchangeWriter{
		w: enc,
		c: ew,
		rg: readerGiver{
			rg:    ew,
			codec: codec,
		},
	}, nil
}

func (x exchangeWriter) Close() error {
	if c, ok := x.w.(msg.Canceler); ok {
		return c.Cancel()
	}
	return errors.Join(x.w.Close(), x.c.Close())
}

func (x exchangeWriter) Encode(v reflect.Value) error {
	return x.w.Encode(v)
}

func (x exchangeWriter) Reader() (Reader, error) {
	if err := x.w.Close(); err != nil {
		return nil, err
	}
	return x.rg.Reader()
}

type exchangeWriterGiver struct {
	ewg   msgio.ExchangeWriterGiver
	codec encoding.Codec
}

func (x exchangeWriterGiver) Writer() (ExchangeWriter, error) {
	w, err := x.ewg.Writer()
	if err != nil {
		return nil, err
	}
	return makeExchangeWriter(w, x.codec)
}

type noCloseWriter struct {
	io.Writer
}

func (x noCloseWriter) Close() error {
	return nil
}

type readerChainer struct {
	rt    ReaderTaker
	codec encoding.Codec
}

func (x *readerChainer) ReaderChain(rt ReaderTaker) error {
	x.rt = rt
	return nil
}

func (x *readerChainer) ReaderTake(r msgio.Reader) error {
	dec, err := x.codec.Decoder(r)
	if err != nil {
		return err
	}
	return x.rt.ReaderTake(dec)
}

type readerGiver struct {
	rg    msgio.ReaderGiver
	codec encoding.Codec
}

func (x readerGiver) Reader() (Reader, error) {
	r, err := x.rg.Reader()
	if err != nil {
		return nil, err
	}
	return x.codec.Decoder(r)
}

type writerGiver struct {
	wg    msgio.WriterGiver
	codec encoding.Codec
}

func (x writerGiver) Writer() (Writer, error) {
	w, err := x.wg.Writer()
	if err != nil {
		return nil, err
	}

	return x.codec.Encoder(w)
}

func ConnOf(c msgio.Conn, codec encoding.Codec) (msg.ConnBlock[Reader, Writer], error) {
	rc := &readerChainer{
		codec: codec,
	}
	x := msg.ConnBlock[Reader, Writer]{
		rc,
		writerGiver{
			wg:    c,
			codec: codec,
		},
	}
	return x, c.ReaderChain(rc)
}

func ExchangeConnOf(ec msgio.ExchangeConn, codec encoding.Codec) (msg.ConnBlock[ExchangeReader, ExchangeWriter], error) {
	erc := &exchangeReaderChainer{
		codec: codec,
	}
	x := msg.ConnBlock[ExchangeReader, ExchangeWriter]{
		erc,
		exchangeWriterGiver{ec, codec},
	}
	return x, ec.ReaderChain(erc)
}

func ExchangeReaderOf(er msgio.ExchangeReader, codec encoding.Codec) (ExchangeReader, error) {
	return makeExchangeReader(er, codec)
}

func ExchangeReaderChainerOf(erc msgio.ExchangeReaderChainer, codec encoding.Codec) (ExchangeReaderChainer, error) {
	o := &exchangeReaderChainer{
		codec: codec,
	}
	return o, erc.ReaderChain(o)
}

func ExchangeWriterOf(ew msgio.ExchangeWriter, codec encoding.Codec) (ExchangeWriter, error) {
	return makeExchangeWriter(ew, codec)
}

func ExchangeWriterGiverOf(ewg msgio.ExchangeWriterGiver, codec encoding.Codec) ExchangeWriterGiver {
	return exchangeWriterGiver{
		ewg:   ewg,
		codec: codec,
	}
}

// The returned value is also a [io/msg.ReaderTaker]
func ReaderChainerOf(rc msgio.ReaderChainer, codec encoding.Codec) (ReaderChainer, error) {
	o := &readerChainer{
		codec: codec,
	}
	return o, rc.ReaderChain(o)
}

func ReaderGiverOf(rg msgio.ReaderGiver, codec encoding.Codec) ReaderGiver {
	return readerGiver{
		rg:    rg,
		codec: codec,
	}
}

func WriterGiverOf(wg msgio.WriterGiver, codec encoding.Codec) WriterGiver {
	return writerGiver{
		wg:    wg,
		codec: codec,
	}
}
