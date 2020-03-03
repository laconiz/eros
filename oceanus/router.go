package oceanus

import (
	"github.com/laconiz/eros/oceanus/proto"
)

func NewRouter() *Router {
	return &Router{nodes: map[proto.NodeID]Node{}, balancers: map[proto.NodeType]*Balancer{}}
}

type Router struct {
	nodes     map[proto.NodeID]Node
	balancers map[proto.NodeType]*Balancer
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

func (router *Router) Remove(node Node) {
	if _, ok := router.nodes[node.Info().ID]; ok {
		delete(router.nodes, node.Info().ID)
		if balancer, ok := router.balancers[node.Info().Type]; ok {
			balancer.Remove(node)
		}
	}
}

func (router *Router) Expired(typo proto.NodeType) {
	if balancer, ok := router.balancers[typo]; ok {
		balancer.Expired()
	}
}
