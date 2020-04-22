package process

func newInvoker(handler interface{}) (*Invoker, error) {

	type initializer interface {
		OnInit()
	}
	if _, ok := handler.(initializer); ok {

	}

}

type Invoker struct {
	fn
}

func (i *Invoker) init() {

}

func (i *Invoker) onMail() {

}

func (i *Invoker) Destroy() {

}
