package oceanus

import (
	"github.com/laconiz/eros/log"
)

type Process interface {
	Run()
}

var logger = log.Std("oceanus")
