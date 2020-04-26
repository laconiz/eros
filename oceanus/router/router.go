package router

import (
	"github.com/laconiz/eros/oceanus/abstract"
	"github.com/laconiz/eros/oceanus/proto"
)

type ID = proto.NodeID
type Type = proto.NodeType
type Key = proto.NodeKey

type Mesh = abstract.Mesh
type Node = abstract.Node

// 创建路由器
func New() abstract.Router {
	return &Router{
		nodes: map[ID]Node{},
		buses: map[Type]*Bus{},
	}
}

// 路由器
type Router struct {
	nodes map[ID]Node   // 节点列表
	buses map[Type]*Bus // 总线列表
}

// 添加节点
func (rt *Router) Insert(node Node) {

	id := node.Info().ID
	// 删除相同ID节点
	rt.Remove(id)
	// 插入节点列表
	rt.nodes[id] = node

	typo := node.Info().Type
	// 查询总线
	bus, ok := rt.buses[typo]
	// 创建总线
	if !ok {
		bus = NewBus()
		rt.buses[typo] = bus
	}
	// 将节点插入总线
	bus.Insert(node)
}

// 移除节点
func (rt *Router) Remove(id ID) {

	// 查询节点
	node, ok := rt.nodes[id]
	if !ok {
		return
	}

	// 移除节点
	delete(rt.nodes, id)
	typo := node.Info().Type
	// 从总线移除节点
	rt.buses[typo].Remove(node)
}

// 设置总线状态过期
func (rt *Router) Expired(typo Type) {

	// 查询总线
	bus, ok := rt.buses[typo]
	if !ok {
		return
	}

	// 设置总线过期
	bus.Expired()
}

// 根据ID查询节点
func (rt *Router) ByID(id ID) Node {
	return rt.nodes[id]
}

// 根据ID列表查询节点列表
func (rt *Router) ByIDList(list []ID) []Node {

	var nodes []Node

	for _, id := range list {

		node, ok := rt.nodes[id]
		if !ok {
			continue
		}

		nodes = append(nodes, node)
	}

	return nodes
}

// 根据KEY查询节点
func (rt *Router) ByKey(typo Type, key Key) Node {

	// 查询总线
	bus, ok := rt.buses[typo]
	if !ok {
		return nil
	}

	return bus.ByKey(key)
}

// 根据KEY列表查询节点列表
func (rt *Router) ByKeys(typo Type, keys []Key) []Node {

	// 查询总线
	bus, ok := rt.buses[typo]
	if !ok {
		return nil
	}

	return bus.ByKeys(keys)
}

// 根据负载查询节点
func (rt *Router) ByLoad(typo Type) Node {

	// 查询总线
	bus, ok := rt.buses[typo]
	if !ok {
		return nil
	}

	return bus.ByLoad()
}

// 根据TYPE查询节点列表
func (rt *Router) ByType(typo Type) []Node {

	// 查询总线
	bus, ok := rt.buses[typo]
	if !ok {
		return nil
	}

	return bus.Nodes()
}
