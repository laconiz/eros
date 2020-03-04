package logis

import (
	"github.com/laconiz/eros/logis/context"
	"time"
)

const (
	Module = "module"
)

type Log struct {
	Level   Level            // 等级
	Error   string           // 错误信息
	Message string           // 文本内容
	Time    time.Time        // 时间
	Value   interface{}      // 即时数据
	Context *context.Context // 上下文
}

type Logger interface {
	Level(level Level) Logger
	Field(key string, value interface{}) Logger
	Fields(fields context.Fields) Logger
	Err(err error) Logger
	Data(value interface{}) Logger
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}
