package local

import "github.com/laconiz/eros/oceanus/proto"

func NewInvoker(v interface{}) *Invoker {
	return &Invoker{}
}

type Invoker struct {
}

func (inv *Invoker) init() {

}

func (inv *Invoker) mail(mail *proto.Mail) {

}

func (inv *Invoker) destroy() {

}
