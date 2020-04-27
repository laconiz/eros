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
	"os"
	"os/signal"
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
	ticker       *time.Ticker     // 心跳计时器
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

	proc.logger.Info("watch consul")
	go proc.synchronizer.Run()

	proc.logger.Info("started")

	proc.ticker = time.NewTicker(time.Second * 10)
	go func() {

	}()
}

func (proc *Process) stop(watcher *consul.Plan) {

	proc.logger.Info("stopping")

	proc.mutex.Lock()
	defer proc.mutex.Unlock()

	// 关闭定时器
	if proc.ticker != nil {
		proc.ticker.Stop()
	}

	// 关闭同步器
	proc.synchronizer.Stop()
	proc.logger.Info("synchronizer stopped")

	// 清理远程网格
	for _, mesh := range proc.remotes {
		mesh.Destroy()
	}
	proc.remotes = Remotes{}
	proc.logger.Info("remote meshes destroyed")

	proc.logger.Info("destroy remote meshes")
	for _, mesh := range proc.remotes {
		mesh.Send(&proto.MeshQuit{Mesh: proc.local.Info()})
		mesh.Destroy()
	}
	proc.remotes = map[proto.MeshID]*remote.Mesh{}

	proc.logger.Info("clear connectors")
	for _, connector := range proc.connectors {
		connector.Stop()
	}
	proc.connectors = map[proto.MeshID]network.Connector{}

	proc.logger.Info("unwatch consul")
	watcher.Stop()

	proc.logger.Info("deregister from consul")
	if err := proc.deregister(); err != nil {
		proc.logger.Err(err).Error("deregister error")
	}

	proc.logger.Info("stop acceptor")
	proc.acceptor.Stop()

	proc.logger.Info("destroy local mesh")
	proc.local.Destroy()

	proc.logger.Info("stopped")
}
