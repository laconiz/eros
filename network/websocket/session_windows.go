package websocket

import "github.com/laconiz/eros/network"

func (ses *Session) invoke(e *network.Event) {
	ses.config.Invoker.Invoke(e)
}
