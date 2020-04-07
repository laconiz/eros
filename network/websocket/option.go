// 配置信息

package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/network/cipher"
	"github.com/laconiz/eros/network/encoder"
	"github.com/laconiz/eros/network/invoker"
	"math"
	"net/http"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------

const module = "websocket"

// ---------------------------------------------------------------------------------------------------------------------
// 连接配置

type SessionOption struct {
	ReadLimit int64           // 读取数据流大小限制
	Timeout   time.Duration   // 超时时间
	QueueLen  int             // 写缓冲区长度
	Encoder   encoder.Maker   // 编码器
	Cipher    cipher.Maker    // 加密器
	Invoker   invoker.Invoker // 消息调用器
}

func (option *SessionOption) parse() {

	if option.ReadLimit <= 0 {
		option.ReadLimit = 1024 * 128
	}

	if option.Timeout <= 0 {
		option.Timeout = time.Second * 15
	}

	if option.QueueLen <= 0 {
		option.QueueLen = 32
	}

	if option.Cipher == nil {
		option.Cipher = cipher.NewEmptyMaker()
	}

	if option.Encoder == nil {
		option.Encoder = encoder.NewNameMaker()
	}

	if option.Invoker == nil {
		option.Invoker = invoker.NewSocketInvoker()
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// 侦听器配置

type AcceptorOption struct {
	Name     string                   // 名称
	Addr     string                   // 侦听地址
	Verify   func(*gin.Context) error // 连接验证
	Upgrader *websocket.Upgrader      // 连接升级
	Level    logis.Level              // 日志等级
	Session  SessionOption            // 连接配置
}

func (option *AcceptorOption) parse() {

	if option.Name == "" {
		option.Name = "acceptor"
	}

	if option.Verify == nil {
		option.Verify = func(*gin.Context) error {
			return nil
		}
	}

	if option.Upgrader == nil {
		option.Upgrader = &websocket.Upgrader{HandshakeTimeout: time.Second * 3, EnableCompression: true}
	}

	if !option.Level.Valid() {
		option.Level = logis.INFO
	}

	option.Session.parse()
}

// ---------------------------------------------------------------------------------------------------------------------
// 连接器配置

type ConnectorOption struct {
	Name      string            // 名称
	Addr      string            // 地址
	Reconnect bool              // 自动重连
	Header    http.Header       // 请求header
	Dialer    *websocket.Dialer // 连接器
	Delays    []time.Duration   // 重连延迟
	Level     logis.Level       // 日志等级
	Session   SessionOption     // session配置
}

func (option *ConnectorOption) parse() {

	if option.Name == "" {
		option.Name = "connector"
	}

	if option.Dialer == nil {
		option.Dialer = &websocket.Dialer{HandshakeTimeout: time.Second * 3}
	}

	if len(option.Delays) == 0 {
		option.Delays = []time.Duration{
			time.Millisecond * 10,
			time.Millisecond * 500,
			time.Millisecond * 1200,
			time.Millisecond * 3600,
			time.Millisecond * 9000,
			time.Millisecond * 15000,
		}
	}

	if !option.Level.Valid() {
		option.Level = logis.INFO
	}

	if option.Session.ReadLimit <= 0 {
		option.Session.ReadLimit = math.MaxInt64
	}

	option.Session.parse()
}
