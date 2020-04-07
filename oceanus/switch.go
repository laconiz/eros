package oceanus

import (
	"container/list"
	"fmt"
)

// ---------------------------------------------------------------------------------------------------------------------

type Bus struct {
	typo    Type
	nodes   map[Key]Node
	expired bool
	list *list.List
}

func (b *Bus) RouteByKey(key Key, mail *Mail) error {
	if node, ok := b.nodes[key]; ok {
		return node.Mail(mail)
	}
	return fmt.Errorf("can not find node by key %v on bus %v", key, b.typo)
}

func (b *Bus) RouteByKeys(keys []Key, mail *Mail) error {

}

func (b *Bus) Route(mail *Mail) error {

	if b.expired {

	}
	b.expired = true

	if
}

func (b *Bus) Broadcast(mail *Mail) error {

	groups := map[MeshID]*Group{}

	for _, node := range b.nodes {

		mesh := node.Mesh()

		group := groups[mesh.Info().ID]
		if group == nil {
			group = &Group{mesh: mesh}
			groups[mesh.Info().ID] = group
		}

		group.nodes = append(group.nodes, node)
	}

	for _, group := range groups {

		var receivers []NodeID
		for _, node := range group.nodes {
			receivers = append(receivers, node.Info().ID)
		}

		group.mesh.Mail()
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type Switch struct {
	nodes map[NodeID]Node
	buses map[Type]*Bus
}

// ---------------------------------------------------------------------------------------------------------------------

func (s *Switch) SendByID(id NodeID, mail *Mail) error {

	node, ok := s.nodes[id]
	if !ok {
		return fmt.Errorf("can not find node by id %v", id)
	}

	return node.Mail(mail)
}

func (s *Switch) SendByIDs(list []NodeID, msg interface{}) error {

	groups := map[MeshID]*Group{}

	for _, id := range list {

		node, ok := s.nodes[id]
		if !ok {
			continue
		}

		mesh := node.Mesh()

		group := groups[mesh.Info().ID]
		if group == nil {
			group = &Group{mesh: mesh}
			groups[mesh.Info().ID] = group
		}

		group.nodes = append(group.nodes, node)
	}

}

func (s *Switch) RouteByKey(typo Type, key Key, mail *Mail) error {

	bus, ok := s.buses[typo]
	if !ok {
		return fmt.Errorf("can not find bus by type %v", typo)
	}

	return bus.RouteByKey(key, mail)
}

func (s *Switch) RouteByKeys(typo Type, keys []Key, mail *Mail) error {

	bus, ok := s.buses[typo]
	if !ok {
		return fmt.Errorf("can not find bus by type %v", typo)
	}

	return bus.RouteByKeys(keys, mail)
}

func (s *Switch) Route(typo Type, mail *Mail) error {

	bus, ok := s.buses[typo]
	if !ok {
		return fmt.Errorf("can not find bus by type %v", typo)
	}

	return bus.Route(mail)
}

func (s *Switch) Broadcast(typo Type, mail *Mail) error {

	bus, ok := s.buses[typo]
	if !ok {
		return fmt.Errorf("can not find bus by type %v", typo)
	}

	return bus.Broadcast(mail)
}

// ---------------------------------------------------------------------------------------------------------------------

type Group struct {
	mesh  Mesh
	nodes []Node
}
