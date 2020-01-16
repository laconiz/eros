package network

import "github.com/laconiz/eros/message"

type HandlerFunc func(*Event)

type Event struct {
	ID  message.ID
	Msg interface{}
	Ses Session
}

// ---------------------------------------------------------------------------------------------------------------------

type Connected struct {
}

type Disconnected struct {
}

type ConnectFailed struct {
}

var (
	ConnectedMetaID     message.ID
	DisconnectedMetaID  message.ID
	ConnectFailedMetaID message.ID
)

func init() {

	meta, err := message.Register(Connected{}, message.JsonCodec())
	if err != nil {
		panic(err)
	}
	ConnectedMetaID = meta.ID()

	meta, err = message.Register(Disconnected{}, message.JsonCodec())
	if err != nil {
		panic(err)
	}
	DisconnectedMetaID = meta.ID()

	meta, err = message.Register(ConnectFailed{}, message.JsonCodec())
	if err != nil {
		panic(err)
	}
	ConnectFailedMetaID = meta.ID()
}
