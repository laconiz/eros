package network

type Invoker interface {
	Invoke(*Event)
}

type defaultInvoker struct {
}

func (inv *defaultInvoker) Invoke(e *Event) {

}

var DefaultInvoker = &defaultInvoker{}
