package hyperion

import "github.com/sirupsen/logrus"

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
}
