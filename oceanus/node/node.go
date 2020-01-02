package node

import "github.com/laconiz/eros/oceanus/proto"

type ID string

type State struct {
	Version uint32
}

type Info struct {
	ID    ID
	Addr  string
	State *State
}

type Node interface {
	Info() *Info
	Send(*proto.Message) error
	Stop()
}
