package local

import (
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/router"
	"github.com/laconiz/eros/queue"
)

type Node struct {
	node  *proto.Node
	mesh  *Mesh
	queue *queue.Queue
}

func (n *Node) Info() *proto.Node {
	return n.node
}

func (n *Node) Mesh() router.Mesh {
	return n.mesh
}

func (n *Node) Push(message *proto.Message) error {
	return n.queue.Add(message)
}

func (n *Node) Close() {
	n.queue.Close()
}
