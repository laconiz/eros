package oceanus

import (
	"github.com/laconiz/eros/codec"
	"github.com/laconiz/eros/network"
)

// 网络节点数据
type NodeListMsg struct {
	Nodes []NodeInfo
}

// 节点加入网络消息
type NodeJoinMsg struct {
	Node NodeInfo
}

// 节点退出网络消息
type NodeExitMsg struct {
	Node NodeInfo
}

// 通道加入列表
type ChannelJoinMsg struct {
	Channels ChannelInfo
}

// 通道退出
type ChannelQuitMsg struct {
	Channel ChannelInfo
}

func init() {
	network.RegisterMeta(NodeListMsg{}, codec.Json())
	network.RegisterMeta(NodeJoinMsg{}, codec.Json())
	network.RegisterMeta(NodeExitMsg{}, codec.Json())
	network.RegisterMeta(ChannelJoinMsg{}, codec.Json())
	network.RegisterMeta(ChannelQuitMsg{}, codec.Json())
}
