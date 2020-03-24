package httpis

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network"
	"net/http"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------

func NewAcceptor(option *AcceptorOption) (*Acceptor, error) {

	option.parse()

	logger := logisor.Level(option.Level).
		Field(logis.Module, module).
		Field(network.FieldName, option.Name)

	engine := gin.New()
	engine.Use(gin.Recovery())

	invoker := NewInvoker(logger).
		Params(option.Params...).
		Creators(option.Creators...)

	if err := invoker.Register(engine, option.Nodes); err != nil {
		return nil, err
	}

	return &Acceptor{
		state: network.Stopped,
		listener: &http.Server{
			Addr:    option.Addr,
			Handler: engine,
		},
		logger: logger,
	}, nil
}

// ---------------------------------------------------------------------------------------------------------------------

type Acceptor struct {
	state    network.State // 状态
	listener *http.Server  // 侦听器
	logger   logis.Logger  // 日志接口
	mutex    sync.RWMutex
}

// ---------------------------------------------------------------------------------------------------------------------

func (acceptor *Acceptor) State() network.State {

	acceptor.mutex.RLock()
	defer acceptor.mutex.RUnlock()

	return acceptor.state
}

// ---------------------------------------------------------------------------------------------------------------------

func (acceptor *Acceptor) Run() {

	acceptor.mutex.Lock()
	defer acceptor.mutex.Unlock()

	if acceptor.state != network.Stopped {
		return
	}
	acceptor.state = network.Running

	acceptor.listener = &http.Server{
		Addr:    acceptor.listener.Addr,
		Handler: acceptor.listener.Handler,
	}

	acceptor.logger.Data(acceptor.listener.Addr).Info("start")

	go func() {

		err := acceptor.listener.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			acceptor.logger.Err(err).Error("listen error")
		}

		acceptor.mutex.Lock()
		defer acceptor.mutex.Unlock()

		acceptor.state = network.Stopped
		acceptor.logger.Info("stopped")
	}()
}

func (acceptor *Acceptor) Stop() {

	acceptor.mutex.Lock()
	defer acceptor.mutex.Unlock()

	if acceptor.state != network.Running {
		return
	}
	acceptor.state = network.Closing

	acceptor.logger.Info("shutting down")

	context, _ := context.WithTimeout(context.Background(), time.Second*2)
	if err := acceptor.listener.Shutdown(context); err != nil {
		acceptor.logger.Err(err).Error("shutdown error")
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func init() {
	gin.SetMode(gin.ReleaseMode)
}
