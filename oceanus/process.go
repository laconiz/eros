package oceanus

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/laconiz/eros/consul"
	"github.com/laconiz/eros/json"
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/tcp"
	"github.com/laconiz/eros/oceanus/local"
	"github.com/laconiz/eros/oceanus/remote"
	uuid "github.com/satori/go.uuid"
	"math/big"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

var namespace = uuid.Must(uuid.FromString("4f31b82c-ca02-432c-afbf-8148c81ccaa2"))

// 生成一个进程
func NewProcess(addr string) Process {
	id := MeshID(uuid.NewV3(namespace, addr).String())
	mesh := &MeshInfo{ID: id, Addr: addr}
	router := NewRouter()
	return &process{
		mesh:       local.NewMesh(mesh, router),
		net:        remote.NewNet(router),
		acceptor:   nil,
		connectors: map[string]network.Connector{},
	}
}

type process struct {

	// 本地网格
	mesh *local.Mesh

	// 远程网格
	net *remote.Net

	// 网格服务端接口
	acceptor network.Acceptor
	// 网格客户端接口
	connectors map[string]network.Connector

	//
	mutex sync.RWMutex
}

// 同步网格连接
// TODO 评估非正常退出网格时未及时清理的网格信息的清理工作
func (p *process) syncMeshConnections(meshes []*MeshInfo) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 获取本地网格的权值
	ap, err := addrPower(p.mesh.Info().Addr)
	if err != nil {
		return
	}

	// 同步网格列表
	for _, mesh := range meshes {
		// 本地节点
		if mesh.Addr == p.mesh.Info().Addr {
			continue
		}
		// 已建立连接
		if _, ok := p.connectors[mesh.Addr]; ok {
			continue
		}
		// 获取网格权值
		mp, err := addrPower(mesh.Addr)
		if err != nil {
			continue
		}
		// 连接负载均衡
		// 当网格数量较多时接近于一半服务器连接一半客户端连接
		if (ap > mp && (ap-mp)%2 == 0) ||
			(ap < mp && (mp-ap)%2 != 0) {
			continue
		}
		// 创建客户端连接
		logger.Infof("connect to mesh: %+v", mesh)
		// 连接配置信息
		conf := tcp.ConnectorConfig{
			Name:      fmt.Sprintf("oceanus.connector.%s", mesh.Addr),
			Addr:      mesh.Addr,
			Reconnect: true,
			Session: tcp.SessionConfig{
				ReadTimeout:  time.Second * 11,
				WriteTimeout: time.Second * 11,
				LogLevel:     log.Warn,
				QueueLen:     64,
				Invoker:      p.newInvoker(),
			},
		}
		// 建立连接
		connector := tcp.NewConnector(conf)
		go connector.Run()
		p.connectors[mesh.Addr] = connector
	}
}

// 销毁
func (p *process) destroy() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// TODO 通知所有节点关闭

	// 同步网格离线消息
	p.net.Broadcast(&MeshQuit{MeshInfo: p.mesh.Info()})
}

// 同步网格状态
func (p *process) notifyState() {
	p.mutex.RLock()
	p.mutex.RUnlock()
	p.net.Broadcast(&MeshJoin{MeshInfo: p.mesh.Info()})
}

// 监听服务端接口
func (p *process) runAcceptor() {
	// 监听信息配置
	conf := tcp.AcceptorConfig{
		Name: "oceanus.acceptor",
		Addr: p.mesh.Info().Addr,
		Session: tcp.SessionConfig{
			ReadTimeout:  time.Second * 11,
			WriteTimeout: time.Second * 11,
			LogLevel:     log.Warn,
			QueueLen:     64,
			Invoker:      p.newInvoker(),
		},
	}
	// 监听端口
	p.acceptor = tcp.NewAcceptor(conf)
	go p.acceptor.Run()
}

