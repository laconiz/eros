package main

import (
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/tcp"
	"github.com/laconiz/eros/oceanus"
)

func main() {

	channels := map[string]*oceanus.ChannelInfo{}

	invoker := network.NewStdInvoker()

	conf := tcp.AcceptorConfig{
		Name: "lookup",
		Addr: ":4369",
		Session: tcp.SessionConfig{
			Invoker: invoker,
		},
	}

	acceptor := tcp.NewAcceptor(conf)

	invoker.Register(network.Connected{}, func(event *network.Event) {

		var list []*oceanus.ChannelInfo
		for _, channel := range channels {
			list = append(list, channel)
		}

		event.Session.Send(oceanus.Channels{Channels: list})
	})

	invoker.Register(oceanus.Channels{}, func(event *network.Event) {

	})

	invoker.Register(oceanus.ChannelStates{}, func(event *network.Event) {

	})

	acceptor.Run()
}
