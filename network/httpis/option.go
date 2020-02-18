package httpis

import "github.com/laconiz/eros/network/invoker"

type AcceptorOption struct {
	Name     string
	Addr     string
	Nodes    []*invoker.Node
	Params   []interface{}
	Creators []interface{}
}

func (o *AcceptorOption) parse() {
	if o.Name == "" {
		o.Name = "acceptor"
	}
	if o.Addr == "" {
		o.Addr = "0.0.0.0:8080"
	}
}
