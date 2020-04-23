package oceanus

import "github.com/laconiz/eros/network/message"

// ---------------------------------------------------------------------------------------------------------------------

type MeshID string

type MeshInfo struct {
	ID   MeshID `json:"id"`
	Addr string `json:"addr"`
}

// ---------------------------------------------------------------------------------------------------------------------

type NodeID string
type NodeType string
type NodeKey string

type NodeInfo struct {
	ID   NodeID   `json:"id"`
	Type NodeType `json:"type"`
	Key  NodeKey  `json:"key"`
}

// ---------------------------------------------------------------------------------------------------------------------

type Version int64

type MeshState struct {
	Version Version `json:"version"`
	Power   int64   `json:"power"`
	Limit   int64   `json:"limit"`
}

// ---------------------------------------------------------------------------------------------------------------------

type MailID string

type RpcID string

const emptyRpcID = ""

type Mail struct {
	ID     MailID                 `json:"id"`
	Header map[string]interface{} `json:"header"`
	From   []*NodeInfo            `json:"from,omitempty"`
	Type   NodeType               `json:"type,omitempty"`
	To     []*NodeInfo            `json:"to,omitempty"`
	Reply  MailID                 `json:"reply,omitempty"`
	User   int64                  `json:"user,omitempty"`
	Body   []byte                 `json:"body"`
}

// ---------------------------------------------------------------------------------------------------------------------

type MeshJoin struct {
	Mesh *MeshInfo `json:"mesh"`
}

type MeshQuit struct {
}

type NodeJoin struct {
	Nodes []*NodeInfo `json:"nodes"`
	State MeshState   `json:"state"`
}

type NodeQuit struct {
	Nodes []NodeID  `json:"nodes"`
	State MeshState `json:"state"`
}

// ---------------------------------------------------------------------------------------------------------------------

func init() {
	message.Register(Mail{}, message.JsonCodec())
	message.Register(MeshJoin{}, message.JsonCodec())
	message.Register(MeshQuit{}, message.JsonCodec())
	message.Register(NodeJoin{}, message.JsonCodec())
	message.Register(NodeQuit{}, message.JsonCodec())
}
