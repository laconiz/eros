package gateway

import (
	"github.com/laconiz/eros/iris/proto"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/websocket"
)

func New(addr string) *Thread {

	g := newGateway()
	t := &Thread{gateway: g}

	invoker := network.NewStdInvoker()

	invoker.Register(network.Connected{}, func(event *network.Event) {

	})

	invoker.Register(proto.UserAuthREQ{}, func(event *network.Event) {

	})

	invoker.Register(network.Disconnected{}, func(event *network.Event) {

	})

	conf := websocket.AcceptorConfig{
		Name:   "gateway",
		Addr:   addr,
		Verify: nil,
		Session: websocket.SessionConfig{
			Encoder: &encoder{},
			Invoker: invoker,
		},
	}

	t.acceptor = websocket.NewAcceptor(conf)

	return t
}

type Thread struct {
	gateway  *Gateway
	acceptor *websocket.Acceptor
}

func (t *Thread) OnStart() {
	go t.acceptor.Start()
}

func (t *Thread) OnMessage() {

}

func (t *Thread) OnStop() {

}
