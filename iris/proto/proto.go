package proto

import (
	"github.com/laconiz/eros/codec"
	"github.com/laconiz/eros/network"
)

type Ping struct {
}

type Pong struct {
}

type KickReason uint32

const (
	KickReasonServerClosed = iota
)

type Kick struct {
	Reason KickReason
}

func init() {
	network.RegisterMeta(Ping{}, codec.Json())
	network.RegisterMeta(Pong{}, codec.Json())
	network.RegisterMeta(Kick{}, codec.Json())
}
