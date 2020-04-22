package local

import (
	"errors"
	"fmt"
	"github.com/laconiz/eros/oceanus/proto"
	"github.com/laconiz/eros/oceanus/router"
)

// ---------------------------------------------------------------------------------------------------------------------

type Progress interface {
}

// ---------------------------------------------------------------------------------------------------------------------

func NewMesh(info *proto.Mesh, state *proto.State, router *router.Router) *Mesh {
	return &Mesh{
		info:   info,
		state:  state,
		nodes:  map[proto.NodeID]*Node{},
		types:  map[proto.NodeType]int64{},
		router: router,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type Mesh struct {
	info   *proto.Mesh
	state  *proto.State
	nodes  map[proto.NodeID]*Node
	types  map[proto.NodeType]int64
	router *router.Router
}

// ---------------------------------------------------------------------------------------------------------------------

func (mesh *Mesh) Info() *proto.Mesh {
	return mesh.info
}

func (mesh *Mesh) State() (*proto.State, bool) {
	return mesh.state, true
}

func (mesh *Mesh) Push(mail *proto.Mail) error {
	for _, receiver := range mail.Receivers {
		if node, ok := mesh.nodes[receiver.ID]; ok {
			node.Mail(mail)
		}
	}
	return nil
}

func (mesh *Mesh) Nodes() []*proto.Node {
	var nodes []*proto.Node
	for _, node := range mesh.nodes {
		nodes = append(nodes, node.Info())
	}
	return nodes
}

// ---------------------------------------------------------------------------------------------------------------------

func (mesh *Mesh) UpdateState(state *proto.State) {
	mesh.state = state
	mesh.Expired()
}

func (mesh *Mesh) Expired() {
	for typo, count := range mesh.types {
		if count > 0 {
			mesh.router.Expired(typo)
		}
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func (mesh *Mesh) Create(pn *proto.Node, h interface{}) *Node {

	if _, ok := mesh.nodes[pn.ID]; ok {
		return nil
	}

	node := NewNode(pn, mesh, h)
	mesh.nodes[pn.ID] = node

	mesh.router.Insert(node)
	mesh.types[pn.Type]++

	return node
}

func (mesh *Mesh) Delete(id proto.NodeID) {

	node, ok := mesh.nodes[id]
	if !ok {
		return
	}

	node.Destroy()
	delete(mesh.nodes, id)

	mesh.types[node.info.Type]--
	mesh.router.Remove([]*proto.Node{node.info})
}

// ---------------------------------------------------------------------------------------------------------------------

func (mesh *Mesh) Destroy() {

	for _, node := range mesh.nodes {
		node.Destroy()
		mesh.router.Remove(node.info)
	}

	mesh.nodes = map[proto.NodeID]*Node{}
	mesh.types = map[proto.NodeType]int64{}
}
