package invoker

import "github.com/laconiz/eros/network"

type Invoker interface {
	Invoke(event *network.Event)
}
