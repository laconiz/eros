package remote

import (
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/router"
)

type Node struct {
	// 节点信息
	node *proto.Node
	// 网格
	mesh *Mesh
	// 均衡器
	hub *router.Hub
}

// 节点信息
func (n *Node) Info() *proto.Node {
	return n.node
}

// 向节点发送数据
func (n *Node) Push(message *proto.Message) error {
	return n.mesh.Push(message)
}

// 节点所属网格
func (n *Node) Mesh() router.Mesh {
	return n.mesh
}
