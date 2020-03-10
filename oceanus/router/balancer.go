package router

import (
	"github.com/laconiz/eros/oceanus/proto"
)

func newBalancer() *Balancer {
	return &Balancer{nodes: map[proto.NodeKey]Node{}}
}

type Balancer struct {
	expired  bool                   // 是否过期
	nodes    map[proto.NodeKey]Node // 节点列表
	balances []Node                 // 均衡列表
}

// 插入节点
func (b *Balancer) Insert(node Node) {
	b.Expired()
	b.nodes[node.Info().Key] = node
}

// 删除节点
func (b *Balancer) Remove(node Node) {
	stored, ok := b.nodes[node.Info().Key]
	if ok && stored.Info().ID == node.Info().ID {
		b.Expired()
		delete(b.nodes, node.Info().Key)
	}
}

// 设置均衡器过期
func (b *Balancer) Expired() {
	b.expired = true
}

// 重新均衡
func (b *Balancer) rebalance() {

}

// 发送消息
func (b *Balancer) Balance(mail *proto.Mail) error {

	if b.expired {
		b.rebalance()
	}
	b.expired = false

	return nil
}
