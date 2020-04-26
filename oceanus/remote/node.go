// 远程节点

package remote

import (
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/router"
)

// 创建远程节点
func newNode(info *proto.Node, mesh *Mesh) *Node {
	return &Node{info: info, mesh: mesh}
}

// 远程节点
type Node struct {
	info *proto.Node // 节点信息
	mesh *Mesh       // 节点所属网格
}

// 节点信息
func (node *Node) Info() *proto.Node {
	return node.info
}

// 节点所属网格
func (node *Node) Mesh() router.Mesh {
	return node.mesh
}

// 发送邮件
func (node *Node) Mail(mail *proto.Mail) error {
	return node.mesh.Mail(mail)
}

// 销毁节点
func (node *Node) Destroy() {

}
