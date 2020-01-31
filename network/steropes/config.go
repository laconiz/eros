package steropes

type AcceptorOption struct {
	Name      string
	Addr      string
	Node      *Node
	Params    []interface{}
	Functions []interface{}
}

func (o *AcceptorOption) make() {
	if o.Name == "" {
		o.Name = "acceptor"
	}
	if o.Addr == "" {
		o.Addr = "0.0.0.0:8080"
	}
	if o.Node == nil {
		o.Node = &Node{}
	}
}
