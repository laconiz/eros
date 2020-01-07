package oceanus

import "github.com/laconiz/eros/queue"

// 本地线程
type Thread struct {
	channel *Channel
	process *Process
	queue   *queue.Queue
	invoker Invoker
}

func (t *Thread) Info() *Channel {
	return t.channel
}

func (t *Thread) Push(message *Message) error {
	return t.queue.Add(message)
}

func (t *Thread) Run() {

	var closed bool
	var messages []interface{}

	for !closed {
		messages, closed = t.queue.Pick()
		for _, message := range messages {
			t.invoker.OnMessage(message.(*Message))
		}
	}
}
