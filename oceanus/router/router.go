package router

import (
	"github.com/laconiz/eros/oceanus/proto"
)

func NewRouter() *Router {
	return &Router{nodes: map[proto.NodeID]Node{}, balancers: map[proto.NodeType]*Balancer{}}
}

type Router struct {
	nodes     map[proto.NodeID]Node
	balancers map[proto.NodeType]*Balancer
	packer    Packer
}

func (router *Router) Insert(node Node) {

	router.nodes[node.Info().ID] = node

	balancer, ok := router.balancers[node.Info().Type]
	if !ok {
		balancer = newBalancer()
		router.balancers[node.Info().Type] = balancer
	}

	balancer.Insert(node)
}

func (router *Router) Remove(list []*proto.Node) {

	for _, info := range list {

		if node, ok := router.nodes[info.ID]; ok {
			delete(router.nodes, info.ID)
			router.balancers[info.Type].Remove(node)
		}
	}
}

func (router *Router) Expired(typo proto.NodeType) {
	if balancer, ok := router.balancers[typo]; ok {
		balancer.Expired()
	}
}

func (router *Router) RouteByID(id proto.NodeID, mail *proto.Mail) {

	node, ok := router.nodes[id]
	if !ok {
		return
	}

	mail = mail.Copy()
	mail.Receivers = []*proto.Node{node.Info()}
	node.Push(mail)
}

func (router *Router) RouteByKey(typo proto.NodeType, key proto.NodeKey, mail *proto.Mail) {

	balancer, ok := router.balancers[typo]
	if !ok {
		return
	}

	balancer.Key(key, mail)
}

func (router *Router) RandByType(typo proto.NodeType, msg interface{}) {

	balancer, ok := router.balancers[typo]
	if !ok {
		return
	}

	balancer.Random()
}
