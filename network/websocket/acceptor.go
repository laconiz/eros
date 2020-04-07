package websocket

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/session"
	"net/http"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------

func NewAcceptor(option *AcceptorOption) *Acceptor {

	option.parse()

	logger := logisor.Level(option.Level).
		Field(logis.Module, module).
		Field(network.FieldName, option.Name)

	acceptor := &Acceptor{
		state:    network.Stopped,
		option:   option,
		sessions: session.NewManager(),
		logger:   logger,
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.GET("/ws", func(context *gin.Context) {
		acceptor.upgrade(context)
	})

	acceptor.listener = &http.Server{
		Addr:    option.Addr,
		Handler: engine,
	}

	return acceptor
}

// ---------------------------------------------------------------------------------------------------------------------
// websocket侦听器

type Acceptor struct {
	state    network.State    // 状态
	option   *AcceptorOption  // 配置信息
	listener *http.Server     // 侦听器
	sessions *session.Manager // 连接管理器
	logger   logis.Logger     // 日志接口
	mutex    sync.Mutex
}

// ---------------------------------------------------------------------------------------------------------------------

func (acceptor *Acceptor) State() network.State {

	acceptor.mutex.Lock()
	defer acceptor.mutex.Unlock()

	return acceptor.state
}

func (acceptor *Acceptor) Count() int64 {
	return acceptor.sessions.Count()
}

// ---------------------------------------------------------------------------------------------------------------------

func (acceptor *Acceptor) Run() {

	acceptor.mutex.Lock()
	defer acceptor.mutex.Unlock()

	if acceptor.state != network.Stopped {
		return
	}
	acceptor.state = network.Running

	option := acceptor.option

	acceptor.listener = &http.Server{
		Addr:    option.Addr,
		Handler: acceptor.listener.Handler,
	}
	acceptor.logger.Data(option.Addr).Info("start")

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

	acceptor.sessions.Range(func(session session.Session) bool {
		session.Close()
		return true
	})

	ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
	if err := acceptor.listener.Shutdown(ctx); err != nil {
		acceptor.logger.Err(err).Error("stop error")
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func (acceptor *Acceptor) upgrade(context *gin.Context) {

	logger := acceptor.logger
	option := acceptor.option

	if err := option.Verify(context); err != nil {
		logger.Err(err).Warn("verify failed")
		return
	}

	conn, err := option.Upgrader.Upgrade(context.Writer, context.Request, context.Request.Header)
	if err != nil {
		logger.Err(err).Warn("upgrade error")
		return
	}

	addr := context.Request.RemoteAddr
	if addr == "" {
		addr = conn.RemoteAddr().String()
	}

	session := newSession(conn, addr, &option.Session, logger)
	acceptor.sessions.Insert(session)
	go session.run(func(session *Session) {
		acceptor.sessions.Remove(session)
	})
}

// ---------------------------------------------------------------------------------------------------------------------

func (acceptor *Acceptor) Broadcast(msg interface{}) {
	acceptor.sessions.Range(func(session session.Session) bool {
		session.Send(msg)
		return true
	})
}

func (acceptor *Acceptor) BroadcastRaw(raw []byte) {
	acceptor.sessions.Range(func(session session.Session) bool {
		session.SendRaw(raw)
		return true
	})
}

// ---------------------------------------------------------------------------------------------------------------------

func init() {
	gin.SetMode(gin.ReleaseMode)
}
