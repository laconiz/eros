// socket服务器

package socket

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/session"
	"net"
	"strings"
	"sync"
)

// ---------------------------------------------------------------------------------------------------------------------

func NewAcceptor(option *AcceptorOption) *Acceptor {

	option.parse()

	logger := logisor.Module(module).
		Level(option.Level).
		Field(network.FieldName, option.Name)

	return &Acceptor{
		option:   option,
		sessions: session.NewManager(),
		logger:   logger,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type Acceptor struct {
	option   *AcceptorOption  // 配置
	listener net.Listener     // 监听器
	sessions *session.Manager // 连接管理器
	logger   logis.Logger     // 日志接口
	mutex    sync.RWMutex
}

// ---------------------------------------------------------------------------------------------------------------------

func (acceptor *Acceptor) State() network.State {

	acceptor.mutex.Lock()
	defer acceptor.mutex.Unlock()

	if acceptor.listener != nil {
		return network.Running
	}
	return network.Stopped
}

func (acceptor *Acceptor) Count() int64 {
	return acceptor.sessions.Count()
}

// ---------------------------------------------------------------------------------------------------------------------

func (acceptor *Acceptor) Run() {

	acceptor.mutex.Lock()
	defer acceptor.mutex.Unlock()

	if acceptor.listener != nil {
		return
	}

	option := acceptor.option

	listener, err := net.Listen("tcp", option.Addr)
	if err != nil {
		acceptor.logger.Err(err).Error("listen error")
		return
	}
	acceptor.listener = listener

	acceptor.logger.Data(option.Addr).Info("start")

	go func() {

		acceptor.accept(listener)

		acceptor.mutex.Lock()
		defer acceptor.mutex.Unlock()

		if acceptor.listener == listener {
			acceptor.listener = nil
		}

		acceptor.logger.Info("stopped")
	}()
}

func (acceptor *Acceptor) Stop() {

	acceptor.mutex.Lock()
	defer acceptor.mutex.Unlock()

	if acceptor.listener == nil {
		return
	}

	acceptor.sessions.Range(func(session session.Session) bool {
		session.Close()
		return true
	})

	acceptor.listener.Close()
}

// ---------------------------------------------------------------------------------------------------------------------

func (acceptor *Acceptor) accept(listener net.Listener) {

	logger := acceptor.logger
	option := acceptor.option

	for {

		conn, err := listener.Accept()

		if err == nil {

			session := newSession(conn, &option.Session, logger)
			acceptor.sessions.Insert(session)

			go session.run(func(session *Session) {
				acceptor.sessions.Remove(session)
			})

			continue
		}

		const strClosed = "use of closed network connection"
		if !strings.Contains(err.Error(), strClosed) {
			logger.Err(err).Error("acceptor error")
		}

		return
	}
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
