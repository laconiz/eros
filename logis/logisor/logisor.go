package logisor

import (
	"github.com/laconiz/eros/database/consul/consulor"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/context"
	"github.com/laconiz/eros/logis/formatter"
	"github.com/laconiz/eros/logis/hook"
	"os"
)

func Module(module string) logis.Logger {
	return logger.Field(logis.Module, module)
}

func Level(level logis.Level) logis.Logger {
	return logger.Level(level)
}

func Field(key string, value interface{}) logis.Logger {
	return logger.Field(key, value)
}

func Fields(fields context.Fields) logis.Logger {
	return logger.Fields(fields)
}

func Data(value interface{}) logis.Logger {
	return logger.Data(value)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

var logger logis.Logger

const (
	path   = "eros/logisor"
	module = "logisor"
)

type Option struct {
	Text bool
}

func init() {

	option := &Option{}
	if err := consulor.KV().Load(path, option); err != nil {
		logger = hook.NewHook(formatter.Text()).
			Add(logis.DEBUG, os.Stdout).
			Entry()
		return
	}
}
