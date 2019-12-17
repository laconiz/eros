package main

import (
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/tcp"
	"github.com/laconiz/eros/network/tcp/examples"
)

func main() {

	invoker := network.NewStdInvoker()

	invoker.Register(examples.REQ{}, func(event *network.Event) {
		event.Session.Send(&examples.ACK{Int: event.Msg.(*examples.REQ).Int})
	})

	conf := tcp.AcceptorConfig{Addr: ":12313"}
	conf.Session.Invoker = invoker

	acceptor := tcp.NewAcceptor(conf)
	acceptor.Run()
}
