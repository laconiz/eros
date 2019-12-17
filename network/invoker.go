package network

import "fmt"

type Invoker interface {
	Invoke(*Event)
}

type StdInvoker struct {
	handlers map[MessageID]func(*Event)
}

func (i *StdInvoker) Invoke(e *Event) {
	if handler, ok := i.handlers[e.Meta.ID()]; ok {
		handler(e)
	}
}

func (i *StdInvoker) Register(msg interface{}, handler func(*Event)) error {

	meta := MetaByMsg(msg)
	if meta == nil {
		return fmt.Errorf("invalid message: %#v", msg)
	}

	if _, ok := i.handlers[meta.ID()]; ok {
		return fmt.Errorf("conflict meta: %v", meta)
	}

	i.handlers[meta.ID()] = handler
	return nil
}

func NewStdInvoker() *StdInvoker {
	return &StdInvoker{handlers: map[MessageID]func(*Event){}}
}
