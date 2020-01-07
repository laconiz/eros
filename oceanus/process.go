// 本地进程

package oceanus

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/tcp"
	uuid "github.com/satori/go.uuid"
	"net"
	"os"
	"sync"
	"time"
)

func NewProcess() *Process {

	return &Process{
		Node: &Node{
			ID:    uuid.NewV1().String(),
			Addr:  os.Args[1],
			State: State{},
		},
		threads:    map[string]*Thread{},
		burls:      map[string]*Burl{},
		courses:    map[string]*Course{},
		routers:    map[string]*Router{},
		acceptor:   nil,
		connectors: map[string]network.Connector{},
	}
}

type Process struct {

	// 本地节点
	*Node

	// 本地线程
	threads map[string]*Thread

	// 远程节点
	burls map[string]*Burl
	// 远程线程
	courses map[string]*Course

	// 分析器
	routers map[string]*Router

	// TCP接口
	acceptor   network.Acceptor
	connectors map[string]network.Connector

	mutex sync.RWMutex
}

// 推送本地消息
func (p *Process) Push(message *Message) error {

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	for _, id := range message.Receivers {
		if thread, ok := p.threads[id]; ok {
			thread.Push(message)
		}
	}

	return nil
}

//
func (p *Process) SyncConnectors(nodes []*Node) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, node := range nodes {

		// 本地节点
		if node.ID == p.ID {
			continue
		}

		// 已存在节点
		if _, ok := p.burls[node.ID]; ok {
			continue
		}

		// 已存在连接
		if _, ok := p.connectors[node.Addr]; ok {
			continue
		}

		net.ParseIP()
	}
}

// 删除远程通道
func (p *Process) destroyCourse(course *Course) {
	delete(p.courses, course.channel.ID)
	course.router.remove(course)
	delete(course.burl.courses, course.channel.ID)
}

// 同步节点
func (p *Process) OnNodeJoin(nodes []*Node) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, node := range nodes {

		// 本地节点
		if node.ID == p.Node.ID {
			continue
		}

		// 已存在节点
		if _, ok := p.burls[node.ID]; ok {
			continue
		}

		p.burls[node.ID] = p.NewBurl(node)
	}
}

// 删除节点
func (p *Process) OnNodeQuit(nodes []*Node) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, node := range nodes {

		burl, ok := p.burls[node.ID]
		if !ok {
			continue
		}

		// 清理通道信息
		for _, course := range burl.courses {
			p.destroyCourse(course)
		}

		// 删除节点
		delete(p.burls, node.ID)
	}
}

// 同步通道
func (p *Process) OnRouteJoin(channels []*Channel) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, channel := range channels {

		// 通道已存在
		if _, ok := p.courses[channel.ID]; ok {
			continue
		}

		// 节点不存在
		burl, ok := p.burls[channel.Node]
		if !ok {
			continue
		}

		// 分析器
		router, ok := p.routers[channel.Type]
		if !ok {
			router = NewRouter(channel.Type)
			p.routers[channel.Type] = router
		}

		course := &Course{
			channel: channel,
			burl:    burl,
			router:  router,
		}

		p.courses[channel.ID] = course
		burl.courses[channel.ID] = course
		router.add(course)
	}
}

// 回收通道
func (p *Process) OnRouteQuit(channels []*Channel) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, channel := range channels {
		if course, ok := p.courses[channel.ID]; ok {
			p.destroyCourse(course)
		}
	}
}

func (p *Process) NewBurl(node *Node) *Burl {

	burl := &Burl{
		Node:    node,
		conn:    nil,
		courses: map[string]*Course{},
	}

	invoker := network.NewStdInvoker()

	// 连接成功
	invoker.Register(network.Connected{}, func(event *network.Event) {

		p.mutex.Lock()
		defer p.mutex.Unlock()

		burl.connected = true

		// 推送当前节点
		event.Session.Send(&NodeJoinMsg{Node: p.Node})

		// 推送当前通道
		var channels []*Channel
		for _, thread := range p.threads {
			channels = append(channels, thread.Channel())
		}
		event.Session.Send(&ChannelJoinMsg{Channels: channels})
	})

	// 断开连接
	// 更新状态信息
	invoker.Register(network.Disconnected{}, func(event *network.Event) {

		p.mutex.Lock()
		defer p.mutex.Unlock()

		burl.connected = false
	})

	// 推送通道加入
	invoker.Register(ChannelJoinMsg{}, func(event *network.Event) {
		p.OnRouteJoin(event.Msg.(*ChannelJoinMsg).Channels)
	})

	// 推送通道退出
	invoker.Register(ChannelQuitMsg{}, func(event *network.Event) {
		p.OnRouteQuit(event.Msg.(*ChannelQuitMsg).Channels)
	})

	return burl
}

func (p *Process) RunAcceptor() {

	invoker := network.NewStdInvoker()

	invoker.Register(NodeJoinMsg{}, func(event *network.Event) {
		p.OnNodeJoin(event.Msg.(*NodeJoinMsg).Nodes)
	})

	invoker.Register(NodeQuitMsg{}, func(event *network.Event) {
		p.OnNodeQuit(event.Msg.(*NodeQuitMsg).Node)
	})

	invoker.Register(ChannelJoinMsg{}, func(event *network.Event) {

	})

}

func (p *Process) Run() error {

	id := uuid.NewV1().String()

	plan, err := watch.Parse(map[string]interface{}{
		"type":    "service",
		"service": "oceanus",
	})
	if err != nil {
		return fmt.Errorf("parse watch plan error: %w", err)
	}

	plan.Handler = func(_ uint64, result interface{}) {

		if entries, ok := result.([]*api.ServiceEntry); ok {

			var nodes []*Node

			for _, entry := range entries {

				service := entry.Service
				addr := fmt.Sprintf("%s:%d", service.Address, service.Port)

				nodes = append(nodes, &Node{
					ID:    entry.Service.ID,
					Addr:  addr,
					State: State{},
				})
			}
		}
	}

	client := &api.Client{}

	if err := client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Kind:              "",
		ID:                "",
		Name:              "",
		Tags:              nil,
		Port:              0,
		Address:           "",
		TaggedAddresses:   nil,
		EnableTagOverride: false,
		Meta:              nil,
		Weights:           nil,
		Check:             nil,
		Checks:            nil,
		Proxy:             nil,
		Connect:           nil,
	}); err != nil {

	}

	defer client.Agent().ServiceDeregister(id)

}
