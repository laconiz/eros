package proxy

import (
	"github.com/laconiz/eros/network/session"
	"github.com/laconiz/eros/oceanus/proto"
)

type UserProxy struct {
	session session.Session
}

func (proxy *UserProxy) OnMessage(proto.Mail) {
	// proxy.session.SendRaw()
}
