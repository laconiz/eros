package node

import (
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/tcp"
	"github.com/laconiz/eros/oceanus/config"
	"github.com/laconiz/eros/oceanus/proto"
)

type remote struct {
	info *Info
	conn network.Connector
}

func (r *remote) Info() *Info {
	return r.info
}

func (r *remote) Send(message *proto.Message) error {
	r.conn.Send(message)
	return nil
}

func (r *remote) Stop() {
	r.conn.Stop()
}

func NewRemote(info *Info) Node {

	conf := tcp.ConnectorConfig{
		Name:      "oceanus",
		Addr:      info.Addr,
		Reconnect: true,
		Session:   config.Session,
	}

	conn := tcp.NewConnector(conf)
	go conn.Run()

	return &remote{info: info, conn: conn}
}
