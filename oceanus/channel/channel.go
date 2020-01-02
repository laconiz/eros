package channel

import (
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/oceanus/node"
)

type ID string

type Type string

type Key string

type Info struct {
	ID   ID
	Type Type
	Key  Key
	Node node.Info
}

type Channel interface {
	Info() *Info
	Node() node.Node
	Push(*oceanus.Message) error
}
