package msg

import (
	"github.com/blitz-frost/encoding"
	"github.com/blitz-frost/msg"
)

type Conn = msg.Conn[Reader, Writer]

type ExchangeConn = msg.Conn[ExchangeReader, ExchangeWriter]

type ExchangeReader interface {
	Reader
	WriterGiver
}

type ExchangeReaderChainer = msg.ReaderChainer[ExchangeReader]

type ExchangeReaderTaker = msg.ReaderTaker[ExchangeReader]

type ExchangeWriter interface {
	Writer
	ReaderGiver
}

type ExchangeWriterGiver = msg.WriterGiver[ExchangeWriter]

type Reader = encoding.Decoder

type ReaderChainer = msg.ReaderChainer[Reader]

type ReaderGiver = msg.ReaderGiver[Reader]

type ReaderTaker = msg.ReaderTaker[Reader]

type Writer = encoding.Encoder

type WriterGiver = msg.WriterGiver[Writer]
