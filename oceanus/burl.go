// 远程节点

package oceanus

import "github.com/laconiz/eros/network"

type Burl struct {
	// 节点信息
	node *Node
	// 节点通道
	courses map[string]*Course
	// 节点连接
	conn network.Connector
	// 连接状态
	connected bool
}

func (b *Burl) Info() *Node {
	return b.node
}

func (b *Burl) Push(message *Message) error {
	return b.conn.Send(message)
}

func (b *Burl) Connected() bool {
	return b.connected
}
