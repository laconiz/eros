package process

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/laconiz/eros/database/consul"
	"github.com/laconiz/eros/database/consul/consulor"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/oceanus/local"
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/remote"
	uuid "github.com/satori/go.uuid"
	"os"
	"os/signal"
	"sync"
	"time"
)

var namespace = uuid.Must(uuid.FromString("4f31b82c-ca02-432c-afbf-8148c81ccaa2"))

// 创建一个进程
func New(addr string) (*Process, error) {
	// 相同IP和端口具有相同的UUID
	id := proto.MeshID(uuid.NewV3(namespace, addr).String())
	// 路由器
	router := oceanus.NewRouter()
	// 本地网格权值
	power, err := addrPower(addr)
	if err != nil {
		return nil, fmt.Errorf("get addr power error: %w", err)
	}
	// 本地网格
	info := &proto.Mesh{ID: id, Addr: addr, Power: power}
	// 进程信息
	process := &Process{
		local:      local.NewMesh(info, &proto.State{}, router),
		remotes:    map[proto.MeshID]*remote.Mesh{},
		acceptor:   nil,
		connectors: map[proto.MeshID]network.Connector{},
		router:     router,
		log:        nil,
	}
	// 网络侦听器
	process.acceptor = process.NewAcceptor(addr)
	// 日志接口
	process.log = logisor.Field(logis.Module, "oceanus")
	//
	return process, nil
}

type Process struct {
	local      *local.Mesh                        // 本地网格数据
	remotes    map[proto.MeshID]*remote.Mesh      // 远程网格列表
	acceptor   network.Acceptor                   // 网络侦听器
	connectors map[proto.MeshID]network.Connector // 网络连接器列表
	router     *oceanus.Router                    // 路由器
	log        logis.Logger                       // 日志接口
	mutex      sync.RWMutex
}

// 同步连接信息
func (process *Process) syncConnections(meshes map[string]*proto.Mesh) {

	process.mutex.Lock()
	defer process.mutex.Unlock()

	local := process.local.Info()

	// 遍历所有网格节点
	for _, mesh := range meshes {
		// 本地节点
		if mesh.Addr == local.Addr {
			continue
		}
		// 已建立连接
		if _, ok := process.connectors[mesh.ID]; ok {
			continue
		}
		// 将由对方网格发起连接
		if (local.Power > mesh.Power && (local.Power-mesh.Power)%2 == 0) ||
			(local.Power < mesh.Power && (mesh.Power-local.Power)%2 != 0) {
			continue
		}
		// 创建网格连接
		process.connectors[mesh.ID] = process.NewConnector(mesh.Addr)
		process.log.Infof("connect to: %+v", mesh)
	}

	for id, connector := range process.connectors {

		if _, ok := meshes[string(id)]; ok {
			continue
		}

		connector.Stop()
		delete(process.connectors, id)
	}
}

//
func (process *Process) Run() {

	// 本地网格信息
	info := process.local.Info()
	process.log.Data(info).Info("start")
	defer process.local.Destroy()

	// 运行侦听器
	process.acceptor.Run()
	defer process.acceptor.Stop()
	// 注册网格信息
	key := string(prefix + info.ID)
	if err := consulor.KV().Store(key, info); err != nil {
		process.log.Errorf("register mesh error: %v", err)
		return
	}
	defer func() {
		if err := consulor.KV().Delete(key); err != nil {
			process.log.Errorf("deregister mesh error: %v", err)
		}
	}()

	// 监视网格列表
	watcher, err := consulor.Watcher().Keyprefix(prefix, process.OnWatcher)
	if err != nil {
		process.log.Errorf("watch meshes error: %v", err)
		return
	}
	go watcher.Run()
	defer watcher.Stop()

	process.Loop()
}

// 网格轮询
func (process *Process) Loop() {

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, os.Kill)

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {

		case <-exit:

			process.Destroy()
			process.log.Info("exit")
			return

		case <-ticker.C:

			process.log.Info("notify state")
			process.NotifyState()
		}
	}
}

// 网格列表监视回调
func (process *Process) OnWatcher(_ uint64, pairs interface{}) {

	meshes := map[string]*proto.Mesh{}

	err := consul.ParsePairs(prefix, pairs.(api.KVPairs), &meshes, false)
	if err != nil {
		process.log.Errorf("parse mesh error: %v", err)
		return
	}

	process.syncConnections(meshes)
}

// 退出
func (process *Process) Destroy() {

	// 通知远端网格销毁本地网格信息
	for _, mesh := range process.remotes {
		mesh.Send(&proto.MeshQuit{Mesh: process.local.Info()})
	}
}

// 通知远端网格本地网格状态
func (process *Process) NotifyState() {

	state, _ := process.local.State()

	for _, mesh := range process.remotes {
		if err := mesh.Send(state); err != nil {
			process.log.Data(mesh.Info()).Warnf("notify state error: %v", err)
		}
	}
}

const prefix = "oceanus/"