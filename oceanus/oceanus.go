package oceanus

import "github.com/laconiz/eros/oceanus/process"

type Oceanus interface {

	//
	Create(typo NodeType, key NodeKey) (NodeID, error)
	//
	Destroy(NodeID)

	// 向指定节点发送消息
	SendByID(id NodeID, msg interface{}) error
	// 向指定节点列表发送消息
	SendByIDs([]NodeID, interface{}) error
	// 指定路由规则发送消息
	SendByKey(typo NodeType, key NodeKey, msg interface{}) error
	// 指定路由规则列表发送消息
	SendByKeys([]NodeKey, interface{}) error

	// 负载均衡发送
	Route(typo NodeType, msg interface{}) error
	// 广播消息
	Broadcast(typo NodeType, msg interface{}) error

	// 指定节点ID的RPC调用
	CallByID(id NodeID, req interface{}, resp interface{}) error
	// 指定路由规则的RPC调用
	CallByKey(typo NodeType, key NodeKey, req interface{}, resp interface{}) error
	// 负载均衡的RPC调用
	CallByRoute(typo NodeType, req interface{}, resp interface{}) error
}

var progress = &process.Process{}

func Create(typo NodeType, key NodeKey) (NodeID, error) {
	return NodeID(0), nil
}
