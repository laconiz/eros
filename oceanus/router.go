package oceanus

import (
	"container/list"
	"fmt"
)

// ---------------------------------------------------------------------------------------------------------------------

type Mesh interface {
	Info() *MeshInfo
	State() (*MeshState, bool)
	Mail(*Mail) error
}

type Node interface {
	Info() *NodeInfo
	Mesh() Mesh
	Mail(*Mail) error
}

// ---------------------------------------------------------------------------------------------------------------------

type Bus struct {
	typo    NodeType
	nodes   map[NodeKey]Node
	expired bool
	list    *list.List
}

func (bus *Bus) RouteByKey(key NodeKey, mail *Mail) error {

	if node, ok := bus.nodes[key]; ok {
		return node.Mail(mail)
	}

	return fmt.Errorf("can not find node by key %v on bus %v", key, bus.typo)
}

func (bus *Bus) RouteByKeys(keys []NodeKey, mail *Mail) error {

	hubs := Hubs{}

	for _, key := range keys {
		if node, ok := bus.nodes[key]; ok {
			hubs.Trunk(node)
		}
	}

	hubs.Mail(mail)
	return nil
}

func (bus *Bus) Route(mail *Mail) error {

	if bus.expired {

	}
	bus.expired = true

	if
}

func (bus *Bus) Broadcast(mail *Mail) error {

	hubs := Hubs{}

	for _, node := range bus.nodes {
		hubs.Trunk(node)
	}

	hubs.Mail(mail)
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

type Router struct {
	nodes map[NodeID]Node
	buses map[NodeType]*Bus
}

// ---------------------------------------------------------------------------------------------------------------------

func (rt *Router) SendByID(id NodeID, mail *Mail) error {

	node, ok := rt.nodes[id]
	if !ok {
		return fmt.Errorf("can not find node by id %v", id)
	}

	return node.Mail(mail)
}

func (rt *Router) SendByIDs(list []NodeID, msg interface{}) error {

	hubs := Hubs{}

	for _, id := range list {

		node, ok := rt.nodes[id]
		if !ok {
			continue
		}

		hubs.Trunk(node)
	}
}

func (rt *Router) RouteByKey(typo NodeType, key NodeKey, mail *Mail) error {

	bus, ok := rt.buses[typo]
	if !ok {
		return fmt.Errorf("can not find bus by type %v", typo)
	}

	return bus.RouteByKey(key, mail)
}

func (rt *Router) RouteByKeys(typo NodeType, keys []NodeKey, mail *Mail) error {

	bus, ok := rt.buses[typo]
	if !ok {
		return fmt.Errorf("can not find bus by type %v", typo)
	}

	return bus.RouteByKeys(keys, mail)
}

func (rt *Router) Route(typo NodeType, mail *Mail) error {

	bus, ok := rt.buses[typo]
	if !ok {
		return fmt.Errorf("can not find bus by type %v", typo)
	}

	return bus.Route(mail)
}

func (rt *Router) Broadcast(typo NodeType, mail *Mail) error {

	bus, ok := rt.buses[typo]
	if !ok {
		return fmt.Errorf("can not find bus by type %v", typo)
	}

	return bus.Broadcast(mail)
}

// ---------------------------------------------------------------------------------------------------------------------

type Hub struct {
	mesh Mesh
	nodes []*NodeInfo
}

type Hubs map[MeshID]*Hub

func (hubs Hubs) Trunk(node Node) {

	mesh := node.Mesh()

	hub, ok := hubs[mesh.Info().ID]
	if !ok {
		hub = &Hub{mesh: mesh}
		hubs[mesh.Info().ID] = hub
	}

	hub.nodes = append(hub.nodes, node.Info())
}

func (hubs Hubs) Mail(mail *Mail) {

	for _, hub := range hubs {
		m := *mail
		m.To = hub.nodes
		hub.mesh.Mail(&m)
	}
}