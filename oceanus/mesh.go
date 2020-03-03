package oceanus

import "github.com/laconiz/eros/oceanus/proto"

type Mesh interface {
	Info() (info *proto.Mesh)
	Push(message *proto.Mail) error
	State() (state *proto.State, valid bool)
}
