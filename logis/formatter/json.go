package formatter

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/utils/json"
)

func Json() *JsonFormatter {
	return &JsonFormatter{timeLayout: DefaultTimeLayout}
}

type JsonLog struct {
	Level   logis.Grade `json:"level"`
	Time    string      `json:"time"`
	Message string      `json:"message"`
	Value   interface{} `json:"data"`
	Context []byte      `json:"context"`
}

type JsonFormatter struct {
	timeLayout  string
	levelLayout map[logis.Level]string
}

// 设置时间序列化规则
func (formatter *JsonFormatter) TimeLayout(layout string) *JsonFormatter {
	formatter.timeLayout = layout
	return formatter
}

// 序列化日志
func (formatter *JsonFormatter) Format(log *logis.Log) ([]byte, error) {
	// 序列化时间戳
	var layout string
	if formatter.timeLayout != "" {
		layout = log.Time.Format(formatter.timeLayout)
	} else {
		layout = log.Time.String()
	}
	// 序列化日志
	return json.Marshal(&JsonLog{
		Level:   log.Level.Grade(),
		Time:    layout,
		Message: log.Message,
		Value:   log.Value,
		Context: log.Context.Json(),
	})
}
