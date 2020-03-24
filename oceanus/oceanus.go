package oceanus

import "github.com/laconiz/eros/oceanus/proto"

type ID = proto.NodeID
type Key = proto.NodeKey
type Type = proto.NodeType

type Oceanus interface {

	// 向指定节点发送消息
	SendByID(id ID, msg interface{}) error
	// 向指定节点列表发送消息
	SendByIDs([]ID, interface{}) error
	// 指定路由规则发送消息
	SendByKey(typo Type, key Key, msg interface{}) error
	// 指定路由规则列表发送消息
	SendByKeys([]Key, interface{}) error

	// 负载均衡发送
	Route(typo Type, msg interface{}) error
	// 广播消息
	Broadcast(typo Type, msg interface{}) error

	// 指定节点ID的RPC调用
	CallByID(id ID, req interface{}, resp interface{}) error
	// 指定路由规则的RPC调用
	CallByKey(typo Type, key Key, req interface{}, resp interface{}) error
	// 负载均衡的RPC调用
	CallByRoute(typo Type, req interface{}, resp interface{}) error
}
