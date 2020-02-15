// socket客户端

package socket

import (
	"errors"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network"
	"net"
	"sync"
	"time"
)

// 生成一个socket客户端
func NewConnector(opt ConnOption) network.Connector {
	opt.parse()
	return &Connector{
		opt: opt,
		log: logisor.Fields(logis.Fields{
			logis.Module:      module,
			network.FieldName: opt.Name,
			network.FieldAddr: opt.Addr,
		}),
	}
}

// socket客户端
type Connector struct {
	opt   ConnOption   // 配置
	ses   *Session     // 连接
	times int          // 重连次数
	log   logis.Logger // 日志接口
	mutex sync.RWMutex
}

// 启动客户端
func (con *Connector) Run() {
	con.mutex.Lock()
	defer con.mutex.Unlock()
	if con.ses != nil {
		return
	}
	con.connect()
}

// 停止客户端
func (con *Connector) Stop() {
	con.mutex.Lock()
	defer con.mutex.Unlock()
	con.opt.Reconnect = false
	if con.ses != nil {
		con.ses.Close()
		con.ses = nil
	}
}

// 客户端状态
func (con *Connector) State() network.State {
	con.mutex.RLock()
	defer con.mutex.RUnlock()
	if con.ses != nil {
		return network.Running
	}
	return network.Stopped
}

// 客户端连接状态
func (con *Connector) Connected() bool {
	con.mutex.RLock()
	defer con.mutex.RUnlock()
	return con.ses != nil
}

// 创建连接
func (con *Connector) connect() {

	opt := con.opt
	conn, err := net.Dial("tcp", opt.Addr)
	ses := newSession(conn, &opt.Session, con.log)

	if err != nil {
		con.log.Errorf("dial error: %v", err)
		go con.reconnect()
		go ses.invoke(network.NewConnectFailedEvent())
		return
	}

	con.ses = ses
	con.times = 0
	go ses.run(con.onSesClose)
}

// session关闭回调
func (con *Connector) onSesClose(ses *Session) {
	con.mutex.Lock()
	if con.ses == ses {
		con.ses = nil
	}
	con.mutex.Unlock()
	con.reconnect()
}

// 重连客户端
func (con *Connector) reconnect() {

	con.mutex.Lock()
	defer con.mutex.Unlock()

	if !con.opt.Reconnect || con.ses != nil {
		return
	}

	delay := con.opt.Delays[len(con.opt.Delays)-1]
	if con.times < len(con.opt.Delays) {
		delay = con.opt.Delays[con.times]
	}
	con.times++

	go func() {
		con.log.Infof("reconnect after %v", delay)
		<-time.After(delay)
		con.mutex.Lock()
		defer con.mutex.Unlock()
		if !con.opt.Reconnect || con.ses != nil {
			return
		}
		con.connect()
	}()
}

// 发送消息
func (con *Connector) Send(msg interface{}) error {
	con.mutex.Lock()
	defer con.mutex.Unlock()
	if con.ses != nil {
		return con.ses.Send(msg)
	}
	return errDisconnected
}

// 发送字节流
func (con *Connector) SendRaw(raw []byte) error {
	con.mutex.RLock()
	defer con.mutex.RUnlock()
	if con.ses != nil {
		return con.ses.SendRaw(raw)
	}
	return errDisconnected
}

var errDisconnected = errors.New("disconnected")
