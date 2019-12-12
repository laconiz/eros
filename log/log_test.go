package log

import (
	"testing"
)

func TestLog(t *testing.T) {

	log := Std("test")

	log.Info(1, 1.1, "hello")
	log.Infof("%v", map[int]string{1: "a", 2: "b"})
	log.Error("error stack:")
	log.Info("hello world")
}
