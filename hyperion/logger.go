package hyperion

import (
	"github.com/sirupsen/logrus"
	"os"
)

const (
	Module = "module"
)

type Entry = logrus.Entry

func NewEntry(module string) *Entry {
	return logger.WithField(Module, module)
}

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)
	logger.SetFormatter(&Formatter{})
}
