package process

import (
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/tcp"
	"github.com/laconiz/eros/oceanus"
	"time"
)

func (p *Process) startSync() {

	inv := network.NewStdInvoker()

	inv.Register(oceanus.SyncMessage{}, func(event *network.Event) {
		p.Sync(event.Msg.(*oceanus.SyncMessage))
	})

	conf := tcp.ConnectorConfig{
		Name:      "oceanus.sync",
		Addr:      "",
		Reconnect: true,
		Session: tcp.SessionConfig{
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
			LogLevel:     log.Warn,
			QueueLen:     64,
			Invoker:      nil,
			EncoderMaker: nil,
		},
	}

}
