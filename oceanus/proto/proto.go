package proto

import "github.com/laconiz/eros/network/message"

type NodeID string
type NodeType string
type NodeKey string

type Node struct {
	ID   NodeID   `json:"id"`   // 节点ID
	Type NodeType `json:"type"` // 节点类型
	Key  NodeKey  `json:"key"`  // 节点KEY
	Mesh MeshID   `json:"mesh"` // 网格ID
}

type MeshID string

type Mesh struct {
	ID    MeshID `json:"id"`    // 网格ID
	Addr  string `json:"addr"`  // 网格地址
	Power uint64 `json:"power"` // 网格权值
}

type State struct {
	Version uint32 `json:"ver"`
}

type MeshJoin struct {
	Mesh  *Mesh   `json:"mesh"`
	State *State  `json:"state"`
	Nodes []*Node `json:"nodes, omitempty"`
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
	message.Register(State{}, message.JsonCodec())
	message.Register(MeshJoin{}, message.JsonCodec())
	message.Register(MeshQuit{}, message.JsonCodec())
	message.Register(NodeJoin{}, message.JsonCodec())
	message.Register(NodeQuit{}, message.JsonCodec())
}
