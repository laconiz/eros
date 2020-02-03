package atlas

import (
	"net"
	"strings"
	"sync"

	"github.com/laconiz/eros/hyperion"
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/epimetheus"
	"github.com/laconiz/eros/network/incremental"
)

type Acceptor struct {
	state      network.State
	option     AcceptorOption      // 配置
	listener   net.Listener        // 监听器
	sessionMgr *network.SessionMgr // session管理器
	logger     *hyperion.Entry     // 日志
	mutex      sync.Mutex
}

func (a *Acceptor) run() net.Listener {

	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.state != network.Stopped {
		return nil
	}

	listener, err := net.Listen("tcp", a.option.Addr)
	if err != nil {
		a.logger.WithError(err).Error("listen error")
		return nil
	}

	a.listener = listener
	return listener
}

const contentListenerClosed = "use of closed network connection"

func (a *Acceptor) Run() {

	listener := a.run()
	if listener == nil {
		return
	}

	a.logger.Info(epimetheus.ContentStarted)

	for {

		conn, err := listener.Accept()
		if err != nil || !strings.Contains(err.Error(), contentListenerClosed) {
			a.logger.WithError(err).Error("accept error")
			break
		}

		sesID := (network.SessionID)(incremental.Get())
		ses := newSession()
	}

	// 检查状态
	a.mutex.Lock()
	if a.listener != nil {
		a.mutex.Unlock()
		return
	}

	// 监听端口
	listener, err := net.Listen("atlas", a.conf.Addr)
	if err != nil {
		a.logger.Errorf("listen at %s error: %v", a.conf.Addr, err)
		return
	}

	// 设置状态
	a.logger.Infof("running at: %s", a.conf.Addr)
	a.listener = listener
	a.mutex.Unlock()

	for {

		// 建立连接
		conn, err := listener.Accept()
		if err != nil && !strings.Contains(err.Error()) {
			a.logger.WithError(err).Error("accept error")
			break
		}

		// 执行session
		id := a.sessionMgr.NewID()
		ses := newSession(a.conf.Name, id, conn, &a.conf.Session)
		go ses.run(a.onSessionClose)
	}

	// 重置状态
	a.mutex.Lock()
	a.listener = nil
	a.mutex.Unlock()
}

func (a *Acceptor) Stop() {

	a.mutex.Lock()
	a.mutex.Unlock()

	a.sessionMgr.Range(func(session network.Session) bool {
		session.Close()
		return true
	})

	if a.listener != nil {
		a.listener.Close()
	}
}

func (a *Acceptor) State() network.State {

	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.listener != nil {
		return network.Running
	}
	return network.Stopped
}

func (a *Acceptor) Count() int64 {
	return a.sessionMgr.Count()
}

func (a *Acceptor) Broadcast(msg interface{}) {
	a.sessionMgr.Range(func(ses network.Session) bool {
		ses.Send(msg)
		return true
	})
}

func (a *Acceptor) BroadcastStream(stream []byte) {
	a.sessionMgr.Range(func(ses network.Session) bool {
		ses.SendStream(stream)
		return true
	})
}

func (a *Acceptor) onSessionClose(ses *Session) {
	a.sessionMgr.Del(ses)
}

func NewAcceptor(conf AcceptorOption) network.Acceptor {

	conf.make()

	return &Acceptor{
		conf:       conf,
		sessionMgr: network.NewSessionMgr(),
		logger:     log.Std(conf.Name),
	}
}
