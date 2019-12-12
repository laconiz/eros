package examples

import "github.com/laconiz/eros/network"

type REQ struct {
	ID uint
}

type ACK struct {
	ID uint
}

func init() {
	network.RegisterMeta(REQ{}, network.JsonCodec)
	network.RegisterMeta(ACK{}, network.JsonCodec)
}
