package logis

import "github.com/laconiz/eros/logis/context"

const (
	Module = "module"
)

type Logger interface {
	Field(key string, value interface{}) Logger
	Fields(fields context.Fields) Logger
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

func NewEmpty() Logger {
	return &Empty{}
}

type Empty struct {
}

func (empty *Empty) Field(string, interface{}) Logger {
	return empty
}

func (empty *Empty) Fields(context.Fields) Logger {
	return empty
}

func (empty *Empty) Data(interface{}) Logger {
	return empty
}

func (empty *Empty) Debug(...interface{}) {

}

func (empty *Empty) Debugf(string, ...interface{}) {

}

func (empty *Empty) Info(...interface{}) {

}

func (empty *Empty) Infof(string, ...interface{}) {

}

func (empty *Empty) Warn(...interface{}) {

}

func (empty *Empty) Warnf(string, ...interface{}) {

}

func (empty *Empty) Error(...interface{}) {

}

func (empty *Empty) Errorf(string, ...interface{}) {

}

func (empty *Empty) Fatal(...interface{}) {

}

func (empty *Empty) Fatalf(string, ...interface{}) {

}
