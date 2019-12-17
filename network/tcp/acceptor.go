package tcp

import (
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	"net"
	"sync"
)

type Acceptor struct {
	conf       AcceptorConfig      // 配置
	listener   net.Listener        // 监听器
	sessionMgr *network.SessionMgr // session管理器
	logger     *log.Logger         // 日志
	mutex      sync.Mutex
}

func (acc *Acceptor) Run() {

	// 检查状态
	acc.mutex.Lock()
	if acc.listener != nil {
		acc.mutex.Unlock()
		return
	}

	// 监听端口
	listener, err := net.Listen("tcp", acc.conf.Addr)
	if err != nil {
		acc.logger.Errorf("listen at %s error: %v", acc.conf.Addr, err)
		return
	}

	// 设置状态
	acc.logger.Infof("running at: %s", acc.conf.Addr)
	acc.listener = listener
	acc.mutex.Unlock()

	for {

		// 建立连接
		conn, err := listener.Accept()
		if err != nil {
			acc.logger.Errorf("accept error: %v", err)
			break
		}

		// 执行session
		id := acc.sessionMgr.NewID()
		ses := newSession(acc.conf.Name, id, conn, &acc.conf.Session)
		go ses.run(acc.onSessionClose)
	}

	// 重置状态
	acc.mutex.Lock()
	acc.listener = nil
	acc.mutex.Unlock()
}

func (acc *Acceptor) Stop() {

	acc.mutex.Lock()
	acc.mutex.Unlock()

	if acc.listener != nil {
		acc.listener.Close()
	}
}

func (acc *Acceptor) State() network.State {

	acc.mutex.Lock()
	defer acc.mutex.Unlock()

	if acc.listener != nil {
		return network.Running
	}
	return network.Stopped
}

func (acc *Acceptor) Count() int64 {
	return acc.sessionMgr.Count()
}

func (acc *Acceptor) onSessionClose(ses *Session) {
	acc.sessionMgr.Del(ses)
}

func NewAcceptor(conf AcceptorConfig) network.Acceptor {

	conf.make()

	return &Acceptor{
		conf:       conf,
		sessionMgr: network.NewSessionMgr(),
		logger:     log.Std(conf.Name),
	}
}
