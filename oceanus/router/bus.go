package router

import (
	"container/list"
)

// 创建路由总线
func NewBus() *Bus {
	return &Bus{
		nodes: map[Key]Node{},
		list:  list.New(),
	}
}

// 路由总线
type Bus struct {
	nodes   map[Key]Node // 节点列表
	expired bool         // 负载均衡队列过期标志
	list    *list.List   // 负载均衡队列
}

// 设置负载均衡过期
func (bus *Bus) Expired() {
	bus.expired = true
}

// 插入节点
func (bus *Bus) Insert(node Node) {
	bus.expired = true
	key := node.Info().Key
	bus.nodes[key] = node
}

// 删除节点
func (bus *Bus) Remove(node Node) {

	info := node.Info()
	id := info.ID
	key := info.Key

	// 查询保存节点
	stored, ok := bus.nodes[key]
	if !ok {
		return
	}

	// 比对节点ID并删除
	if id == stored.Info().ID {
		bus.expired = true
		delete(bus.nodes, key)
	}
}

// 根据KEY查询节点
func (bus *Bus) ByKey(key Key) Node {
	return bus.nodes[key]
}

// 根据KEY列表查询节点列表
func (bus *Bus) ByKeys(keys []Key) []Node {

	var nodes []Node

	for _, key := range keys {

		node, ok := bus.nodes[key]
		if !ok {
			continue
		}

		nodes = append(nodes, node)
	}

	return nodes
}

// 根据负载查询节点
func (bus *Bus) ByLoad() Node {

	// 重新负载均衡
	if bus.expired {
		bus.list = balance(bus.nodes)
		bus.expired = false
	}

	list := bus.list
	// 节点列表为空
	if list.Len() == 0 {
		return nil
	}

	// 删除头节点
	front := list.Front()
	node := list.Remove(front)
	// 移动节点至尾部
	list.PushBack(node)

	return node.(Node)
}

// 查询节点列表
func (bus *Bus) Nodes() []Node {

	var nodes []Node

	for _, node := range bus.nodes {
		nodes = append(nodes, node)
	}

	return nodes
}
