package examples

import "github.com/laconiz/eros/network"

type REQ struct {
	Int int64
}

type ACK struct {
	Int int64
}

func init() {
	network.RegisterMeta(REQ{}, network.JsonCodec)
	network.RegisterMeta(ACK{}, network.JsonCodec)
}
