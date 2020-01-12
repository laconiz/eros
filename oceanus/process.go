// 本地进程

package oceanus

import (
	"github.com/laconiz/eros/network"
	uuid "github.com/satori/go.uuid"
	"os"
	"sync"
)

func NewProcess() *Process {

	addr := os.Args[1]
	id := uuid.NewV3(uuid.NamespaceURL, addr).String()

	return &Process{
		Node: &Node{
			ID:    id,
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

// 状态
func (p *Process) Connected() bool {
	return true
}

// 同步节点
func (p *Process) OnNodeJoin(node *Node, session network.Session) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	burl, ok := p.burls[node.ID]
	if !ok {
		burl = NewBurl(node)
		burl.session = session
		p.burls[node.ID] = burl
		logger.Infof("node join: %+v", node)
	} else {
		burl.node = node
		burl.session = session
		logger.Infof("node update: %+v", node)
	}

	// 更新状态
	for _, course := range burl.courses {
		course.router.expired()
	}
}

// 删除节点
func (p *Process) OnNodeQuit(node *Node) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	burl, ok := p.burls[node.ID]
	if !ok {
		return
	}

	// 回收节点所属通道
	for _, course := range burl.courses {
		delete(p.courses, course.channel.ID)
		course.router.remove(course)
		delete(course.burl.courses, course.channel.ID)
	}

	// 删除节点
	delete(p.burls, node.ID)
}

// 节点连接断开
func (p *Process) onNodeDisconnected(node *Node) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	burl, ok := p.burls[node.ID]
	if !ok {
		return
	}

	for _, course := range burl.courses {
		course.router.expired()
	}

	burl.session = nil
}

// 同步通道
func (p *Process) addCourse(channels []*Channel) {

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
			logger.Warnf("can not find channel[%+v]'s node", channel)
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

		// 添加通道
		p.courses[channel.ID] = course
		burl.courses[channel.ID] = course
		router.add(course)
	}
}

// 回收通道
func (p *Process) removeCourse(channel *Channel) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if course, ok := p.courses[channel.ID]; ok {
		delete(p.courses, course.channel.ID)
		course.router.remove(course)
		delete(course.burl.courses, course.channel.ID)
	}
}

// 将本节点信息和通道列表同步至指定连接
func (p *Process) syncProcess(session network.Session) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	session.Send(&NodeJoinMsg{Node: p.Node})
	var channels []*Channel
	for _, thread := range p.threads {
		channels = append(channels, thread.Info())
	}
	session.Send(&ChannelJoinMsg{Channels: channels})
}

// 关闭节点
func (p *Process) destroy() {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	// 通知线程退出
	for _, thread := range p.threads {
		thread.Quit()
	}
	// 广播节点离线消息
	p.broadcast(&NodeQuitMsg{Node: p.Node})
}

// 向其他所有节点同步状态兼心跳功能
func (p *Process) notifyState() {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	p.broadcast(&NodeJoinMsg{Node: p.Node})
}

// 推送消息至所有节点
func (p *Process) broadcast(msg interface{}) {
	for _, burl := range p.burls {
		if burl.session != nil {
			burl.session.Send(msg)
		}
	}
}
