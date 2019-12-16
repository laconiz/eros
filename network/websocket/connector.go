package websocket

import (
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	"sync"
	"time"
)

type Connector struct {
	state        network.State   // 客户端状态
	session      *Session        // websocket连接
	config       ConnectorConfig // 配置
	connectTimes int             // 重连次数
	log          *log.Logger     // 日志
	mutex        sync.Mutex
}

// 启动客户端连接
func (cnt *Connector) Start() {

	cnt.mutex.Lock()
	defer cnt.mutex.Unlock()

	// 修改状态
	if cnt.state != network.Stopped {
		return
	}
	cnt.state = network.Running

	// 建立连接
	if cnt.session == nil {
		cnt.connect()
	}
}

// 连接websocket
func (cnt *Connector) connect() {

	config := cnt.config

	conn, _, err := cnt.config.Dialer.Dial(config.Addr, config.Header)

	// 连接失败
	if err != nil {

		cnt.log.Warnf("dial error: %v", err)

		// 检查重连
		go cnt.reconnect()

		go func() {
			ses := newSession(connectorSessionID, cnt.config.Name, "", nil, &cnt.config.Session)
			ses.invoke(&network.Event{
				Meta: network.MetaConnectFailed,
				Msg:  &network.ConnectFailed{},
			})
		}()

		return
	}

	ses := newSession(connectorSessionID, cnt.config.Name, cnt.config.Addr, conn, &cnt.config.Session)
	cnt.session = ses
	go ses.run(cnt.onSessionClose)
}

// 检查重连
func (cnt *Connector) reconnect() {

	cnt.mutex.Lock()
	defer cnt.mutex.Unlock()

	// 检查参数
	if cnt.state != network.Running || !cnt.config.Reconnect || cnt.session != nil {
		return
	}

	// 延迟重连
	duration := reconnectDuration[len(reconnectDuration)-1]
	if cnt.connectTimes < len(reconnectDuration) {
		duration = reconnectDuration[cnt.connectTimes]
	}

	go func() {

		// 等待时间
		<-time.After(duration)

		cnt.mutex.Lock()
		defer cnt.mutex.Unlock()

		// 检查状态
		if cnt.state != network.Running || !cnt.config.Reconnect || cnt.session != nil {
			return
		}

		// 增加重连计数
		cnt.connectTimes++
		// 重新连接
		cnt.connect()
	}()
}

// 关闭连接
func (cnt *Connector) Stop() {

	cnt.mutex.Lock()
	defer cnt.mutex.Unlock()

	// 修改状态
	if cnt.state != network.Running {
		return
	}
	cnt.state = network.Stopped

	// 关闭旧连接
	if cnt.session != nil {
		cnt.session.Close()
	}

	// 强行重置连接以快速重新开始
	cnt.session = nil
}

func (cnt *Connector) onSessionClose(ses *Session) {

	// 清理连接
	cnt.mutex.Lock()
	// 异步回调 检查老连接断开时已建立新连接
	if cnt.session == ses {
		cnt.session = nil
	}
	cnt.mutex.Unlock()

	// 检查重连
	cnt.reconnect()
}

func NewConnector(config ConnectorConfig) *Connector {

	// 检查配置
	config.make()

	return &Connector{
		state:   network.Stopped,
		session: nil,
		config:  config,
		log:     log.Std(config.Name),
	}
}

// session id
const connectorSessionID = 0

// 重连时间
var reconnectDuration = []time.Duration{
	time.Millisecond * 10,
	time.Second,
	time.Second * 3,
	time.Second * 6,
	time.Second * 10,
	time.Second * 15,
	time.Second * 21,
	time.Second * 30,
}
