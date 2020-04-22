package oceanus

import "github.com/laconiz/eros/network/queue"

func newThread() {

}

type Thread struct {
	info  *NodeInfo
	mesh  *Progress
	queue *queue.Queue
}

func (t *Thread) Info() *NodeInfo {
	return t.info
}

func (t *Thread) Mesh() Mesh {
	return t.mesh
}

func (t *Thread) Mail(mail *Mail) error {
	return t.queue.Add(mail)
}
