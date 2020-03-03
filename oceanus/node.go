package oceanus

import "github.com/laconiz/eros/oceanus/proto"

type Node interface {
	Info() *proto.Node
	Mesh() Mesh
	Push(message *proto.Mail) error
}
