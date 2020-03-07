package process

import (
	"encoding/hex"
	"fmt"
	"github.com/laconiz/eros/database/consul"
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

	id := proto.MeshID(hex.EncodeToString(uuid.NewV3(namespace, addr).Bytes()))

	router := oceanus.NewRouter()

	power, err := addrPower(addr)
	if err != nil {
		return nil, fmt.Errorf("get addr power error: %w", err)
	}

	info := &proto.Mesh{ID: id, Addr: addr, Power: power}

	return &Process{
		local:      local.NewMesh(info, &proto.State{}, router),
		remotes:    map[proto.MeshID]*remote.Mesh{},
		connectors: map[proto.MeshID]network.Connector{},
		router:     router,
		logger:     logisor.Field(logis.Module, "oceanus"),
	}, nil
}

type Process struct {
	local      *local.Mesh                        // 本地网格数据
	remotes    map[proto.MeshID]*remote.Mesh      // 远程网格列表
	acceptor   network.Acceptor                   // 网络侦听器
	connectors map[proto.MeshID]network.Connector // 网络连接器列表
	router     *oceanus.Router                    // 路由器
	logger     logis.Logger                       // 日志接口
	mutex      sync.RWMutex
}

func (process *Process) Run() {

	info := process.local.Info()
	process.logger.Data(info).Info("starting")

	process.logger.Info("start acceptor")
	process.acceptor = process.NewAcceptor(info.Addr)
	process.acceptor.Run()

	process.logger.Info("register to consul")
	if err := process.register(); err != nil {
		process.logger.Err(err).Error("register error")
		return
	}

	process.logger.Info("watch consul")
	watcher := process.watcher()
	go watcher.Run()

	process.logger.Info("started")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, os.Kill)

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			process.broadcastState()
		case <-exit:
			goto BREAK
		}
	}

BREAK:
	process.stop(watcher)
}

func (process *Process) stop(watcher *consul.Plan) {

	process.logger.Info("stopping")

	process.mutex.Lock()
	defer process.mutex.Unlock()

	process.logger.Info("destroy remote meshes")
	for _, mesh := range process.remotes {
		mesh.Send(&proto.MeshQuit{Mesh: process.local.Info()})
		mesh.Destroy()
	}
	process.remotes = map[proto.MeshID]*remote.Mesh{}

	process.logger.Info("clear connectors")
	for _, connector := range process.connectors {
		connector.Stop()
	}
	process.connectors = map[proto.MeshID]network.Connector{}

	process.logger.Info("unwatch consul")
	watcher.Stop()

	process.logger.Info("deregister from consul")
	if err := process.deregister(); err != nil {
		process.logger.Err(err).Error("deregister error")
	}

	process.logger.Info("stop acceptor")
	process.acceptor.Stop()

	process.logger.Info("destroy local mesh")
	process.local.Destroy()

	process.logger.Info("stopped")
}
