package oceanus

import (
	"github.com/laconiz/eros/oceanus/proto"
)

type Node interface {
	Info() *proto.Node
	Send(*proto.Message) error
}

type mesh struct {
	*proto.Node
}

func (m *mesh) Info() *proto.Node {
	return m.Node
}

// TODO
func (m *mesh) Send(msg *proto.Message) error {
	return nil
}
