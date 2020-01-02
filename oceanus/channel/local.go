package channel

import (
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/oceanus/node"
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/queue"
)

type Local struct {
	info   *Info
	node   node.Node
	queue  *queue.Queue
	thread oceanus.Thread
}

func (l *Local) Node() node.Node {
	return l.node
}

func (l *Local) Push(message *proto.Message) error {
	return l.queue.Add(message)
}

func (l *Local) Run() {

	var exited bool
	var messages []interface{}

	for !exited {
		messages, exited = l.queue.Pick()
		for _, message := range messages {
			l.thread.OnMessage(message.(*proto.Message))
		}
	}
}

func NewLocal(info *Info, node node.Node, thread oceanus.Thread) *Local {
	return &Local{
		info:   info,
		node:   node,
		queue:  queue.New(64),
		thread: thread,
	}
}
