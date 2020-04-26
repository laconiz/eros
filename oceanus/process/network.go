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

func (proc *Process) NewInvoker() invoker.Invoker {

	invoker := invoker.NewSocketInvoker()

	invoker.Register(network.Connected{}, proc.OnConnected)
	invoker.Register(network.Disconnected{}, proc.OnDisconnected)
	invoker.Register(proto.Mail{}, proc.OnMail)
	invoker.Register(proto.State{}, proc.OnState)

	invoker.Register(proto.MeshJoin{}, proc.OnMeshJoin)
	invoker.Register(proto.MeshQuit{}, proc.OnMeshQuit)
	invoker.Register(proto.NodeJoin{}, proc.onNodeJoin)
	invoker.Register(proto.NodeQuit{}, proc.onNodeQuit)

	return invoker
}

func (proc *Process) NewSessionOption() socket.SessionOption {
	return socket.SessionOption{
		Timeout: time.Second * 11,
		Queue:   64,
		Invoker: proc.NewInvoker(),
		Encoder: encoder.NewNameMaker(),
		Cipher:  cipher.NewEmptyMaker(),
		Reader:  reader.NewSizeMaker(),
	}
}

func (proc *Process) NewAcceptor(addr string) network.Acceptor {

	option := &socket.AcceptorOption{
		Name:    "oceanus.acceptor",
		Addr:    addr,
		Session: proc.NewSessionOption(),
		Level:   logis.WARN,
	}

	return socket.NewAcceptor(option)
}

func (proc *Process) NewConnector(addr string) network.Connector {

	option := socket.ConnectorOption{
		Name:      "oceanus.connector",
		Addr:      addr,
		Reconnect: true,
		Delays:    []time.Duration{time.Millisecond},
		Session:   proc.NewSessionOption(),
		Level:     logis.WARN,
	}

	connector := socket.NewConnector(option)
	return connector
}

func (proc *Process) broadcast(msg interface{}) {

	for _, mesh := range proc.remotes {
		if err := mesh.Send(msg); err != nil {
			proc.logger.Data(msg).Warn("send message failed")
		}
	}
}
