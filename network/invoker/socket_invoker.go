package invoker

import (
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/message"
)

func NewSocketInvoker() *SocketInvoker {
	return &SocketInvoker{handlers: map[message.ID]network.HandlerFunc{}}
}

type SocketInvoker struct {
	handlers map[message.ID]network.HandlerFunc
}

func (invoker *SocketInvoker) Register(msg interface{}, handler network.HandlerFunc) {

	meta, ok := message.MetaByMessage(msg)
	if !ok {
		return
	}

	invoker.handlers[meta.ID()] = handler
}

func (invoker *SocketInvoker) Invoke(event *network.Event) {
	if handler, ok := invoker.handlers[event.Meta.ID()]; ok {
		handler(event)
	}
}
