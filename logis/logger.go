package logis

const (
	Module = "module"
)

type Logger interface {
	Field(key string, value interface{}) Logger
	Fields(fields Fields) Logger
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
	return &empty{}
}

type empty struct {
}

func (e *empty) Field(string, interface{}) Logger {
	return e
}

func (e *empty) Fields(Fields) Logger {
	return e
}

func (e *empty) Debug(...interface{}) {

}

func (e *empty) Debugf(string, ...interface{}) {

}

func (e *empty) Info(...interface{}) {

}

func (e *empty) Infof(string, ...interface{}) {

}

func (e *empty) Warn(...interface{}) {

}

func (e *empty) Warnf(string, ...interface{}) {

}

func (e *empty) Error(...interface{}) {

}

func (e *empty) Errorf(string, ...interface{}) {

}

func (e *empty) Fatal(...interface{}) {

}

func (e *empty) Fatalf(string, ...interface{}) {

}
