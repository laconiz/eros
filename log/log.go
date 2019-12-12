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

type Log struct {
	name  []byte
	level Level
	file  File
}

func (l *Log) Logf(lv Level, format string, args ...interface{}) {

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

func (l *Log) Log(level Level, args ...interface{}) {
	l.Logf(level, "", args...)
}

func (l *Log) Debug(args ...interface{}) {
	l.Log(Debug, args...)
}

func (l *Log) Debugf(format string, args ...interface{}) {
	l.Logf(Debug, format, args...)
}

func (l *Log) Info(args ...interface{}) {
	l.Log(Info, args...)
}

func (l *Log) Infof(format string, args ...interface{}) {
	l.Logf(Info, format, args...)
}

func (l *Log) Warn(args ...interface{}) {
	l.Log(Warn, args...)
}

func (l *Log) Warnf(format string, args ...interface{}) {
	l.Logf(Warn, format, args...)
}

func (l *Log) Error(args ...interface{}) {
	l.Log(Error, args...)
}

func (l *Log) Errorf(format string, args ...interface{}) {
	l.Logf(Error, format, args...)
}

func (l *Log) Fatal(args ...interface{}) {
	l.Log(Fatal, args...)
	os.Exit(-1)
}

func (l *Log) Fatalf(format string, args ...interface{}) {
	l.Logf(Fatal, format, args...)
	os.Exit(-1)
}

func New(name string, level Level, file File) *Log {
	return &Log{
		name:  []byte(fmt.Sprintf(" %s $ ", name)),
		level: level,
		file:  file,
	}
}

func Std(name string) *Log {
	return New(name, Debug, os.Stdout)
}
