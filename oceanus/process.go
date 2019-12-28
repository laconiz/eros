package oceanus

import (
	"errors"
	"fmt"
	"github.com/laconiz/eros/network/tcp"
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/queue"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"sync"
)

type Process struct {

	// 同步版本
	version uint64

	// 本地节点
	node Node
	// 节点信息
	nodes map[string]Node
	// 通道信息
	channels map[string]Channel
	// 通道分组
	groups map[string]map[string]Channel

	// 监听器
	acceptor tcp.Acceptor

	//
	mutex sync.RWMutex

	//
	rand rand.Rand
}

// 同步通道模型
func (p *Process) Sync(states *proto.ChannelStates) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, uuid := range states.Delete {

		// 查找通道信息
		channel, ok := p.channels[uuid]
		if !ok {
			continue
		}

		// 删除通道信息
		info := channel.Info()
		delete(p.channels, uuid)
		delete(p.groups[info.Type], info.Key)
	}

	for _, info := range states.Update {

		// 通道信息已存在
		if _, ok := p.channels[info.UUID]; ok {
			continue
		}

		// 创建通道组
		group, ok := p.groups[info.Type]
		if !ok {
			group = map[string]Channel{}
			p.groups[info.Type] = group
		}

		// 创建远程节点
		node, ok := p.nodes[info.Node.UUID]
		if !ok {
			node = &mesh{Node: info.Node}
			p.nodes[info.Node.UUID] = node
		}

		// 创建远程通道记录
		channel := newAccess(info, node)
		p.channels[info.UUID] = channel
		p.groups[info.Type][info.Key] = channel
	}
}

func (p *Process) Run(typo, key string, handler Handler) error {

	// 检查参数
	if typo == "" {
		return errors.New("empty type for thread")
	}
	if key == "" {
		return errors.New("empty key for thread")
	}
	if handler == nil {
		return errors.New("nil handler for thread")
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 创建通道组
	group, ok := p.groups[typo]
	if !ok {
		group = map[string]Channel{}
		p.groups[typo] = group
	}

	// 冲突的通道名
	if _, ok := group[key]; ok {
		return fmt.Errorf("conflict channel: %s.%s", typo, key)
	}

	// 创建通道
	thread := &thread{
		Channel: &proto.Channel{
			UUID: uuid.NewV1().String(),
			Type: typo,
			Key:  key,
			Node: p.node.Info(),
		},
		Peer:    p,
		Queue:   queue.New(128),
		Handler: handler,
	}

	// 写入通道记录
	group[thread.Key] = thread
	p.channels[thread.UUID] = thread

	// 启动携程
	go thread.run()

	return nil
}

func (p *Process) direct(typo, key string, msg *proto.Message) error {

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	group, ok := p.groups[typo]
	if !ok {
		return fmt.Errorf("can not find type: %s", typo)
	}

	channel, ok := group[key]
	if !ok {
		return fmt.Errorf("can not find channel: %s.%s", typo, key)
	}

	return channel.Push(msg)
}

func (p *Process) multicast(typo string, keys []string, message *proto.Message) error {

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	group, ok := p.groups[typo]
	if !ok || len(group) == 0 {
		return fmt.Errorf("%w: type[%s]", ErrNotFound, typo)
	}

	for _, key := range keys {

		channel, ok := group[key]
		if !ok {
			continue
		}

		// 本地通道直接投递
		if channel.Local() {
			msg := message.Copy()
			msg.Receiver = channel.Info()
			channel.Push(msg)
			continue
		}

	}

	return nil
}

func (p *Process) broadcast(typo string, message *proto.Message) error {

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	group, ok := p.groups[typo]
	if !ok || len(group) == 0 {
		return fmt.Errorf("can not find type: %s", typo)
	}

	for _, channel := range group {
		msg := message.Copy()
		msg.Receiver = channel.Info()
		channel.Push(msg)
	}

	return nil
}

func (p *Process) balance(typo string, message *proto.Message) error {

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	group, ok := p.groups[typo]
	if !ok || len(group) == 0 {
		return fmt.Errorf("%w: type[%s]", ErrNotFound, typo)
	}

	index := 0
	slice := make([]Channel, len(group))
	for _, channel := range group {
		slice[index] = channel
		index++
	}

	channel := slice[rand.Intn(len(slice))]
	msg := message.Copy()
	msg.Receiver = channel.Info()

	return channel.Push(msg)
}

var (
	ErrNotFound = errors.New("can not find channel")
)
