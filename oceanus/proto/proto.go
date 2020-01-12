package proto

import (
	"github.com/laconiz/eros/codec"
	"github.com/laconiz/eros/network"
)

type MeshID string

type State struct {
	Version uint64
}

type Mesh struct {
	ID    MeshID
	Addr  string
	State State
}

type NodeID string

type NodeType string

type NodeKey string

type Node struct {
	ID   NodeID
	Type NodeType
	Key  NodeKey
	Mesh MeshID
}

type MsgID string

type MsgType uint32

type Message struct {
	ID        MsgID
	Sender    []Node
	Receivers []Node
	Type      MsgType
	Body      []byte
}

type MeshJoin struct {
	*Mesh
}

type MeshQuit struct {
	*Mesh
}

type NodeJoin []*Node

type NodeQuit []*Node

func init() {
	network.RegisterMeta(Message{}, codec.Json())
	network.RegisterMeta(MeshJoin{}, codec.Json())
	network.RegisterMeta(MeshQuit{}, codec.Json())
	network.RegisterMeta(NodeJoin{}, codec.Json())
	network.RegisterMeta(NodeQuit{}, codec.Json())
}
