package oceanus

import (
	"github.com/laconiz/eros/codec"
	"github.com/laconiz/eros/network"
)

type State struct {
	Version uint32
}

type Node struct {
	ID    string
	Addr  string
	State State
}

type Channel struct {
	ID   string
	Type string
	Key  string
	Node string
}

type Message struct {
	// 全局唯一ID
	// 消息追溯
	ID string
	// 发送队列列表
	// 追溯消息来源
	// 路径健康检查
	Sender []*Channel
	// 接收队列
	Receivers []string
	// 消息源数据
	Body []byte
}

// 节点加入网络消息
type NodeJoinMsg struct {
	Node *Node
}

// 节点退出网络消息
type NodeQuitMsg struct {
	Node *Node
}

// 通道加入列表
type ChannelJoinMsg struct {
	Channels []*Channel
}

// 通道退出
type ChannelQuitMsg struct {
	Channels []*Channel
}

func init() {
	network.RegisterMeta(Message{}, codec.Json())
	network.RegisterMeta(NodeJoinMsg{}, codec.Json())
	network.RegisterMeta(NodeQuitMsg{}, codec.Json())
	network.RegisterMeta(ChannelJoinMsg{}, codec.Json())
	network.RegisterMeta(ChannelQuitMsg{}, codec.Json())
}
