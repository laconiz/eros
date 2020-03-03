package hook

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/logis/formatter"
	"io"
)

type EmptyHook interface {
	Set(writer *Writer) *Hook
	Add(level logis.Level, writer io.Writer) *Hook
}

// 生成一个日志钩子
func NewHook(formatter formatter.Formatter) EmptyHook {
	return &Hook{level: logis.INVALID, formatter: formatter}
}

// 日志钩子
type Hook struct {
	level     logis.Level
	formatter formatter.Formatter
	writers   []*Writer
}

// 调用钩子
func (hook *Hook) Hook(log *logis.Log) {
	// 序列化日志
	raw, err := hook.formatter.Format(log)
	if err != nil {
		formatError(log, err)
		return
	}
	// 写日志
	for _, writer := range hook.writers {
		writer.Write(log.Level, raw)
	}
}

// 添加一个日志写入器
func (hook *Hook) Set(writer *Writer) *Hook {
	if writer != nil && writer.level.Valid() {
		if !hook.level.Enable(writer.level) {
			hook.level = writer.level
		}
		hook.writers = append(hook.writers, writer)
	}
	return hook
}

// 添加一个日志写入器
func (hook *Hook) Add(level logis.Level, writer io.Writer) *Hook {
	return hook.Set(&Writer{level: level, writer: writer})
}

// 将钩子构造成日志钩子列表
func (hook *Hook) Strap() *Strap {
	return &Strap{level: hook.level, hooks: []*Hook{hook}}
}

// 生成日志入口
func (hook *Hook) Entry() logis.Logger {
	return hook.Strap().Entry()
}
