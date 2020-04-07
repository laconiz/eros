package socket

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/network/cipher"
	"github.com/laconiz/eros/network/encoder"
	"github.com/laconiz/eros/network/invoker"
	"github.com/laconiz/eros/network/socket/reader"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------

const module = "socket"

// ---------------------------------------------------------------------------------------------------------------------
// 连接配置

type SessionOption struct {
	Queue   int             // 写缓冲区长度
	Timeout time.Duration   // 超时
	Reader  reader.Maker    // 包装器
	Cipher  cipher.Maker    // 加密器
	Encoder encoder.Maker   // 编码器
	Invoker invoker.Invoker // 消息调用器
}

func (opt *SessionOption) parse() {

	if opt.Timeout <= 0 {
		opt.Timeout = time.Second * 15
	}

	if opt.Queue <= 0 {
		opt.Queue = 32
	}

	if opt.Invoker == nil {
		opt.Invoker = invoker.NewSocketInvoker()
	}

	if opt.Encoder == nil {
		opt.Encoder = encoder.NewNameMaker()
	}

	if opt.Cipher == nil {
		opt.Cipher = cipher.NewIndexMaker()
	}

	if opt.Reader == nil {
		opt.Reader = reader.NewSizeMaker()
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// 侦听器配置

type AcceptorOption struct {
	Name    string        // 名称
	Addr    string        // 侦听地址
	Level   logis.Level   // 日志等级
	Session SessionOption // 连接配置
}

func (option *AcceptorOption) parse() {

	if option.Name == "" {
		option.Name = "acceptor"
	}

	if !option.Level.Valid() {
		option.Level = logis.INFO
	}

	option.Session.parse()
}

// ---------------------------------------------------------------------------------------------------------------------
// 连接器配置

type ConnectorOption struct {
	Name      string          // 名称
	Addr      string          // 连接地址
	Reconnect bool            // 是否重连
	Delays    []time.Duration // 重连延迟
	Level     logis.Level     // 日志等级
	Session   SessionOption   // session配置
}

func (option *ConnectorOption) parse() {

	if option.Name == "" {
		option.Name = "connector"
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

	option.Session.parse()
}
