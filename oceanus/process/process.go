package process

import (
	"errors"
	"fmt"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/oceanus/channel"
	"github.com/laconiz/eros/oceanus/node"
	"github.com/laconiz/eros/queue"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"sync"
)

type Process struct {

	// 本地节点信息
	info *node.Info

	// 节点信息
	nodes map[node.ID]node.Node
	// 通道信息
	channels map[channel.ID]channel.Channel
	// 通道分组
	groups map[channel.Type]map[channel.Key]channel.Channel
	// 版本控制通道
	balances map[channel.Type]*balance

	// 监听器
	acceptor network.Acceptor

	//
	mutex sync.RWMutex

	//
	rand rand.Rand
}

// 本地节点信息
func (p *Process) Info() *node.Info {
	return p.info
}

// 本地节点数据推送
func (p *Process) Push(message *oceanus.Message) error {

	// 本地节点之间的推送会占用锁
	go func() {

		p.mutex.RLock()
		defer p.mutex.RUnlock()

		for _, id := range message.Receivers {
			if c, ok := p.channels[id]; ok {
				c.Push(message)
			}
		}
	}()

	return nil
}

// 同步通道模型
func (p *Process) Sync(sync *oceanus.SyncMessage) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, uuid := range sync.Deleted {

		// 查找通道信息
		channel, ok := p.channels[uuid]
		if !ok {
			continue
		}

		// 删除通道信息
		info := channel.Info()
		delete(p.channels, uuid)
		delete(p.groups[info.Type], info.Key)

		// 更新负载均衡通道
		if b, ok := p.balances[info.Type]; ok {
			b.Replace(p.groups[info.Type])
		}
	}

	for _, info := range sync.Updated {

		// 通道信息已存在
		if _, ok := p.channels[info.ID]; ok {
			continue
		}

		// 创建通道组
		group, ok := p.groups[info.Type]
		if !ok {
			group = map[channel.Key]channel.Channel{}
			p.groups[info.Type] = group
		}

		// 创建节点
		n, ok := p.nodes[info.Node.ID]
		if !ok {
			n = node.NewRemote(&info.Node)
			p.nodes[info.Node.ID] = n
		}

		// 创建通道记录
		channel := channel.NewRemote(info, n)
		p.channels[info.ID] = channel
		p.groups[info.Type][info.Key] = channel

		// 更新负载均衡通道
		if b, ok := p.balances[info.Type]; ok {
			b.Replace(p.groups[info.Type])
		}
	}
}

func (p *Process) NewThread(typo, key string, handler oceanus.Handler) error {

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
		group = map[string]oceanus.Channel{}
		p.groups[typo] = group
	}

	// 冲突的通道名
	if _, ok := group[key]; ok {
		return fmt.Errorf("conflict channel: %s.%s", typo, key)
	}

	// 创建通道
	thread := &oceanus.thread{
		Channel: &proto.Channel{
			UUID: uuid.NewV1().String(),
			Type: typo,
			Key:  key,
			Node: p.node,
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
