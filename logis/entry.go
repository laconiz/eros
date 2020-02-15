package logis

import (
	"fmt"
	"github.com/laconiz/eros/utils/json"
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
	Data    Fields          `json:"-"`
	RawData json.RawMessage `json:"field,omitempty"`
	Hooks   *Hooks          `json:"-"`
}

func (e *Entry) copy() *Entry {
	n := &Entry{Data: Fields{}, RawData: nil, Hooks: e.Hooks}
	for key, value := range e.Data {
		n.Data[key] = value
	}
	return n
}

func (e *Entry) Field(key string, value interface{}) Logger {
	n := e.copy()
	n.Data[key] = value
	n.RawData = nil
	return n
}

func (e *Entry) Fields(fields Fields) Logger {
	n := e.copy()
	for key, value := range fields {
		n.Data[key] = value
	}
	n.RawData = nil
	return n
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
	if e.Hooks.enable(level) {
		e.log(level, fmt.Sprint(args...))
	}
}

func (e *Entry) Logf(level Level, format string, args ...interface{}) {
	if e.Hooks.enable(level) {
		e.log(level, fmt.Sprintf(format, args...))
	}
}

func (e *Entry) log(level Level, message string) {
	e.Hooks.fire(&Log{Entry: e, Level: level, Time: time.Now(), Message: message})
}

func (e *Entry) ParseField() error {
	if e.RawData == nil {
		raw, err := json.Marshal(e.Data)
		if err != nil {
			return err
		}
		e.RawData = raw
	}
	return nil
}
