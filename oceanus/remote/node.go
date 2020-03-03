// 远程节点

package remote

import (
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/oceanus/proto"
)

func newNode(info *proto.Node, mesh *Mesh) *Node {
	return &Node{info: info, mesh: mesh}
}

type Node struct {
	info *proto.Node // 节点信息
	mesh *Mesh       // 所属网格
}

// 节点信息
func (node *Node) Info() *proto.Node {
	return node.info
}

// 节点所属网格
func (node *Node) Mesh() oceanus.Mesh {
	return node.mesh
}

// 向节点发送数据
func (node *Node) Push(message *proto.Mail) error {
	return node.mesh.Push(message)
}

// 销毁节点
func (node *Node) Destroy() {

}
