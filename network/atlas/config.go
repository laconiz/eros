package atlas

import (
	"time"

	"github.com/laconiz/eros/holder/message"
	"github.com/laconiz/eros/network"
)

type SessionOption struct {
	// 读超时
	ReadTimeout time.Duration
	// 写超时
	WriteTimeout time.Duration
	// 写缓冲区长度
	QueueLen int
	// 回调接口
	Invoker network.Invoker
	// 编码器
	Encoder message.Encoder
}

func (conf *SessionOption) make() {

	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = time.Second * 15
	}

	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = time.Second * 15
	}

	if conf.QueueLen <= 0 {
		conf.QueueLen = 32
	}

	if conf.Invoker == nil {
		conf.Invoker = network.NewStdInvoker()
	}

	if conf.EncoderMaker == nil {
		conf.EncoderMaker = &StdEncoderMaker{}
	}
}

type AcceptorOption struct {
	Name    string
	Addr    string
	Session SessionOption
}

func (conf *AcceptorOption) make() {

	if conf.Name == "" {
		conf.Name = "acceptor"
	}

	conf.Session.make()
}

type ConnectorConfig struct {
	Name      string
	Addr      string
	Reconnect bool
	Session   SessionOption
}

func (conf *ConnectorConfig) make() {

	if conf.Name == "" {
		conf.Name = "connector"
	}

	conf.Session.make()
}
