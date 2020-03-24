package local

import (
	"github.com/laconiz/eros/network/queue"
	"github.com/laconiz/eros/oceanus/proto"
)

func newNode(info *proto.Node, mesh *Mesh, invoker interface{}) (*Node, error) {
	return &Node{info: info, mesh: mesh, queue: queue.New(queueLen)}, nil
}

type Node struct {
	info  *proto.Node  // 节点信息
	mesh  Mesh         // 所属网格
	queue *queue.Queue // 消息队列
}

// 节点信息
func (node *Node) Info() *proto.Node {
	return node.info
}

// 节点所属网格
func (node *Node) Mesh() Mesh {
	return node.mesh
}

// 向节点发送数据
func (node *Node) Push(message *proto.Mail) error {
	return node.queue.Add(message)
}

// 销毁节点
func (node *Node) Destroy() {

}

const queueLen = 64
