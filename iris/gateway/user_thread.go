package gateway

import (
	"github.com/laconiz/eros/oceanus"
)

type UserThread struct {
	gateway *Gateway
}

func (t *UserThread) OnStart() {

}

func (t *UserThread) OnMessage(message *oceanus.Message) {
	t.gateway.Send(message)
}

func (t *UserThread) OnStop() {

}
