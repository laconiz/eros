package hook

import (
	"fmt"
	"github.com/laconiz/eros/logis"
	"os"
)

func formatError(log *logis.Log, err error) {
	os.Stderr.WriteString(fmt.Sprintf("format log[%+v] error: %v", log, err))
}

func writeError(raw []byte, err error) {
	os.Stderr.WriteString(fmt.Sprintf("write log[%s] error: %v", string(raw), err))
}
