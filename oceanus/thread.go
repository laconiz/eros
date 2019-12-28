package oceanus

import (
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/queue"
)

type Thread interface {
	Reply(msg interface{})
}

type Handler func(*proto.Message)

type thread struct {
	*proto.Channel
	Peer    *Process
	Queue   *queue.Queue
	Handler Handler
}

func (t *thread) Info() *proto.Channel {
	return t.Channel
}

func (t *thread) Node() Node {
	return t.Peer.Node
}

func (t *thread) Local() bool {
	return true
}

func (t *thread) Push(message *proto.Message) error {
	return t.Queue.Add(message)
}

func (t *thread) run() {

	var exited bool
	var messages []interface{}

	for !exited {
		messages, exited = t.Queue.Pick()
		for _, message := range messages {
			t.Handler(message.(*proto.Message))
		}
	}
}

func newThread(info *proto.Channel, process *Process) *thread {
	return nil
}
