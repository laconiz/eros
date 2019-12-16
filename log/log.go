package log

import (
	"bytes"
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

type File interface {
	Write(b []byte) (n int, err error)
}

type Level uint8

const (
	Debug Level = iota
	Info
	Warn
	Error
	Fatal
)

var levelName = map[Level][]byte{
	Debug: []byte("[DEBUG] "),
	Info:  []byte("[INFO ] "),
	Warn:  []byte("[WARN ] "),
	Error: []byte("[ERROR] "),
	Fatal: []byte("[FATAL] "),
}

const (
	timeFormatter = "2006-01-02 15:04:05.000"
)

type Logger struct {
	name  []byte
	level Level
	file  File
}

func (l *Logger) Logf(lv Level, format string, args ...interface{}) {

	// 不输出级别不足的日志
	if lv < l.level {
		return
	}

	var buf bytes.Buffer
	buf.Write(levelName[lv])
	buf.WriteString(time.Now().Format(timeFormatter))
	buf.Write(l.name)

	// 序列化
	var raw string
	if format == "" {
		raw = fmt.Sprintln(args...)
	} else {
		raw = fmt.Sprintf(format, args...)
	}
	buf.WriteString(raw)
	if len(raw) > 0 && raw[len(raw)-1] != '\n' {
		buf.WriteByte('\n')
	}

	if lv >= Error {
		buf.Write(debug.Stack())
	}

	l.file.Write(buf.Bytes())
}

func (l *Logger) Log(level Level, args ...interface{}) {
	l.Logf(level, "", args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.Log(Debug, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logf(Debug, format, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.Log(Info, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logf(Info, format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.Log(Warn, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logf(Warn, format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Log(Error, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logf(Error, format, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.Log(Fatal, args...)
	os.Exit(-1)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Logf(Fatal, format, args...)
	os.Exit(-1)
}

func (l *Logger) SetLevel(lv Level) {
	l.level = lv
}

func New(name string, level Level, file File) *Logger {
	return &Logger{
		name:  []byte(fmt.Sprintf(" %s $ ", name)),
		level: level,
		file:  file,
	}
}

func Std(name string) *Logger {
	return New(name, Debug, os.Stdout)
}
