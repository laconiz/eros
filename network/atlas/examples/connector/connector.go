package main

import (
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/atlas"
	"github.com/laconiz/eros/network/atlas/examples"
	"sync/atomic"
	"time"
)

func main() {

	flag := int64(0)

	invoker := network.NewStdInvoker()

	// invoker.Register(network.Connected{}, func(event *network.Event) {
	// 	event.Session.Send(examples.REQ{Int: flag})
	// })
	//
	// invoker.Register(examples.ACK{}, func(event *network.Event) {
	// 	flag = event.Msg.(*examples.ACK).Int + 1
	// 	event.Session.Send(examples.REQ{Int: flag})
	// })

	conf := atlas.ConnectorConfig{
		Addr:      "192.168.10.108:12313",
		Reconnect: true,
	}
	conf.Session.Invoker = invoker
	conf.Session.LogLevel = log.Warn

	connector := atlas.NewConnector(conf)
	connector.Run()

	go func() {
		for {
			connector.Send(examples.REQ{Int: atomic.AddInt64(&flag, 1)})
		}
	}()

	last := flag
	for {
		select {
		case <-time.After(time.Second):
			if connector.Connected() {
				log.Std("main").Info(flag - last)
			}

			last = flag
		}
	}
}
