package proto

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
