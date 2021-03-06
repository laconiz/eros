package main

import (
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/invoker"
	"github.com/laconiz/eros/network/socket"
	"github.com/laconiz/eros/network/socket/example"
	"time"
)

var times uint64

func NewAcceptor() *socket.Acceptor {

	invoker := invoker.NewSocketInvoker()
	invoker.Register(example.REQ{}, func(event *network.Event) {
		req := event.Msg.(*example.REQ)
		event.Ses.Send(&example.ACK{Time: req.Time})
		times++
	})

	option := &socket.AcceptorOption{
		Addr: example.Addr,
		Session: socket.SessionOption{
			Invoker: invoker,
		},
	}

	return socket.NewAcceptor(option)
}

func main() {

	acceptor := NewAcceptor()
	acceptor.Run()

	for {
		<-time.After(time.Second)
		logger.Info(times)
		times = 0
	}
}

var logger = logisor.Module("socket.example.acceptor")
