package local

import (
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/router"
)

func NewMesh(info *proto.Mesh, state *proto.State, router *router.Router) *Mesh {
	return &Mesh{
		info:   info,
		state:  state,
		nodes:  map[proto.NodeID]*Node{},
		types:  map[proto.NodeType]int64{},
		router: router,
	}
}

type Mesh struct {
	info   *proto.Mesh
	state  *proto.State
	nodes  map[proto.NodeID]*Node
	types  map[proto.NodeType]int64
	router *router.Router
}

func (mesh *Mesh) Info() *proto.Mesh {
	return mesh.info
}

func (mesh *Mesh) State() (*proto.State, bool) {
	return mesh.state, true
}

func (mesh *Mesh) Push(message *proto.Mail) error {
	for _, receiver := range message.Receivers {
		if node, ok := mesh.nodes[receiver.ID]; ok {
			node.Push(message)
		}
	}
	return nil
}

func (mesh *Mesh) Nodes() []*proto.Node {
	var nodes []*proto.Node
	for _, node := range mesh.nodes {
		nodes = append(nodes, node.Info())
	}
	return nodes
}

func (mesh *Mesh) UpdateState(state *proto.State) {
	mesh.state = state
	mesh.Expired()
}

func (mesh *Mesh) Expired() {
	for typo, count := range mesh.types {
		if count > 0 {
			mesh.router.Expired(typo)
		}
	}
}

func (mesh *Mesh) Insert(info *proto.Node, invoker interface{}) (*Node, error) {
	// 删除旧节点
	mesh.Remove(info.ID)
	// 新建一个节点
	node, err := newNode(info, mesh, invoker)
	if err != nil {
		return nil, err
	}
	// 写入节点列表
	mesh.nodes[info.ID] = node
	// 插入路由器
	mesh.router.Insert(node)
	// 更新节点类型统计列表
	mesh.types[info.Type]++
	// 返回数据
	return node, nil
}

// 销毁一个节点
func (mesh *Mesh) Remove(id proto.NodeID) {
	// 查询节点
	if node, ok := mesh.nodes[id]; ok {
		// 销毁节点
		node.Destroy()
		// 删除节点列表数据
		delete(mesh.nodes, id)
		// 从路由器删除节点
		mesh.router.Remove(node)
		// 更新节点类型统计列表
		mesh.types[node.Info().Type]--
	}
}

// 销毁一个网格
func (mesh *Mesh) Destroy() {
	// 遍历节点
	for _, node := range mesh.nodes {
		// 销毁节点
		node.Destroy()
		// 从路由器中删除节点
		mesh.router.Remove(node)
	}
	// 重置节点列表
	mesh.nodes = map[proto.NodeID]*Node{}
	// 重置节点类型统计列表
	mesh.types = map[proto.NodeType]int64{}
}
