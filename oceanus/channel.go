package oceanus

import "github.com/laconiz/eros/oceanus/proto"

type Channel interface {
	Info() *proto.Channel
	Node() Node
	Local() bool
	Push(*proto.Message) error
}

type access struct {
	channel *proto.Channel
	node    Node
}

func (a *access) Info() *proto.Channel {
	return a.channel
}

func (a *access) Node() Node {
	return a.node
}

func (a *access) Local() bool {
	return false
}

func (a *access) Push(msg *proto.Message) error {
	return a.node.Send(msg)
}

func newAccess(info *proto.Channel, node Node) Channel {
	return &access{channel: info, node: node}
}
