package process

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/cipher"
	"github.com/laconiz/eros/network/encoder"
	"github.com/laconiz/eros/network/invoker"
	"github.com/laconiz/eros/network/socket"
	"github.com/laconiz/eros/network/socket/reader"
	"github.com/laconiz/eros/oceanus/proto"
	"time"
)

func (process *Process) NewInvoker() invoker.Invoker {

	invoker := invoker.NewSocketInvoker()

	invoker.Register(network.Connected{}, process.OnConnected)
	invoker.Register(network.Disconnected{}, process.OnDisconnected)
	invoker.Register(proto.Mail{}, process.OnMail)
	invoker.Register(proto.State{}, process.OnState)

	invoker.Register(proto.MeshJoin{}, process.OnMeshJoin)
	invoker.Register(proto.MeshQuit{}, process.OnMeshQuit)
	invoker.Register(proto.NodeJoin{}, process.onNodeJoin)
	invoker.Register(proto.NodeQuit{}, process.onNodeQuit)

	return invoker
}

func (process *Process) NewSessionOption() socket.SessionOption {
	return socket.SessionOption{
		Timeout: time.Second * 11,
		Queue:   64,
		Invoker: process.NewInvoker(),
		Encoder: encoder.NewNameMaker(),
		Cipher:  cipher.NewEmptyMaker(),
		Reader:  reader.NewSizeMaker(),
	}
}

func (process *Process) NewAcceptor(addr string) network.Acceptor {

	option := socket.AcceptorOption{
		Name:    "oceanus.acceptor",
		Addr:    addr,
		Session: process.NewSessionOption(),
		Level:   logis.WARN,
	}

	return socket.NewAcceptor(option)
}

func (process *Process) NewConnector(addr string) network.Connector {

	option := socket.ConnectorOption{
		Name:      "oceanus.connector",
		Addr:      addr,
		Reconnect: true,
		Delays:    []time.Duration{time.Millisecond},
		Session:   process.NewSessionOption(),
		Level:     logis.WARN,
	}

	connector := socket.NewConnector(option)
	return connector
}

func (process *Process) broadcast(msg interface{}) {

	for _, mesh := range process.remotes {
		if err := mesh.Send(msg); err != nil {
			process.logger.Data(msg).Warn("send message failed")
		}
	}
}
