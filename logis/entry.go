package logis

import (
	"fmt"
	"time"
)

type Log struct {
	*Entry
	Level   Level     `json:"level"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

// ---------------------------------------------------------------------------------------------------------------------

type Fields map[string]interface{}

// 日志入口
type Entry struct {
	Context    *Context   `json:"-"`
	ContextRaw ContextRaw `json:"field,omitempty"`
	Strap      *Strap     `json:"-"`
}

func (e *Entry) Field(key string, value interface{}) Logger {
	return &Entry{Context: NewContext(e.Context).Field(key, value), Strap: e.Strap}
}

func (e *Entry) Fields(fields Fields) Logger {
	return &Entry{Context: NewContext(e.Context).Fields(fields), Strap: e.Strap}
}

func (e *Entry) Debug(args ...interface{}) {
	e.Log(DEBUG, fmt.Sprint(args...))
}

func (e *Entry) Debugf(format string, args ...interface{}) {
	e.Log(DEBUG, fmt.Sprintf(format, args...))
}

func (e *Entry) Info(args ...interface{}) {
	e.Log(INFO, fmt.Sprint(args...))
}

func (e *Entry) Infof(format string, args ...interface{}) {
	e.Log(INFO, fmt.Sprintf(format, args...))
}

func (e *Entry) Warn(args ...interface{}) {
	e.Log(WARN, fmt.Sprint(args...))
}

func (e *Entry) Warnf(format string, args ...interface{}) {
	e.Log(WARN, fmt.Sprintf(format, args...))
}

func (e *Entry) Error(args ...interface{}) {
	e.Log(ERROR, fmt.Sprint(args...))
}

func (e *Entry) Errorf(format string, args ...interface{}) {
	e.Log(ERROR, fmt.Sprintf(format, args...))
}

func (e *Entry) Fatal(args ...interface{}) {
	e.Log(FATAL, fmt.Sprint(args...))
}

func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.Log(FATAL, fmt.Sprintf(format, args...))
}

func (e *Entry) Log(level Level, args ...interface{}) {
	if e.Strap.Enable(level) {
		e.log(level, fmt.Sprint(args...))
	}
}

func (e *Entry) Logf(level Level, format string, args ...interface{}) {
	if e.Strap.Enable(level) {
		e.log(level, fmt.Sprintf(format, args...))
	}
}

func (e *Entry) log(level Level, message string) {
	e.Strap.Hook(&Log{Entry: e, Level: level, Time: time.Now(), Message: message})
}
