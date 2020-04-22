package router

import (
	"container/list"
	"fmt"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/oceanus/proto"
	"sort"
)

func newBalancer() *Bus {
	return &Bus{nodes: map[proto.NodeKey]Node{}}
}

type Bus struct {
	expired  bool
	nodes    map[proto.NodeKey]Node
	elements list.List
	logger   logis.Logger
}

// 插入节点
func (balancer *Bus) Insert(node Node) {
	balancer.Expired()
	balancer.nodes[node.Info().Key] = node
}

// 删除节点
func (balancer *Bus) Remove(node Node) {
	stored, ok := balancer.nodes[node.Info().Key]
	if ok && stored.Info().ID == node.Info().ID {
		balancer.Expired()
		delete(balancer.nodes, node.Info().Key)
	}
}

// 设置均衡器过期
func (balancer *Bus) Expired() {
	balancer.expired = true
}

// 重新均衡
func (balancer *Bus) rebalance() {

	loads := Loads{}

	for _, node := range balancer.nodes {

		state, ok := node.Mesh().State()
		if !ok {
			continue
		}

	}

	sort.Sort(loads)

	balancer.balances = loads.Nodes()
}

// 发送消息
func (balancer *Bus) Balance(mail *proto.Mail) error {

	if balancer.expired {
		balancer.rebalance()
	}
	balancer.expired = false

	return nil
}

// 随机发送消息
func (balancer *Bus) Random(raw *proto.Mail) error {

	node, ok := balancer.nodes[raw.Type]
	if !ok {
		return nil
	}

	mail.Copy()
}

//
func (balancer *Bus) Key(key proto.NodeKey, mail *proto.Mail) error {

	node, ok := balancer.nodes[key]
	if !ok {
		return fmt.Errorf("cannot find node by key %v", key)
	}

	mail.Receivers = []*proto.Node{node.Info()}
	return node.Push(mail)
}

//
func (balancer *Bus) Broadcast(origin *proto.Mail) error {

	group := map[Mesh][]Node{}
	for _, node := range balancer.nodes {
		mesh := node.Mesh()
		group[mesh] = append(group[mesh], node)
	}

	for mesh, nodes := range group {

		var receivers []*proto.Node
		for _, node := range nodes {
			receivers = append(receivers, node.Info())
		}
	}

	return nil
}

//
