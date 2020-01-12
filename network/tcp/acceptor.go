package tcp

import (
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	"net"
	"sync"
)

type Acceptor struct {
	state      network.State       // 当前状态
	conf       AcceptorConfig      // 配置
	listener   net.Listener        // 监听器
	sessionMgr *network.SessionMgr // session管理器
	// logger     *log.Logger         // 日志
	logger *log.Entry
	mutex  sync.Mutex
}

// 运行服务器接口
func (a *Acceptor) Run() {

	// 检查状态
	a.mutex.Lock()
	if a.state != network.Stopped {
		a.mutex.Unlock()
		return
	}

	a.state = network.Running

	// 监听端口
	listener, err := net.Listen("tcp", a.conf.Addr)
	if err != nil {
		//a.logger.Errorf("listen at %s error: %v", a.conf.Addr, err)
		return
	}

	// 设置状态
	a.logger.Infof("running at: %s", a.conf.Addr)
	a.listener = listener
	a.mutex.Unlock()

	for {

		// 建立连接
		conn, err := listener.Accept()
		if err != nil {
			a.logger.Errorf("accept error: %v", err)
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

	a.sessionMgr.Range(func(session network.Session) bool {
		session.Close()
		return true
	})

	a.mutex.Lock()
	a.mutex.Unlock()

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

func NewAcceptor(conf AcceptorConfig) network.Acceptor {

	conf.make()

	return &Acceptor{
		conf:       conf,
		sessionMgr: network.NewSessionMgr(),
		logger:     log.Std(conf.Name),
	}
}
