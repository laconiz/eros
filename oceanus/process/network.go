package process

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/cipher"
	"github.com/laconiz/eros/network/encoder"
	"github.com/laconiz/eros/network/invoker"
	"github.com/laconiz/eros/network/socket"
	"github.com/laconiz/eros/network/socket/reader"
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/remote"
	"time"
)

func (proc *Process) networkInvoker() invoker.Invoker {

	const key = "mesh"

	// 获取绑定信息
	bind := func(event *network.Event) (*proto.Mesh, bool) {
		value, ok := event.Ses.Load(key)
		if !ok {
			return nil, false
		}
		return value.(*proto.Mesh), true
	}

	invoker := invoker.NewSocketInvoker()

	// 连接建立
	invoker.Register(network.Connected{}, func(event *network.Event) {

		proc.mutex.RLock()
		defer proc.mutex.RUnlock()

		mesh := proc.local
		state, _ := mesh.State()

		// 发送网格信息
		event.Ses.Send(&proto.MeshJoin{Mesh: mesh.Info()})
		// 发送网格状态
		event.Ses.Send(state)
		// 发送节点列表
		event.Ses.Send(&proto.NodeJoin{Nodes: mesh.Nodes()})

		proc.logger.Data(event.Ses.Addr()).Info("connected")
	})

	// 连接断开
	invoker.Register(network.Disconnected{}, func(event *network.Event) {

		// 获取绑定信息
		info, ok := bind(event)
		if !ok {
			return
		}

		proc.mutex.Lock()
		defer proc.mutex.Unlock()

		// 更新连接信息
		if mesh, ok := proc.remotes[info.ID]; ok {
			mesh.UpdateSession(event.Ses)
			proc.logger.Data(info).Info("disconnected")
		}
	})

	// 网格状态
	invoker.Register(proto.State{}, func(event *network.Event) {

		// 获取绑定信息
		info, ok := bind(event)
		if !ok {
			return
		}

		state := event.Msg.(*proto.State)

		proc.mutex.Lock()
		defer proc.mutex.Unlock()

		type State struct {
			Mesh  *proto.Mesh
			State *proto.State
		}

		// 更新状态
		if mesh, ok := proc.remotes[info.ID]; ok {
			mesh.UpdateState(state)
			// data := &State{Mesh: info, State: state}
			// proc.logger.Data(data).Info("state updated")
		}
	})

	// 收到邮件
	invoker.Register(proto.Mail{}, func(event *network.Event) {

		mail := event.Msg.(*proto.Mail)

		proc.mutex.RLock()
		defer proc.mutex.RUnlock()

		// RPC RESPONSE
		if len(mail.To) == 0 && mail.Reply != proto.EmptyRpcID {
			if ch, ok := proc.channels[mail.Reply]; ok {
				ch <- mail
			}
		} else {
			// 普通邮件
			proc.local.Mail(mail)
		}
	})

	// 网格加入
	invoker.Register(proto.MeshJoin{}, func(event *network.Event) {

		// 设置绑定信息
		info := event.Msg.(*proto.MeshJoin).Mesh
		event.Ses.Store(key, info)

		proc.mutex.Lock()
		defer proc.mutex.Unlock()

		// 查询网格记录
		mesh, ok := proc.remotes[info.ID]
		if !ok {
			// 新建网格
			mesh = remote.New(info, proc)
			proc.remotes[info.ID] = mesh
		}

		// 更新连接信息
		mesh.UpdateSession(event.Ses)
		proc.logger.Data(info).Info("mesh join")
	})

	// 网格退出
	invoker.Register(proto.MeshQuit{}, func(event *network.Event) {

		// 获取绑定信息
		info, ok := bind(event)
		if !ok {
			return
		}

		proc.mutex.Lock()
		defer proc.mutex.Unlock()

		// 销毁网格
		if mesh, ok := proc.remotes[info.ID]; ok {
			mesh.Destroy()
			delete(proc.remotes, info.ID)
			proc.logger.Data(info).Info("mesh quit")
		}
	})

	// 节点加入
	invoker.Register(proto.NodeJoin{}, func(event *network.Event) {

		// 获取绑定信息
		info, ok := bind(event)
		if !ok {
			return
		}

		msg := event.Msg.(*proto.NodeJoin)

		proc.mutex.Lock()
		defer proc.mutex.Unlock()

		if mesh, ok := proc.remotes[info.ID]; ok {

			// 插入节点
			mesh.Insert(msg.Nodes)

			// 更新网格状态
			if msg.State != nil {
				mesh.UpdateState(msg.State)
			}

			proc.logger.Data(msg).Info("nodes join")
		}
	})

	// 节点退出
	invoker.Register(proto.NodeQuit{}, func(event *network.Event) {

		// 获取绑定信息
		info, ok := bind(event)
		if !ok {
			return
		}

		msg := event.Msg.(*proto.NodeQuit)

		proc.mutex.Lock()
		defer proc.mutex.Unlock()

		if mesh, ok := proc.remotes[info.ID]; ok {

			// 移除节点
			mesh.Remove(msg.Nodes)

			// 更新网格状态
			if msg.State != nil {
				mesh.UpdateState(msg.State)
			}

			proc.logger.Data(msg).Info("nodes quit")
		}
	})

	return invoker
}

func (proc *Process) NewAcceptor(addr string) network.Acceptor {

	return socket.NewAcceptor(&socket.AcceptorOption{
		Name: "oceanus.acceptor",
		Addr: addr,
		Session: socket.SessionOption{
			Timeout: time.Second * 11,
			Queue:   64,
			Invoker: proc.invoker,
			Encoder: encoder.NewNameMaker(),
			Cipher:  cipher.NewEmptyMaker(),
			Reader:  reader.NewSizeMaker(),
		},
		Level: logis.WARN,
	})
}

func (proc *Process) NewConnector(addr string) network.Connector {

	return socket.NewConnector(&socket.ConnectorOption{
		Name:      "oceanus.connector",
		Addr:      addr,
		Reconnect: true,
		Delays: []time.Duration{
			time.Millisecond,
			time.Millisecond * 100,
			time.Millisecond * 500,
			time.Second,
		},
		Session: socket.SessionOption{
			Timeout: time.Second * 11,
			Queue:   64,
			Invoker: proc.invoker,
			Encoder: encoder.NewNameMaker(),
			Cipher:  cipher.NewEmptyMaker(),
			Reader:  reader.NewSizeMaker(),
		},
		Level: logis.WARN,
	})
}
