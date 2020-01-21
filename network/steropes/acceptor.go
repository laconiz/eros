package steropes

import (
	"net/http"
	"sync"

	"github.com/laconiz/eros/hyperion"
	"github.com/laconiz/eros/network"
)

const (
	Module    = "steropes"
	FieldName = "name"
	FieldAddr = "addr"
)

func NewAcceptor(option AcceptorOption) *Acceptor {

	logger := hyperion.NewEntry(Module).
		WithField(FieldName, option.Name).
		WithField(FieldAddr, option.Addr)

	acceptor := &Acceptor{
		state:  network.Stopped,
		server: nil,
		logger: logger,
	}

	return acceptor
}

type Acceptor struct {
	state  network.State
	server *http.Server
	logger *hyperion.Entry
	mutex  sync.Mutex
}

func (a *Acceptor) Run() {

	if !func() bool {

		a.mutex.Lock()
		defer a.mutex.Unlock()

		if a.state != network.Stopped {
			return false
		}

		return false

	}() {
		return
	}

	a.server = &http.Server{Addr: a.server.Addr, Handler: a.server.Handler}

	a.logger.WithField(FieldAddr, a.server.Addr).Infof("start listening")
	err := a.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		a.logger.WithError(err).Error("stopped")
	} else {
		a.logger.Info("stopped")
	}
}

func (a *Acceptor) Stop() {

}
