package remote

import (
	"github.com/laconiz/eros/oceanus"
)

type Node struct {
	// 节点信息
	node *oceanus.NodeInfo
	// 网格
	mesh *Mesh
	// 均衡器
	hub *oceanus.Balancer
}

// 节点信息
func (n *Node) Info() *oceanus.NodeInfo {
	return n.node
}

// 向节点发送数据
func (n *Node) Push(message *oceanus.Message) error {
	return n.mesh.Push(message)
}

// 节点所属网格
func (n *Node) Mesh() oceanus.MeshInfo {
	return n.mesh
}

func newNode(info *oceanus.NodeInfo, mesh *Mesh, router *oceanus.Router) *Node {
	node := &Node{node: info, mesh: mesh}
	node.hub = router.Add(node)
	return node
}
