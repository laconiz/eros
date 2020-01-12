// 远程节点

package oceanus

import (
	"fmt"
	"github.com/laconiz/eros/network"
)

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

	if b.session == nil {
		return fmt.Errorf("node %v is avaliable", b.node)
	}

	return b.session.Send(message)
}

func (b *Burl) Connected() bool {
	return b.session != nil
}

func (b *Burl) Update(node *Node, session network.Session) {

	b.node = node
	b.session = session

	for _, course := range b.courses {
		course.Expired()
	}
}

func (b *Burl) destroy() {
	for _, course := range b.courses {
		course.destroy()
	}
	b.courses = map[string]*Course{}
}

func NewBurl(node *Node) *Burl {
	return &Burl{
		node:    node,
		courses: map[string]*Course{},
		session: nil,
	}
}
