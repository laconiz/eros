package formatter

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/utils/json"
)

func Json() *JsonFormatter {
	return &JsonFormatter{timeLayout: DefaultTimeLayout}
}

type jsonLog struct {
	Level   logis.Grade `json:"level"`
	Error   string      `json:"error"`
	Time    string      `json:"time"`
	Message string      `json:"message"`
	Value   interface{} `json:"data"`
	Context []byte      `json:"context"`
}

type JsonFormatter struct {
	timeLayout string
}

// 设置时间序列化规则
func (formatter *JsonFormatter) TimeLayout(layout string) *JsonFormatter {
	formatter.timeLayout = layout
	return formatter
}

// 序列化日志
func (formatter *JsonFormatter) Format(log *logis.Log) ([]byte, error) {

	var layout string
	if formatter.timeLayout != "" {
		layout = log.Time.Format(formatter.timeLayout)
	} else {
		layout = log.Time.String()
	}

	return json.Marshal(&jsonLog{
		Level:   log.Level.Grade(),
		Error:   log.Error,
		Time:    layout,
		Message: log.Message,
		Value:   log.Value,
		Context: log.Context.Json(),
	})
}
