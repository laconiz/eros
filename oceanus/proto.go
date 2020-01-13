package oceanus

import (
	"github.com/laconiz/eros/codec"
	"github.com/laconiz/eros/network"
)

// 网格ID UUID.V3
// 由指定的命名空间和连接地址生成, 相同的服务地址生成的ID相同
type MeshID string

// 网格状态
// 路由消息时应根据网格状态进行负载均衡
// TODO 当前服务器负载信息
// TODO 当前网格节点统计
type State struct {
	// 版本号
	Version uint64
}

// 网格信息
type MeshInfo struct {
	// 网格ID
	ID MeshID
	// 监听地址
	Addr string
	// 网格状态
	State State
}

// 节点ID UUID.V1
// 生成每个节点时保证每个节点ID都不一致
type NodeID string

// 节点类型
// 用于消息路由
type NodeType string

// 节点主键
// 用于消息路由
// 对于某些分布式等价节点 主键会由UUID.V1生成唯一主键
type NodeKey string

// 节点信息
type NodeInfo struct {
	// 节点ID
	ID NodeID
	// 节点类型
	Type NodeType
	// 节点主键
	Key NodeKey
	// 节点所属网格ID
	Mesh MeshID
}

// 消息ID UUID.V1
// 每个消息都应该自动生成唯一ID
type MessageID string

// 消息类型
// 用于消息处理时逻辑处理
type MsgType uint32

// 路由消息
type Message struct {
	// 消息ID
	ID MessageID
	// 消息调用节点
	// 展示完整的消息调用过程
	// RPC消息回调
	Senders []NodeInfo
	// 消息接收节点
	Receivers []NodeInfo
	// 消息类型
	Type MsgType
	// 消息体
	Body []byte
}

// 网格加入 & 更新消息
// 当底层网络连接成功时自动向其他网格同步
// TODO 定时向其他网格同步当前网格状态
// TODO 同时该消息也附带了心跳功能
// TODO 当网格数量较大且状态信息较多时, 考虑流量
// TODO 当状态并没有实际改变时也会出发其他网格的路由规则更新
// TODO　状态更新不及时
// TODO 是否加入心跳消息
// TODO 智能状态管理 当状态未改变时不触发路由规则更新
type MeshJoin struct {
	*MeshInfo
}

// 网格离线消息
// 当前网格离线时自动向其他网格推送
// TODO 当与其他网格网络连接出现问题时, 网格无法回收, 如何解决?
// TODO 是否可以使用网格列表进行同步判定, 定期回收失效网格
type MeshQuit struct {
	*MeshInfo
}

// 节点加入消息
// 当底层网络连接成功时自动向其他网格同步
// 当新节点生成时自动向其他网格同步
type NodeJoin struct {
	Nodes []*NodeInfo
}

// 节点离线消息
// 当节点离线时自动向其他网格同步
type NodeQuit struct {
	Nodes []*NodeInfo
}

func init() {
	network.RegisterMeta(Message{}, codec.Json())
	network.RegisterMeta(MeshJoin{}, codec.Json())
	network.RegisterMeta(MeshQuit{}, codec.Json())
	network.RegisterMeta(NodeJoin{}, codec.Json())
	network.RegisterMeta(NodeQuit{}, codec.Json())
}
