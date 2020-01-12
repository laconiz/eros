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

func (m *Mesh) update() {

}

func (m *Mesh) insertNode() {

}

func (m *Mesh) removeNode() {

}
