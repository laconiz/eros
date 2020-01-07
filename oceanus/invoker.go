package oceanus

type Invoker interface {
	OnMessage(*Message)
}
