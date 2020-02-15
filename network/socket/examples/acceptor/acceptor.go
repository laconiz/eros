package main

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/socket"
	"github.com/laconiz/eros/network/socket/examples"
	"sync/atomic"
	"time"
)

func main() {

	flag := int64(0)

	invoker := network.NewStdInvoker()
	invoker.Register(examples.REQ{}, func(event *network.Event) {
		flag++
	})

	opt := socket.AccOption{
		Addr:    examples.Addr,
		Session: socket.SesOption{Invoker: invoker},
	}
	acc := socket.NewAcceptor(opt)
	acc.Run()

	last := flag

	lt := time.NewTicker(time.Second)
	defer lt.Stop()
	bt := time.NewTicker(time.Second * 10)
	defer bt.Stop()

	for {
		select {
		case <-lt.C:
			log.Infof("sessions: %d count: %d", acc.Count(), atomic.LoadInt64(&flag)-last)
			last = flag
		case <-bt.C:
			acc.Broadcast(&examples.ACK{})
		}
	}
}

var log = logisor.Field(logis.Module, "main")
