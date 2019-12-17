package main

import (
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/tcp"
	"github.com/laconiz/eros/network/tcp/examples"
	"log"
	"time"
)

func main() {

	flag := int64(0)

	invoker := network.NewStdInvoker()

	invoker.Register(network.Connected{}, func(event *network.Event) {
		event.Session.Send(examples.REQ{Int: flag})
	})

	invoker.Register(examples.ACK{}, func(event *network.Event) {
		flag = event.Msg.(*examples.ACK).Int + 1
		event.Session.Send(examples.REQ{Int: flag})
	})

	conf := tcp.ConnectorConfig{
		Addr:      "192.168.10.106:12313",
		Reconnect: true,
	}
	conf.Session.Invoker = invoker

	connector := tcp.NewConnector(conf)
	connector.Run()

	last := flag
	for {
		select {
		case <-time.After(time.Second):
			if connector.Connected() {
				log.Println(flag-last, flag, last)
			}

			last = flag
		}
	}
}
