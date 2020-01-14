package oceanus

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
	//
	node *LocalNode
	//
	process *process
}

func (m *Message) Reply(msg interface{}) {

}

func (m *Message) Create() {

}
