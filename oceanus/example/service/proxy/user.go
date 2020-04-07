package proxy

import (
	"github.com/laconiz/eros/network/session"
	"github.com/laconiz/eros/oceanus"
)

type UserProxy struct {
	session session.Session
}

func (proxy *UserProxy) OnMessage(oceanus.Mail) {
	// proxy.session.SendRaw()
}
