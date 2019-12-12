package websocket

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
	state      network.State       // 服务器状态
	sessionMgr *network.SessionMgr // 连接管理器
	server     *http.Server        // HTTP服务
	config     AcceptorConfig      // 配置
	log        *log.Log            // 日志
	mutex      sync.Mutex
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

	// 通知所有连接关闭
	acc.sessionMgr.Range(func(ses network.Session) bool {
		ses.Close()
		return true
	})

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

// 当前连接数量
func (acc *Acceptor) Count() int64 {
	return acc.sessionMgr.Count()
}

// 升级HTTP请求
func (acc *Acceptor) upgrade(context *gin.Context) {

	// 请求合法性检测
	data, err := acc.config.Verify(context)
	if err != nil {
		acc.log.Warnf("verify failed: %v", err)
		return
	}

	// 升级连接
	conn, err := acc.config.Upgrader.Upgrade(context.Writer, context.Request, context.Request.Header)
	if err != nil {
		acc.log.Errorf("upgrade error: %v", err)
		return
	}

	// 客户端真实IP地址
	addr := context.Request.RemoteAddr
	if addr == "" {
		addr = conn.RemoteAddr().String()
	}

	// 构造session
	session := newSession(acc.sessionMgr.NewID(), acc.config.Name, addr, conn, &acc.config.Session)

	// 客户端携带信息
	for key, value := range data {
		session.Set(key, value)
	}

	go session.run(acc.onSessionClose)
}

func (acc *Acceptor) onSessionClose(ses *Session) {
	acc.sessionMgr.Del(ses)
}

func NewAcceptor(config AcceptorConfig) *Acceptor {

	config.make()

	acc := &Acceptor{
		state:      network.Stopped,
		sessionMgr: network.NewSessionMgr(),
		config:     config,
		log:        log.Std(config.Name),
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.GET("/ws", func(context *gin.Context) {
		acc.upgrade(context)
	})

	acc.server = &http.Server{
		Addr:    config.Addr,
		Handler: engine,
	}

	return acc
}
