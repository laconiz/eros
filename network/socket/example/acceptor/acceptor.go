package main

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/socket"
	"github.com/laconiz/eros/network/socket/example"
	"os"
	"sync/atomic"
	"time"
)

func main() {

	flag := int64(0)

	invoker := network.NewStdInvoker()
	invoker.Register(example.REQ{}, func(event *network.Event) {
		flag++
	})

	opt := socket.AccOption{
		Addr:    example.Addr,
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
			acc.Broadcast(&example.ACK{})
		}
	}
}

var log = logis.NewHook(logis.NewTextFormatter()).AddWriter(logis.DEBUG, os.Stdout).Entry()
