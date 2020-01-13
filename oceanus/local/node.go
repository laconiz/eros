package local

import (
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/queue"
)

type Node struct {
	node  *oceanus.NodeInfo
	mesh  *Mesh
	hub   *oceanus.Balancer
	queue *queue.Queue
}

func (n *Node) Info() *oceanus.NodeInfo {
	return n.node
}

func (n *Node) Mesh() oceanus.MeshInfo {
	return n.mesh
}

func (n *Node) Push(message *oceanus.Message) error {
	return n.queue.Add(message)
}

func (n *Node) Close() {
	n.queue.Close()
}

func newNode(info *oceanus.NodeInfo, mesh *Mesh, router *oceanus.Router) *Node {
	node := &Node{
		node:  info,
		mesh:  mesh,
		queue: queue.New(64),
	}
	node.hub = router.Add(node)
	return node
}
