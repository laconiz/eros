package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/laconiz/eros/network"
	"math"
	"net/http"
	"time"
)

type SessionConfig struct {
	ReadLimit     int64           // 读取数据流大小限制
	ReadTimeout   time.Duration   // 读取超时
	WriteTimeout  time.Duration   // 写入超时
	WriteQueueLen int             // 写入队列大小
	Encoder       Encoder         // 编码器
	Invoker       network.Invoker // 消息调用器
}

func (c *SessionConfig) make() {

	if c.ReadLimit <= 0 {
		c.ReadLimit = 1024 * 32
	}

	if c.ReadTimeout <= 0 {
		c.ReadTimeout = time.Second * 15
	}

	if c.WriteTimeout <= 0 {
		c.WriteTimeout = time.Second * 15
	}

	if c.WriteQueueLen <= 0 {
		c.WriteQueueLen = 32
	}

	if c.Encoder == nil {
		c.Encoder = NameEncoder
	}

	if c.Invoker == nil {
		c.Invoker = network.NewStdInvoker()
	}
}

type AcceptorConfig struct {
	Name     string
	Addr     string
	Verify   func(*gin.Context) (map[interface{}]interface{}, error)
	Upgrader *websocket.Upgrader // 连接升级
	Session  SessionConfig
}

func (config *AcceptorConfig) make() {

	if config.Name == "" {
		config.Name = "acceptor"
	}

	if config.Verify == nil {
		config.Verify = verify
	}

	if config.Upgrader == nil {
		config.Upgrader = &websocket.Upgrader{
			HandshakeTimeout:  time.Second * 3,
			EnableCompression: true,
		}
	}

	config.Session.make()
}

func verify(_ *gin.Context) (map[interface{}]interface{}, error) {
	return map[interface{}]interface{}{}, nil
}

type ConnectorConfig struct {
	Name      string            // 日志名
	Addr      string            // 地址
	Header    http.Header       // 请求header
	Reconnect bool              // 自动重连
	Dialer    *websocket.Dialer // 连接器
	Session   SessionConfig     // session配置
}

func (config *ConnectorConfig) make() {

	if config.Name == "" {
		config.Name = "connector"
	}

	if config.Dialer == nil {
		config.Dialer = &websocket.Dialer{
			HandshakeTimeout: time.Second * 3,
			// EnableCompression: true,
		}
	}

	// 客户端不限制数据大小
	if config.Session.ReadLimit <= 0 {
		config.Session.ReadLimit = math.MaxInt64
	}

	config.Session.make()
}
