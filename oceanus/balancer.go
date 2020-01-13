package oceanus

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

func (b *Balancer) Type() NodeType {
	return b.typo
}

func (b *Balancer) expired() {
	b.dirty = true
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

//
func (b *Balancer) Send(message *Message) {

}

// 均衡消息
func (b *Balancer) Balance(message *Message) {

}
