package logis

import (
	"fmt"
	"github.com/laconiz/eros/logis/context"
	"time"
)

func NewEntry(strap Strap) *Entry {
	return &Entry{context: context.New(nil), strap: strap}
}

type Strap interface {
	Enable(Level) bool
	Hook(log *Log)
}

// 日志入口
type Entry struct {
	level   Level            // 最低等级
	err     error            // 错误信息
	value   interface{}      // 即时数据
	context *context.Context // 上下文
	strap   Strap            // 钩子列表
}

func (entry *Entry) Level(level Level) Logger {
	copy := *entry
	copy.level = level
	return &copy
}

func (entry *Entry) Field(key string, value interface{}) Logger {
	return entry.Fields(context.Fields{key: value})
}

func (entry *Entry) Fields(fields context.Fields) Logger {
	copy := *entry
	copy.context = context.New(entry.context).Fields(fields)
	return &copy
}

func (entry *Entry) Err(err error) Logger {
	copy := *entry
	copy.err = err
	return &copy
}

func (entry *Entry) Data(value interface{}) Logger {
	copy := *entry
	copy.value = value
	return &copy
}

func (entry *Entry) Debug(args ...interface{}) {
	entry.Log(DEBUG, args...)
}

func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.Logf(DEBUG, format, args...)
}

func (entry *Entry) Info(args ...interface{}) {
	entry.Log(INFO, args...)
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.Logf(INFO, format, args...)
}

func (entry *Entry) Warn(args ...interface{}) {
	entry.Log(WARN, args...)
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	entry.Logf(WARN, format, args...)
}

func (entry *Entry) Error(args ...interface{}) {
	entry.Log(ERROR, args...)
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.Logf(ERROR, format, args...)
}

func (entry *Entry) Fatal(args ...interface{}) {
	entry.Log(FATAL, args...)
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.Logf(FATAL, format, args...)
}

func (entry *Entry) Log(level Level, args ...interface{}) {

	if !entry.level.Enable(level) {
		return
	}

	if !entry.strap.Enable(level) {
		return
	}

	entry.log(level, fmt.Sprint(args...))
}

func (entry *Entry) Logf(level Level, format string, args ...interface{}) {

	if !entry.level.Enable(level) {
		return
	}

	if !entry.strap.Enable(level) {
		return
	}

	entry.log(level, fmt.Sprintf(format, args...))
}

func (entry *Entry) log(level Level, message string) {

	log := &Log{
		Level:   level,
		Message: message,
		Time:    time.Now(),
		Value:   entry.value,
		Context: entry.context,
	}

	if entry.err != nil {
		log.Error = entry.err.Error()
	}

	entry.strap.Hook(log)
}