// 网格通信回调接口
func (p *process) newInvoker() network.Invoker {

	const key = "mesh"

	invoker := network.NewStdInvoker()

	// 节点消息
	invoker.Register(Message{}, func(event *network.Event) {
		p.mutex.RLock()
		defer p.mutex.RUnlock()
		p.mesh.Push(event.Msg.(*Message))
	})

	// 网格上线
	invoker.Register(MeshJoin{}, func(event *network.Event) {
		msg := event.Msg.(*MeshJoin)
		logger.Infof("mesh join: %+v", msg.MeshInfo)
		p.mutex.Lock()
		defer p.mutex.Unlock()
		p.net.AddMesh(msg.MeshInfo, event.Session)
	})

	// 网格离线
	invoker.Register(MeshQuit{}, func(event *network.Event) {
		msg := event.Msg.(*MeshQuit)
		logger.Infof("mesh quit: %+v", msg.MeshInfo)
		p.mutex.Lock()
		defer p.mutex.Unlock()
		p.net.RemoveMesh(msg.ID)
	})

	// 节点上线
	invoker.Register(NodeJoin{}, func(event *network.Event) {
		p.mutex.Lock()
		defer p.mutex.Unlock()
		for _, node := range *event.Msg.(*NodeJoin) {
			p.net.InsertNode(node)
		}
	})

	// 节点离线
	invoker.Register(NodeQuit{}, func(event *network.Event) {
		p.mutex.Lock()
		defer p.mutex.Unlock()
		for _, node := range *event.Msg.(*NodeQuit) {
			p.net.RemoveNode(node)
		}
	})

	// 连接成功时 想远程网格同步本地网格和节点数据
	invoker.Register(network.Connected{}, func(event *network.Event) {
		p.mutex.RLock()
		defer p.mutex.RUnlock()
		event.Session.Send(&MeshJoin{MeshInfo: p.mesh.Info()})
	})

	// 连接断开时 根据连接上附加的网格数据更新网格状态
	invoker.Register(network.Disconnected{}, func(event *network.Event) {
		if mesh := event.Session.Get(key); mesh != nil {
			p.mutex.Lock()
			defer p.mutex.Unlock()
			p.net.AddMesh(mesh.(*MeshInfo), event.Session)
		}
	})

	return invoker
}

// 运行网格
func (p *process) Run() {

	const prefixKey = "oceanus/"

	mesh := p.mesh.Info()
	logger.Infof("mesh: %+v", mesh)

	// 服务端监听
	p.runAcceptor()
	defer p.acceptor.Stop()

	// 注册网格信息
	key := fmt.Sprintf("%s%v", prefixKey, mesh.ID)
	if err := consul.KV().Store(key, mesh); err != nil {
		logger.Errorf("register mesh error: %v", err)
	}
	defer func() {
		if err := consul.KV().Delete(key); err != nil {
			logger.Errorf("deregister mesh error: %v", err)
		}
	}()

	// 监视网格列表
	watcher, err := consul.NewKeyPrefixWatcher(prefixKey, func(pairs api.KVPairs) {
		// 反序列化网格列表
		var meshes []*MeshInfo
		for _, pair := range pairs {
			mesh := &MeshInfo{}
			if err := json.Unmarshal(pair.Value, mesh); err == nil {
				meshes = append(meshes, mesh)
			}
		}
		// 同步网格列表
		p.syncMeshConnections(meshes)
	})
	if err != nil {
		logger.Errorf("create watcher error: %v", err)
		return
	}
	go watcher.Run()
	defer watcher.Stop()

	// 监听进程退出信号
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, os.Kill)

	// TODO 添加退出信号
	for {
		select {
		case signal := <-exit:
			logger.Infof("exit signal received: %v", signal)
			p.destroy()
			return
		case <-time.After(time.Second * 5):
			p.notifyState()
		}
	}
}

func (p *process) NewTypeNode(typo NodeType) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	node := &NodeInfo{
		ID:   randomNodeID(),
		Type: typo,
		Key:  randomNodeKey(),
		Mesh: p.mesh.Info().ID,
	}

	p.mesh.InsertNode(node)
}

func (p *process) NewKeyNode(typo NodeType, key NodeKey) {

	id := uuid.NewV1().String()

	p.mutex.Lock()
	defer p.mutex.Unlock()

	node := &NodeInfo{
		ID:   NodeID(id),
		Type: typo,
		Key:  key,
		Mesh: p.mesh.Info().ID,
	}

	p.mesh.InsertNode(node)
}

// 生成一个随机的节点ID
func randomNodeID() NodeID {
	return NodeID(uuid.NewV1().String())
}

// 生成一个随机的节点KEY
func randomNodeKey() NodeKey {
	return NodeKey(uuid.NewV1().String())
}

// 获取一个IPV4地址的权值
func addrPower(addr string) (uint64, error) {
	// 分离IP和端口
	ap := strings.Split(addr, ":")
	if len(ap) != 2 {
		return 0, errors.New("invalid addr format")
	}
	// 解析IP地址
	ip := net.ParseIP(ap[0])
	if ip == nil {
		return 0, errors.New("invalid ip address")
	}
	// 反序列化端口号
	port, err := strconv.ParseUint(ap[1], 10, 64)
	if err != nil || port > 65535 {
		return 0, errors.New("invalid port address")
	}
	// 计算权值
	power := big.NewInt(0).SetBytes(ip.To4()).Uint64()
	return port<<32 | power, nil
}
