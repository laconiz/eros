package oceanus

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/network/queue"
)

// ---------------------------------------------------------------------------------------------------------------------

type LMesh struct {
	info    *MeshInfo
	state   *MeshState
	nodes   map[NodeID]*LNode
	types   map[NodeType]int32
	logger  logis.Logger
	process *Process
}

func (lm *LMesh) Info() *MeshInfo {
	return lm.info
}

func (lm *LMesh) State() (*MeshState, bool) {
	return lm.state, true
}

func (lm *LMesh) Nodes() []*NodeInfo {

	var nodes []*NodeInfo
	for _, node := range lm.nodes {
		nodes = append(nodes, node.Info())
	}

	return nodes
}

func (lm *LMesh) Mail(mail *Mail) error {

	for _, ni := range mail.To {

		node, ok := lm.nodes[ni.ID]
		if !ok {
			continue
		}

		node.Mail(mail)
	}

	return nil
}

func (lm *LMesh) Destroy() {

	for id, node := range lm.nodes {
		node.Destroy()
	}

	lm.nodes = map[NodeID]*LNode{}
	lm.types = map[NodeType]int32{}
}

//
func (lm *LMesh) Create(ni NodeInfo, handler interface{}) *LNode {

}

//
func (lm *LMesh) Delete(id NodeID) {

}

// ---------------------------------------------------------------------------------------------------------------------

type LNode struct {
	info    *NodeInfo
	mesh    *LMesh
	queue   *queue.Queue
	invoker *Invoker
}

func (ln *LNode) Info() *NodeInfo {
	return ln.info
}

func (ln *LNode) Mesh() Mesh {
	return ln.mesh
}

func (ln *LNode) Mail(mail *Mail) error {
	return ln.queue.Add(mail)
}

func (ln *LNode) Destroy() {
	ln.queue.Close()
}

func (ln *LNode) run() {

	defer func() {
		if err := recover(); err != nil {
			ln.mesh.logger.Data(err).Error("invoke error")
		}
	}()

	ln.invoker.init()

	for {

		events, closed := ln.queue.Pick()

		for _, event := range events {
			mail := event.(*Mail)
			ln.invoker.mail(mail)
		}

		if closed {
			break
		}
	}

	ln.invoker.destroy()
}

// ---------------------------------------------------------------------------------------------------------------------

type Invoker struct {
}

func (inv *Invoker) init() {

}

func (inv *Invoker) mail(mail *Mail) {

}

func (inv *Invoker) destroy() {

}
