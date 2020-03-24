package proto

import "github.com/laconiz/eros/network/message"

type NodeID string
type NodeType string
type NodeKey string

type Node struct {
	ID   NodeID   `json:"id"`
	Type NodeType `json:"type"`
	Key  NodeKey  `json:"key"`
}

type MeshID string

type Mesh struct {
	ID   MeshID `json:"id"`
	Addr string `json:"addr"`
}

type MeshJoin struct {
	Mesh  *Mesh  `json:"mesh"`
	State *State `json:"state"`
}

type MeshQuit struct {
	Mesh *Mesh `json:"mesh"`
}

type NodeJoin struct {
	Nodes []*Node `json:"nodes"`
}

type NodeQuit struct {
	Nodes []*Node `json:"nodes"`
}

func init() {
	message.Register(MeshJoin{}, message.JsonCodec())
	message.Register(MeshQuit{}, message.JsonCodec())
	message.Register(NodeJoin{}, message.JsonCodec())
	message.Register(NodeQuit{}, message.JsonCodec())
}
