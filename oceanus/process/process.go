package process

import (
	"fmt"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/oceanus/local"
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/remote"
	uuid "github.com/satori/go.uuid"
	"sync"
)

var namespace = uuid.Must(uuid.FromString("4f31b82c-ca02-432c-afbf-8148c81ccaa2"))

// 创建一个进程
func New(addr string) (*Process, error) {
	// 相同IP和端口具有相同的UUID
	id := proto.MeshID(string(uuid.NewV3(namespace, addr).Bytes()))
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
		connectors: map[string]network.Connector{},
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
	local      *local.Mesh                   // 本地网格数据
	remotes    map[proto.MeshID]*remote.Mesh // 远程网格列表
	acceptor   network.Acceptor              // 网络侦听器
	connectors map[string]network.Connector  // 网络连接器列表
	router     *oceanus.Router               // 路由器
	log        logis.Logger                  // 日志接口
	mutex      sync.RWMutex
}

// 同步连接信息
func (process *Process) syncConnections(meshes []*proto.Mesh) {
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
		if _, ok := process.connectors[mesh.Addr]; ok {
			continue
		}
		// 将由对方网格发起连接
		if (local.Power > mesh.Power && (local.Power-mesh.Power)%2 == 0) ||
			(local.Power < mesh.Power && (mesh.Power-local.Power)%2 != 0) {
			continue
		}
		// 创建网格连接
		process.connectors[mesh.Addr] = process.NewConnector(mesh.Addr)
		process.log.Infof("connect to: %+v", mesh)
	}
}

//
func (process *Process) Run() {
	// 数据列表前缀
	const prefixKey = "oceanus/"
	// 数据
}
