package process

import (
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/remote"
)

// 网络连接时发送网格数据
func (process *Process) OnConnected(event *network.Event) {
	process.mutex.RLock()
	defer process.mutex.RUnlock()
	mesh := process.local
	state, _ := mesh.State()
	event.Ses.Send(&proto.MeshJoin{Mesh: mesh.Info(), State: state, Nodes: mesh.Nodes()})
}

// 网络断开时更新网格状态
func (process *Process) OnDisconnected(event *network.Event) {
	if mesh, ok := process.boundMesh(event); ok {
		process.mutex.Lock()
		defer process.mutex.Unlock()
		mesh.UpdateSession(nil)
	}
}

// 处理消息
func (process *Process) OnMail(event *network.Event) {
	process.mutex.RLock()
	defer process.mutex.RUnlock()
	if err := process.local.Push(event.Msg.(*proto.Mail)); err != nil {
		process.log.Warnf("recv mail error: %v", err)
	}
}

// 插入网格
func (process *Process) OnMeshJoin(event *network.Event) {
	msg := event.Msg.(*proto.MeshJoin)
	event.Ses.Set(sessionKey, msg.Mesh)
	process.mutex.Lock()
	defer process.mutex.Unlock()
	mesh, ok := process.remotes[msg.Mesh.ID]
	if !ok {
		mesh = remote.NewMesh(msg.Mesh, msg.State, process.router)
		process.remotes[msg.Mesh.ID] = mesh
	}
	mesh.UpdateSession(event.Ses)
	mesh.Insert(msg.Nodes)
}

// 移除网格
func (process *Process) OnMeshQuit(event *network.Event) {
	if mesh, ok := process.boundMesh(event); ok {
		process.mutex.Lock()
		defer process.mutex.Unlock()
		mesh.Destroy()
		delete(process.remotes, mesh.Info().ID)
	}
}

// 插入节点
func (process *Process) OnNodeJoin(event *network.Event) {
	if mesh, ok := process.boundMesh(event); ok {
		process.mutex.Lock()
		defer process.mutex.Unlock()
		mesh.Insert(event.Msg.(*proto.NodeJoin).Nodes)
	}
}

// 移除节点
func (process *Process) OnNodeQuit(event *network.Event) {
	if mesh, ok := process.boundMesh(event); ok {
		process.mutex.Lock()
		defer process.mutex.Unlock()
		mesh.Remove(event.Msg.(*proto.NodeQuit).Nodes)
	}
}

// 获取连接绑定
func (process *Process) boundMesh(event *network.Event) (*remote.Mesh, bool) {
	if value, ok := event.Ses.Get(sessionKey); ok {
		mesh, ok := process.remotes[value.(*proto.Mesh).ID]
		return mesh, ok
	}
	return nil, false
}

const sessionKey = "mesh"
