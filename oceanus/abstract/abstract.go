package abstract

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/oceanus/proto"
)

type Mesh interface {
	Info() *proto.Mesh
	State() (*proto.State, bool)
	Mail(*proto.Mail) error
}

type Node interface {
	Info() *proto.Node
	Mesh() Mesh
	Mail(*proto.Mail) error
}

type Router interface {
	Insert(Node)
	Remove(proto.NodeID)
	Expired(proto.NodeType)
	ByID(proto.NodeID) Node
	ByIDList([]proto.NodeID) []Node
	ByKey(proto.NodeType, proto.NodeKey) Node
	ByKeys(proto.NodeType, []proto.NodeKey) []Node
	ByLoad(proto.NodeType) Node
	ByType(proto.NodeType) []Node
}

type Process interface {
	Local() Mesh
	Router() Router
	Logger() logis.Logger
	NewConnector(string) network.Connector
}
