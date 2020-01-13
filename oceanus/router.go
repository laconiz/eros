package oceanus

type Mesh interface {
	Info() *MeshInfo
	Push(*Message) error
	Connected() bool
}

type Node interface {
	Info() *NodeInfo
	Mesh() Mesh
	Push(*Message) error
}

type Router struct {
	// 所有节点
	nodes map[NodeID]Node
	// 均衡器
	balancers map[NodeType]*Balancer
}

// 添加一个节点
func (r *Router) Insert(node Node) *Balancer {
	// 设置精确查找字典
	info := node.Info()
	r.nodes[info.ID] = node
	// 获取均衡器
	balancer, ok := r.balancers[info.Type]
	if !ok {
		balancer = NewBalancer()
		r.balancers[info.Type] = balancer
	}
	// 插入均衡器
	balancer.Insert(node)
	return balancer
}

// 删除一个节点
func (r *Router) Remove(id NodeID) {
	if node, ok := r.nodes[id]; ok {
		delete(r.nodes, id)
		info := node.Info()
		// 从均衡器中删除
		if balancer, ok := r.balancers[info.Type]; ok {
			balancer.Remove(info)
		}
	}
}

// 将指定类型的均衡器设置为过期状态
func (r *Router) Expired(typo NodeType) {
	if balancer, ok := r.balancers[typo]; ok {
		balancer.Expired()
	}
}

func NewRouter() *Router {
	return &Router{
		nodes: map[NodeID]Node{},
	}
}
