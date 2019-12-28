package proto

// 节点信息
type Node struct {
	// 唯一ID
	UUID string
	// 节点地址
	Addr string
}

// 通道信息
type Channel struct {
	// 唯一ID
	UUID string
	// 通道类型
	// 消息路由分组
	Type string
	// 逻辑ID
	Key string
	// 所属节点
	Node *Node
}

// channel更新消息
type ChannelStates struct {
	// 更新的channel
	Update []*Channel `json:"omitempty"`
	// 删除的channel
	Delete []string `json:"omitempty"`
}

// 转发消息
type Message struct {
	// 全局唯一ID
	// 消息追溯
	UUID string
	// 发送队列列表
	// 追溯消息来源
	// 路径健康检查
	Sender []*Channel
	// 接收队列
	Receiver *Channel
	// 消息源数据
	Body []byte
}

func (m *Message) Copy() *Message {
	return &Message{
		UUID:     m.UUID,
		Sender:   m.Sender,
		Receiver: m.Receiver,
		Body:     m.Body,
	}
}
