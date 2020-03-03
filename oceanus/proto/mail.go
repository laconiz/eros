package proto

import "github.com/laconiz/eros/network/message"

type MailID string

// 路由消息
type Mail struct {
	// 消息ID
	ID MailID
	// 消息类型
	Type NodeType
	// 消息调用节点
	// 展示完整的消息调用过程
	// RPC消息回调
	Senders []Node
	// 消息接收节点
	Receivers []Node
	// 消息体
	Body []byte
}

func init() {
	message.Register(Mail{}, message.JsonCodec())
}
