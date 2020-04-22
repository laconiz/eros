package proto

import (
	"github.com/laconiz/eros/network/message"
)

type MailID string

type Mail struct {
	ID   MailID   `json:"id"`
	From []*Node  `json:"from,omitempty"`
	Type NodeType `json:"type,omitempty"` // 广播消息节点类型
	To   []*Node  `json:"to,omitempty"`   // 接收节点
	User int64    `json:"user,omitempty"`
	Body []byte   `json:"body"`
}

func init() {
	message.Register(Mail{}, message.JsonCodec())
}
