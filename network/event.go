package network

import (
	"github.com/laconiz/eros/holder/message"
	"github.com/laconiz/eros/network/session"
)

type HandlerFunc func(*Event)

type Event struct {
	Meta message.Meta
	Msg  interface{}
	Ses  session.Session
}

// ---------------------------------------------------------------------------------------------------------------------

func NewConnectedEvent(ses session.Session) *Event {
	return &Event{Meta: metaConnected, Msg: &Connected{}, Ses: ses}
}

func NewDisconnectedEvent(ses session.Session) *Event {
	return &Event{Meta: metaDisconnected, Msg: &Disconnected{}, Ses: ses}
}

func NewConnectFailedEvent() *Event {
	return &Event{Meta: metaConnectFailed, Msg: &ConnectFailed{}}
}

type Connected struct {
}

type Disconnected struct {
}

type ConnectFailed struct {
}

var (
	metaConnected     message.Meta
	metaDisconnected  message.Meta
	metaConnectFailed message.Meta
)

func init() {

	var err error

	if metaConnected, err = message.Register(Connected{}, message.JsonCodec()); err != nil {
		panic(err)
	}

	if metaDisconnected, err = message.Register(Disconnected{}, message.JsonCodec()); err != nil {
		panic(err)
	}

	if metaConnectFailed, err = message.Register(ConnectFailed{}, message.JsonCodec()); err != nil {
		panic(err)
	}
}
