package router

import "github.com/laconiz/eros/oceanus/proto"

type Mesh interface {
	Info() *proto.Mesh
	Push(*proto.Mail) error
	State() (*proto.State, bool)
}

type Node interface {
	Info() *proto.Node
	Mesh() Mesh
	Push(*proto.Mail) error
}

type Load struct {
	node Node
	load int
}

type Loads []*Load

func (loads Loads) Len() int {
	return len(loads)
}

func (loads Loads) Less(i, j int) bool {
	return loads[i].load < loads[j].load
}

func (loads Loads) Swap(i, j int) {
	loads[i], loads[j] = loads[j], loads[i]
}

func (loads Loads) Nodes() []Node {
	var list []Node
	for _, load := range loads {
		list = append(list, load.node)
	}
	return list
}
