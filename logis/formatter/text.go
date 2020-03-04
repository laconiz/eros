package formatter

import (
	"bytes"
	"fmt"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/utils/json"
)

func Text() *TextFormatter {
	return &TextFormatter{
		timeLayout: DefaultTimeLayout,
		levelLayout: map[logis.Level]string{
			logis.DEBUG: "[DEBUG]",
			logis.INFO:  "[INFO ]",
			logis.WARN:  "[WARN ]",
			logis.ERROR: "[ERROR]",
			logis.FATAL: "[FATAL]",
		},
	}
}

type TextFormatter struct {
	timeLayout  string
	levelLayout map[logis.Level]string
}

// 设置时间序列化规则
func (f *TextFormatter) TimeLayout(layout string) *TextFormatter {
	f.timeLayout = layout
	return f
}

// 设置消息等级文本
func (f *TextFormatter) LevelLayout(strings map[logis.Level]string) *TextFormatter {
	if strings != nil {
		f.levelLayout = strings
	}
	return f
}

// 序列化日志
func (f *TextFormatter) Format(log *logis.Log) ([]byte, error) {

	var buf bytes.Buffer

	// level
	buf.WriteString(f.levelLayout[log.Level])

	// time
	buf.WriteByte(' ')
	buf.WriteString(log.Time.Format(f.timeLayout))

	// context
	if context := log.Context.Json(); len(context) > 0 {
		buf.WriteString(" context:")
		buf.Write(context)
	}

	// message
	buf.WriteString(" $ ")
	buf.WriteString(log.Message)

	// value
	var value string
	if log.Value != nil {
		if raw, err := json.Marshal(log.Value); err != nil {
			value = string(raw)
		} else {
			value = fmt.Sprintf("%#v", log.Value)
		}
	}

	if log.Error != "" && value != "" {

		buf.WriteString(": error=")
		buf.WriteString(log.Error)
		buf.WriteString(" value=")
		buf.WriteString(value)

	} else if log.Error != "" {

		buf.WriteString(": ")
		buf.WriteString(log.Error)

	} else if value != "" {

		buf.WriteString(": ")
		buf.WriteString(value)
	}

	buf.WriteByte('\n')
	return buf.Bytes(), nil
}
