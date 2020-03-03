package remote

import (
	"errors"
	"github.com/laconiz/eros/network/session"
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/oceanus/proto"
)

func NewMesh(info *proto.Mesh, state *proto.State, router *oceanus.Router) *Mesh {
	return &Mesh{
		info:    info,
		state:   state,
		nodes:   map[proto.NodeID]*Node{},
		types:   map[proto.NodeType]int64{},
		router:  router,
		session: nil,
	}
}

type Mesh struct {
	info    *proto.Mesh              // 网格信息
	state   *proto.State             // 网格状态
	nodes   map[proto.NodeID]*Node   // 节点列表
	types   map[proto.NodeType]int64 // 网格节点类型统计
	router  *oceanus.Router          // 路由器
	session session.Session          // 网络连接
}

// 网格信息
func (mesh *Mesh) Info() *proto.Mesh {
	return mesh.info
}

// 网格状态
func (mesh *Mesh) State() (*proto.State, bool) {
	return mesh.state, mesh.session != nil
}

// 向网格发送数据
func (mesh *Mesh) Push(message *proto.Mail) error {
	if mesh.session != nil {
		return mesh.session.Send(message)
	}
	return errors.New("invalid session")
}

// 更新网格状态
func (mesh *Mesh) UpdateState(state *proto.State) {
	mesh.state = state
	mesh.Expired()
}

// 更新网格连接
func (mesh *Mesh) UpdateSession(session session.Session) {
	mesh.session = session
	mesh.Expired()
}

// 设置网格过期
func (mesh *Mesh) Expired() {
	for typo, count := range mesh.types {
		if count > 0 {
			mesh.router.Expired(typo)
		}
	}
}

// 插入一个节点
func (mesh *Mesh) Insert(list []*proto.Node) {
	mesh.Remove(list)
	for _, info := range list {
		node := newNode(info, mesh)
		mesh.nodes[info.ID] = node
		mesh.router.Insert(node)
		mesh.types[info.Type]++
	}
}

// 销毁一个节点
func (mesh *Mesh) Remove(list []*proto.Node) {
	for _, info := range list {
		if node, ok := mesh.nodes[info.ID]; ok {
			node.Destroy()
			delete(mesh.nodes, info.ID)
			mesh.router.Remove(node)
			mesh.types[info.Type]--
		}
	}
}

// 销毁一个网格
func (mesh *Mesh) Destroy() {
	for _, node := range mesh.nodes {
		node.Destroy()
		mesh.router.Remove(node)
	}
	mesh.nodes = map[proto.NodeID]*Node{}
	mesh.types = map[proto.NodeType]int64{}
}
