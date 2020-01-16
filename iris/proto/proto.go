package proto

import (
	"github.com/laconiz/eros/message"
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
	network.RegisterMeta(Ping{}, message.Json())
	network.RegisterMeta(Pong{}, message.Json())
	network.RegisterMeta(Kick{}, message.Json())
}
