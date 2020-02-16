package invoker

import (
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/utils/ioc"
)

type Option struct {
	Params   []interface{}
	Creators []interface{}
}

type Handler func(*network.Event)

type Invoker interface {
	Invoke(event *network.Event)
}

func NewIOCInvoker() *IOCInvoker {
	return &IOCInvoker{ioc.New().Function()}
}

type IOCInvoker struct {
	squirt *ioc.Squirt
}

func (i *IOCInvoker) Invoke(event *network.Event) {

}

func (i *IOCInvoker) Register() {

}

func (i *IOCInvoker) RegisterEx(msg interface{}, handler interface{}) {

}

func (i *IOCInvoker) Parse(handler interface{}) (Handler, error) {

	return nil, nil
}
