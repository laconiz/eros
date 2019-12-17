package tcp

import (
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	"time"
)

type SessionConfig struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	LogLevel     log.Level
	QueueLen     int
	Invoker      network.Invoker
	EncoderMaker EncoderMaker
}

func (conf *SessionConfig) make() {

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

type AcceptorConfig struct {
	Name    string
	Addr    string
	Session SessionConfig
}

func (conf *AcceptorConfig) make() {

	if conf.Name == "" {
		conf.Name = "acceptor"
	}

	conf.Session.make()
}

type ConnectorConfig struct {
	Name      string
	Addr      string
	Reconnect bool
	Session   SessionConfig
}

func (conf *ConnectorConfig) make() {

	if conf.Name == "" {
		conf.Name = "connector"
	}

	conf.Session.make()
}
