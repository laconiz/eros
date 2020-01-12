package router

import "github.com/laconiz/eros/oceanus/proto"

type Mesh interface {
	Info() *proto.Mesh
	Push(*proto.Message) error
}

type Node interface {
	Info() *proto.Node
	Mesh() Mesh
	Push(*proto.Message) error
}

type Router struct {
	hubs  map[proto.NodeType]*Hub
	nodes map[proto.NodeID]Node
}

func (r *Router) Add(node Node) *Hub {

}

func (r *Router) Remove(id proto.NodeID) {

}

func (r Router) Hub(typo proto.NodeType) *Hub {
	hub, ok := r.hubs[typo]
	if !ok {
		hub := newHub(typo)
		r.hubs[typo] = hub
	}
	return hub
}

func (r *Router) Expired(typo proto.NodeType) {

}

func NewRouter() *Router {
	return &Router{
		nodes: map[proto.NodeID]Node{},
	}
}