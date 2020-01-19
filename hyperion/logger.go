package hyperion

import "github.com/sirupsen/logrus"

const (
	Module = "module"
)

func NewEntry(module string) *logrus.Entry {
	return logger.WithField(Module, module)
}

var logger *logrus.Logger

func init() {
	logger = logrus.New()
}
