package proto

import "github.com/laconiz/eros/network/message"

type MailID string

type Headers map[string]interface{}

type RpcID MailID

const EmptyRpcID RpcID = ""

type Mail struct {
	ID      MailID   `json:"id"`
	Headers Headers  `json:"header,omitempty"`
	From    []*Node  `json:"from,omitempty"`
	Type    NodeType `json:"type,omitempty"`
	To      []*Node  `json:"to,omitempty"`
	Reply   RpcID    `json:"reply,omitempty"`
	Body    []byte   `json:"body"`
}

func init() {
	message.Register(Mail{}, message.JsonCodec())
}
