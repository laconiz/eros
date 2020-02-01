package hyperion

import (
	"bytes"

	"github.com/laconiz/eros/utils/json"
)

type Formatter struct {
}

func (f *Formatter) Format(entry *Entry) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(entry.Level.String())
	buf.WriteByte(' ')
	buf.WriteString(entry.Time.Format(timeLayout))
	raw, _ := json.Marshal(entry.Data)
	buf.WriteByte(' ')
	buf.Write(raw)
	buf.WriteByte(' ')
	buf.WriteString(entry.Message)
	if len(entry.Message) == 0 || entry.Message[len(entry.Message)-1] != '\n' {
		buf.WriteByte('\n')
	}
	// buf.WriteString("\n")
	return buf.Bytes(), nil
}

const timeLayout = "2006-01-02 15:04:05.999"
