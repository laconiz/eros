package httpis

import (
	"context"
	"github.com/laconiz/eros/network/invoker"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/network"
)

const module = "httpis"

func NewAcceptor(opt AcceptorOption, log logis.Logger) (*Acceptor, error) {

	opt.parse()

	engine := gin.New()
	engine.Use(gin.Recovery())

	log = log.Fields(logis.Fields{logis.Module: module, network.FieldName: opt.Name})

	invoker := invoker.NewGinInvoker(log).Params(opt.Params...).Creators(opt.Creators...)
	if err := invoker.Register(engine, opt.Nodes); err != nil {
		return nil, err
	}

	server := &http.Server{Addr: opt.Addr, Handler: engine}
	return &Acceptor{state: network.Stopped, server: server, log: log}, nil
}

type Acceptor struct {
	state  network.State
	server *http.Server
	log    logis.Logger
	mutex  sync.RWMutex
}

func (acceptor *Acceptor) Run() {

	acceptor.mutex.Lock()
	defer acceptor.mutex.Unlock()

	if acceptor.state != network.Stopped {
		return
	}
	acceptor.state = network.Running

	server := &http.Server{
		Addr:    acceptor.server.Addr,
		Handler: acceptor.server.Handler,
	}
	acceptor.server = server

	go func() {

		acceptor.log.Infof("listening at %s", server.Addr)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			acceptor.log.Errorf("listen error: %v", err)
		} else {
			acceptor.log.Info("stopped")
		}

		acceptor.mutex.Lock()
		defer acceptor.mutex.Unlock()

		if server == acceptor.server {
			acceptor.state = network.Stopped
		}
	}()
}

func (acceptor *Acceptor) Stop() {

	acceptor.mutex.Lock()
	defer acceptor.mutex.Unlock()

	if acceptor.state != network.Running {
		return
	}

	err := acceptor.server.Shutdown(context.Background())
	if err != nil {
		acceptor.log.Errorf("shutdown error: %v", err)
	}
}

func (acceptor *Acceptor) State() network.State {
	acceptor.mutex.RLock()
	defer acceptor.mutex.RUnlock()
	return acceptor.state
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}
