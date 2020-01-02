package channel

import (
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/oceanus/node"
	"github.com/laconiz/eros/oceanus/proto"
)

type remote struct {
	info *Info
	node node.Node
}

func (r *remote) Info() *Info {
	return r.info
}

func (r *remote) Node() node.Node {
	return r.node
}

func (r *remote) Push(message *proto.Message) error {
	return r.node.Send(message)
}

func NewRemote(info *Info, node node.Node) Channel {
	return &remote{info: info, node: node}
}
