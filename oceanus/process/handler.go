package process

import (
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/remote"
)

// 网络连接时发送网格数据
func (proc *Process) OnConnected(event *network.Event) {

	proc.mutex.RLock()
	defer proc.mutex.RUnlock()

	mesh := proc.local
	state, _ := mesh.State()
	event.Ses.Send(&proto.MeshJoin{
		Mesh:  mesh.Info(),
		State: state,
		Nodes: mesh.Nodes(),
	})

	proc.logger.Info("join to remote")
}

// 网络断开时更新网格状态
func (proc *Process) OnDisconnected(event *network.Event) {

	value, ok := event.Ses.Get(sessionKey)
	if !ok {
		return
	}
	info := value.(*proto.Mesh)

	proc.mutex.Lock()
	defer proc.mutex.Unlock()

	if mesh, ok := proc.remotes[info.ID]; ok {
		mesh.UpdateSession(nil)
	}
}

// 处理消息
func (proc *Process) OnMail(event *network.Event) {

	mail := event.Msg.(*proto.Mail)

	// PROXY RESPONSE
	if len(mail.To) == 0 && mail.Reply != proto.EmptyRpcID {

		proc.mutex.RLock()
		defer proc.mutex.RUnlock()

		if ch, ok := proc.channels[mail.Reply]; ok {
			ch <- mail
		}
	}

	proc.mutex.RLock()
	defer proc.mutex.RUnlock()
	if err := proc.local.Push(event.Msg.(*proto.Mail)); err != nil {
		proc.logger.Warnf("recv mail error: %v", err)
	}
}

// 网格状态
func (proc *Process) OnState(event *network.Event) {

	value, ok := event.Ses.Get(sessionKey)
	if !ok {
		return
	}
	info := value.(*proto.Mesh)

	msg := event.Msg.(*proto.State)

	proc.mutex.Lock()
	defer proc.mutex.Unlock()

	if mesh, ok := proc.remotes[info.ID]; ok {
		mesh.UpdateState(msg)
		data := &proto.MeshJoin{Mesh: info, State: msg}
		proc.logger.Data(data).Info("remote mesh update")
	}
}

// 插入网格
func (proc *Process) OnMeshJoin(event *network.Event) {

	msg := event.Msg.(*proto.MeshJoin)
	event.Ses.Set(sessionKey, msg.Mesh)

	proc.mutex.Lock()
	defer proc.mutex.Unlock()

	mesh, ok := proc.remotes[msg.Mesh.ID]
	if !ok {
		mesh = remote.NewMesh(msg.Mesh, msg.State, proc.router)
		proc.remotes[msg.Mesh.ID] = mesh
	}

	mesh.UpdateSession(event.Ses)
	mesh.Insert(msg.Nodes)

	proc.logger.Data(msg).Info("remote mesh join")
}

// 移除网格
func (proc *Process) OnMeshQuit(event *network.Event) {

	value, ok := event.Ses.Get(sessionKey)
	if !ok {
		return
	}
	info := value.(*proto.Mesh)

	proc.mutex.Lock()
	defer proc.mutex.Unlock()

	if mesh, ok := proc.remotes[info.ID]; ok {
		mesh.Destroy()
		delete(proc.remotes, info.ID)
		proc.logger.Data(info).Info("remote mesh quit")
	}

	if connector, ok := proc.connectors[info.ID]; ok {
		connector.Stop()
		delete(proc.connectors, info.ID)
		proc.logger.Data(info.ID).Info("connector stopped")
	}
}

// 插入节点
func (proc *Process) onNodeJoin(event *network.Event) {

	value, ok := event.Ses.Get(sessionKey)
	if !ok {
		return
	}
	info := value.(*proto.Mesh)

	msg := event.Msg.(*proto.NodeJoin)

	proc.mutex.Lock()
	defer proc.mutex.Unlock()

	if mesh, ok := proc.remotes[info.ID]; ok {
		mesh.Insert(msg.Nodes)
		proc.logger.Data(msg).Info("remote node join")
	}
}

// 移除节点
func (proc *Process) onNodeQuit(event *network.Event) {

	value, ok := event.Ses.Get(sessionKey)
	if !ok {
		return
	}
	info := value.(*proto.Mesh)

	msg := event.Msg.(*proto.NodeQuit)

	proc.mutex.Lock()
	defer proc.mutex.Unlock()

	if mesh, ok := proc.remotes[info.ID]; ok {
		mesh.Remove(msg.Nodes)
		proc.logger.Data(msg).Info("remote node quit")
	}
}

// 广播状态
func (proc *Process) broadcastState() {

	proc.mutex.RLock()
	defer proc.mutex.RUnlock()

	state, _ := proc.local.State()
	proc.broadcast(state)
}

const sessionKey = "mesh"
