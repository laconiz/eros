package local

import (
	"github.com/laconiz/eros/network/queue"
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/router"
)

// ---------------------------------------------------------------------------------------------------------------------

const queueLen = 64

func NewNode(pn *proto.Node, m *Mesh, h interface{}) *Node {
	return &Node{
		info:    pn,
		mesh:    m,
		invoker: NewInvoker(h),
		queue:   queue.New(queueLen),
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type Node struct {
	info    *proto.Node  // 节点信息
	mesh    *Mesh        // 所属网格
	invoker *Invoker     // 调用器
	queue   *queue.Queue // 消息队列
}

// ---------------------------------------------------------------------------------------------------------------------

func (n *Node) Info() *proto.Node {
	return n.info
}

func (n *Node) Mesh() router.Mesh {
	return n.mesh
}

func (n *Node) Mail(mail *proto.Mail) error {
	return n.queue.Add(mail)
}

// ---------------------------------------------------------------------------------------------------------------------

func (n *Node) Destroy() {
	n.queue.Close()
}

// ---------------------------------------------------------------------------------------------------------------------

func (n *Node) run() {

	n.invoker.init()

	for {

		msgs, closed := n.queue.Pick()

		for _, msg := range msgs {
			mail := msg.(*proto.Mail)
			n.invoker.onMail(mail)
		}

		if closed {
			break
		}
	}

	n.invoker.destroy()
}
