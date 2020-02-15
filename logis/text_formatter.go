package logis

import (
	"bytes"
	"fmt"
)

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

func (f *TextFormatter) SetTimeLayout(layout string) *TextFormatter {
	f.timeLayout = layout
	return f
}

func (f *TextFormatter) SetLevelString(strings map[Level]string) *TextFormatter {
	if strings != nil {
		f.levelString = strings
	}
	return f
}

func (f *TextFormatter) Format(log *Log) ([]byte, error) {

	if err := parseData(log); err != nil {
		return nil, fmt.Errorf("parse fields error: %v", err)
	}

	var buf bytes.Buffer

	buf.WriteString(f.levelString[log.Level])

	buf.WriteByte(' ')
	buf.WriteString(log.Time.Format(f.timeLayout))

	if log.RawData != nil {
		buf.WriteByte(' ')
		buf.Write(log.RawData)
	}

	buf.WriteString(" $ ")
	buf.WriteString(log.Message)
	buf.WriteByte('\n')

	return buf.Bytes(), nil
}
