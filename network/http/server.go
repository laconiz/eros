package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	"net/http"
	"sync"
	"time"
)

type Acceptor struct {
	state  network.State // 服务器状态
	server *http.Server  // HTTP服务
	log    *tlog.Log     // 日志
	mutex  sync.Mutex
}

// 启动服务器
func (acc *Acceptor) Start() {

	acc.mutex.Lock()
	defer acc.mutex.Unlock()

	// 修改状态
	if acc.state != network.Stopped {
		return
	}
	acc.state = network.Running

	// 重置server解决shutdown无法重新listen的问题
	acc.server = &http.Server{
		Addr:    acc.server.Addr,
		Handler: acc.server.Handler,
	}

	go func() {

		acc.log.Infof("start at %s", acc.server.Addr)

		// 开始监听
		err := acc.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			acc.log.Errorf("stopped by error: %s", err)
		}

		acc.mutex.Lock()
		defer acc.mutex.Unlock()

		// 更新服务器状态
		acc.state = network.Stopped
		acc.log.Info("stopped")
	}()
}

// 停止服务器
func (acc *Acceptor) Stop() {

	acc.mutex.Lock()
	defer acc.mutex.Unlock()

	// 修改状态
	if acc.state != network.Running {
		return
	}
	acc.state = network.Closing

	acc.log.Info("shutting down")

	// 关闭HTTP服务
	ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
	if err := acc.server.Shutdown(ctx); err != nil {
		acc.log.Errorf("stop error: %v", err)
	}
}

// 当前服务器状态
func (acc *Acceptor) State() network.State {
	acc.mutex.Lock()
	defer acc.mutex.Unlock()
	return acc.state
}

func NewAcceptor(config AcceptorConfig) *Acceptor {

	config.Load()

	return &Acceptor{
		state: network.Stopped,
		server: &http.Server{
			Addr:    config.Addr,
			Handler: config.Engine,
		},
		log: tlog.Std(config.Name),
	}
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}
