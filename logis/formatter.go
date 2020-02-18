package logis

import (
	"bytes"
	"github.com/laconiz/eros/utils/json"
)

type Formatter interface {
	Format(*Log) ([]byte, error)
}

// ---------------------------------------------------------------------------------------------------------------------
// JSON格式化

type JsonFormatter struct {
}

// 序列化日志
func (f *JsonFormatter) Format(log *Log) ([]byte, error) {
	if log.ContextRaw == nil {
		log.ContextRaw = ParseContext(log)
	}
	return json.Marshal(log)
}

// ---------------------------------------------------------------------------------------------------------------------
// 文本格式化

func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		timeLayout: "2006-01-02 15:04:05.000",
		levelString: map[Level]string{
			DEBUG: "[DEBUG]",
			INFO:  "[INFO ]",
			WARN:  "[WARN ]",
			ERROR: "[ERROR]",
			FATAL: "[FATAL]",
		},
	}
}

type TextFormatter struct {
	timeLayout  string
	levelString map[Level]string
}

// 设置时间序列化规则
func (f *TextFormatter) SetTimeLayout(layout string) *TextFormatter {
	f.timeLayout = layout
	return f
}

// 设置消息等级文本
func (f *TextFormatter) SetLevelString(strings map[Level]string) *TextFormatter {
	if strings != nil {
		f.levelString = strings
	}
	return f
}

// 序列化日志
func (f *TextFormatter) Format(log *Log) ([]byte, error) {

	var buf bytes.Buffer

	buf.WriteString(f.levelString[log.Level])

	buf.WriteByte(' ')
	buf.WriteString(log.Time.Format(f.timeLayout))

	buf.WriteString(" context:")
	buf.Write(ParseContext(log))

	buf.WriteString(" $ ")
	buf.WriteString(log.Message)
	buf.WriteByte('\n')

	return buf.Bytes(), nil
}

// ---------------------------------------------------------------------------------------------------------------------

// 序列化字段
func ParseContext(log *Log) []byte {
	if log.ContextRaw == nil {
		log.ContextRaw = log.Context.Json()
	}
	return log.ContextRaw
}
