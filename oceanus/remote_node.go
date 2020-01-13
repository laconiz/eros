// 远程节点

package oceanus

type RemoteNode struct {
	// 节点信息
	info *NodeInfo
	// 网格
	mesh *RemoteMesh
	// 均衡器
	balancer *Balancer
}

// 节点信息
func (n *RemoteNode) Info() *NodeInfo {
	return n.info
}

// 向节点发送数据
func (n *RemoteNode) Push(message *Message) error {
	return n.mesh.Push(message)
}

// 节点所属网格
func (n *RemoteNode) Mesh() Mesh {
	return n.mesh
}

func NewRemoteNode(info *NodeInfo, mesh *RemoteMesh, router *Router) *RemoteNode {
	node := &RemoteNode{info: info, mesh: mesh}
	node.balancer = router.Insert(node)
	return node
}
