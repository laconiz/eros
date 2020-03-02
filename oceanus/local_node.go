// 本地节点

package oceanus

import (
	"github.com/laconiz/eros/network/queue"
)

func NewLocalNode(info *NodeInfo, mesh *LocalMesh, router *Router) *LocalNode {
	node := &LocalNode{
		info:  info,
		mesh:  mesh,
		queue: queue.New(64),
	}
	node.balancer = router.Insert(node)
	return node
}

type LocalNode struct {
	info     *NodeInfo
	mesh     *LocalMesh
	balancer *Balancer
	queue    *queue.Queue
}

// 节点信息
func (n *LocalNode) Info() *NodeInfo {
	return n.info
}

// 节点所属网格
func (n *LocalNode) Mesh() Mesh {
	return n.mesh
}

// 向节点推送信息
func (n *LocalNode) Push(message *Message) error {
	return n.queue.Add(message)
}

// 销毁节点
func (n *LocalNode) Destroy() {

}
