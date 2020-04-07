// 连接器

package websocket

import (
	"errors"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/logisor"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/utils/mathe"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------

func NewConnector(option *ConnectorOption) *Connector {

	option.parse()

	logger := logisor.Level(option.Level).
		Field(logis.Module, module).
		Field(network.FieldName, option.Name)

	return &Connector{option: option, logger: logger}
}

// ---------------------------------------------------------------------------------------------------------------------

type Connector struct {
	option    *ConnectorOption // 配置
	session   *Session         // 连接
	reconnect bool             // 是否重连
	times     int              // 重连次数
	logger    logis.Logger     // 日志
	mutex     sync.RWMutex
}

// ---------------------------------------------------------------------------------------------------------------------

func (connector *Connector) Connected() bool {

	connector.mutex.RLock()
	defer connector.mutex.RUnlock()

	return connector.session != nil
}

// ---------------------------------------------------------------------------------------------------------------------

func (connector *Connector) Run() {

	connector.mutex.Lock()
	defer connector.mutex.Unlock()

	if connector.session != nil {
		return
	}

	connector.reconnect = connector.option.Reconnect
	connector.connect()
}

func (connector *Connector) Stop() {

	connector.mutex.Lock()
	defer connector.mutex.Unlock()

	connector.reconnect = false

	if connector.session != nil {
		connector.session.Close()
		connector.session = nil
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func (connector *Connector) connect() {

	option := connector.option

	connector.logger.Data(option.Addr).Info("connect")

	conn, _, err := option.Dialer.Dial(option.Addr, option.Header)
	session := newSession(conn, option.Addr, &option.Session, connector.logger)

	if err != nil {
		connector.logger.Err(err).Error("dial error")
		go connector.delay()
		go session.invoke(network.NewConnectFailedEvent())
		return
	}

	connector.session = session
	connector.times = 0

	go session.run(func(session *Session) {

		connector.mutex.Lock()
		defer connector.mutex.Unlock()

		if connector.session == session {
			connector.session = nil
			go connector.delay()
		}
	})
}

func (connector *Connector) delay() {

	connector.mutex.Lock()
	defer connector.mutex.Unlock()

	if !connector.reconnect || connector.session != nil {
		return
	}

	delays := connector.option.Delays
	delay := delays[mathe.MinInt(connector.times, len(delays)-1)]
	connector.times++

	go func() {

		<-time.After(delay)
		connector.logger.Data(delay.String()).Info("reconnect")

		connector.mutex.Lock()
		defer connector.mutex.Unlock()

		if !connector.reconnect || connector.session != nil {
			return
		}

		connector.connect()
	}()
}

// ---------------------------------------------------------------------------------------------------------------------

func (connector *Connector) Send(msg interface{}) error {

	connector.mutex.RLock()
	defer connector.mutex.RUnlock()

	if connector.session != nil {
		return connector.session.Send(msg)
	}
	return errors.New("disconnected")
}

func (connector *Connector) SendRaw(raw []byte) error {

	connector.mutex.RLock()
	defer connector.mutex.RUnlock()

	if connector.session != nil {
		return connector.session.SendRaw(raw)
	}
	return errors.New("disconnected")
}
