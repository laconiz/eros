// 远程节点

package oceanus

import "github.com/laconiz/eros/network"

type Burl struct {
	// 节点信息
	node *Node
	// 节点通道
	courses map[string]*Course
	// 连接信息
	session network.Session
}

func (b *Burl) Info() *Node {
	return b.node
}

func (b *Burl) Push(message *Message) error {
	return b.session.Send(message)
}
