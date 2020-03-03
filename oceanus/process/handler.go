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
	state, _ := process.local.State()
	// 发送网格数据
	event.Ses.Send(&proto.MeshJoin{
		Mesh:  process.local.Info(),
		State: state,
		Nodes: process.local.Nodes(),
	})
}

// 网络断开时更新网格状态
func (process *Process) OnDisconnected(event *network.Event) {
	if value, ok := event.Ses.Get(sessionKey); ok {
		// 连接绑定网格信息
		info := value.(*proto.Mesh)
		process.mutex.Lock()
		defer process.mutex.Unlock()
		// 查找网格
		if mesh, ok := process.remotes[info.ID]; ok {
			// 更新网格连接
			mesh.UpdateSession(nil)
		}
	}
}

// 处理消息
func (process *Process) OnMail(event *network.Event) {
	msg := event.Msg.(*proto.Mail)
	process.mutex.RLock()
	defer process.mutex.RUnlock()
	// 派发消息
	if err := process.local.Push(msg); err != nil {
		process.log.Warnf("recv mail error: %v", err)
	}
}

// 插入网格
func (process *Process) OnMeshJoin(event *network.Event) {
	msg := event.Msg.(*proto.MeshJoin)
	process.mutex.Lock()
	defer process.mutex.Unlock()
	event.Ses.Set(sessionKey, msg.Mesh)
	// 查找网格
	mesh, ok := process.remotes[msg.Mesh.ID]
	if !ok {
		// 新建网格
		mesh = remote.NewMesh(msg.Mesh, msg.State, process.router)
		process.remotes[msg.Mesh.ID] = mesh
	}
	// 更新网格连接
	mesh.UpdateSession(event.Ses)
	// 插入节点
	for _, node := range msg.Nodes {
		mesh.Insert(node)
	}
}

// 移除网格
func (process *Process) OnMeshQuit(event *network.Event) {
	msg := event.Msg.(*proto.MeshQuit)
	process.mutex.Lock()
	defer process.mutex.Unlock()
	if mesh, ok := process.remotes[msg.Mesh.ID]; ok {
		// 移除网格
		mesh.Destroy()
		delete(process.remotes, msg.Mesh.ID)
	}
}

// 插入节点
func (process *Process) OnNodeJoin(event *network.Event) {
	msg := event.Msg.(*proto.NodeJoin)
	if value, ok := event.Ses.Get(sessionKey); ok {
		// 连接绑定网格信息
		info := value.(*proto.Mesh)
		// 查找网格
		if mesh, ok := process.remotes[info.ID]; ok {
			process.mutex.Lock()
			defer process.mutex.Unlock()
			// 插入节点
			for _, node := range msg.Nodes {
				mesh.Insert(node)
			}
		}
	}
}

// 移除节点
func (process *Process) OnNodeQuit(event *network.Event) {
	msg := event.Msg.(*proto.NodeQuit)
	if value, ok := event.Ses.Get(sessionKey); ok {
		// 连接绑定网格信息
		info := value.(*proto.Mesh)
		// 查找网格
		if mesh, ok := process.remotes[info.ID]; ok {
			process.mutex.Lock()
			defer process.mutex.Unlock()
			// 移除节点
			for _, node := range msg.Nodes {
				mesh.Remove(node.ID)
			}
		}
	}
}

const sessionKey = "mesh"
