// socket服务器

package socket

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/context"
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/session"
	"net"
	"strings"
	"sync"
)

// 生成一个socket服务器
func NewAcceptor(opt AccOption) *Acceptor {
	opt.parse()
	return &Acceptor{
		option:   opt,
		sessions: session.NewManager(),
		log: logisor.Fields(context.Fields{
			logis.Module:      module,
			network.FieldName: opt.Name,
		}),
	}
}

// socket服务器
type Acceptor struct {
	option   AccOption        // 配置
	listener net.Listener     // 监听器
	sessions *session.Manager // session管理器
	log      logis.Logger     // 日志
	mutex    sync.RWMutex
}

const errClosed = "use of closed network connection"

// 启动服务器
func (acc *Acceptor) Run() {

	acc.mutex.Lock()
	defer acc.mutex.Unlock()
	if acc.listener != nil {
		return
	}

	// 创建侦听器
	listener, err := net.Listen("tcp", acc.option.Addr)
	if err != nil {
		acc.log.Errorf("listen error: %v", err)
		return
	}
	acc.listener = listener

	acc.log.Infof("listen at: %s", acc.option.Addr)

	go func() {

		for {

			conn, err := listener.Accept()
			if err != nil {
				if !strings.Contains(err.Error(), errClosed) {
					acc.log.Errorf("accept error: %v", err)
				}
				break
			}

			ses := newSession(conn, &acc.option.Session, acc.log)
			acc.sessions.Insert(ses)
			go ses.run(acc.onSessionClose)
		}

		acc.mutex.Lock()
		defer acc.mutex.Unlock()
		if acc.listener == listener {
			acc.listener = nil
		}
		acc.log.Info("stopped")
	}()
}

// 停止服务器
func (acc *Acceptor) Stop() {

	acc.mutex.Lock()
	defer acc.mutex.Unlock()

	acc.sessions.Range(func(ses session.Session) bool {
		ses.Close()
		return true
	})

	if acc.listener != nil {
		acc.listener.Close()
	}
}

// 服务器状态
func (acc *Acceptor) State() network.State {
	acc.mutex.Lock()
	defer acc.mutex.Unlock()
	if acc.listener != nil {
		return network.Running
	}
	return network.Stopped
}

// 服务器连接数量
func (acc *Acceptor) Count() int64 {
	return acc.sessions.Count()
}

// 广播消息
func (acc *Acceptor) Broadcast(msg interface{}) {
	acc.sessions.Range(func(ses session.Session) bool {
		ses.Send(msg)
		return true
	})
}

// 广播字节流
func (acc *Acceptor) BroadcastRaw(stream []byte) {
	acc.sessions.Range(func(ses session.Session) bool {
		ses.SendRaw(stream)
		return true
	})
}

// session关闭回调
func (acc *Acceptor) onSessionClose(ses *Session) {
	acc.sessions.Remove(ses)
}
