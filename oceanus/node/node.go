package node

import "github.com/laconiz/eros/oceanus/routing"

// ---------------------------------------------------------------------------------------------------------------------

type ID string

type Info struct {
	ID      ID                `json:"id"`
	Routing []routing.Routing `json:"routing"`
}

// ---------------------------------------------------------------------------------------------------------------------

type Local struct {
	info *Info
}

func (l *Local) Info() *Info {
	return l.info
}

func (l *Local) Mesh() *oceanus.Mesh {

}
