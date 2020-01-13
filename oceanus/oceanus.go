package oceanus

import (
	"github.com/laconiz/eros/log"
)

type Process interface {
	Run()
	// NewThread(NodeType, NodeKey, interface{}) error
}

type Thread interface {
	Call(interface{}) (interface{}, error)
	Stop()
}

var logger = log.Std("oceanus")
