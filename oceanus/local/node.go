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

type Node struct {
	info    *proto.Node
	mesh    *Mesh
	queue   *queue.Queue
	invoker *Invoker
}

func (node *Node) Info() *proto.Node {
	return node.info
}

func (node *Node) Mesh() router.Mesh {
	return node.mesh
}

func (node *Node) Mail(mail *proto.Mail) error {
	return node.queue.Add(mail)
}

func (node *Node) Destroy() {
	node.queue.Close()
}

func (node *Node) run() {

	defer func() {
		if err := recover(); err != nil {
			node.mesh.logger.Data(err).Error("invoke error")
		}
	}()

	node.invoker.init()

	for {

		events, closed := node.queue.Pick()

		for _, event := range events {
			mail := event.(*proto.Mail)
			node.invoker.mail(mail)
		}

		if closed {
			break
		}
	}

	node.invoker.destroy()
}
