package logis

import (
	"fmt"
	"io"
	"os"
)

// 日志钩子列表
type Strap struct {
	level Level
	hooks []*Hook
}

// 添加一个日志钩子
func (strap *Strap) AddHook(hook *Hook) *Strap {
	strap.hooks = append(strap.hooks, hook)
	strap.level = MinLevel(strap.level, hook.level)
	return strap
}

// 生成日志入口
func (strap *Strap) Entry() *Entry {
	return &Entry{Context: NewContext(nil), Strap: strap}
}

// 是否有日志钩子需要调用
func (strap *Strap) Enable(level Level) bool {
	return strap.level <= level
}

// 调用日志钩子
func (strap *Strap) Hook(log *Log) {
	for _, hook := range strap.hooks {
		if hook.Enable(log.Level) {
			hook.Hook(log)
		}
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// 生成一个日志钩子
func NewHook(formatter Formatter) *Hook {
	return &Hook{level: INVALID, formatter: formatter, writers: []*Writer{}}
}

// 日志钩子
type Hook struct {
	level     Level
	formatter Formatter
	writers   []*Writer
}

// 日志钩子是否需要调用
func (h *Hook) Enable(level Level) bool {
	return h.level <= level
}

// 调用钩子
func (h *Hook) Hook(log *Log) {

	bytes, err := h.formatter.Format(log)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("format entry[%+v] error: %v", log, err))
		return
	}

	for _, writer := range h.writers {
		if log.Level >= writer.level {
			if _, err := writer.writer.Write(bytes); err != nil {
				os.Stderr.WriteString(fmt.Sprintf("write bytes[%s] error: %v", string(bytes), err))
			}
		}
	}
}

// 添加一个日志写入器
func (h *Hook) AddWriter(level Level, writer io.Writer) *Hook {
	if writer != nil && level.Valid() {
		h.level = MinLevel(h.level, level)
		h.writers = append(h.writers, &Writer{level: level, writer: writer})
	}
	return h
}

// 将钩子构造成日志钩子列表
func (h *Hook) Strap() *Strap {
	return &Strap{level: h.level, hooks: []*Hook{h}}
}

// 生成日志入口
func (h *Hook) Entry() *Entry {
	return h.Strap().Entry()
}

// ---------------------------------------------------------------------------------------------------------------------

// 日志写入器
type Writer struct {
	level  Level
	writer io.Writer
}
