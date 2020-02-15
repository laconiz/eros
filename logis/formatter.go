package logis

import (
	"fmt"
	"github.com/laconiz/eros/utils/json"
)

type Formatter interface {
	Format(*Log) ([]byte, error)
}

func parseData(log *Log) error {
	if log.RawData == nil && len(log.Data) > 0 {
		raw, err := json.Marshal(log.Data)
		if err != nil {
			return err
		}
		log.RawData = raw
	}
	return nil
}

type JsonFormatter struct {
}

func (f *JsonFormatter) Format(log *Log) ([]byte, error) {
	if err := log.ParseField(); err != nil {
		return nil, fmt.Errorf("parse fields error: %v", err)
	}
	return json.Marshal(log)
}
