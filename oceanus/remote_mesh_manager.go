// 远程网格管理器

package oceanus

import (
	"github.com/laconiz/eros/network"
)

type RemoteMeshMgr struct {
	// 路由器
	router *Router
	// 网格列表
	meshes map[MeshID]*RemoteMesh
}

// 插入或更新一个网格
func (n *RemoteMeshMgr) InsertMesh(info *MeshInfo, session network.Session) {
	mesh, ok := n.meshes[info.ID]
	if !ok {
		// 网格不存在, 生成网格
		mesh = NewRemoteMesh(info, session, n.router)
		n.meshes[info.ID] = mesh
		logger.Infof("mesh join: %+v", info)
	} else {
		// 网格存在, 更新网格信息
		mesh.Update(info, session)
		logger.Infof("mesh update: %+v", info)
	}
}

// 删除一个网格
func (n *RemoteMeshMgr) RemoveMesh(id MeshID) {
	if mesh, ok := n.meshes[id]; ok {
		mesh.Destroy()
		delete(n.meshes, id)
		logger.Infof("mesh quit: %+v", mesh)
	}
}

// 插入一个节点
func (n *RemoteMeshMgr) InsertNode(info *NodeInfo) {
	if mesh, ok := n.meshes[info.Mesh]; ok {
		mesh.InsertNode(info)
	}
}

// 移除一个节点
func (n *RemoteMeshMgr) RemoveNode(info *NodeInfo) {
	if mesh, ok := n.meshes[info.Mesh]; ok {
		mesh.RemoveNode(info.ID)
	}
}

// 向所有节点广播消息
func (n *RemoteMeshMgr) Broadcast(msg interface{}) {
	for _, mesh := range n.meshes {
		if mesh.session != nil {
			mesh.session.Send(msg)
		}
	}
}

// 生成一个网格管理器
func NewRemoteMeshMgr(router *Router) *RemoteMeshMgr {
	return &RemoteMeshMgr{
		router: router,
		meshes: map[MeshID]*RemoteMesh{},
	}
}
