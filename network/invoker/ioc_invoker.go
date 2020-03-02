package invoker

import (
	"fmt"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/message"
	"github.com/laconiz/eros/network/session"
	"github.com/laconiz/eros/utils/ioc"
	"reflect"
)

type Invoker interface {
	Invoke(event *network.Event)
}

// ---------------------------------------------------------------------------------------------------------------------

func NewSocketIOCInvoker(log logis.Logger) *SocketIOCInvoker {
	return &SocketIOCInvoker{squirt: ioc.New(), log: log, handlers: map[message.ID][]network.HandlerFunc{}}
}

type SocketIOCInvoker struct {
	squirt   *ioc.Squirt
	log      logis.Logger
	handlers map[message.ID][]network.HandlerFunc
}

func (i *SocketIOCInvoker) Params(params ...interface{}) *SocketIOCInvoker {
	i.squirt.Params(params...)
	return i
}

func (i *SocketIOCInvoker) Creators(creators ...interface{}) *SocketIOCInvoker {
	i.squirt.Creators(creators...)
	return i
}

func (i *SocketIOCInvoker) Register(handlers ...interface{}) error {

	for _, handler := range handlers {

		args, err := i.squirt.UnknownArgs(handler, i.dynamicParams()...)
		if err != nil {
			return err
		}

		if len(args) > 1 {
			return fmt.Errorf("too many args in handler: %v", args)
		}

		if len(args) == 1 {
			return i.RegisterEx(args[0], handler)
		}
	}

	return nil
}

func (i *SocketIOCInvoker) RegisterEx(typo reflect.Type, handler interface{}) error {

	meta, ok := message.MetaByType(typo)
	if !ok {
		return fmt.Errorf("invalid message type: %v", typo)
	}

	squirt := i.squirt.Copy().Creator(typo, func(event *network.Event) (interface{}, error) {
		return event.Msg, nil
	})

	invoker, err := squirt.Handle(handler, i.dynamicParams()...)
	if err != nil {
		return err
	}

	i.handlers[meta.ID()] = append(i.handlers[meta.ID()], func(event *network.Event) {
		invoker(event, event.Ses)
	})
	return nil
}

func (i *SocketIOCInvoker) dynamicParams() []interface{} {
	return []interface{}{&network.Event{}, (session.Session)(&session.EmptySession{})}
}
