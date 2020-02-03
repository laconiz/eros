package atlas

import (
	"errors"
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	"net"
	"sync"
	"time"
)

type Connector struct {
	state          network.State   // 状态
	conf           ConnectorConfig // 配置
	ses            *Session        // 连接
	reconnectTimes int             // 重连次数
	logger         *log.Logger
	mutex          sync.Mutex
}

func (c *Connector) Run() {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.state != network.Stopped {
		return
	}
	c.state = network.Running

	if c.ses == nil {
		c.connect()
	}
}

func (c *Connector) Stop() {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.ses != nil {
		c.ses.Close()
	}
}

func (c *Connector) State() network.State {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.state
}

func (c *Connector) Connected() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.ses != nil
}

func (c *Connector) connect() {

	conf := c.conf

	// 连接
	conn, err := net.Dial("atlas", conf.Addr)
	ses := newSession(conf.Name, connectorSessionID, conn, &conf.Session)

	if err != nil {

		c.logger.Warnf("dial error: %v", err)

		// 重连
		go c.reconnect()

		// 连接失败回调
		go ses.invoke(&network.Event{
			Meta: network.MetaConnectFailed,
			Msg:  &network.ConnectFailed{},
		})

	} else {

		c.reconnectTimes = 0
		c.ses = ses
		go ses.run(c.onSesClose)
	}
}

func (c *Connector) onSesClose(ses *Session) {

	c.mutex.Lock()
	if c.ses == ses {
		c.ses = nil
	}
	c.mutex.Unlock()

	c.reconnect()
}

func (c *Connector) reconnect() {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.state != network.Running || !c.conf.Reconnect || c.ses != nil {
		return
	}

	// 延迟重连
	duration := reconnectDuration[len(reconnectDuration)-1]
	if c.reconnectTimes < len(reconnectDuration) {
		duration = reconnectDuration[c.reconnectTimes]
	}

	go func() {

		<-time.After(duration)

		c.mutex.Lock()
		defer c.mutex.Unlock()

		if c.state != network.Running || !c.conf.Reconnect || c.ses != nil {
			return
		}

		c.reconnectTimes++
		c.connect()
	}()
}

func (c *Connector) Send(msg interface{}) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.ses != nil {
		return c.ses.Send(msg)
	}

	return errors.New("disconnected")
}

func NewConnector(conf ConnectorConfig) network.Connector {

	conf.make()

	return &Connector{
		state:  network.Stopped,
		conf:   conf,
		logger: log.Std(conf.Name),
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
}
