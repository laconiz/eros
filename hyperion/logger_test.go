package hyperion

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestFormatter_Format(t *testing.T) {
	logrus.New().Info("hello")
	ln := logrus.New()
	ln.SetFormatter(&Formatter{})
	ln.Info("world")
}
