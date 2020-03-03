package logis

import (
	"fmt"
	"github.com/laconiz/eros/logis/context"
	"time"
)

type Strap interface {
	Enable(Level) bool
	Hook(log *Log)
}

// 日志入口
type Entry struct {
	Value   interface{}      // 即时数据
	Context *context.Context // 上下文
	Strap   Strap            // 钩子列表
}

func (entry *Entry) Field(key string, value interface{}) Logger {
	return entry.Fields(context.Fields{key: value})
}

func (entry *Entry) Fields(fields context.Fields) Logger {
	return &Entry{
		Value:   entry.Value,
		Context: context.New(entry.Context).Fields(fields),
		Strap:   entry.Strap,
	}
}

func (entry *Entry) Data(value interface{}) Logger {
	return &Entry{Value: value, Context: entry.Context, Strap: entry.Strap}
}

func (entry *Entry) Debug(args ...interface{}) {
	entry.Log(DEBUG, fmt.Sprint(args...))
}

func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.Log(DEBUG, fmt.Sprintf(format, args...))
}

func (entry *Entry) Info(args ...interface{}) {
	entry.Log(INFO, fmt.Sprint(args...))
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.Log(INFO, fmt.Sprintf(format, args...))
}

func (entry *Entry) Warn(args ...interface{}) {
	entry.Log(WARN, fmt.Sprint(args...))
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	entry.Log(WARN, fmt.Sprintf(format, args...))
}

func (entry *Entry) Error(args ...interface{}) {
	entry.Log(ERROR, fmt.Sprint(args...))
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.Log(ERROR, fmt.Sprintf(format, args...))
}

func (entry *Entry) Fatal(args ...interface{}) {
	entry.Log(FATAL, fmt.Sprint(args...))
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.Log(FATAL, fmt.Sprintf(format, args...))
}

func (entry *Entry) Log(level Level, args ...interface{}) {
	if entry.Strap.Enable(level) {
		entry.log(level, fmt.Sprint(args...))
	}
}

func (entry *Entry) Logf(level Level, format string, args ...interface{}) {
	if entry.Strap.Enable(level) {
		entry.log(level, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) log(level Level, message string) {
	entry.Strap.Hook(&Log{
		Level:   level,
		Message: message,
		Time:    time.Now(),
		Value:   entry.Value,
		Context: entry.Context,
	})
}
