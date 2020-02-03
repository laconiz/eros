package steropes

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/laconiz/eros/hyperion"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/epimetheus"
	"github.com/laconiz/eros/utils/ioc"
)

const module = "steropes"

func NewAcceptor(option AcceptorOption) (*Acceptor, error) {
	option.make()
	logger := hyperion.NewEntry(module).WithField(epimetheus.FieldName, option.Name)
	squirt := ioc.New().Params(option.Params...).Functions(option.Functions...)
	engine := gin.New()
	engine.Use(gin.Recovery())
	if err := handleNode(engine, option.Node, squirt, logger); err != nil {
		return nil, err
	}
	return &Acceptor{
		state:  network.Stopped,
		server: &http.Server{Addr: option.Addr, Handler: engine},
		logger: logger,
	}, nil
}

type Acceptor struct {
	state  network.State
	server *http.Server
	logger *hyperion.Entry
	mutex  sync.RWMutex
}

func (a *Acceptor) Run() {
	server := a.run()
	if server == nil {
		return
	}
	a.logger.Infof(epimetheus.ContentStarted)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		a.logger.WithError(err).Error(epimetheus.ContentStopped)
	} else {
		a.logger.Info(epimetheus.ContentStopped)
	}
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if server == a.server {
		a.state = network.Stopped
	}
}

func (a *Acceptor) run() *http.Server {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.state != network.Stopped {
		return nil
	}
	a.state = network.Running
	a.server = &http.Server{Addr: a.server.Addr, Handler: a.server.Handler}
	return a.server
}

func (a *Acceptor) Stop() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.state != network.Running {
		return
	}
	a.state = network.Closing
	a.logger.Info(epimetheus.ContentStopping)
	if err := a.server.Shutdown(context.Background()); err != nil {
		a.logger.WithError(err).Error(epimetheus.ContentStopError)
	}
}

func (a *Acceptor) State() network.State {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.state
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}
