package process

import (
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/cipher"
	"github.com/laconiz/eros/network/encoder"
	"github.com/laconiz/eros/network/invoker"
	"github.com/laconiz/eros/network/socket"
	"github.com/laconiz/eros/network/socket/packer"
	"github.com/laconiz/eros/oceanus/proto"
	"time"
)

func (process *Process) NewInvoker() invoker.Invoker {
	invoker := invoker.NewSocketInvoker()
	invoker.Register(network.Connected{}, process.OnConnected)
	invoker.Register(network.Disconnected{}, process.OnDisconnected)
	invoker.Register(proto.Mail{}, process.OnMail)
	invoker.Register(proto.MeshJoin{}, process.OnMeshJoin)
	invoker.Register(proto.MeshJoin{}, process.OnMeshQuit)
	invoker.Register(proto.MeshJoin{}, process.OnNodeJoin)
	invoker.Register(proto.MeshJoin{}, process.OnNodeQuit)
	return invoker
}

func (process *Process) NewSessionOption() socket.SesOption {
	return socket.SesOption{
		Timeout:  time.Second * 11,
		QueueLen: 64,
		Invoker:  nil,
		Encoder:  encoder.NewIDMaker(),
		Cipher:   cipher.NewEmptyMaker(),
		Packer:   packer.NewSizeMaker(),
	}
}

func (process *Process) NewAcceptor(addr string) network.Acceptor {
	option := socket.AccOption{
		Name:    "oceanus.acceptor",
		Addr:    addr,
		Session: process.NewSessionOption(),
	}
	return socket.NewAcceptor(option)
}

func (process *Process) NewConnector(addr string) network.Connector {
	option := socket.ConnOption{
		Name:      "oceanus.connector",
		Addr:      addr,
		Reconnect: true,
		Delays:    []time.Duration{time.Millisecond},
		Session:   process.NewSessionOption(),
	}
	connector := socket.NewConnector(option)
	go connector.Run()
	return connector
}
