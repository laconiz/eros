package local

import (
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/router"
)

type Mesh struct {
	// 网格信息
	mesh *proto.Mesh
	// 节点列表
	nodes map[proto.NodeID]*Node
	// 记录当前网格所拥有的节点的均衡器数量然后直接更新
	// 以避免当节点数量过多时遍历节点列表设置均衡器过期
	types map[proto.NodeType]int64
	// 路由器
	router *router.Router
}

func (m *Mesh) Info() *proto.Mesh {
	return m.mesh
}

func (m *Mesh) Push(message *proto.Message) error {
	for _, receiver := range message.Receivers {
		if node, ok := m.nodes[receiver.ID]; ok {
			node.Push(message)
		}
	}
	return nil
}

func (m *Mesh) update(info *proto.Mesh) {
	m.mesh = info
	for typo, count := range m.types {
		if count > 0 {
			m.router.Expired(typo)
		}
	}
}

func (m *Mesh) Nodes() []*proto.Node {
	var nodes []*proto.Node
	for _, node := range m.nodes {
		nodes = append(nodes, node.Info())
	}
	return nodes
}

func (m *Mesh) InsertNode(info *proto.Node) *Node {
	m.RemoveNode(info.ID)
	node := newNode(info, m, m.router)
	m.nodes[info.ID] = node
	m.types[info.Type]++
	return node
}

func (m *Mesh) RemoveNode(id proto.NodeID) *Node {
	if node, ok := m.nodes[id]; ok {
		delete(m.nodes, id)
		m.router.Remove(id)
		m.types[node.Info().Type]--
		return node
	}
	return nil
}
