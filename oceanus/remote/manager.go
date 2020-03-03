package remote

import (
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/oceanus/proto"
)

func NewManager() *Manager {

}

type Manager struct {
	meshes map[proto.MeshID]*Mesh // 网格列表
	router *oceanus.Router        // 路由器
}

func (manager *Manager) InsertMesh() {

}

func (manager *Manager) RemoveMesh() {

}

func (manager *Manager) InsertNode() {

}

func (manager *Manager) RemoveNode() {

}
