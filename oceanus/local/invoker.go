package local

import "github.com/laconiz/eros/oceanus/proto"

func NewInvoker(v interface{}) *Invoker {
	return &Invoker{}
}

type Invoker struct {
}

func (i *Invoker) init() {

}

func (i *Invoker) onMail(m *proto.Mail) {

}

func (i *Invoker) destroy() {

}
