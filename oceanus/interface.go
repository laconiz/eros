package oceanus

import (
	"github.com/laconiz/eros/oceanus/channel"
	"github.com/laconiz/eros/oceanus/node"
)

type MessageID string

type Message struct {
	// 全局唯一ID
	// 消息追溯
	ID MessageID
	// 发送队列列表
	// 追溯消息来源
	// 路径健康检查
	Sender []*channel.Info
	// 接收队列
	Receivers []channel.ID
	// 消息源数据
	Body []byte
}

type NodeSyncMessage struct {
	Deleted []node.ID
	Updated []*node.Info
}

type ChannelSyncMessage struct {
	Deleted []channel.ID
	Updated []*channel.Info
}

type SyncMessage struct {
	Deleted []channel.ID
	Updated []*channel.Info
}
