package logis

import (
	"fmt"
	"io"
	"os"
)

// 生成一个日志钩子列表
func NewHooks(formatter Formatter, writers map[Level]io.Writer) *Hooks {
	hook := NewHook(formatter)
	for level, writer := range writers {
		hook.AddWriter(level, writer)
	}
	return hook.Hooks()
}

// 日志钩子列表
type Hooks struct {
	level Level
	hooks []*Hook
}

// 添加一个日志钩子
func (h *Hooks) AddHook(hook *Hook) *Hooks {
	h.hooks = append(h.hooks, hook)
	h.level = MinLevel(h.level, hook.level)
	return h
}

// 生成日志入口
func (h *Hooks) Entry(key string, value interface{}) *Entry {
	return &Entry{Data: Fields{key: value}, Hooks: h}
}

// 是否有日志钩子需要调用
func (h *Hooks) enable(level Level) bool {
	return h.level <= level
}

// 调用日志钩子
func (h *Hooks) fire(log *Log) {
	for _, hook := range h.hooks {
		if hook.enable(log.Level) {
			hook.fire(log)
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
func (h *Hook) enable(level Level) bool {
	return h.level <= level
}

// 调用钩子
func (h *Hook) fire(log *Log) {

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
func (h *Hook) Hooks() *Hooks {
	return &Hooks{level: h.level, hooks: []*Hook{h}}
}

// 生成日志入口
func (h *Hook) Entry(key string, value interface{}) *Entry {
	return &Entry{Data: Fields{key: value}, Hooks: h.Hooks()}
}

// ---------------------------------------------------------------------------------------------------------------------

// 日志写入器
type Writer struct {
	level  Level
	writer io.Writer
}
