package remote

import (
	"errors"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/router"
)

type Mesh struct {
	// 网格信息
	mesh *proto.Mesh
	// 节点列表
	nodes map[proto.NodeID]*Node
	// 当state发生变化时 应触发对应的均衡器过期
	// 记录当前网格所拥有的节点的均衡器数量然后直接更新
	// 以避免当节点数量过多时遍历节点列表设置均衡器过期
	types map[proto.NodeType]int64
	// 路由器
	router *router.Router
	// 网络连接
	session network.Session
}

// 网格信息
func (m *Mesh) Info() *proto.Mesh {
	return m.mesh
}

// 向网格发送数据
func (m *Mesh) Push(message *proto.Message) error {
	if m.session != nil {
		return m.session.Send(message)
	}
	return errors.New("session disconnected")
}

// 更新网格数据并将网格引用的均衡器过期
func (m *Mesh) update(mesh *proto.Mesh, session network.Session) {
	m.mesh = mesh
	m.session = session
	for typo, count := range m.types {
		if count > 0 {
			m.router.Expired(typo)
		}
	}
}

// 插入一个节点并将其插入路由器
// 路由器在插拔节点时会自动过期 所以不再进行手动过期
func (m *Mesh) insertNode(info *proto.Node) {
	m.removeNode(info.ID)
	m.nodes[info.ID] = newNode(info, m, m.router)
	m.types[info.Type]++
}

// 移除一个节点并将其从路由器中移除
// 路由器在插拔节点时会自动过期 所以不再进行手动过期
func (m *Mesh) removeNode(id proto.NodeID) {
	if node, ok := m.nodes[id]; ok {
		delete(m.nodes, id)
		m.router.Remove(id)
		m.types[node.Info().Type]--
	}
}

// 销毁网格并将网格中的所有节点从路由器中移除
// 路由器在插拔节点时会自动过期 所以不再进行手动过期
func (m *Mesh) destroy() {
	for _, node := range m.nodes {
		m.router.Remove(node.Info().ID)
	}
	m.nodes = map[proto.NodeID]*Node{}
	m.types = map[proto.NodeType]int64{}
}

// 生成一个网格
func newMesh(mesh *proto.Mesh, session network.Session, router *router.Router) *Mesh {
	return &Mesh{
		mesh:    mesh,
		nodes:   map[proto.NodeID]*Node{},
		types:   map[proto.NodeType]int64{},
		router:  router,
		session: session,
	}
}
