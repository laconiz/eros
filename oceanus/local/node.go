package local

import (
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/router"
	"github.com/laconiz/eros/queue"
)

type Node struct {
	node  *proto.Node
	mesh  *Mesh
	hub   *router.Hub
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

func newNode(info *proto.Node, mesh *Mesh, router *router.Router) *Node {
	node := &Node{
		node:  info,
		mesh:  mesh,
		queue: queue.New(64),
	}
	node.hub = router.Add(node)
	return node
}
