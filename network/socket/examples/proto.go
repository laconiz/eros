package examples

import (
	"github.com/laconiz/eros/holder/message"
)

const Addr = "127.0.0.1:1024"

type REQ struct {
	Int int64
}

type ACK struct {
	Int int64
}

func init() {
	message.Register(REQ{}, message.JsonCodec())
	message.Register(ACK{}, message.JsonCodec())
}
