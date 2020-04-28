package process

import (
	"github.com/laconiz/eros/database/consul"
	"github.com/laconiz/eros/database/consul/consulor"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/encoder"
	"github.com/laconiz/eros/network/invoker"
	"github.com/laconiz/eros/oceanus/abstract"
	"github.com/laconiz/eros/oceanus/local"
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/remote"
	"github.com/laconiz/eros/oceanus/router"
	"sync"
	"time"
)

type RpcID = proto.RpcID
type MeshID = proto.MeshID
type Router = abstract.Router
type Remotes = map[MeshID]*remote.Mesh
type Channels = map[RpcID]chan interface{}

const module = "oceanus"

// 创建进程
func New(addr string, encoder encoder.Encoder) (*Process, error) {

	proc := &Process{
		remotes:  Remotes{},
		router:   router.New(),
		logger:   logisor.Module(module),
		encoder:  encoder,
		channels: Channels{},
		signal:   make(chan bool, 1),
	}

	// 创建本地网格
	info := &proto.Mesh{
		ID:   MeshID(NewNamespaceUUID(addr, namespace)),
		Addr: addr,
	}
	proc.local = local.New(info, proc)

	// 创建网络回调接口
	proc.invoker = proc.networkInvoker()

	// 创建网络侦听器
	proc.acceptor = proc.NewAcceptor(addr)

	// 创建同步器
	synchronizer, err := consulor.Watcher().Prefix(prefix, proc.synchronize)
	if err != nil {
		return nil, err
	}
	proc.synchronizer = synchronizer

	return proc, nil
}

type Process struct {
	local        *local.Mesh      // 本地网格数据
	remotes      Remotes          // 远程网格列表
	router       Router           // 路由器
	invoker      invoker.Invoker  // 网络回调接口
	acceptor     network.Acceptor // 网络侦听器
	logger       logis.Logger     // 日志接口
	encoder      encoder.Encoder  // 邮件编码器
	channels     Channels         // RPC CHANNEL列表
	synchronizer *consul.Plan     // 网格同步器
	signal       chan bool        // 退出信号
	mutex        sync.RWMutex
}

func (proc *Process) Local() abstract.Mesh {
	return proc.local
}

func (proc *Process) Router() abstract.Router {
	return proc.router
}

func (proc *Process) Logger() logis.Logger {
	return proc.logger
}

func (proc *Process) Run() {

	info := proc.local.Info()
	proc.logger.Data(info).Info("starting")

	proc.logger.Info("start acceptor")
	proc.acceptor = proc.NewAcceptor(info.Addr)
	proc.acceptor.Run()

	proc.logger.Info("register to consul")
	err := consulor.KV().Store(prefix+string(info.ID), info)
	if err != nil {
		proc.logger.Err(err).Error("register error")
		return
	}

	proc.logger.Info("start synchronizer")
	go proc.synchronizer.Run()

	proc.logger.Info("start tick")
	go proc.tick()

	proc.logger.Info("started")
}

func (proc *Process) Stop() {

	proc.logger.Info("stopping")

	proc.mutex.Lock()
	defer proc.mutex.Unlock()

	info := proc.local.Info()

	// 关闭定时器
	proc.logger.Info("stop tick")
	proc.signal <- true

	// 关闭同步器
	proc.logger.Info("stop synchronizer")
	proc.synchronizer.Stop()

	// 清理远程网格
	proc.logger.Info("destroy remote meshes")
	msg := &proto.MeshQuit{ID: info.ID}
	for _, mesh := range proc.remotes {
		// 发送网格退出消息
		mesh.Send(msg)
		// 清理网格
		mesh.Destroy()
	}
	proc.remotes = Remotes{}

	proc.logger.Info("stop synchronizer")
	proc.synchronizer.Stop()

	proc.logger.Info("deregister from consul")
	err := consulor.KV().Delete(prefix + string(info.ID))
	if err != nil {
		proc.logger.Err(err).Error("deregister error")
	}

	proc.logger.Info("stop acceptor")
	proc.acceptor.Stop()

	proc.logger.Info("destroy local mesh")
	proc.local.Destroy()

	proc.logger.Info("stopped")
}

func (proc *Process) tick() {

	sender := func() {

		proc.mutex.RLock()
		defer proc.mutex.RUnlock()

		state, _ := proc.local.State()

		for _, mesh := range proc.remotes {
			mesh.Send(state)
		}
	}

	ticker := time.NewTicker(time.Second * 10)

	for {
		select {
		case <-proc.signal:
			return
		case <-ticker.C:
			sender()
		}
	}
}

func (proc *Process) broadcast(msg interface{}) {

	proc.mutex.RLock()
	defer proc.mutex.RUnlock()

	for _, mesh := range proc.remotes {
		mesh.Send(msg)
	}
}
