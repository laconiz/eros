// socket客户端

package socket

import (
	"errors"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/utils/mathe"
	"net"
	"sync"
	"time"
)

// 生成一个socket客户端
func NewConnector(option ConnectorOption) network.Connector {
	option.parse()
	logger := logisor.Level(option.Level).Field(logis.Module, module).Field(network.FieldName, option.Name)
	return &Connector{option: option, logger: logger}
}

// socket客户端
type Connector struct {
	option    ConnectorOption // 配置
	session   *Session        // 连接
	reconnect bool            // 是否重连
	times     int             // 重连次数
	logger    logis.Logger    // 日志接口
	mutex     sync.RWMutex
}

// 启动客户端
func (connector *Connector) Run() {

	connector.mutex.Lock()
	defer connector.mutex.Unlock()

	if connector.session != nil {
		return
	}

	connector.reconnect = connector.option.Reconnect
	connector.connect()
}

// 停止客户端
func (connector *Connector) Stop() {

	connector.mutex.Lock()
	defer connector.mutex.Unlock()

	connector.reconnect = false

	if connector.session != nil {
		connector.session.Close()
		connector.session = nil
	}
}

// 客户端状态
func (connector *Connector) State() network.State {

	connector.mutex.RLock()
	defer connector.mutex.RUnlock()

	if connector.session != nil {
		return network.Running
	}
	return network.Stopped
}

// 客户端连接状态
func (connector *Connector) Connected() bool {
	connector.mutex.RLock()
	defer connector.mutex.RUnlock()
	return connector.session != nil
}

// 创建连接
func (connector *Connector) connect() {

	opt := connector.option
	conn, err := net.Dial("tcp", opt.Addr)
	session := newSession(conn, &opt.Session, connector.logger)

	if err != nil {
		connector.logger.Err(err).Error("dial error")
		go connector.delayConnect()
		go session.invoke(network.NewConnectFailedEvent())
		return
	}

	connector.session = session
	connector.times = 0
	go session.run(connector.onSesClose)
}

// session关闭回调
func (connector *Connector) onSesClose(session *Session) {

	connector.mutex.Lock()
	defer connector.mutex.Unlock()

	if connector.session == session {
		connector.session = nil
		go connector.delayConnect()
	}
}

// 重连客户端
func (connector *Connector) delayConnect() {

	connector.mutex.Lock()
	defer connector.mutex.Unlock()

	if !connector.reconnect {
		return
	}
	if connector.session != nil {
		return
	}

	option := connector.option
	delay := option.Delays[mathe.MinInt(connector.times, len(option.Delays)-1)]
	connector.times++

	go func() {

		<-time.After(delay)
		connector.logger.Data(delay.String()).Info("reconnect")

		connector.mutex.Lock()
		defer connector.mutex.Unlock()

		if !connector.reconnect {
			return
		}
		if connector.session != nil {
			return
		}

		connector.connect()
	}()
}

// 发送消息
func (connector *Connector) Send(msg interface{}) error {

	connector.mutex.Lock()
	defer connector.mutex.Unlock()

	if connector.session != nil {
		return connector.session.Send(msg)
	}
	return errors.New("disconnected")
}

// 发送字节流
func (connector *Connector) SendRaw(raw []byte) error {

	connector.mutex.RLock()
	defer connector.mutex.RUnlock()

	if connector.session != nil {
		return connector.session.SendRaw(raw)
	}
	return errors.New("disconnected")
}
