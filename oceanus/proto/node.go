package proto

import "github.com/laconiz/eros/network/message"

type NodeID string
type NodeType string
type NodeKey string

const (
	EmptyNodeType NodeType = ""
)

type Node struct {
	ID   NodeID   `json:"id"`
	Type NodeType `json:"type"`
	Key  NodeKey  `json:"key"`
}

type NodeJoin struct {
	Nodes []*Node `json:"nodes"`
	State *State  `json:"state"`
}

type NodeQuit struct {
	Nodes []*Node `json:"nodes"`
	State *State  `json:"state"`
}

func init() {
	message.Register(NodeJoin{}, message.JsonCodec())
	message.Register(NodeQuit{}, message.JsonCodec())
}
