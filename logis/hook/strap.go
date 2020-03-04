package hook

import (
	"github.com/laconiz/eros/logis"
)

// 创建一个钩子列表
func NewStrap() *Strap {
	return &Strap{level: logis.INVALID}
}

// 日志钩子列表
type Strap struct {
	level logis.Level
	hooks []*Hook
}

// 添加一个日志钩子
func (strap *Strap) Add(hook *Hook) *Strap {
	if !strap.level.Enable(hook.level) {
		strap.level = hook.level
	}
	strap.hooks = append(strap.hooks, hook)
	return strap
}

// 生成日志入口
func (strap *Strap) Entry() logis.Logger {
	return logis.NewEntry(strap)
}

// 是否有日志钩子需要调用
func (strap *Strap) Enable(level logis.Level) bool {
	return strap.level.Enable(level)
}

// 调用日志钩子
func (strap *Strap) Hook(log *logis.Log) {
	for _, hook := range strap.hooks {
		if hook.level.Enable(log.Level) {
			hook.Hook(log)
		}
	}
}
