package hook

import (
	"github.com/laconiz/eros/logis"
	"io"
)

func NewWriter(level logis.Level, writer io.Writer) *Writer {
	return &Writer{level: level, writer: writer}
}

// 日志写入器
type Writer struct {
	level  logis.Level
	writer io.Writer
}

func (writer *Writer) Write(level logis.Level, raw []byte) {
	// 不需要写入的等级
	if !writer.level.Enable(level) {
		return
	}
	// 写入日志
	if _, err := writer.writer.Write(raw); err != nil {
		writeError(raw, err)
	}
}
