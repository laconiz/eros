package config

import (
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network/tcp"
	"time"
)

var Session = tcp.SessionConfig{
	ReadTimeout:  time.Minute,
	WriteTimeout: time.Minute,
	LogLevel:     log.Warn,
	QueueLen:     64,
}
