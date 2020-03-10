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
