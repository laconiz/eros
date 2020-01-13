package remote

import (
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/oceanus"
)

type Net struct {
	router *oceanus.Router
	meshes map[oceanus.MeshID]*Mesh
}

// 插入或更新一个网格
func (n *Net) AddMesh(info *oceanus.MeshInfo, session network.Session) {
	mesh, ok := n.meshes[info.ID]
	if !ok {
		mesh = newMesh(info, session, n.router)
		n.meshes[info.ID] = mesh
	} else {
		mesh.update(info, session)
	}
}

// 删除一个网格
func (n *Net) RemoveMesh(id oceanus.MeshID) {
	if mesh, ok := n.meshes[id]; ok {
		mesh.destroy()
		delete(n.meshes, id)
	}
}

// 插入一个节点
func (n *Net) InsertNode(info *oceanus.NodeInfo) {
	if mesh, ok := n.meshes[info.Mesh]; ok {
		mesh.insertNode(info)
	}
}

// 移除一个节点
func (n *Net) RemoveNode(info *oceanus.NodeInfo) {
	if mesh, ok := n.meshes[info.Mesh]; ok {
		mesh.removeNode(info.ID)
	}
}

// 向所有节点广播消息
func (n *Net) Broadcast(msg interface{}) {
	for _, mesh := range n.meshes {
		if mesh.session != nil {
			mesh.session.Send(msg)
		}
	}
}

// 生成一个网格管理器
func NewNet(router *oceanus.Router) *Net {
	return &Net{
		router: router,
		meshes: map[oceanus.MeshID]*Mesh{},
	}
}
