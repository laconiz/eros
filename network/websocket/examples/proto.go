package examples

import "github.com/laconiz/eros/network"

type REQ struct {
	ID uint64
}

type ACK struct {
	ID uint64
}

func init() {
	network.RegisterMeta(REQ{}, network.JsonCodec)
	network.RegisterMeta(ACK{}, network.JsonCodec)
}
