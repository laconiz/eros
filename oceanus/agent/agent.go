package agent

import (
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/oceanus/local"
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/remote"
	uuid "github.com/satori/go.uuid"
	"sync"
)

type Agent struct {

	// 本地网格
	mesh *local.Mesh

	// 远程网格
	net *remote.Net

	//
	mutex sync.RWMutex
}

func (a *Agent) newInvoker() network.Invoker {

	const key = "mesh"

	invoker := network.NewStdInvoker()

	// 节点消息
	invoker.Register(proto.Message{}, func(event *network.Event) {
		a.mutex.RLock()
		defer a.mutex.RUnlock()
		a.mesh.Push(event.Msg.(*proto.Message))
	})

	// 网格上线
	invoker.Register(proto.MeshJoin{}, func(event *network.Event) {
		a.mutex.Lock()
		defer a.mutex.Unlock()
		a.net.AddMesh(event.Msg.(*proto.MeshJoin).Mesh, event.Session)
	})

	// 网格离线
	invoker.Register(proto.MeshQuit{}, func(event *network.Event) {
		a.mutex.Lock()
		defer a.mutex.Unlock()
		a.net.RemoveMesh(event.Msg.(*proto.MeshQuit).ID)
	})

	// 节点上线
	invoker.Register(proto.NodeJoin{}, func(event *network.Event) {
		a.mutex.Lock()
		defer a.mutex.Unlock()
		for _, node := range *event.Msg.(*proto.NodeJoin) {
			a.net.InsertNode(node)
		}
	})

	// 节点离线
	invoker.Register(proto.NodeQuit{}, func(event *network.Event) {
		a.mutex.Lock()
		defer a.mutex.Unlock()
		for _, node := range *event.Msg.(*proto.NodeQuit) {
			a.net.RemoveNode(node)
		}
	})

	// 连接成功时 想远程网格同步本地网格和节点数据
	invoker.Register(network.Connected{}, func(event *network.Event) {
		a.mutex.RLock()
		defer a.mutex.RUnlock()
		event.Session.Send(&proto.MeshJoin{Mesh: a.mesh.Info()})
	})

	// 连接断开时 根据连接上附加的网格数据更新网格状态
	invoker.Register(network.Disconnected{}, func(event *network.Event) {
		if mesh := event.Session.Get(key); mesh != nil {
			a.mutex.Lock()
			defer a.mutex.Unlock()
			a.net.AddMesh(mesh.(*proto.Mesh), event.Session)
		}
	})

	return invoker
}

func (a *Agent) NewTypeNode(typo proto.NodeType) {

	a.mutex.Lock()
	defer a.mutex.Unlock()

	node := &proto.Node{
		ID:   randomNodeID(),
		Type: typo,
		Key:  randomNodeKey(),
		Mesh: a.mesh.Info().ID,
	}

	a.mesh.InsertNode(node)
}

func (a *Agent) NewKeyNode(typo proto.NodeType, key proto.NodeKey) {

	id := uuid.NewV1().String()

	a.mutex.Lock()
	defer a.mutex.Unlock()

	node := &proto.Node{
		ID:   proto.NodeID(id),
		Type: typo,
		Key:  key,
		Mesh: a.mesh.Info().ID,
	}

	a.mesh.InsertNode(node)
}

func randomNodeID() proto.NodeID {
	return proto.NodeID(uuid.NewV1().String())
}

func randomNodeKey() proto.NodeKey {
	return proto.NodeKey(uuid.NewV1().String())
}
