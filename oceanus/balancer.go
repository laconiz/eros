package oceanus

// 创建一个均衡器
func NewBalancer() *Balancer {
	return &Balancer{nodes: map[NodeKey]Node{}}
}

type Balancer struct {
	typo NodeType
	// 当前均衡器是否过期
	dirty bool
	// 当前节点列表
	nodes map[NodeKey]Node
	// 均衡列表
	// TODO 当前为随机选中, 需实现更优化的算法
	balances []Node
}

// 将均衡器状态设置未过期
func (b *Balancer) Expired() {
	b.dirty = true
}

// 插入一个节点
func (b *Balancer) Insert(node Node) {
	info := node.Info()
	b.nodes[info.Key] = node
	b.Expired()
}

// 删除一个节点
// KEY有可能重复, 删除节点时需判定节点ID是否一致
func (b *Balancer) Remove(info *NodeInfo) {
	node, ok := b.nodes[info.Key]
	if ok && node.Info().ID == info.ID {
		delete(b.nodes, info.Key)
		b.Expired()
	}
}

// 重新均衡均衡器
func (b *Balancer) rebalance() {

	b.dirty = false
	b.balances = nil

	for _, node := range b.nodes {

		mesh := node.Mesh()
		// 网格未连接
		if mesh.Connected() {
			continue
		}
		b.balances = append(b.balances, node)
	}
}

// //
// func (b *Balancer) Send(message *Message) {
//
// }
//
// // 均衡消息
// func (b *Balancer) Balance(message *Message) {
//
// }
