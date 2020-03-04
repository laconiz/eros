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

// 生成一个socket服务器
func NewAcceptor(option AcceptorOption) *Acceptor {
	option.parse()
	logger := logisor.Level(option.Level).Field(logis.Module, module).Field(network.FieldName, option.Name)
	return &Acceptor{option: option, sessions: session.NewManager(), logger: logger}
}

// socket服务器
type Acceptor struct {
	option   AcceptorOption   // 配置
	listener net.Listener     // 监听器
	sessions *session.Manager // session管理器
	logger   logis.Logger     // 日志
	mutex    sync.RWMutex
}

// 启动服务器
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

		for {

			conn, err := listener.Accept()
			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					acceptor.logger.Err(err).Error("accept error")
				}
				break
			}

			ses := newSession(conn, &acceptor.option.Session, acceptor.logger)
			acceptor.sessions.Insert(ses)
			go ses.run(acceptor.onSessionClose)
		}

		acceptor.mutex.Lock()
		defer acceptor.mutex.Unlock()

		if acceptor.listener == listener {
			acceptor.listener = nil
		}

		acceptor.logger.Info("stopped")
	}()
}

// 停止服务器
func (acceptor *Acceptor) Stop() {

	acceptor.mutex.Lock()
	defer acceptor.mutex.Unlock()

	acceptor.sessions.Range(func(ses session.Session) bool {
		ses.Close()
		return true
	})

	if acceptor.listener != nil {
		acceptor.listener.Close()
	}
}

// 服务器状态
func (acceptor *Acceptor) State() network.State {

	acceptor.mutex.Lock()
	defer acceptor.mutex.Unlock()

	if acceptor.listener != nil {
		return network.Running
	}
	return network.Stopped
}

// 服务器连接数量
func (acceptor *Acceptor) Count() int64 {
	return acceptor.sessions.Count()
}

// 广播消息
func (acceptor *Acceptor) Broadcast(msg interface{}) {
	acceptor.sessions.Range(func(ses session.Session) bool {
		ses.Send(msg)
		return true
	})
}

// 广播字节流
func (acceptor *Acceptor) BroadcastRaw(stream []byte) {
	acceptor.sessions.Range(func(ses session.Session) bool {
		ses.SendRaw(stream)
		return true
	})
}

// session关闭回调
func (acceptor *Acceptor) onSessionClose(ses *Session) {
	acceptor.sessions.Remove(ses)
}
