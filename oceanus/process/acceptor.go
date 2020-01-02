package process

import (
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/tcp"
	"github.com/laconiz/eros/oceanus"
	"time"
)

func (p *Process) Run() {

	invoker := network.NewStdInvoker()

	invoker.Register(&oceanus.Message{}, func(event *network.Event) {
		p.Push(event.Msg.(*oceanus.Message))
	})

	conf := tcp.AcceptorConfig{
		Name: "oceanus",
		Addr: ":4370",
		Session: tcp.SessionConfig{
			ReadTimeout:  time.Minute,
			WriteTimeout: time.Minute,
			LogLevel:     log.Warn,
			QueueLen:     64,
			Invoker:      invoker,
		},
	}

	p.acceptor = tcp.NewAcceptor(conf)

	go p.acceptor.Run()
}
