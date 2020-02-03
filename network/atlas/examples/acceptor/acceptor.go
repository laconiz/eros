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

	invoker := network.NewStdInvoker()

	flag := int64(0)

	invoker.Register(examples.REQ{}, func(event *network.Event) {
		atomic.StoreInt64(&flag, event.Msg.(*examples.REQ).Int)
		// event.Session.Send(&examples.ACK{Int: event.Msg.(*examples.REQ).Int})
	})

	conf := atlas.AcceptorOption{Addr: ":12313"}
	conf.Session.Invoker = invoker
	conf.Session.LogLevel = log.Warn

	acceptor := atlas.NewAcceptor(conf)
	go acceptor.Run()

	last := flag
	for {
		select {
		case <-time.After(time.Second):
			log.Std("main").Info(flag - last)
			last = flag
		}
	}
}
