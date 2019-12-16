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
	Invoker      network.Invoker
}

func (conf *SessionConfig) make() {

	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = time.Second * 15
	}

	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = time.Second * 15
	}
}
