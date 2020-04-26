package oceanus

import (
	"encoding/hex"
	"github.com/laconiz/eros/oceanus/process"
	"github.com/laconiz/eros/oceanus/proto"
	uuid "github.com/satori/go.uuid"
)

type ID = proto.NodeID
type Type = proto.NodeType
type Key = proto.NodeKey

type Message struct {
	Header  proto.Headers
	Message interface{}
}

type Oceanus interface {

	//
	Create(typo proto.NodeType, key proto.NodeKey) (proto.NodeID, error)
	//
	Destroy(proto.NodeID)

	// 向指定节点发送消息
	SendByID(id proto.NodeID, msg interface{}) error
	// 向指定节点列表发送消息
	SendByIDs([]proto.NodeID, interface{}) error
	// 指定路由规则发送消息
	SendByKey(typo proto.NodeType, key proto.NodeKey, msg interface{}) error
	// 指定路由规则列表发送消息
	SendByKeys([]proto.NodeKey, interface{}) error

	// 负载均衡发送
	Route(typo proto.NodeType, msg interface{}) error
	// 广播消息
	Broadcast(typo proto.NodeType, msg interface{}) error

	// 指定节点ID的RPC调用
	CallByID(id proto.NodeID, req interface{}, resp interface{}) error
	// 指定路由规则的RPC调用
	CallByKey(typo proto.NodeType, key proto.NodeKey, req interface{}, resp interface{}) error
	// 负载均衡的RPC调用
	CallByRoute(typo proto.NodeType, req interface{}, resp interface{}) error
}

var progress = &process.Process{}

func NewUUID() string {
	return hex.EncodeToString(uuid.NewV1().Bytes())
}
