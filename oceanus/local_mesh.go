// 本地网格

package oceanus

// 生成一个本地网格
func NewLocalMesh(info *MeshInfo, router *Router) *LocalMesh {
	return &LocalMesh{
		info:   info,
		nodes:  map[NodeID]*LocalNode{},
		types:  map[NodeType]int64{},
		router: router,
	}
}

type LocalMesh struct {
	// 网格信息
	info *MeshInfo
	// 节点列表
	nodes map[NodeID]*LocalNode
	// 记录当前网格所拥有的节点的均衡器数量然后直接更新
	// 以避免当节点数量过多时遍历节点列表设置均衡器过期
	types map[NodeType]int64
	// 路由器
	router *Router
}

// 网格信息
func (m *LocalMesh) Info() *MeshInfo {
	return m.info
}

// 推送节点消息
func (m *LocalMesh) Push(message *Message) error {
	for _, receiver := range message.Receivers {
		if node, ok := m.nodes[receiver.ID]; ok {
			node.Push(message)
		}
	}
	return nil
}

// 本地节点永远在线
func (m *LocalMesh) Connected() bool {
	return true
}

// 更新网格信息
func (m *LocalMesh) Update(info *MeshInfo) {
	m.info = info
	for typo, count := range m.types {
		if count > 0 {
			m.router.Expired(typo)
		}
	}
}

// 获取网格拥有的节点信息列表
func (m *LocalMesh) Nodes() []*NodeInfo {
	var nodes []*NodeInfo
	for _, node := range m.nodes {
		nodes = append(nodes, node.Info())
	}
	return nodes
}

func (m *LocalMesh) InsertNode(info *NodeInfo) *LocalNode {
	m.RemoveNode(info.ID)
	node := NewLocalNode(info, m, m.router)
	m.nodes[info.ID] = node
	m.types[info.Type]++
	return node
}

func (m *LocalMesh) RemoveNode(id NodeID) *LocalNode {
	if node, ok := m.nodes[id]; ok {
		delete(m.nodes, id)
		m.router.Remove(id)
		m.types[node.Info().Type]--
		return node
	}
	return nil
}

// 销毁网格
func (m *LocalMesh) Destroy() {
	for _, node := range m.nodes {
		node.Destroy()
	}
}
