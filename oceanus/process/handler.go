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
	event.Ses.Send(&proto.MeshJoin{
		Mesh:  mesh.Info(),
		State: state,
		Nodes: mesh.Nodes(),
	})

	process.logger.Info("join to remote")
}

// 网络断开时更新网格状态
func (process *Process) OnDisconnected(event *network.Event) {

	value, ok := event.Ses.Get(sessionKey)
	if !ok {
		return
	}
	info := value.(*proto.Mesh)

	process.mutex.Lock()
	defer process.mutex.Unlock()

	if mesh, ok := process.remotes[info.ID]; ok {
		mesh.UpdateSession(nil)
	}
}

// 处理消息
func (process *Process) OnMail(event *network.Event) {
	process.mutex.RLock()
	defer process.mutex.RUnlock()
	if err := process.local.Push(event.Msg.(*proto.Mail)); err != nil {
		process.logger.Warnf("recv mail error: %v", err)
	}
}

// 网格状态
func (process *Process) OnState(event *network.Event) {

	value, ok := event.Ses.Get(sessionKey)
	if !ok {
		return
	}
	info := value.(*proto.Mesh)

	msg := event.Msg.(*proto.State)

	process.mutex.Lock()
	defer process.mutex.Unlock()

	if mesh, ok := process.remotes[info.ID]; ok {
		mesh.UpdateState(msg)
		data := &proto.MeshJoin{Mesh: info, State: msg}
		process.logger.Data(data).Info("remote mesh update")
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

	process.logger.Data(msg).Info("remote mesh join")
}

// 移除网格
func (process *Process) OnMeshQuit(event *network.Event) {

	value, ok := event.Ses.Get(sessionKey)
	if !ok {
		return
	}
	info := value.(*proto.Mesh)

	process.mutex.Lock()
	defer process.mutex.Unlock()

	if mesh, ok := process.remotes[info.ID]; ok {
		mesh.Destroy()
		delete(process.remotes, info.ID)
		process.logger.Data(info).Info("remote mesh quit")
	}

	if connector, ok := process.connectors[info.ID]; ok {
		connector.Stop()
		delete(process.connectors, info.ID)
		process.logger.Data(info.ID).Info("connector stopped")
	}
}

// 插入节点
func (process *Process) onNodeJoin(event *network.Event) {

	value, ok := event.Ses.Get(sessionKey)
	if !ok {
		return
	}
	info := value.(*proto.Mesh)

	msg := event.Msg.(*proto.NodeJoin)

	process.mutex.Lock()
	defer process.mutex.Unlock()

	if mesh, ok := process.remotes[info.ID]; ok {
		mesh.Insert(msg.Nodes)
		process.logger.Data(msg).Info("remote node join")
	}
}

// 移除节点
func (process *Process) onNodeQuit(event *network.Event) {

	value, ok := event.Ses.Get(sessionKey)
	if !ok {
		return
	}
	info := value.(*proto.Mesh)

	msg := event.Msg.(*proto.NodeQuit)

	process.mutex.Lock()
	defer process.mutex.Unlock()

	if mesh, ok := process.remotes[info.ID]; ok {
		mesh.Remove(msg.Nodes)
		process.logger.Data(msg).Info("remote node quit")
	}
}

// 广播状态
func (process *Process) broadcastState() {

	process.mutex.RLock()
	defer process.mutex.RUnlock()

	state, _ := process.local.State()

	for _, mesh := range process.remotes {
		if err := mesh.Send(state); err != nil {
			process.logger.Err(err).Data(mesh.Info()).Warn("send state failed")
		}
	}
}

const sessionKey = "mesh"
